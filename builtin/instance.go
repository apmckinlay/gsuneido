package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	InstanceMethods = Methods{
		"Copy": method0(func(this Value) Value {
			return this.(*SuInstance).Copy()
		}),
		"Delete": methodRaw("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args []Value) Value {
				if all := getNamed(as, args, SuStr("all")); all == True {
					this.(*SuInstance).Clear()
				} else {
					iter := NewArgsIter(as, args)
					for {
						k, v := iter()
						if k != nil || v == nil {
							break
						}
						this.(*SuInstance).Delete(v)
					}
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
