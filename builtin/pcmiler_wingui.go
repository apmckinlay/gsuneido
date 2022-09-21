// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"github.com/apmckinlay/gsuneido/builtin/goc"
	"github.com/apmckinlay/gsuneido/builtin/heap"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/str"
	"golang.org/x/exp/maps"
	"golang.org/x/sys/windows"
)

// Interface to PC*Miler DLL

type suPcmsrv64 struct {
	staticClass[suPcmsrv64]
}

var _ = Global.Builtin("Pcmsrv64", &suPcmsrv64{})

func (*suPcmsrv64) String() string {
	return "Pcmsrv64 /* builtin class */"
}

func (*suPcmsrv64) Lookup(_ *Thread, method string) Callable {
	return pcmMethods[method]
}

var pcmMethods = methods()

var _ = staticMethod(pcm_Members, "()")

func pcm_Members() Value {
	return SuObjectOfStrs(maps.Keys(pcmMethods))
}

var pcmsrv = windows.NewLazyDLL("pcmsrv64.dll")

// dll long pcmsrv32:PCMSAbout(string which, buffer buf, long bufsize)
var pcmsAbout = pcmsrv.NewProc("PCMSAbout")
var _ = staticMethod(pcm_Version, "()")

func pcm_Version() Value {
	if pcmsrv.Load() != nil {
		return False
	}
	defer heap.FreeTo(heap.CurSize())
	const buflen = 200
	buf := heap.Alloc(buflen)
	goc.Syscall3(pcmsAbout.Addr(),
		uintptr(heap.CopyStr("ProductVersion")),
		uintptr(buf),
		buflen)
	s := heap.GetStrZ(buf, int(buflen))
	return SuStr(str.BeforeFirst(s, "."))
}

// dll long pcmsrv32:PCMSAddStop(long tripId, [in] string stop)
var pcmsAddStop = pcmsrv.NewProc("PCMSAddStop")
var _ = staticMethod(pcm_PCMAddStop, "(tripId, stop)")

func pcm_PCMAddStop(a, b Value) Value {
	defer heap.FreeTo(heap.CurSize())
	rtn := goc.Syscall2(pcmsAddStop.Addr(),
		intArg(a),
		uintptr(stringArg(b)))
	return int32Ret(rtn)
}

// dll long pcmsrv32:PCMSCalculate(long tripId)
var pcmsCalculate = pcmsrv.NewProc("PCMSCalculate")
var _ = staticMethod(pcm_PCMCalculate, "(tripId)")

func pcm_PCMCalculate(a Value) Value {
	rtn := goc.Syscall1(pcmsCalculate.Addr(),
		intArg(a))
	return int32Ret(rtn)
}

// dll long pcmsrv32:PCMSCalcTrip(long tripId, [in] string orig, [in] string dest)
var pcmsCalcTrip = pcmsrv.NewProc("PCMSCalcTrip")
var _ = staticMethod(pcm_PCMCalcTrip, "(tripId, orig, dest)")

func pcm_PCMCalcTrip(a, b, c Value) Value {
	defer heap.FreeTo(heap.CurSize())
	rtn := goc.Syscall3(pcmsCalcTrip.Addr(),
		intArg(a),
		uintptr(stringArg(b)),
		uintptr(stringArg(c)))
	return int32Ret(rtn)
}

// dll long pcmsrv32:PCMSCloseServer(long server)
var pcmsCloseServer = pcmsrv.NewProc("PCMSCloseServer")

var _ = staticMethod(pcm_PCMCloseServer, "(server)")

func pcm_PCMCloseServer(a Value) Value {
	rtn := goc.Syscall1(pcmsCloseServer.Addr(),
		intArg(a))
	return int32Ret(rtn)
}

// dll long pcmsrv32:PCMSDeleteTrip(long tripId)
var pcmsDeleteTrip = pcmsrv.NewProc("PCMSDeleteTrip")
var _ = staticMethod(pcm_PCMDeleteTrip, "(tripId)")

func pcm_PCMDeleteTrip(a Value) Value {
	rtn := goc.Syscall1(pcmsDeleteTrip.Addr(),
		intArg(a))
	return int32Ret(rtn)
}

// dll long pcmsrv32:PCMSGetMatch(long tripId,
// long index, string buffer, long bufLen)
var pcmsGetMatch = pcmsrv.NewProc("PCMSGetMatch")
var _ = staticMethod(pcm_GetMatch, "(tripId, i)")

func pcm_GetMatch(a, b Value) Value {
	defer heap.FreeTo(heap.CurSize())
	const buflen = 200
	buf := heap.Alloc(buflen)
	goc.Syscall4(pcmsGetMatch.Addr(),
		intArg(a),
		intArg(b),
		uintptr(buf),
		buflen)
	return SuStr(heap.GetStrZ(buf, int(buflen)))
}

// dll long pcmsrv32:PCMSGetRpt(long tripId, long rpt, string buffer, long bufLen)
var pcmsGetRpt = pcmsrv.NewProc("PCMSGetRpt")

// dll long pcmsrv32:PCMSNumRptBytes(long tripId, long rpt)
var pcmsNumRptBytes = pcmsrv.NewProc("PCMSNumRptBytes")

const PCM_RPT_STATE = 1

var _ = staticMethod(pcm_GetRpt, "(tripId, rpt)")

func pcm_GetRpt(a, b Value) Value {
	size := goc.Syscall2(pcmsNumRptBytes.Addr(),
		intArg(a),
		PCM_RPT_STATE)
	if size <= 0 {
		return False
	}
	defer heap.FreeTo(heap.CurSize())
	buf := heap.Alloc(size)
	goc.Syscall4(pcmsGetRpt.Addr(),
		intArg(a),
		intArg(b),
		uintptr(buf),
		size)
	return SuStr(heap.GetStrZ(buf, int(size)))
}

// dll long pcmsrv32:PCMSNewTrip(long serverId)
var pcmsNewTrip = pcmsrv.NewProc("PCMSNewTrip")
var _ = staticMethod(pcm_PCMNewTrip, "(tripId)")

func pcm_PCMNewTrip(a Value) Value {
	rtn := goc.Syscall1(pcmsNewTrip.Addr(),
		intArg(a))
	return int32Ret(rtn)
}

// dll long pcmsrv32:PCMSNumMatches(long tripId)
var pcmsNumMatches = pcmsrv.NewProc("PCMSNumMatches")
var _ = staticMethod(pcm_PCMNumMatches, "(tripId)")

func pcm_PCMNumMatches(a Value) Value {
	rtn := goc.Syscall1(pcmsNumMatches.Addr(),
		intArg(a))
	return int32Ret(rtn)
}

// dll long pcmsrv32:PCMSOpenServer(long hInstance, long hwnd)
var pcmsOpenServer = pcmsrv.NewProc("PCMSOpenServer")
var _ = staticMethod(pcm_PCMOpenServer, "(tripId, hwnd)")

func pcm_PCMOpenServer(a, b Value) Value {
	rtn := goc.Syscall2(pcmsOpenServer.Addr(),
		intArg(a),
		intArg(b))
	return int32Ret(rtn)
}

// dll long pcmsrv32:PCMSIsValid(long serverId)
var pcmsIsValid = pcmsrv.NewProc("PCMSIsValid")
var _ = staticMethod(pcm_PCMIsValid, "(serverId)")

func pcm_PCMIsValid(a Value) Value {
	rtn := goc.Syscall1(pcmsIsValid.Addr(),
		intArg(a))
	return int32Ret(rtn)
}

// dll void pcmsrv32:PCMSSetBordersOpen(long tripId, bool open)
var pcmsSetBordersOpen = pcmsrv.NewProc("PCMSSetBordersOpen")
var _ = staticMethod(pcm_PCMSetBordersOpen, "(tripId, open)")

func pcm_PCMSetBordersOpen(a, b Value) Value {
	goc.Syscall2(pcmsSetBordersOpen.Addr(),
		intArg(a),
		boolArg(b))
	return nil
}

// dll void pcmsrv32:PCMSSetCalcType(long tripId, long routeType)
var pcmsSetCalcType = pcmsrv.NewProc("PCMSSetCalcType")
var _ = staticMethod(pcm_PCMSetCalcType, "(tripId, routeType)")

func pcm_PCMSetCalcType(a, b Value) Value {
	goc.Syscall2(pcmsSetCalcType.Addr(),
		intArg(a),
		intArg(b))
	return nil
}

// dll long pcmsrv32:PCMSLookup(long tripId, [in] string placeName, long easyMatch)
var pcmsLookup = pcmsrv.NewProc("PCMSLookup")
var _ = staticMethod(pcm_PCMLookup, "(tripId, placeName, easyMatch)")

func pcm_PCMLookup(a, b, c Value) Value {
	defer heap.FreeTo(heap.CurSize())
	rtn := goc.Syscall3(pcmsLookup.Addr(),
		intArg(a),
		uintptr(stringArg(b)),
		intArg(c))
	return int32Ret(rtn)
}

// dll void pcmsrv32:PCMSSetMiles(long tripId)
var pcmsSetMiles = pcmsrv.NewProc("PCMSSetMiles")
var _ = staticMethod(pcm_PCMSetMiles, "(tripId)")

func pcm_PCMSetMiles(a Value) Value {
	rtn := goc.Syscall1(pcmsSetMiles.Addr(),
		intArg(a))
	return int32Ret(rtn)
}

// dll void pcmsrv32:PCMSSetCustomMode(long tripId, bool onOff)
var pcmsSetCustomMode = pcmsrv.NewProc("PCMSSetCustomMode")
var _ = staticMethod(pcm_PCMSetCustomMode, "(tripId, onOff)")

func pcm_PCMSetCustomMode(a, b Value) Value {
	goc.Syscall2(pcmsSetCustomMode.Addr(),
		intArg(a),
		boolArg(b))
	return nil
}

// dll void pcmsrv32:PCMSSetCalcTypeEx(long tripId, long routeType,
// long optFlags, long vehType)
var pcmsSetCalcTypeEx = pcmsrv.NewProc("PCMSSetCalcTypeEx")
var _ = staticMethod(pcm_PCMSetCalcTypeEx, "(tripId, routeType, optFlags, vehType)")

func pcm_PCMSetCalcTypeEx(a, b, c, d Value) Value {
	rtn := goc.Syscall4(pcmsSetCalcTypeEx.Addr(),
		intArg(a),
		intArg(b),
		intArg(c),
		intArg(d))
	return int32Ret(rtn)
}

// dll void pcmsrv32:PCMSSetVehicleType(long tripId, bool onOff)
var pcmsSetVehicleType = pcmsrv.NewProc("PCMSSetVehicleType")
var _ = staticMethod(pcm_PCMSetVehicleType, "(tripId, onOff)")

func pcm_PCMSetVehicleType(a, b Value) Value {
	goc.Syscall2(pcmsSetVehicleType.Addr(),
		intArg(a),
		boolArg(b))
	return nil
}

// dll void pcmsrv32:PCMSSetRouteLevel(long trip, bool UseStreets)
var pcmsSetRouteLevel = pcmsrv.NewProc("PCMSSetRouteLevel")
var _ = staticMethod(pcm_PCMSetRouteLevel, "(tripId, onOff)")

func pcm_PCMSetRouteLevel(a, b Value) Value {
	goc.Syscall2(pcmsSetRouteLevel.Addr(),
		intArg(a),
		boolArg(b))
	return nil
}
