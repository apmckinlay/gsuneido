package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

// methods common to SuClass and SuInstance

func init() {
	BaseMethods = Methods{
		"Base": method0(func(this Value) Value {
			return base(this, func(v Value, _ *MemBase) Value { return v })
		}),
		"Base?": method1("(class)", func(this, c Value) Value {
			return base(this, func(v Value, _ *MemBase) Value {
				if v == c {
					return True
				}
				return nil
			})
		}),
		//TODO Eval, Eval2
		"GetDefault": methodRaw("(member, default)", // methodRaw to get thread
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				args = t.Args(&paramSpecGetDef, as)
				if x := this.Get(t, args[0]); x != nil {
					return x
				}
				if args[1].Type() == types.Block {
					return t.CallWithArgs(args[1])
				}
				return args[1]
			}),
		"Member?": method1("(string)", func(this, arg Value) Value {
			m := IfStr(arg)
			result := this.(Findable).Finder(func(v Value, mb *MemBase) Value {
				if _, ok := mb.Data[m]; ok {
					return True
				}
				return nil
			})
			if result == nil {
				result = False
			}
			return result
		}),
		"Members": method0(func(this Value) Value {
			//TODO all:
			return this.(Findable).Finder(func(_ Value, mb *MemBase) Value {
				return mb.Members()
			})
		}),
		"Size": method0(func(this Value) Value {
			return this.(Findable).Finder(func(_ Value, mb *MemBase) Value {
				return IntVal(mb.Size())
			})
		}),
	}
}

// base skips the first
func base(x Value, fn func(Value, *MemBase) Value) Value {
	first := true
	return nilToFalse(x.(Findable).Finder(func(v Value, mb *MemBase) Value {
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
