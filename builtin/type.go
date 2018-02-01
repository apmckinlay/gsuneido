package builtin

import (
	. "github.com/apmckinlay/gsuneido/base"
	"github.com/apmckinlay/gsuneido/interp/global"
)

// type is a builtin function that returns a value's type as a string
func SuType(ct Context, self Value, args ...Value) Value {
	value := args[0]
	return SuStr(value.TypeName())
}

var _ = builtinFunction("Type", SuType)

func builtinFunction(name string,
	fn func(ct Context, self Value, args ...Value) Value) int {
	f := Builtin{}
	// TODO parse a string to get params spec
	f.Nparams = 1
	f.Flags = []Flag{0}
	f.Strings = []string{"value"}
	f.fn = SuType
	global.Add(name, &f)
	return 0
}
