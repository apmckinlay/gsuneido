package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

var _ = builtin("Object(@args)",
	func(_ *Thread, args []Value) Value {
		return args[0]
	})

// NOTE: ObjectMethods are shared with SuRecord

func init() {
	ObjectMethods = Methods{
		"Add": methodRaw("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args []Value) Value {
				ob := ToContainer(this)
				iter := NewArgsIter(as, args)
				if at := getNamed(as, args, SuStr("at")); at != nil {
					if i, ok := at.IfInt(); ok {
						addAt(ob, i, iter)
					} else {
						putAt(t, ob, at, iter)
					}
				} else {
					addAt(ob, ob.ListSize(), iter)
				}
				return this
			}),
		"Assocs": methodRaw("(list = true, named = true)",
			func(_ *Thread, as *ArgSpec, this Value, args []Value) Value {
				list, named := iterWhich(as, args)
				return NewSuSequence(IterAssocs(ToContainer(this), list, named))
			}),
		"BinarySearch": method("(value, block = false)",
			func(t *Thread, this Value, args []Value) Value {
				ob := ToContainer(this).ToObject()
				if args[1] == False {
					return IntVal(ob.BinarySearch(args[0]))
				}
				return IntVal(ob.BinarySearch2(t, args[0], args[1]))
			}),
		"Copy": method0(func(this Value) Value {
			return ToContainer(this).Copy()
		}),
		"Delete": methodRaw("(@args)",
			obDelete),
		"Erase": methodRaw("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args []Value) Value {
				ob := ToContainer(this)
				iter := NewArgsIter(as, args)
				for {
					k, v := iter()
					if k != nil || v == nil {
						break
					}
					ob.Erase(t, v)
				}
				return this
			}),
		"Eval": methodRaw("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args []Value) Value {
				return nilToEmptyStr(EvalAsMethod(t, as, this, args))
			}),
		"Eval2": methodRaw("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args []Value) Value {
				ob := &SuObject{}
				if result := EvalAsMethod(t, as, this, args); result != nil {
					ob.Add(result)
				}
				return ob
			}),
		"Find": method1("(value)", func(this Value, val Value) Value {
			return ToContainer(this).ToObject().Find(val)
		}),
		"GetDefault": method("(member, block)",
			func(t *Thread, this Value, args []Value) Value {
				ob := ToContainer(this)
				if x := ob.GetIfPresent(t, args[0]); x != nil {
					return x
				}
				if args[1].Type() == types.Block {
					return t.Call(args[1])
				}
				return args[1]
			}),
		"Iter": method0(func(this Value) Value {
			return SuIter{Iter: IterValues(ToContainer(this), true, true)}
		}),
		"Join": method1("(separator='')", func(this Value, arg Value) Value {
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
		}),
		"Members": methodRaw("(list = true, named = true)",
			func(_ *Thread, as *ArgSpec, this Value, args []Value) Value {
				list, named := iterWhich(as, args)
				return NewSuSequence(IterMembers(ToContainer(this), list, named))
			}),
		"Member?": method1("(member)", func(this Value, val Value) Value {
			return SuBool(ToContainer(this).HasKey(val))
		}),
		"Readonly?": method0(func(this Value) Value {
			return SuBool(ToContainer(this).IsReadOnly())
		}),
		"Reverse!": method0(func(this Value) Value {
			ToContainer(this).ToObject().Reverse()
			return this
		}),
		"Set_default": method1("(value=nil)", func(this Value, val Value) Value {
			ToContainer(this).ToObject().SetDefault(val)
			return this
		}),
		"Set_readonly": method0(func(this Value) Value {
			ToContainer(this).SetReadOnly()
			return this
		}),
		"Size": method2("(list=false,named=false)",
			func(this, arg1, arg2 Value) Value {
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
			}),
		"Sort!": method("(block = false)",
			func(t *Thread, this Value, args []Value) Value {
				ToContainer(this).ToObject().Sort(t, args[0])
				return this
			}),
		"Unique!": method0(func(this Value) Value {
			ToContainer(this).ToObject().Unique()
			return this
		}),
		"Values": methodRaw("(list = true, named = true)",
			func(_ *Thread, as *ArgSpec, this Value, args []Value) Value {
				list, named := iterWhich(as, args)
				return NewSuSequence(IterValues(ToContainer(this), list, named))
			}),
	}
	ObjectMethods["LowerBound"] = ObjectMethods["BinarySearch"]
}

func obDelete(t *Thread, as *ArgSpec, this Value, args []Value) Value {
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
			ob.Delete(t, v)
		}
	}
	return this
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

func putAt(t *Thread, ob Container, at Value, iter ArgsIter) {
	k, v := iter()
	if k != nil || v == nil {
		return
	}
	if k, v := iter(); k == nil && v != nil {
		panic("can only Add multiple values to un-named or numeric positions")
	}
	ob.Put(t, at, v)
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
