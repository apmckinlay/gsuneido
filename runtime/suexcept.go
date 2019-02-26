package runtime

import (
	"runtime"
)

type SuExcept struct {
	SuStr
	Callstack *SuObject
}

func NewSuExcept(t *Thread, s SuStr) *SuExcept {
	return &SuExcept{SuStr: s, Callstack: CallStack(t)}
}

// SuValue interface ------------------------------------------------

func (*SuExcept) TypeName() string {
	return "Except"
}

// SuExceptMethods is initialized by the builtin package
var SuExceptMethods Methods

func (*SuExcept) Lookup(method string) Value {
	if m := SuExceptMethods[method]; m != nil {
		return m
	}
	return StringMethods[method]
}

// callstack --------------------------------------------------------

// NOTE: it might be more efficient
// to capture the call stack in an internal format (not an SuObject)
// and only build the SuObject if required

// callstack captures the call stack at the point of recover
// the Go stack will be as of the panic
// the Suneido frame pointer will be unwound
// but the frames themselves will still be intact
func CallStack(t *Thread) *SuObject {
	cs := &SuObject{}
	nframes := countSuneidoFrames()
	for i := t.fp + nframes - 1; i >= 0; i-- {
		fr := t.frames[i]
		call := &SuObject{}
		call.Put(SuStr("fn"), fr.fn)
		locals := &SuObject{}
		for i, v := range fr.locals {
			if v != nil {
				locals.Put(SuStr(fr.fn.Names[i]), v)
			}
		}
		call.Put(SuStr("locals"), locals)
		cs.Add(call)
	}
	return cs
}

const framer = "github.com/apmckinlay/gsuneido/runtime.(*Thread).Run"

// countSuneidoFrames uses the Go call stack
// to count how many frames are active
func countSuneidoFrames() int {
	pc := make([]uintptr, 100)
	n := runtime.Callers(1, pc)
	nframes := 0
	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, more := frames.Next()
		if frame.Function == framer {
			nframes++
		}
		if !more {
			break
		}
	}
	return nframes
}
