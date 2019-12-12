package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

// methods common to SuClass and SuInstance

func init() {
	BaseMethods = Methods{
		"Base": method("()", func(t *Thread, this Value, args []Value) Value {
			return base(t, this, func(v Value, _ *MemBase) Value { return v })
		}),
		"Base?": method("(class)", func(t *Thread, this Value, args []Value) Value {
			return nilToFalse(this.(Findable).Finder(t,
				func(v Value, _ *MemBase) Value {
					if v == args[0] {
						return True
					}
					return nil
				}))
		}),
		"Eval": methodRaw("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args []Value) Value {
				return EvalAsMethod(t, as, this, args)
			}),
		"Eval2": methodRaw("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args []Value) Value {
				ob := &SuObject{}
				if result := EvalAsMethod(t, as, this, args); result != nil {
					ob.Add(result)
				}
				return ob
			}),
		"GetDefault": method("(member, block)",
			func(t *Thread, this Value, args []Value) Value {
				if x := this.Get(t, args[0]); x != nil {
					return x
				}
				if args[1].Type() == types.Block {
					return t.Call(args[1])
				}
				return args[1]
			}),
		"Member?": method("(string)", func(t *Thread, this Value, arg []Value) Value {
			m := ToStr(arg[0])
			result := this.(Findable).Finder(t, func(v Value, mb *MemBase) Value {
				if _, ok := mb.Data[m]; ok {
					return True
				}
				return nil
			})
			return nilToFalse(result)
		}),
		"Members": method("(all = false)", func(t *Thread, this Value, args []Value) Value {
			if args[0] == True {
				args[0] = nil
			}
			list := NewSuObject()
			this.(Findable).Finder(t, func(v Value, mb *MemBase) Value {
				mb.AddMembersTo(list)
				return args[0]
			})
			list.Sort(nil, False)
			list.Unique()
			return list
		}),
		"Size": method("()", func(t *Thread, this Value, args []Value) Value {
			return this.(Findable).Finder(t, func(_ Value, mb *MemBase) Value {
				return IntVal(mb.Size())
			})
		}),
	}
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
