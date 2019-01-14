package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

func init() {
	ObjectMethods = Methods{
		"Add": rawmethod("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				// TODO handle at: and @args
				ob := ToObject(this)
				for i := 0; i < as.Unnamed(); i++ {
					ob.Add(args[i])
				}
				return this
			}),
		"Members": method0(func(this Value) Value { // TODO sequence
			ob := ToObject(this)
			mems := new(SuObject)
			it := ob.MapIter()
			for {
				key, _ := it()
				if key == nil {
					break
				}
				mems.Add(key)
			}
			return mems
		}),
		"Set_default": method1("(value=nil)", func(this Value, val Value) Value {
			ToObject(this).SetDefault(val)
			return this
		}),
		"Set_readonly": method0(func(this Value) Value {
			ToObject(this).SetReadOnly()
			return this
		}),
		"Size": method0(func(this Value) Value {
			ob := ToObject(this)
			return IntToValue(ob.Size())
		}),
		"Sort!": method0(func(this Value) Value { // TODO override Lt
			ToObject(this).Sort()
			return this
		}),
		// TODO more methods
	}
}
