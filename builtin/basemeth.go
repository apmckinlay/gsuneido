// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

// methods common to SuClass and SuInstance

var _ = exportMethods(&BaseMethods)

var _ = method(base_Base, "()")

func base_Base(t *Thread, this Value, args []Value) Value {
	return base(t, this, func(v Value, _ *MemBase) Value { return v })
}

var _ = method(base_Eval, "(@args)")

func base_Eval(t *Thread, as *ArgSpec, this Value, args []Value) Value {
	return EvalAsMethod(t, as, this, args)
}

var _ = method(base_Eval2, "(@args)")

func base_Eval2(t *Thread, as *ArgSpec, this Value, args []Value) Value {
	ob := &SuObject{}
	if result := EvalAsMethod(t, as, this, args); result != nil {
		ob.Add(result)
	}
	return ob
}

var _ = method(base_GetDefault, "(member, block)")

func base_GetDefault(t *Thread, this Value, args []Value) Value {
	if x := this.Get(t, args[0]); x != nil {
		return x
	}
	if args[1].Type() == types.Block {
		return t.Call(args[1])
	}
	return args[1]
}

var _ = method(base_MemberQ, "(string)")

func base_MemberQ(t *Thread, this Value, arg []Value) Value {
	m := ToStr(arg[0])
	result := this.(Findable).Finder(t, func(v Value, mb *MemBase) Value {
		if mb.Has(m) {
			return True
		}
		return nil
	})
	return nilToFalse(result)
}

var _ = method(base_Members, "(all = false)")

func base_Members(t *Thread, this Value, args []Value) Value {
	if args[0] == True {
		args[0] = nil
	}
	list := &SuObject{}
	this.(Findable).Finder(t, func(v Value, mb *MemBase) Value {
		mb.AddMembersTo(list)
		return args[0]
	})
	list.Sort(nil, False)
	list.Unique()
	return list
}

var _ = method(base_Size, "()")

func base_Size(t *Thread, this Value, args []Value) Value {
	return this.(Findable).Finder(t, func(_ Value, mb *MemBase) Value {
		return IntVal(mb.Size())
	})
}

// base skips the first
func base(t *Thread, x Value, fn func(Value, *MemBase) Value) Value {
	first := true
	return nilToFalse(x.(Findable).Finder(t, func(v Value, mb *MemBase) Value {
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
