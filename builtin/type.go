package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = AddGlobal("Type", Builtin(suType))

func suType(_ *ArgSpec, args ...Value) Value {
	value := args[0]
	return SuStr(value.TypeName())
}
