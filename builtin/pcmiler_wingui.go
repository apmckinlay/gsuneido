// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"github.com/apmckinlay/gsuneido/builtin/goc"
	"github.com/apmckinlay/gsuneido/builtin/heap"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/str"
	"golang.org/x/sys/windows"
)

// Interface to PC*Miler DLL

var pcmsrv64 = &SuClass{Name: "Pcmsrv64",
	MemBase: MemBase{Data: map[string]Value{}}}

var _ = Global.Builtin("Pcmsrv64", pcmsrv64)

var pcmsrv = windows.NewLazyDLL("pcmsrv64.dll")

// dll long pcmsrv32:PCMSAbout(string which, buffer buf, long bufsize)
var pcmsAbout = pcmsrv.NewProc("PCMSAbout")

var _ = meth0("Version()",
	func() Value {
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
	})

// dll long pcmsrv32:PCMSAddStop(long tripId, [in] string stop)
var pcmsAddStop = pcmsrv.NewProc("PCMSAddStop")

var _ = meth2("PCMAddStop(tripId, stop)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		rtn := goc.Syscall2(pcmsAddStop.Addr(),
			intArg(a),
			uintptr(stringArg(b)))
		return int32Ret(rtn)
	})

// dll long pcmsrv32:PCMSCalculate(long tripId)
var pcmsCalculate = pcmsrv.NewProc("PCMSCalculate")

var _ = meth1("PCMCalculate(tripId)",
	func(a Value) Value {
		rtn := goc.Syscall1(pcmsCalculate.Addr(),
			intArg(a))
		return int32Ret(rtn)
	})

// dll long pcmsrv32:PCMSCalcTrip(long tripId, [in] string orig, [in] string dest)
var pcmsCalcTrip = pcmsrv.NewProc("PCMSCalcTrip")

var _ = meth3("PCMCalcTrip(tripId, orig, dest)",
	func(a, b, c Value) Value {
		defer heap.FreeTo(heap.CurSize())
		rtn := goc.Syscall3(pcmsCalcTrip.Addr(),
			intArg(a),
			uintptr(stringArg(b)),
			uintptr(stringArg(c)))
		return int32Ret(rtn)
	})

// dll long pcmsrv32:PCMSCloseServer(long server)
var pcmsCloseServer = pcmsrv.NewProc("PCMSCloseServer")

var _ = meth1("PCMCloseServer(server)",
	func(a Value) Value {
		rtn := goc.Syscall1(pcmsCloseServer.Addr(),
			intArg(a))
		return int32Ret(rtn)
	})

// dll long pcmsrv32:PCMSDeleteTrip(long tripId)
var pcmsDeleteTrip = pcmsrv.NewProc("PCMSDeleteTrip")

var _ = meth1("PCMDeleteTrip(tripId)",
	func(a Value) Value {
		rtn := goc.Syscall1(pcmsDeleteTrip.Addr(),
			intArg(a))
		return int32Ret(rtn)
	})

// dll long pcmsrv32:PCMSGetMatch(long tripId,
//
//	long index, string buffer, long bufLen)
var pcmsGetMatch = pcmsrv.NewProc("PCMSGetMatch")

var _ = meth2("GetMatch(tripId, i)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		const buflen = 200
		buf := heap.Alloc(buflen)
		goc.Syscall4(pcmsGetMatch.Addr(),
			intArg(a),
			intArg(b),
			uintptr(buf),
			buflen)
		return SuStr(heap.GetStrZ(buf, int(buflen)))
	})

// dll long pcmsrv32:PCMSGetRpt(long tripId, long rpt, string buffer, long bufLen)
var pcmsGetRpt = pcmsrv.NewProc("PCMSGetRpt")

// dll long pcmsrv32:PCMSNumRptBytes(long tripId, long rpt)
var pcmsNumRptBytes = pcmsrv.NewProc("PCMSNumRptBytes")

const PCM_RPT_STATE = 1

var _ = meth2("GetRpt(tripId, rpt)",
	func(a, b Value) Value {
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
	})

// dll long pcmsrv32:PCMSNewTrip(long serverId)
var pcmsNewTrip = pcmsrv.NewProc("PCMSNewTrip")

var _ = meth1("PCMNewTrip(tripId)",
	func(a Value) Value {
		rtn := goc.Syscall1(pcmsNewTrip.Addr(),
			intArg(a))
		return int32Ret(rtn)
	})

// dll long pcmsrv32:PCMSNumMatches(long tripId)
var pcmsNumMatches = pcmsrv.NewProc("PCMSNumMatches")

var _ = meth1("PCMNumMatches(tripId)",
	func(a Value) Value {
		rtn := goc.Syscall1(pcmsNumMatches.Addr(),
			intArg(a))
		return int32Ret(rtn)
	})

// dll long pcmsrv32:PCMSOpenServer(long hInstance, long hwnd)
var pcmsOpenServer = pcmsrv.NewProc("PCMSOpenServer")

var _ = meth2("PCMOpenServer(tripId, hwnd)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(pcmsOpenServer.Addr(),
			intArg(a),
			intArg(b))
		return int32Ret(rtn)
	})

// dll long pcmsrv32:PCMSIsValid(long serverId)
var pcmsIsValid = pcmsrv.NewProc("PCMSIsValid")

var _ = meth1("PCMIsValid(serverId)",
	func(a Value) Value {
		rtn := goc.Syscall1(pcmsIsValid.Addr(),
			intArg(a))
		return int32Ret(rtn)
	})

// dll void pcmsrv32:PCMSSetBordersOpen(long tripId, bool open)
var pcmsSetBordersOpen = pcmsrv.NewProc("PCMSSetBordersOpen")

var _ = meth2("PCMSetBordersOpen(tripId, open)",
	func(a, b Value) Value {
		goc.Syscall2(pcmsSetBordersOpen.Addr(),
			intArg(a),
			boolArg(b))
		return nil
	})

// dll void pcmsrv32:PCMSSetCalcType(long tripId, long routeType)
var pcmsSetCalcType = pcmsrv.NewProc("PCMSSetCalcType")

var _ = meth2("PCMSetCalcType(tripId, routeType)",
	func(a, b Value) Value {
		goc.Syscall2(pcmsSetCalcType.Addr(),
			intArg(a),
			intArg(b))
		return nil
	})

// dll long pcmsrv32:PCMSLookup(long tripId, [in] string placeName, long easyMatch)
var pcmsLookup = pcmsrv.NewProc("PCMSLookup")

var _ = meth3("PCMLookup(tripId, placeName, easyMatch)",
	func(a, b, c Value) Value {
		defer heap.FreeTo(heap.CurSize())
		rtn := goc.Syscall3(pcmsLookup.Addr(),
			intArg(a),
			uintptr(stringArg(b)),
			intArg(c))
		return int32Ret(rtn)
	})

// dll void pcmsrv32:PCMSSetMiles(long tripId)
var pcmsSetMiles = pcmsrv.NewProc("PCMSSetMiles")

var _ = meth1("PCMSetMiles(tripId)",
	func(a Value) Value {
		rtn := goc.Syscall1(pcmsSetMiles.Addr(),
			intArg(a))
		return int32Ret(rtn)
	})

// dll void pcmsrv32:PCMSSetCustomMode(long tripId, bool onOff)
var pcmsSetCustomMode = pcmsrv.NewProc("PCMSSetCustomMode")

var _ = meth2("PCMSetCustomMode(tripId, onOff)",
	func(a, b Value) Value {
		goc.Syscall2(pcmsSetCustomMode.Addr(),
			intArg(a),
			boolArg(b))
		return nil
	})

// dll void pcmsrv32:PCMSSetCalcTypeEx(long tripId, long routeType,
//
//	long optFlags, long vehType)
var pcmsSetCalcTypeEx = pcmsrv.NewProc("PCMSSetCalcTypeEx")

var _ = meth4("PCMSetCalcTypeEx(tripId, routeType, optFlags, vehType)",
	func(a, b, c, d Value) Value {
		rtn := goc.Syscall4(pcmsSetCalcTypeEx.Addr(),
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d))
		return int32Ret(rtn)
	})

// dll void pcmsrv32:PCMSSetVehicleType(long tripId, bool onOff)
var pcmsSetVehicleType = pcmsrv.NewProc("PCMSSetVehicleType")

var _ = meth2("PCMSetVehicleType(tripId, onOff)",
	func(a, b Value) Value {
		goc.Syscall2(pcmsSetVehicleType.Addr(),
			intArg(a),
			boolArg(b))
		return nil
	})

// dll void pcmsrv32:PCMSSetRouteLevel(long trip, bool UseStreets)
var pcmsSetRouteLevel = pcmsrv.NewProc("PCMSSetRouteLevel")

var _ = meth2("PCMSetRouteLevel(tripId, onOff)",
	func(a, b Value) Value {
		goc.Syscall2(pcmsSetRouteLevel.Addr(),
			intArg(a),
			boolArg(b))
		return nil
	})

//-------------------------------------------------------------------

func meth0(s string, f func() Value) bool {
	name, ps := paramSplit(s)
	pcmsrv64.Data[name] =
		&SuBuiltin0{Fn: f, BuiltinParams: BuiltinParams{ParamSpec: ps}}
	return true
}

func meth1(s string, f func(a1 Value) Value) bool {
	name, ps := paramSplit(s)
	pcmsrv64.Data[name] =
		&SuBuiltin1{Fn: f, BuiltinParams: BuiltinParams{ParamSpec: ps}}
	return true
}

func meth2(s string, f func(a1, a2 Value) Value) bool {
	name, ps := paramSplit(s)
	pcmsrv64.Data[name] =
		&SuBuiltin2{Fn: f, BuiltinParams: BuiltinParams{ParamSpec: ps}}
	return true
}

func meth3(s string, f func(a1, a2, a3 Value) Value) bool {
	name, ps := paramSplit(s)
	pcmsrv64.Data[name] =
		&SuBuiltin3{Fn: f, BuiltinParams: BuiltinParams{ParamSpec: ps}}
	return true
}

func meth4(s string, f func(a1, a2, a3, a4 Value) Value) bool {
	name, ps := paramSplit(s)
	pcmsrv64.Data[name] =
		&SuBuiltin4{Fn: f, BuiltinParams: BuiltinParams{ParamSpec: ps}}
	return true
}
