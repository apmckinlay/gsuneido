// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"sync"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/options"
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
				ParamSpec: params("(block, name = false)")}}})
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
	if options.ThreadDisabled {
		return nil
	}
	fn := args[0]
	fn.SetConcurrent()
	t2 := NewThread(th)
	thread_Name(t2, args[1:])
	threads.add(t2)
	go func() {
		defer func() {
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

var threadMethods = methods()

var _ = staticMethod(thread_Name, "(name=false)")

func thread_Name(th *Thread, args []Value) Value {
	if args[0] != False {
		th.Name = str.BeforeFirst(th.Name, " ") + " " + ToStr(args[0])
	}
	return SuStr(th.Name)
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
	total, self, ops, calls := th.StopProfile()
	prof := &SuObject{}
	for name, op := range ops {
		ob := &SuObject{}
		ob.Set(SuStr("name"), SuStr(name))
		ob.Set(SuStr("ops"), IntVal(int(op)))
		ob.Set(SuStr("calls"), IntVal(int(calls[name])))
		ob.Set(SuStr("total"), IntVal(int(total[name])))
		ob.Set(SuStr("self"), IntVal(int(self[name])))
		prof.Add(ob)
	}
	return prof
}

var _ = staticMethod(thread_NewSuneidoGlobal, "()")

func thread_NewSuneidoGlobal(th *Thread, _ []Value) Value {
	th.Suneido.Store(new(SuneidoObject))
	return nil
}

func (d *suThreadGlobal) Get(_ *Thread, key Value) Value {
	m := ToStr(key)
	if fn, ok := threadMethods[m]; ok {
		return fn.(Value)
	}
	if fn, ok := ParamsMethods[m]; ok {
		return NewSuMethod(d, fn.(Value))
	}
	return nil
}

func (d *suThreadGlobal) Lookup(th *Thread, method string) Callable {
	if f, ok := threadMethods[method]; ok {
		return f
	}
	return d.SuBuiltin.Lookup(th, method) // for Params
}

func (d *suThreadGlobal) String() string {
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
