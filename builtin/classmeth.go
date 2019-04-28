package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

func init() {
	ClassMethods = Methods{
		"*new*": methodRaw("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				return this.(*SuClass).New(t, as)
			}),
		"Method?":     method1("(string)", methodQ),
		"MethodClass": method1("(string)", methodClass),
		"Readonly?": method0(func(this Value) Value {
			return True
		}),
	}
}

func methodClass(this, arg Value) Value {
	m := IfStr(arg)
	return nilToFalse(this.(Findable).Finder(func(c Value, mb *MemBase) Value {
		if x, ok := mb.Data[m]; ok && x.Type() == types.Function {
			return c
		}
		return nil
	}))
}

func methodQ(this, arg Value) Value {
	m := IfStr(arg)
	return nilToFalse(this.(Findable).Finder(func(c Value, mb *MemBase) Value {
		if x, ok := mb.Data[m]; ok && x.Type() == types.Function {
			return True
		}
		return nil
	}))
}
