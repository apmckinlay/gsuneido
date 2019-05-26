package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	InstanceMethods = Methods{
		"Copy": method0(func(this Value) Value {
			return this.(*SuInstance).Copy()
		}),
		"Delete": method2("(key = nil, all = false)",
			func(this, key, all Value) Value {
				if all == True {
					this.(*SuInstance).Clear()
				} else {
					this.(*SuInstance).Delete(key)
				}
				return this
			}),
		"Method?": method("(string)",
			func(t *Thread, this Value, args []Value) Value {
				return methodQ(t, this.(*SuInstance).Base(), args)
			}),
		"MethodClass": method("(string)",
			func(t *Thread, this Value, args []Value) Value {
				return methodClass(t, this.(*SuInstance).Base(), args)
			}),
		"Readonly?": method0(func(this Value) Value {
			return False
		}),
	}
}
