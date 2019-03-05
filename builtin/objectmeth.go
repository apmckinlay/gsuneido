package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	ObjectMethods = Methods{
		"Add": methodRaw("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				ob := ToObject(this)
				iter := NewArgsIter(as, args)
				if at := getNamed(as, args, SuStr("at")); at != nil {
					if i, ok := at.IfInt(); ok {
						addAt(ob, i, iter)
					} else {
						putAt(ob, at, iter)
					}
				} else {
					addAt(ob, ob.ListSize(), iter)
				}
				return this
			}),
		"Assocs": method0(func(this Value) Value {
			return NewSuSequence(ToObject(this).IterAssocs())
		}),
		"Iter": method0(func(this Value) Value {
			return SuIter{Iter: ToObject(this).Iter()}
		}),
		"Members": method0(func(this Value) Value {
			return NewSuSequence(ToObject(this).IterMembers())
		}),
		"Set_default": method1("(value=nil)", func(this Value, val Value) Value {
			ToObject(this).SetDefault(val)
			return this
		}),
		"Set_readonly": method0(func(this Value) Value {
			ToObject(this).SetReadOnly()
			return this
		}),
		"Size": method0(func(this Value) Value {
			return IntToValue(ToObject(this).Size())
		}),
		"Sort!": methodRaw("(block = false)", // methodRaw to get thread
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				args = t.Args(&ParamSpecOptionalBlock, as)
				ToObject(this).Sort(t, args[0])
				return this
			}),
		"Values": method0(func(this Value) Value {
			return NewSuSequence(ToObject(this).Iter())
		}),
		// TODO more methods
	}
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

func addAt(ob *SuObject, at int, iter ArgsIter) {
	for {
		k, v := iter()
		if k != nil || v == nil {
			break
		}
		ob.Insert(at, v)
		at++
	}
}

func putAt(ob *SuObject, at Value, iter ArgsIter) {
	k, v := iter()
	if k != nil || v == nil {
		return
	}
	if k, v := iter(); k == nil && v != nil {
		panic("can only Add multiple values to un-named or numeric positions")
	}
	ob.Put(at, v)
}
