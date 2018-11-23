package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

func init() {
	ObjectMethods = Methods{
		"Size": method0(func(self Value) Value {
			ob := self.(*SuObject)
			return SuInt(ob.Size())
		}),
		"Add": rawmethod("(@args)",
			func(t *Thread, self Value, as *ArgSpec, args ...Value) Value {
				// TODO handle at: and @args
				ob := self.(*SuObject)
				for i := 0; i < int(as.Unnamed); i++ {
					ob.Add(args[i])
				}
				return self
			}),
		// TODO more methods
	}
}
