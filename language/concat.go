package language

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

func Concat() {
	fn := func(x, y Value) Value {
		return NewSuConcat().Add(ToStr(x)).Add(ToStr(y))
	}
	Global.Add("Concat", &SuBuiltin2{fn, BuiltinParams{ParamSpec: ParamSpec2}})
}
