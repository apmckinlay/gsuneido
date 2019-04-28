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
		"Method?": method1("(string)", func(this, arg Value) Value {
			return methodQ(this.(*SuInstance).Base(), arg)
		}),
		"MethodClass": method1("(string)", func(this, arg Value) Value {
			return methodClass(this.(*SuInstance).Base(), arg)
		}),
		"Readonly?": method0(func(this Value) Value {
			return False
		}),
	}
}
