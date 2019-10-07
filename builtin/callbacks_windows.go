package builtin

import (
	"fmt"

	. "github.com/apmckinlay/gsuneido/runtime"
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

var callbacks [maxCb]callback

func (cb *callback) callv(args ...Value) uintptr {
	defer func () {
		if e := recover(); e != nil {
			fmt.Println("panic in callback", e)
		}
	}()
	if !cb.active {
		fmt.Println("CALLBACK TO INACTIVE!!!")
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

func NewCallback(fn Value, nargs byte) uintptr {
	for j := range callbacks {
		i := j
		cb := &callbacks[i]
		if !cb.active && (cb.gcb == 0 || cb.nargs == nargs) {
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
			} else {
				fmt.Println("--- reuse callback", i)
			}
			cb.fn = fn
			cb.nargs = nargs
			cb.active = true
			return cb.gcb
		}
	}
	panic("too many callbacks")
}

const clearCallbackDisabled = false

func init() {
	if clearCallbackDisabled {
		fmt.Println("ClearCallback disabled")
	}
}

func ClearCallback(fn Value) bool {
	for _, cb := range callbacks {
		if cb.fn == fn {
			if !clearCallbackDisabled {
				cb.active = false
			}
			// keep the fn in case it gets called soon after clear
			// keep the go callback to reuse
			return true
		}
	}
	return false // not found
}

var _ = builtin1("ClearCallback(fn)", func(fn Value) Value {
	return SuBool(ClearCallback(fn))
})

//TODO may want to delay reuse, e.g. add to tail of free list (per nargs)
