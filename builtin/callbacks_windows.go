// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"fmt"
	"log"
	"sync"

	"github.com/apmckinlay/gsuneido/builtin/goc"
	heap "github.com/apmckinlay/gsuneido/builtin/heapstack"
	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

// WARNING: Not thread-safe.
// Should only be used on main UI thread.

type callback struct {
	fn     Value
	used   bool
	active bool
}

var cb2s [goc.Ncb2s]callback
var cb3s [goc.Ncb3s]callback
var cb4s [goc.Ncb4s]callback

var cbs = [3][]callback{cb2s[:], cb3s[:], cb4s[:]}

func callback2(i, a, b uintptr) uintptr {
	return cb2s[i].callv(
		IntVal(int(a)),
		IntVal(int(b)))
}
func callback3(i, a, b, c uintptr) uintptr {
	return cb3s[i].callv(
		IntVal(int(a)),
		IntVal(int(b)),
		IntVal(int(c)))
}
func callback4(i, a, b, c, d uintptr) uintptr {
	return cb4s[i].callv(
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
			log.Fatalln("callback: heapSize", heapSize, "=>", heap.CurSize(),
				"in", cb.fn, args)
		}
	}()
	if !cb.active {
		log.Println("CALLBACK TO INACTIVE!!!", cb.fn)
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
			log.Println("Error in Handler", e)
			MessageBox(fmt.Sprint(e), "Error in Handler")
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

var cblock sync.Mutex

func NewCallback(fn Value, nargs int) uintptr {
	// need locking because SetTimer(Delayed) can be called from other threads
	cblock.Lock()
	defer cblock.Unlock()
	if fn.Type() == types.Number {
		return uintptr(ToInt(fn))
	}
	callbacks := cbs[nargs-2]
	j := -1
	for i := range callbacks {
		cb := &callbacks[i]
		if !cb.used {
			cb.used = true
			j = i
			break
		}
		if j == -1 && !cb.active { // reuse
			j = i
		}
		if cb.active && cbeq(fn, cb.fn) {
			panic("duplcate callback")
		}
	}
	if j == -1 {
		log.Fatalln("too many callbacks")
	}
	cb := &callbacks[j]
	cb.fn = fn
	cb.active = true
	return goc.GetCallback(nargs, j)
}

// cbeq is identity equality, except for bound methods
// can't just use Equals because it's deep equals for SuInstance
func cbeq(x, y Value) bool {
	if x == y {
		return true
	}
	if mx, ok := x.(*SuMethod); ok {
		return mx.Equal(y)
	}
	return false
}

func ClearCallback(fn Value) bool {
	foundInactive := false
	for _, c := range cbs {
		for i := range c {
			cb := &c[i]
			if !cb.used {
				break
			}
			if cbeq(fn, cb.fn) {
				if cb.active {
					if !options.ClearCallbackDisabled {
						cb.active = false
					}
					// keep the fn in case it gets called soon after clear
					// keep the go callback to reuse
					return true
				}
				foundInactive = true
			}
		}
	}
	if foundInactive {
		log.Println("ClearCallback FOUND INACTIVE", fn)
	} else {
		log.Println("ClearCallback NOT FOUND", fn)
	}
	return false // not found
}

var _ = builtin0("Callbacks()", func() Value {
	ob := NewSuObject()
	for _, c := range cbs {
		for _, cb := range c {
			if !cb.used {
				break
			}
			if cb.active {
				ob.Add(cb.fn)
			}
		}
	}
	return ob
})

var _ = builtin1("ClearCallback(fn)", func(fn Value) Value {
	return SuBool(ClearCallback(fn))
})

//TODO may want to delay reuse, e.g. add to tail of free list (per nargs)
