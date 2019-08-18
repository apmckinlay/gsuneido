package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

const maxCb = 2000 // same as Go's limit (as of 20190813)

var callbacks [maxCb]struct {
	fn     Value
	nargs  byte
	active bool
	gcb    uintptr
}

func NewCallback(fn Value, nargs byte) uintptr {
	for i := range callbacks {
		cb := &callbacks[i]
		if !cb.active && (cb.gcb == 0 || cb.nargs == nargs) {
			if cb.gcb == 0 {
				// create a reusable callback for callbacks[i]
				switch nargs {
				case 1:
					cb.gcb = windows.NewCallback(func(a uintptr) uintptr {
						x := UIThread.Call(callbacks[i].fn, IntVal(int(a)))
						return uintptr(ToInt(x))
					})
				case 2:
					cb.gcb = windows.NewCallback(func(a, b uintptr) uintptr {
						x := UIThread.Call(callbacks[i].fn,
							IntVal(int(a)),
							IntVal(int(b)))
						if x == nil {
							return 0
						}
						return uintptr(ToInt(x))
					})
				case 3:
					cb.gcb = windows.NewCallback(func(a, b, c uintptr) uintptr {
						x := UIThread.Call(callbacks[i].fn,
							IntVal(int(a)),
							IntVal(int(b)),
							IntVal(int(c)))
						if x == nil {
							return 0
						}
						return uintptr(ToInt(x))
					})
				case 4:
					cb.gcb = windows.NewCallback(func(a, b, c, d uintptr) uintptr {
						x := UIThread.Call(callbacks[i].fn,
							IntVal(int(a)),
							IntVal(int(b)),
							IntVal(int(c)),
							IntVal(int(d)))
						if x == nil {
							return 0
						}
						return uintptr(ToInt(x))
					})
				default:
					panic("callback with unsupported nargs")
				}
			}
			cb.fn = fn
			cb.nargs = nargs
			cb.active = true
			return cb.gcb
		}
	}
	panic("too many callback1")
}

func ClearCallback(fn Value) bool {
	for _, cb := range callbacks {
		if cb.fn == fn {
			cb.active = false
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
