// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"runtime"
	"sync"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/str"
)

// See also: call.go

type suThreadGlobal struct {
	SuBuiltin
}

func init() {
	Global.Builtin("Thread", &suThreadGlobal{
		SuBuiltin{Fn: threadCallClass,
			BuiltinParams: BuiltinParams{
				ParamSpec: params("(@args)")}}})
}

type threadList struct {
	list map[int32]*Thread // map so we can remove
	lock sync.Mutex
}

var threads = threadList{list: map[int32]*Thread{}}

func (ts *threadList) add(th *Thread) {
	ts.lock.Lock()
	defer ts.lock.Unlock()
	ts.list[th.Num] = th
}

func (ts *threadList) remove(num int32) {
	ts.lock.Lock()
	defer ts.lock.Unlock()
	delete(ts.list, num)
}

func (ts *threadList) count() int {
	ts.lock.Lock()
	defer ts.lock.Unlock()
	return len(ts.list)
}

func threadCallClass(th *Thread, args []Value) Value {
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
	go func() {
		defer func() {
			t2.Close()
			threads.remove(t2.Num)
			if e := recover(); e != nil {
				LogUncaught(t2, "Thread", e)
			}
		}()
		t2.CallEach(fn, ob)
	}()
	return nil
}

var threadMethods = methods("thread")

var _ = staticMethod(thread_GC, "()")

func thread_GC() Value {
	runtime.GC()
	return nil
}

var _ = staticMethod(thread_Name, "(name=false)")

func thread_Name(th *Thread, args []Value) Value {
	if args[0] != False {
		threadName(th, args[0])
	}
	return SuStr(th.Name)
}

func threadName(th *Thread, name Value) {
	th.Name = str.BeforeFirst(th.Name, " ") + " " + ToStr(name)
}

var _ = staticMethod(thread_Count, "()")

func thread_Count() Value {
	return IntVal(threads.count())
}

var _ = AddInfo("builtin.nThread", threads.count)

var _ = staticMethod(thread_List, "()")

func thread_List() Value {
	ob := &SuObject{}
	threads.lock.Lock()
	defer threads.lock.Unlock()
	for _, t := range threads.list {
		ob.Add(SuStr(t.Name))
	}
	return ob
}

var _ = staticMethod(thread_Sleep, "(ms)")

func thread_Sleep(ms Value) Value {
	time.Sleep(time.Duration(ToInt(ms)) * time.Millisecond)
	return nil
}

var _ = staticMethod(thread_Profile, "(block)")

func thread_Profile(th *Thread, args []Value) Value {
	th.StartProfile()
	defer th.StopProfile()
	th.Call(args[0])
	total, self, calls := th.StopProfile()
	prof := &SuObject{}
	for f, c := range calls {
		ob := &SuObject{}
		ob.Set(SuStr("name"), SuStr(f.String()))
		ob.Set(SuStr("calls"), IntVal(int(c)))
		ob.Set(SuStr("total"), Int64Val(int64(total[f])))
		ob.Set(SuStr("self"), Int64Val(int64(self[f])))
		prof.Add(ob)
	}
	return prof
}

var _ = staticMethod(thread_NewSuneidoGlobal, "()")

func thread_NewSuneidoGlobal(th *Thread, _ []Value) Value {
	th.Suneido.Store(new(SuneidoObject))
	return nil
}

var _ = staticMethod(thread_MainQ, "()")

func thread_MainQ(th *Thread, _ []Value) Value {
	return SuBool(th == MainThread || OnUIThread())
}

var _ = staticMethod(thread_Exit, "()")

func thread_Exit(th *Thread, _ []Value) Value {
	if th == MainThread || OnUIThread() {
		panic("suneido: cannot use Thread.Exit on main thread")
	}
	runtime.Goexit()
	return nil
}

var _ = staticMethod(thread_Members, "()")

func thread_Members() Value {
	return thread_members
}

var thread_members = methodList(threadMethods)

func (tg *suThreadGlobal) Lookup(th *Thread, method string) Value {
	if f, ok := threadMethods[method]; ok {
		return f
	}
	return tg.SuBuiltin.Lookup(th, method) // for Params
}

func (*suThreadGlobal) String() string {
	return "Thread /* builtin class */"
}

// ThreadList is used by HttpStatus
func ThreadList() []string {
	threads.lock.Lock()
	defer threads.lock.Unlock()
	list := make([]string, 0, len(threads.list))
	for _, t := range threads.list {
		list = append(list, t.Name)
	}
	return list
}
