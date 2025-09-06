// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"sync"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/types"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// methods common to SuClass and SuInstance

var _ = exportMethods(&BaseMethods, "base")

var _ = method(base_Base, "()")

func base_Base(th *Thread, this Value, args []Value) Value {
	return base(th, this, func(v Value, _ *MemBase) Value { return v })
}

var _ = method(base_Eval, "(@args)")

func base_Eval(th *Thread, as *ArgSpec, this Value, args []Value) Value {
	return EvalAsMethod(th, as, this, args)
}

var _ = method(base_Eval2, "(@args)")

func base_Eval2(th *Thread, as *ArgSpec, this Value, args []Value) Value {
	ob := &SuObject{}
	if result := EvalAsMethod(th, as, this, args); result != nil {
		ob.Add(result)
	}
	return ob
}

var _ = method(base_GetDefault, "(member, block)")

func base_GetDefault(th *Thread, this Value, args []Value) Value {
	if x := this.Get(th, args[0]); x != nil {
		return x
	}
	if args[1].Type() == types.Block {
		return th.Call(args[1])
	}
	return args[1]
}

var _ = method(base_MemberQ, "(string)")

func base_MemberQ(th *Thread, this Value, arg []Value) Value {
	m := ToStr(arg[0])
	result := this.(Findable).Finder(th, func(v Value, mb *MemBase) Value {
		if mb.Has(m) {
			return True
		}
		return nil
	})
	return nilToFalse(result)
}

var _ = method(base_Members, "(all = false)")

func base_Members(th *Thread, this Value, args []Value) Value {
	if args[0] == True {
		args[0] = nil
	}
	list := &SuObject{}
	this.(Findable).Finder(th, func(v Value, mb *MemBase) Value {
		mb.AddMembersTo(list)
		return args[0]
	})
	list.Sort(nil, False)
	list.Unique()
	return list
}

var _ = method(base_Size, "()")

func base_Size(th *Thread, this Value, args []Value) Value {
	return this.(Findable).Finder(th, func(_ Value, mb *MemBase) Value {
		return IntVal(mb.Size())
	})
}

var _ = method(base_Synchronized, "(block)")

func base_Synchronized(th *Thread, this Value, args []Value) Value {
	name := th.ClassName()
	assert.That(name != "")
	mutVal, ok := classMutexes.Load(name)
	if !ok {
		// multiple threads could get here (race) but only one will store
		mutVal, _ = classMutexes.LoadOrStore(name, MakeMutexT())
	}
	mut := mutVal.(MutexT)
	mut.Lock()
	defer mut.Unlock()
	return th.Call(args[0])
}

// classMutexes holds mutexes for synchronized access by class name.
// This map only grows, never shrinks.
// The assumption is that Synchronized is not heavily used.
var classMutexes sync.Map

// base skips the first
func base(th *Thread, x Value, fn func(Value, *MemBase) Value) Value {
	first := true
	return nilToFalse(x.(Findable).Finder(th, func(v Value, mb *MemBase) Value {
		if first {
			first = false
			return nil
		}
		return fn(v, mb)
	}))
}

func nilToFalse(result Value) Value {
	if result == nil {
		result = False
	}
	return result
}
