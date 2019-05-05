package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

// NOTE: ObjectMethods are shared with SuRecord

func init() {
	ObjectMethods = Methods{
		"Add": methodRaw("(@args)",
			func(_ *Thread, as *ArgSpec, this Value, args ...Value) Value {
				ob := ToObject(this)
				iter := NewArgsIter(as, args)
				if at := getNamed(as, args, SuStr("at")); at != nil {
					if i, ok := at.IfInt(); ok {
						addAt(ob, i, iter)
					} else {
						putAt(ob.Set, at, iter)
					}
				} else {
					addAt(ob, ob.ListSize(), iter)
				}
				return this
			}),
		"Assocs": methodRaw("(list = true, named = true)",
			func(_ *Thread, as *ArgSpec, this Value, args ...Value) Value {
				return NewSuSequence(ToObject(this).IterAssocs(iterWhich(as, args)))
			}),
		"Clear": method0(func(this Value) Value {
			ToObject(this).Clear()
			return nil
		}),
		"Copy": method0(func(this Value) Value {
			return ToObject(this).Copy()
		}),
		"Delete": methodRaw("(@args)",
			func(_ *Thread, as *ArgSpec, this Value, args ...Value) Value {
				ob := ToObject(this)
				if all := getNamed(as, args, SuStr("all")); all == True {
					ob.Clear()
				} else {
					iter := NewArgsIter(as, args)
					for {
						k, v := iter()
						if k != nil || v == nil {
							break
						}
						ob.Delete(v)
					}
				}
				return this
			}),
		"Erase": methodRaw("(@args)",
			func(_ *Thread, as *ArgSpec, this Value, args ...Value) Value {
				ob := ToObject(this)
				iter := NewArgsIter(as, args)
				for {
					k, v := iter()
					if k != nil || v == nil {
						break
					}
					ob.Erase(v)
				}
				return this
			}),
		"Eval": methodRaw("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				return nilToEmptyStr(EvalAsMethod(t, as, this, args))
			}),
		"Eval2": methodRaw("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				ob := &SuObject{}
				if result := EvalAsMethod(t, as, this, args); result != nil {
					ob.Add(result)
				}
				return ob
			}),
		"Find": method1("(value)", func(this Value, val Value) Value {
			k,_ := ToObject(this).Find(val)
			return k
		}),
		"GetDefault": method("(member, block)",
			func(t *Thread, this Value, args ...Value) Value {
				ob := ToObject(this)
				if x := ob.GetDefault(args[0], nil); x != nil {
					return x
				}
				if args[1].Type() == types.Block {
					return t.CallWithArgs(args[1])
				}
				return args[1]
			}),
		"Iter": method0(func(this Value) Value {
			return SuIter{Iter: ToObject(this).IterValues(true, true)}
		}),
		"Join": method1("(separator='')", func(this Value, arg Value) Value {
			ob := ToObject(this)
			separator := ToStr(arg)
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
		"BinarySearch": method("(value, block = false)",
			func(t *Thread, this Value, args ...Value) Value {
				ob := ToObject(this)
				if args[1] == False {
					return IntVal(ob.BinarySearch(args[0]))
				}
				return IntVal(ob.BinarySearch2(t, args[0], args[1]))
			}),
		"Members": methodRaw("(list = true, named = true)",
			func(_ *Thread, as *ArgSpec, this Value, args ...Value) Value {
				return NewSuSequence(ToObject(this).IterMembers(iterWhich(as, args)))
			}),
		"Member?": method1("(member)", func(this Value, val Value) Value {
			return SuBool(ToObject(this).Has(val))
		}),
		"Reverse!": method0(func(this Value) Value {
			ToObject(this).Reverse()
			return this
		}),
		"Set_default": method1("(value=nil)", func(this Value, val Value) Value {
			ToObject(this).SetDefault(val)
			return this
		}),
		"Set_readonly": method0(func(this Value) Value {
			ToObject(this).SetReadOnly()
			return this
		}),
		"Size": method2("(list=false,named=false)",
			func(this, arg1, arg2 Value) Value {
				list := ToBool(arg1)
				named := ToBool(arg2)
				ob := ToObject(this)
				var n int
				if list == named {
					n = ob.Size()
				} else if list {
					n = ob.ListSize()
				} else {
					n = ob.NamedSize()
				}
				return IntVal(n)
			}),
		"Sort!": method("(block = false)",
			func(t *Thread, this Value, args ...Value) Value {
				ToObject(this).Sort(t, args[0])
				return this
			}),
		"Unique!": method0(func(this Value) Value {
			ToObject(this).Unique()
			return this
		}),
		"Values": methodRaw("(list = true, named = true)",
			func(_ *Thread, as *ArgSpec, this Value, args ...Value) Value {
				return NewSuSequence(ToObject(this).IterValues(iterWhich(as, args)))
			}),
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

func putAt(put func(Value, Value), at Value, iter ArgsIter) {
	k, v := iter()
	if k != nil || v == nil {
		return
	}
	if k, v := iter(); k == nil && v != nil {
		panic("can only Add multiple values to un-named or numeric positions")
	}
	put(at, v)
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
