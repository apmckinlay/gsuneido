package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

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
						putAt(ob, at, iter)
					}
				} else {
					addAt(ob, ob.ListSize(), iter)
				}
				return this
			}),
		"Assocs": method0(func(this Value) Value { //TODO list? and named?
			return NewSuSequence(ToObject(this).IterAssocs())
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
				result := EvalAsMethod(t, as, ToObject(this), args)
				if result == nil {
					return EmptyStr
				}
				return result
			}),
		"Eval2": methodRaw("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				ob := &SuObject{}
				if result := EvalAsMethod(t, as, ToObject(this), args); result != nil {
					ob.Add(result)
				}
				return ob
			}),
		"Find": method1("(value)", func(this Value, val Value) Value {
			iter := ToObject(this).Iter2()
			for k, v := iter(); v != nil; k, v = iter() {
				if v.Equal(val) {
					return k
				}
			}
			return False
		}),
		"GetDefault": methodRaw("(member, default)", // methodRaw to get thread
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				args = t.Args(&paramSpecGetDef, as)
				ob := ToObject(this)
				x := ob.GetDefault(args[0], nil)
				if x == nil {
					if args[1].Type() == types.Block {
						x = t.CallWithArgs(args[1])
					} else {
						x = args[1]
					}
				}
				return x
			}),
		"Iter": method0(func(this Value) Value {
			return SuIter{Iter: ToObject(this).Iter()}
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
		"Members": method0(func(this Value) Value { //TODO list? and named?
			return NewSuSequence(ToObject(this).IterMembers())
		}),
		"Member?": method1("(member)", func(this Value, val Value) Value {
			return SuBool(ToObject(this).Has(val))
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
		"Sort!": methodRaw("(block = false)", // methodRaw to get thread
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				args = t.Args(&ParamSpecOptionalBlock, as)
				ToObject(this).Sort(t, args[0])
				return this
			}),
		"Values": method0(func(this Value) Value { //TODO list? and named?
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
	ob.Set(at, v)
}

var paramSpecGetDef = params("(member,block)")
