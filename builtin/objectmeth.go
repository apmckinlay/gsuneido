package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

func init() {
	ObjectMethods = Methods{
		"Size": method0(func(this Value) Value {
			ob := this.(*SuObject)
			return SuInt(ob.Size())
		}),
		"Add": rawmethod("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				// TODO handle at: and @args
				ob := this.(*SuObject)
				for i := 0; i < int(as.Unnamed); i++ {
					ob.Add(args[i])
				}
				return this
			}),
		// TODO more methods
	}
}
