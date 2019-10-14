package builtin

import (
	"fmt"
	"log"

	heap "github.com/apmckinlay/gsuneido/builtin/heapstack"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
	"golang.org/x/sys/windows"
)

// WARNING: Not thread-safe.
// Should only be used on main UI thread.

const maxCb = 2000 // same as Go's limit (as of 20190813)

type callback struct {
	fn     Value
	nargs  byte
	active bool
	gcb    uintptr
}

var callbacks [maxCb]callback

func (cb *callback) call1(a uintptr) uintptr {
	return cb.callv(
		IntVal(int(a)))
}
func (cb *callback) call2(a, b uintptr) uintptr {
	return cb.callv(
		IntVal(int(a)),
		IntVal(int(b)))
}
func (cb *callback) call3(a, b, c uintptr) uintptr {
	return cb.callv(
		IntVal(int(a)),
		IntVal(int(b)),
		IntVal(int(c)))
}
func (cb *callback) call4(a, b, c, d uintptr) uintptr {
	return cb.callv(
		IntVal(int(a)),
		IntVal(int(b)),
		IntVal(int(c)),
		IntVal(int(d)))
}

func (cb *callback) callv(args ...Value) uintptr {
	heapSize := heap.CurSize()
	defer func() {
		if e := recover(); e != nil {
			handler(e)
		}
		if heap.CurSize() != heapSize {
			log.Fatalln("callback: heap", heapSize, "=>", heap.CurSize(),
				"in", cb.fn, args)
		}
	}()
	if !cb.active {
		log.Println("CALLBACK TO INACTIVE!!!")
	}
	x := UIThread.Call(cb.fn, args...)
	if x == nil || x == False {
		return 0
	}
	if x == True {
		return 1
	}
	return uintptr(ToInt(x))
}

func handler(e interface{}) {
	defer func() {
		if e := recover(); e != nil {
			log.Fatalln("error in Handler:", e)
		}
	}()
	// debug.PrintStack()
	// UIThread.PrintStack()
	log.Println("panic in callback:", e, "<<<<<<<<<<<<<<<<")

	se, ok := e.(*SuExcept)
	if !ok {
		s := fmt.Sprint(e) // TODO avoid fmt
		se = NewSuExcept(UIThread, SuStr(s))
	}
	handler := Global.GetName(UIThread, "Handler")
	UIThread.Call(handler, se, Zero, se.Callstack)
}

func NewCallback(fn Value, nargs byte) uintptr {
	if fn.Type() == types.Number {
		return uintptr(ToInt(fn))
	}
	j := -1
	for i := range callbacks {
		cb := &callbacks[i]
		if j == -1 { // haven't found one yet
			if cb.gcb == 0 || // unused
				(!cb.active && cb.nargs == nargs) { // reuse
				j = i
			}
		}
		if cb.active && cb.fn == fn {
			panic("duplcate callback")
		}
	}
	if j == -1 {
		tooManyCallbacks()
	}
	cb := &callbacks[j]
	if cb.gcb == 0 {
		// create a reusable callback for callbacks[i]
		switch nargs {
		case 1:
			cb.gcb = windows.NewCallback(cb.call1)
		case 2:
			cb.gcb = windows.NewCallback(cb.call2)
		case 3:
			cb.gcb = windows.NewCallback(cb.call3)
		case 4:
			cb.gcb = windows.NewCallback(cb.call4)
		default:
			panic("callback with unsupported number of arguments")
		}
	}
	cb.fn = fn
	cb.nargs = nargs
	cb.active = true
	return cb.gcb
}

func tooManyCallbacks() {
	for i, cb := range callbacks {
		if cb.active {
			log.Println(i, cb.fn)
		}
	}
	log.Fatalln("too many callbacks")
}

const clearCallbackDisabled = false

func init() {
	if clearCallbackDisabled {
		fmt.Println("ClearCallback disabled")
	}
}

func ClearCallback(fn Value) bool {
	for i := range callbacks {
		cb := &callbacks[i]
		if cb.gcb == 0 {
			break
		}
		if cb.active && cb.fn == fn {
			if !clearCallbackDisabled {
				cb.active = false
			}
			// keep the fn in case it gets called soon after clear
			// keep the go callback to reuse
			return true
		}
	}
	log.Println("NOT FOUND")
	return false // not found
}

var _ = builtin0("Callbacks()", func() Value {
	ob := NewSuObject()
	for _, cb := range callbacks {
		if cb.active {
			ob.Add(cb.fn)
		}
	}
	return ob
})

var _ = builtin1("ClearCallback(fn)", func(fn Value) Value {
	return SuBool(ClearCallback(fn))
})

//TODO may want to delay reuse, e.g. add to tail of free list (per nargs)
