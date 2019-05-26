package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtinRaw("Display(value)", // raw to get thread
	func(t *Thread, as *ArgSpec, args []Value) Value {
		args = t.Args(&ParamSpec1, as)
		return SuStr(display(t, args[0]))
	})

func display(t *Thread, val Value) string {
	if d, ok := val.(ToStringable); ok {
		return d.ToString(t)
	}
	return val.String()
}
