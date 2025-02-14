// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"log"
	"sync"
	"syscall"
	"time"

	"github.com/apmckinlay/gsuneido/builtin/heap"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/types"
	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/assert"
)

const ncb2s = 32
const ncb3s = 32
const ncb4s = 512

// var ncbs = []int{ncb2s, ncb3s, ncb4s}

// wndProcCb is the single callback shared by WndProc's
var wndProcCb = syscall.NewCallback(wndProcCall)

// hwndToCb maps hwnd's to Suneido callbacks for wndProcCall
var hwndToCb = map[uintptr]Value{}

// WndProcCallback is used by SetWindowProc
func WndProcCallback(hwnd uintptr, fn Value) uintptr {
	hwndToCb[hwnd] = fn
	return wndProcCb
}

const delay = 10

var startTime = time.Now()

// clock ticks every millisecond
func clock() uint32 {
	d := time.Since(startTime)
	return uint32(d / time.Millisecond)
}

type callback struct {
	// gocb is from syscall.NewCallback
	gocb uintptr
	// fn is the current Suneido function for the callback
	fn Value
	// active is set to true when the callback is allocated
	// and set to false when it's cleared
	active bool
	// keepTill records "when" the callback was cleared.
	// We delay reusing slots since calls may happen after clear.
	keepTill uint32
}

var cb2s [ncb2s]callback
var cb3s [ncb3s]callback
var cb4s [ncb4s]callback

var cbs = [3][]callback{cb2s[:], cb3s[:], cb4s[:]}

func newGoCallback(nargs, i int) uintptr {
	switch nargs {
	case 2:
		return syscall.NewCallback(func(a, b uintptr) uintptr {
			return cb2s[i].callv(
				IntVal(int(a)),
				IntVal(int(b)))
		})
	case 3:
		return syscall.NewCallback(func(a, b, c uintptr) uintptr {
			return cb3s[i].callv(
				IntVal(int(a)),
				IntVal(int(b)),
				IntVal(int(c)))
		})
	case 4:
		return syscall.NewCallback(func(a, b, c, d uintptr) uintptr {
			return cb4s[i].callv(
				IntVal(int(a)),
				IntVal(int(b)),
				IntVal(int(c)),
				IntVal(int(d)))
		})
	}
	panic(assert.ShouldNotReachHere())
}

func wndProcCall(a, b, c, d uintptr) uintptr {
	if fn, ok := hwndToCb[a]; ok {
		return call(fn,
			IntVal(int(a)),
			IntVal(int(b)),
			IntVal(int(c)),
			IntVal(int(d)))
	}
	log.Fatalln("FATAL: WndProc callback missing hwnd")
	return 0
}

func (cb *callback) callv(args ...Value) uintptr {
	if !cb.active && cb.keepTill < clock() {
		log.Println("ERROR: callback to inactive", cb.fn,
			"keepTill", cb.keepTill, "clock", clock())
	}
	return call(cb.fn, args...)
}

func call(fn Value, args ...Value) uintptr {
	heapSize := heap.CurSize()
	state := MainThread.GetState()
	defer func() {
		if e := recover(); e != nil {
			handler(e, state)
		}
		if heap.CurSize() != heapSize {
			Fatal("callback: heapSize", heapSize, "=>", heap.CurSize(),
				"in", fn, args)
		}
	}()
	x := MainThread.Call(fn, args...)
	if x == nil || x == False {
		return 0
	}
	if x == True {
		return 1
	}
	return uintptr(ToInt(x))
}

func handler(e any, state ThreadState) {
	if MainThread.InHandler {
		LogUncaught(MainThread, "Handler", e)
		Alert("Error in Handler:", e)
		return
	}
	MainThread.InHandler = true
	defer func() {
		MainThread.InHandler = false
		MainThread.RestoreState(state)
		if e2 := recover(); e2 != nil {
			LogUncaught(MainThread, "Handler", e2)
			Alert("Error in Handler:", e2, "\ncaused by", e)
		}
	}()
	se := ToSuExcept(MainThread, e)
	handler := Global.GetName(MainThread, "Handler")
	MainThread.Call(handler, se, Zero, se.Callstack)
}

var cblock sync.Mutex

func NewCallback(fn Value, nargs int) uintptr {
	if fn.Type() == types.Number {
		return uintptr(ToInt(fn))
	}
	cblock.Lock()
	defer cblock.Unlock()
	clock := clock()
	callbacks := cbs[nargs-2]
	j := -1
	for i := range callbacks {
		cb := &callbacks[i]
		if cb.gocb == 0 {
			if j == -1 {
				cb.gocb = newGoCallback(nargs, i)
				j = i
			}
			break
		}
		if j == -1 && !cb.active && cb.keepTill < clock {
			j = i // reuse
			// don't break so we finish checking for duplicate
		}
		if cb.active && cbeq(fn, cb.fn) {
			panic("duplicate callback")
		}
	}
	if j == -1 {
		// fmt.Println("Last 10 callbacks, clock ", clock)
		// for _, c := range callbacks[ncbs[nargs-2]-10:] {
		// 	fmt.Println(c.active, c.keepTill, c.fn)
		// }
		Fatal("too many callbacks", nargs-2)
	}
	cb := &callbacks[j]
	cb.fn = fn
	cb.active = true
	return cb.gocb
}

// cbeq is identity equality, except for bound methods
// can't just use Equal because it's deep equals for SuInstance
func cbeq(x, y Value) bool {
	if x == y {
		return true
	}
	if mx, ok := x.(*SuMethod); ok {
		return mx.Equal(y)
	}
	return false
}

func clearCallback(fn Value) bool {
	cblock.Lock()
	defer cblock.Unlock()
	foundInactive := false
	for _, c := range cbs {
		for i := range c {
			cb := &c[i]
			if cb.gocb == 0 {
				break
			}
			if cbeq(fn, cb.fn) {
				if cb.active {
					if !options.ClearCallbackDisabled {
						cb.active = false
						cb.keepTill = clock() + delay
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
		log.Println("ERROR: ClearCallback found inactive", fn)
	} else {
		for hwnd, cbfn := range hwndToCb {
			if cbfn == fn {
				delete(hwndToCb, hwnd)
				return true
			}
		}
	}
	// one reason it may not be found is if it was overwritten in hwndToCb
	return false // not found
}

var _ = builtin(Callbacks, "()")

func Callbacks() Value {
	ob := &SuObject{}
	for _, c := range cbs {
		for _, cb := range c {
			if cb.gocb == 0 {
				break
			}
			if cb.active {
				ob.Add(cb.fn)
			}
		}
	}
	return ob
}

func CallbacksCount() int {
	n := 0
	for _, c := range cbs {
		for _, cb := range c {
			if cb.gocb == 0 {
				break
			}
			if cb.active {
				n++
			}
		}
	}
	return n
}

var _ = AddInfo("windows.nCallback", CallbacksCount)

var _ = builtin(ClearCallback, "(fn)")

func ClearCallback(fn Value) Value {
	return SuBool(clearCallback(fn))
}

func WndProcCount() int {
	return len(hwndToCb)
}

var _ = AddInfo("windows.nWndProc", WndProcCount)
