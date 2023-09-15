// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

var _ = builtin(ConcurrentQ, "(value)")

func ConcurrentQ(v Value) Value {
	return IsConcurrent(v)
}

var _ = builtin(Object, "(@args)")

func Object(arg Value) Value {
	return arg
}

// NOTE: ObjectMethods are shared with SuRecord

var _ = exportMethods(&ObjectMethods)

var _ = method(ob_Add, "(@args)")

func ob_Add(th *Thread, as *ArgSpec, this Value, args []Value) Value {
	ob := ToContainer(this)
	iter := NewArgsIter(as, args)
	if at := getNamed(as, args, SuStr("at")); at != nil {
		if i, ok := at.IfInt(); ok {
			addAt(ob, i, iter)
		} else {
			putAt(th, ob, at, iter)
		}
	} else {
		for {
			k, v := iter()
			if k != nil || v == nil {
				break
			}
			ob.Add(v)
		}
	}
	return this
}

var _ = method(ob_Assocs, "(list = true, named = true)")

func ob_Assocs(_ *Thread, as *ArgSpec, this Value, args []Value) Value {
	list, named := iterWhich(as, args)
	return NewSuSequence(IterAssocs(ToContainer(this), list, named))
}

var _ = method(ob_BinarySearch, "(value, block = false)")

func ob_BinarySearch(th *Thread, this Value, args []Value) Value {
	ob := ToContainer(this).ToObject()
	if args[1] == False {
		return IntVal(ob.BinarySearch(args[0]))
	}
	return IntVal(ob.BinarySearch2(th, args[0], args[1]))
}

var _ = method(ob_Copy, "()")

func ob_Copy(this Value) Value {
	return ToContainer(this).Copy()
}

var _ = method(ob_Delete, "(@args)")

func ob_Delete(th *Thread, as *ArgSpec, this Value, args []Value) Value {
	ob := ToContainer(this)
	if all := getNamed(as, args, SuStr("all")); all == True {
		ob.DeleteAll()
	} else {
		iter := NewArgsIter(as, args)
		for {
			k, v := iter()
			if k != nil || v == nil {
				break
			}
			ob.Delete(th, v)
		}
	}
	return this
}

var _ = method(ob_Erase, "(@args)")

func ob_Erase(th *Thread, as *ArgSpec, this Value, args []Value) Value {
	ob := ToContainer(this)
	iter := NewArgsIter(as, args)
	for {
		k, v := iter()
		if k != nil || v == nil {
			break
		}
		ob.Erase(th, v)
	}
	return this
}

var _ = method(ob_Eval, "(@args)")

func ob_Eval(th *Thread, as *ArgSpec, this Value, args []Value) Value {
	return EvalAsMethod(th, as, this, args)
}

var _ = method(ob_Eval2, "(@args)")

func ob_Eval2(th *Thread, as *ArgSpec, this Value, args []Value) Value {
	ob := &SuObject{}
	if result := EvalAsMethod(th, as, this, args); result != nil {
		ob.Add(result)
	}
	return ob
}

var _ = method(ob_Find, "(value)")

func ob_Find(this Value, val Value) Value {
	return ToContainer(this).ToObject().Find(val)
}

var _ = method(ob_GetDefault, "(member, block)")

func ob_GetDefault(th *Thread, this Value, args []Value) Value {
	ob := ToContainer(this)
	if x := ob.GetIfPresent(th, args[0]); x != nil {
		return x
	}
	if args[1].Type() == types.Block {
		return th.Call(args[1])
	}
	return args[1]
}

var _ = method(ob_HasQ, "(value)")

func ob_HasQ(this Value, val Value) Value {
	return SuBool(ToContainer(this).ToObject().Find(val) != False)
}

var _ = method(ob_Iter, "()")

func ob_Iter(this Value) Value {
	return SuIter{Iter: IterValues(ToContainer(this), true, true)}
}

var _ = method(ob_Join, "(separator='')")

func ob_Join(this Value, arg Value) Value {
	ob := ToContainer(this)
	separator := AsStr(arg)
	sb := strings.Builder{}
	sep := ""
	iter := ob.ArgsIter()
	for {
		k, v := iter()
		if k != nil || v == nil {
			break
		}
		sb.WriteString(sep)
		sep = separator
		sb.WriteString(ToStrOrString(v))
	}
	return SuStr(sb.String())
}

var _ = method(ob_Max, "()")

func ob_Max(this Value) Value {
	iter := ToContainer(this).Iter2(true, true)
	_, max := iter()
	if max == nil {
		panic("cannot use Max on empty object")
	}
	for _, v := iter(); v != nil; _, v = iter() {
		if v.Compare(max) > 0 {
			max = v
		}
	}
	return max
}

var _ = method(ob_Members, "(list = true, named = true)")

func ob_Members(_ *Thread, as *ArgSpec, this Value, args []Value) Value {
	list, named := iterWhich(as, args)
	return NewSuSequence(IterMembers(ToContainer(this), list, named))
}

var _ = method(ob_MemberQ, "(member)")

func ob_MemberQ(this Value, val Value) Value {
	return SuBool(ToContainer(this).HasKey(val))
}

var _ = method(ob_Min, "()")

func ob_Min(this Value) Value {
	iter := ToContainer(this).Iter2(true, true)
	_, min := iter()
	if min == nil {
		panic("cannot use Min on empty object")
	}
	for _, v := iter(); v != nil; _, v = iter() {
		if v.Compare(min) < 0 {
			min = v
		}
	}
	return min
}

var _ = method(ob_PopFirst, "()")

func ob_PopFirst(this Value) Value {
	x := ToContainer(this).ToObject().PopFirst()
	if x == nil {
		return this
	}
	return x
}

var _ = method(ob_PopLast, "()")

func ob_PopLast(this Value) Value {
	x := ToContainer(this).ToObject().PopLast()
	if x == nil {
		return this
	}
	return x
}

var _ = method(ob_ReadonlyQ, "()")

func ob_ReadonlyQ(this Value) Value {
	return SuBool(ToContainer(this).IsReadOnly())
}

var _ = method(ob_ReverseX, "()")

func ob_ReverseX(this Value) Value {
	ToContainer(this).ToObject().Reverse()
	return this
}

var _ = method(ob_Set_default, "(value=nil)")

func ob_Set_default(this Value, val Value) Value {
	ToContainer(this).ToObject().SetDefault(val)
	return this
}

var _ = method(ob_Set_readonly, "()")

func ob_Set_readonly(this Value) Value {
	ToContainer(this).SetReadOnly()
	return this
}

var _ = method(ob_Size, "(list=false,named=false)")

func ob_Size(this, arg1, arg2 Value) Value {
	list := ToBool(arg1)
	named := ToBool(arg2)
	ob := ToContainer(this)
	var n int
	if list == named {
		n = ob.ListSize() + ob.NamedSize()
	} else if list {
		n = ob.ListSize()
	} else {
		n = ob.NamedSize()
	}
	return IntVal(n)
}

var _ = method(ob_SortX, "(block = false)")

func ob_SortX(th *Thread, this Value, args []Value) Value {
	ToContainer(this).ToObject().Sort(th, args[0])
	return this
}

var _ = method(ob_UniqueX, "()")

func ob_UniqueX(this Value) Value {
	ToContainer(this).ToObject().Unique()
	return this
}

var _ = method(ob_Values, "(list = true, named = true)")

func ob_Values(_ *Thread, as *ArgSpec, this Value, args []Value) Value {
	list, named := iterWhich(as, args)
	return NewSuSequence(IterValues(ToContainer(this), list, named))
}

func getNamed(as *ArgSpec, args []Value, name Value) Value {
	iter := NewArgsIter(as, args)
	for k, v := iter(); v != nil; k, v = iter() {
		if name.Equal(k) {
			return v
		}
	}
	return nil
}

func addAt(ob Container, at int, iter ArgsIter) {
	for {
		k, v := iter()
		if k != nil || v == nil {
			break
		}
		ob.Insert(at, v)
		at++
	}
}

func putAt(th *Thread, ob Container, at Value, iter ArgsIter) {
	k, v := iter()
	if k != nil || v == nil {
		return
	}
	if k, v := iter(); k == nil && v != nil {
		panic("can only Add multiple values to un-named or numeric positions")
	}
	ob.Put(th, at, v)
}

func iterWhich(as *ArgSpec, args []Value) (list bool, named bool) {
	ai := NewArgsIter(as, args)
	for k, v := ai(); v != nil; k, v = ai() {
		if k == nil && v != nil {
			panic("usage: () or (list:) or (named:)")
		}
		if k.Equal(SuStr("list")) {
			list = true
		} else if k.Equal(SuStr("named")) {
			named = true
		}
	}
	return
}

// EvalAsMethod runs a function as if it were a method of an object
// i.e. object.Eval
func EvalAsMethod(th *Thread, as *ArgSpec, ob Value, args []Value) Value {
	// first argument is function
	k, f := NewArgsIter(as, args)()
	if k != nil || f == nil {
		panic("usage: object.Eval(callable, ...)")
	}
	if m, ok := f.(*SuMethod); ok {
		f = m.GetFn()
	}
	return f.Call(th, ob, as.DropFirst())
}
