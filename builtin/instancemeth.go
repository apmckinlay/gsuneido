package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

func init() {
	InstanceMethods = Methods{
		//TODO Copy
		"Delete": method2("(key = nil, all = false)",
			func(this, key, all Value) Value {
				if all == True {
					this.(*SuInstance).Clear()
				} else {
					this.(*SuInstance).Delete(key)
				}
				return this
			}),
		"Readonly?": method0(func(this Value) Value {
			return False
		}),
	}
}
