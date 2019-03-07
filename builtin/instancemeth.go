package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

func init() {
	InstanceMethods = Methods{
		"Base": method0(func(this Value) Value {
			return this.(*SuInstance).Base()
		}),
		"Members": method0(func(this Value) Value {
			return this.(*SuInstance).Members()
		}),
		"Member?": memberq, // from ClassMethods
		"Size": method0(func(this Value) Value {
			return this.(*SuInstance).Size()
		}),
		// TODO more methods
	}
}
