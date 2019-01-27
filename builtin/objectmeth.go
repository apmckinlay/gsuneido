package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	ObjectMethods = Methods{
		"Add": rawmethod("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				ob := ToObject(this)
				iter := NewArgsIter(as, args)
				if at := getNamed(as, args, SuStr("at")); at != nil {
					if i, ok := ToInt(at); ok {
						addAt(ob, i, iter)
					} else {
						putAt(ob, at, iter)
					}
				} else {
					addAt(ob, ob.ListSize(), iter)
				}
				return this
			}),
		"Iter": method0(func(this Value) Value { // TODO sequence
			ob := ToObject(this)
			return SuIter{ob.IterValues()}
		}),
		"Members": method0(func(this Value) Value { // TODO sequence
			ob := ToObject(this)
			mems := new(SuObject)
			it := ob.MapIter()
			for {
				key, _ := it()
				if key == nil {
					break
				}
				mems.Add(key)
			}
			return mems
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
			ob := ToObject(this)
			return IntToValue(ob.Size())
		}),
		"Sort!": rawmethod("(block = false)", // rawmethod to get thread
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				args = t.Args(&ParamSpecOptionalBlock, as)
				ToObject(this).Sort(t, args[0])
				return this
			}),
		// TODO more methods
	}
}

func ToInt(x Value) (int, bool) {
	if i, ok := SmiToInt(x); ok {
		return i, ok
	}
	if dn, ok := x.(SuDnum); ok {
		return dn.Dnum.ToInt()
	}
	return 0, false
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
