package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

func init() {
	ObjectMethods = Methods{
		"Add": rawmethod("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				// TODO handle at: and @args
				ob := this.(*SuObject)
				for i := 0; i < as.Unnamed(); i++ {
					ob.Add(args[i])
				}
				return this
			}),
		"Members": method0(func(this Value) Value { // TODO sequence
			ob := this.(*SuObject)
			mems := new(SuObject)
			it := ob.MapIter();
			for {
				key,_ := it()
				if key == nil {
					break
				}
				mems.Add(key)
			}
			return mems
		}),
		"Size": method0(func(this Value) Value {
			ob := this.(*SuObject)
			return IntToValue(ob.Size())
		}),
		"Sort!": method0(func(this Value) Value { // TODO override Lt
			this.(*SuObject).Sort()
			return this
		}),
		// TODO more methods
	}
}
