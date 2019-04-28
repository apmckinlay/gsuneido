package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

func init() {
	ClassMethods = Methods{
		"*new*": methodRaw("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				return this.(*SuClass).New(t, as)
			}),
		"Readonly?": method0(func(this Value) Value {
			return True
		}),
	}
}
