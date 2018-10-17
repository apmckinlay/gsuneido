package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

func init() {
	ObjectMethods = Methods{
		method("Size()", func(t *Thread, self Value, args ...Value) Value {
			ob := self.(*SuObject)
			return SuInt(ob.Size())
		}),
		method("Add(value)", func(t *Thread, self Value, args ...Value) Value {
			ob := self.(*SuObject)
			ob.Add(args[0]) // TODO handle multiple arguments (no massage)
			return self
		}),
		// TODO
	}
}
