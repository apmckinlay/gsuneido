package language

import (
	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/runtime"
)

func Def() {
	def := func(nameVal, val Value) Value {
		name := string(nameVal.(SuStr))
		if ss, ok := val.(SuStr); ok {
			val = compile.NamedConstant(name, string(ss))
		}
		TestGlobal(name, val)
		return nil
	}
	AddGlobal("Def", &Builtin2{def, BuiltinParams{ParamSpec: ParamSpec2}})
}
