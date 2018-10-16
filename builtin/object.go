package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = AddGlobal("Object",
	&Builtin{Fn: suObject, ParamSpec: ParamSpec{Nparams: 1, Flags: []Flag{AtParam}}})

func suObject(_ *Thread, args ...Value) Value {
	return args[0]
}
