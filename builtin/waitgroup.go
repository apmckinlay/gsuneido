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

var _ = builtin(WaitGroup, "() :unknown")

func WaitGroup() Value {
	return &suWaitGroup{}
}

var suWaitGroupMethods = methods("wg")

var _ = method(wg_Add, "(inc = 1) :void")

func wg_Add(this Value, a Value) Value {
	wg := this.(*suWaitGroup)
	inc := ToInt(a)
	wg.wg.Add(inc)
	return nil
}

var _ = method(wg_Done, "() :void")

func wg_Done(this Value) Value {
	wg := this.(*suWaitGroup)
	wg.wg.Done()
	return nil
}

var _ = method(wg_Thread, "(@args) :void")

func wg_Thread(th *Thread, this Value, args []Value) Value {
	wg := this.(*suWaitGroup)
	ob := args[0].(*SuObject)
	ob.SetConcurrent()
	var fn Value
	if block := ob.NamedGet(SuStr("block")); block != nil {
		fn = block
		ob.Delete(th, SuStr("block"))
	} else {
		fn = ob.ListGet(0)
		ob.Delete(th, Zero)
	}
	t2 := NewThread(th)
	if name := ob.NamedGet(SuStr("name")); name != nil {
		threadName(t2, name)
		ob.Delete(th, SuStr("name"))
	}
	threads.add(t2)
	wg.wg.Go(func() {
		defer func() {
			t2.Close()
			threads.remove(t2.Num)
			if e := recover(); e != nil {
				LogUncaught(t2, "Thread", e)
			}
		}()
		t2.CallEach(fn, ob)
	})
	return nil
}

var _ = method(wg_Wait, "(secs :number = 10) :true")

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
