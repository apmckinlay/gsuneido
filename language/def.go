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
		Global.TestDef(name, val)
		return nil
	}
	Global.Add("Def", &SuBuiltin2{def, BuiltinParams{ParamSpec: ParamSpec2}})
}
