package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtinRaw("Display(value)", // raw to get thread
	func(t *Thread, as *ArgSpec, args []Value) Value {
		args = t.Args(&ParamSpec1, as)
		return SuStr(Display(t, args[0]))
	})
