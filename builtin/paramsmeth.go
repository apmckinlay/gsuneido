package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

type paramsable interface {
	Params() string
}

func init() {
	ParamsMethods = Methods{
		"Params": method0(func(this Value) Value {
			fn := this.(paramsable)
			return SuStr(fn.Params())
		}),
	}
}
