package language

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

func Concat() {
	fn := func(x, y Value) Value {
		return NewSuConcat().Add(AsStr(x)).Add(AsStr(y))
	}
	Global.Builtin("Concat", &SuBuiltin2{fn, BuiltinParams{ParamSpec: ParamSpec2}})
}
