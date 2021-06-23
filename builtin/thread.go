// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/str"
)

type suThreadGlobal struct {
	SuBuiltin
}

func init() {
	name, ps := paramSplit("Thread(block)")
	Global.Builtin(name, &suThreadGlobal{
		SuBuiltin{Fn: threadCallClass,
			BuiltinParams: BuiltinParams{ParamSpec: *ps}}})
}

type threadList struct {
	list map[int32]*Thread // map so we can remove
	lock sync.Mutex
}

var threads = threadList{list: map[int32]*Thread{}}

func (ts *threadList) add(num int32, t *Thread) {
	ts.lock.Lock()
	defer ts.lock.Unlock()
	ts.list[num] = t
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

func threadCallClass(_ *Thread, args []Value) Value {
	if options.ThreadDisabled {
		return nil
	}
	fn := args[0]
	fn.SetConcurrent()
	t2 := NewThread()

	threads.add(t2.Num, t2)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				log.Println("ERROR in thread:", e)
				t2.PrintStack()
				if InternalError(e) {
					buf := make([]byte, 512)
					n := runtime.Stack(buf, false)
					os.Stderr.Write(buf[:n])
				}
			}
			t2.Close()
			threads.remove(t2.Num)
		}()
		t2.Call(fn)
	}()
	return nil
}

func InternalError(e interface{}) bool {
	if _, ok := e.(runtime.Error); ok {
		return true
	}
	if s, ok := e.(string); ok {
		return strings.HasPrefix(s, "assert failed")
	}
	if se, ok := e.(*SuExcept); ok {
		return strings.HasPrefix(string(se.SuStr), "assert failed")
	}
	return false
}

var threadMethods = Methods{
	"Name": method("(name=false)", func(t *Thread, _ Value, args []Value) Value {
		if args[0] != False {
			t.Name = str.BeforeFirst(t.Name, " ") + " " + ToStr(args[0])
		}
		return SuStr(t.Name)
	}),
	"Count": method0(func(this Value) Value {
		return IntVal(threads.count())
	}),
	"List": method0(func(this Value) Value {
		ob := &SuObject{}
		threads.lock.Lock()
		defer threads.lock.Unlock()
		for _, t := range threads.list {
			ob.Put(nil, SuStr(t.Name), True)
		}
		return ob
	}),
	"Sleep": method1("(ms)", func(this, ms Value) Value {
		time.Sleep(time.Duration(ToInt(ms)) * time.Millisecond)
		return nil
	}),
}

func (d *suThreadGlobal) Lookup(t *Thread, method string) Callable {
	if f, ok := threadMethods[method]; ok {
		return f
	}
	return d.SuBuiltin.Lookup(t, method) // for Params
}

func (d *suThreadGlobal) String() string {
	return "Thread /* builtin class */"
}

var _ = builtin("Scheduled(ms, block)",
	func(_ *Thread, args []Value) Value {
		ms := time.Duration(ToInt(args[0])) * time.Millisecond
		t2 := NewThread()
		block := args[1]
		block.SetConcurrent()
		go func() {
			defer t2.Close()
			time.Sleep(ms)
			t2.Call(block)
		}()
		return nil
	})
