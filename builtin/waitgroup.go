// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"sync"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
)

type suWaitGroup struct {
	ValueBase[suWaitGroup]
	wg sync.WaitGroup
}

var _ = builtin(WaitGroup, "()")

func WaitGroup() Value {
	return &suWaitGroup{}
}

var suWaitGroupMethods = methods("wg")

var _ = method(wg_Add, "()")

func wg_Add(this Value) Value {
	wg := this.(*suWaitGroup)
	wg.wg.Add(1)
	return nil
}

var _ = method(wg_Done, "()")

func wg_Done(this Value) Value {
	wg := this.(*suWaitGroup)
	wg.wg.Done()
	return nil
}

var _ = method(wg_Thread, "(block, name = false)")

func wg_Thread(th *Thread, this Value, args []Value) Value {
	wg := this.(*suWaitGroup)
	wg.wg.Add(1)
	fn := args[0]
	fn.SetConcurrent()
	t2 := NewThread(th)
	thread_Name(t2, args[1:])
	threads.add(t2)
	go func() {
		defer func() {
			wg.wg.Done()
			t2.Close()
			threads.remove(t2.Num)
			if e := recover(); e != nil {
				LogUncaught(t2, "Thread", e)
			}
		}()
		t2.Call(fn)
	}()
	return nil
}

var _ = method(wg_Wait, "(secs = 10)")

func wg_Wait(th *Thread, this Value, args []Value) Value {
	timeout := IfInt(args[0])
	if timeout <= 0 {
		panic("WaitGroup.Wait: timeout must be > 0")
	}
	wg := this.(*suWaitGroup)
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.wg.Wait()
	}()
	select {
	case <-c:
		return True // completed normally
	case <-time.After(time.Duration(timeout) * time.Second):
		th.ReturnThrow = true
		return SuStr("WaitGroup: timeout")
	}
}

// Value implementation

var _ Value = (*suWaitGroup)(nil)

func (wg *suWaitGroup) Equal(other any) bool {
	return wg == other
}

func (*suWaitGroup) Lookup(_ *Thread, method string) Value {
	return suWaitGroupMethods[method]
}

func (*suWaitGroup) SetConcurrent() {
	// ok for concurrent use
}
