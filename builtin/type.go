package builtin

import (
	. "github.com/apmckinlay/gsuneido/base"
	. "github.com/apmckinlay/gsuneido/interp"
	"github.com/apmckinlay/gsuneido/interp/global"
)

var _ = global.Add("Type", Builtin(suType))

func suType(_ *ArgSpec, args ...Value) Value {
	value := args[0]
	return SuStr(value.TypeName())
}
