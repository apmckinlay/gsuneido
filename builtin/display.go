package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtinRaw("Display(value)", // raw to get thread
	func(t *Thread, as *ArgSpec, args ...Value) Value {
		args = t.Args(&ParamSpec1, as)
		return SuStr(display(t, args[0]))
	})

func display(t *Thread, val Value) string {
	if d, ok := val.(displayable); ok {
		return d.Display(t)
	}
	return val.String()
}

type displayable interface {
	Display(*Thread) string
}
