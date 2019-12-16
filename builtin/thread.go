// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"log"
	"time"

	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
)

type SuThreadGlobal struct {
	SuBuiltin
}

func init() {
	name, ps := paramSplit("Thread(block)")
	Global.Builtin(name, &SuThreadGlobal{
		SuBuiltin{threadCallClass, BuiltinParams{ParamSpec: *ps}}})
}

var threads = map[int32]*Thread{}

func threadCallClass(t *Thread, args []Value) Value {
	if options.ThreadDisabled {
		return nil
	}
	fn := args[0]
	fn.SetConcurrent()
	t2 := NewThread()
	t2.Token = t.Dbms().Token()
	threads[t2.Num] = t2 //TODO lock
	go func() {
		defer func() {
			if e := recover(); e != nil {
				log.Println("error in thread:", e)
				t2.PrintStack()
			}
			t2.Close()
			delete(threads, t2.Num) //TODO lock
		}()
		t2.Call(fn)
	}()
	return nil
}

var threadMethods = Methods{
	"Name": method("(name=false)", func(t *Thread, _ Value, args []Value) Value {
		if args[0] != False {
			t.Name = ToStr(args[0])
		}
		return SuStr(t.Name)
	}),
	"Count": method0(func(this Value) Value {
		return IntVal(len(threads))
	}),
	"List": method0(func(this Value) Value {
		ob := NewSuObject()
		for _, t := range threads { //TODO lock
			ob.Put(nil, SuStr(t.Name), True)
		}
		return ob
	}),
	"Sleep": method1("(ms)", func(this, ms Value) Value {
		time.Sleep(time.Duration(1000000 * ToInt(ms)))
		return nil
	}),
}

func (d *SuThreadGlobal) Lookup(t *Thread, method string) Callable {
	if f, ok := threadMethods[method]; ok {
		return f
	}
	return d.SuBuiltin.Lookup(t, method) // for Params
}

func (d *SuThreadGlobal) String() string {
	return "Thread /* builtin class */"
}

var _ = builtin("Scheduled(ms, block)",
	func(t *Thread, args []Value) Value {
		ms := time.Duration(ToInt(args[0])) * time.Millisecond
		t2 := NewThread()
		t2.Token = t.Dbms().Token()
		block := args[1]
		go func() {
			defer t2.Close()
			time.Sleep(ms)
			t2.Call(block)
		}()
		return nil
	})
