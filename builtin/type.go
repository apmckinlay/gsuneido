package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = AddGlobal("Type", &Builtin{Fn: suType, ParamSpec: ParamSpec{Nparams: 1}})

func suType(_ *Thread, args ...Value) Value {
	return SuStr(args[0].TypeName())
}
