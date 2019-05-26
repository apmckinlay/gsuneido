package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

func init() {
	ClassMethods = Methods{
		"*new*": methodRaw("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args []Value) Value {
				return this.(*SuClass).New(t, as)
			}),
		"Method?":     method("(string)", methodQ),
		"MethodClass": method("(string)", methodClass),
		"Readonly?": method0(func(this Value) Value {
			return True
		}),
	}
}

func methodClass(t *Thread, this Value, args []Value) Value {
	m := ToStr(args[0])
	return nilToFalse(this.(Findable).Finder(t, func(c Value, mb *MemBase) Value {
		if x, ok := mb.Data[m]; ok && x.Type() == types.Function {
			return c
		}
		return nil
	}))
}

func methodQ(t *Thread, this Value, args []Value) Value {
	m := ToStr(args[0])
	return nilToFalse(this.(Findable).Finder(t, func(c Value, mb *MemBase) Value {
		if x, ok := mb.Data[m]; ok && x.Type() == types.Function {
			return True
		}
		return nil
	}))
}
