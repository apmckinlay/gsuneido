package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

func init() {
	ObjectMethods = Methods{
		"Add": rawmethod("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				// TODO handle at: and @args
				ob := toObject(this)
				for i := 0; i < as.Unnamed(); i++ {
					ob.Add(args[i])
				}
				return this
			}),
		"Members": method0(func(this Value) Value { // TODO sequence
			ob := toObject(this)
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
		"Set_readonly": method0(func(this Value) Value {
			toObject(this).SetReadOnly()
			return this
		}),
		"Size": method0(func(this Value) Value {
			ob := toObject(this)
			return IntToValue(ob.Size())
		}),
		"Sort!": method0(func(this Value) Value { // TODO override Lt
			toObject(this).Sort()
			return this
		}),
		// TODO more methods
	}
}

func toObject(x Value) *SuObject {
	if ob, ok := x.(*SuObject); ok {
		return ob
	}
	if r, ok := x.(*SuRecord); ok {
		return &r.SuObject
	}
	panic("can't convert " + x.TypeName() + " to object")
}
