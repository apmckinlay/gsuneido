// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	InstanceMethods = Methods{
		"Base?": method("(class)", func(t *Thread, this Value, args []Value) Value {
			instance := this.(*SuInstance)
			class := instance.Base()
			if class == args[0] {
				return True
			}
			return nilToFalse(class.Finder(t,
				func(v Value, _ *MemBase) Value {
					if v == args[0] {
						return True
					}
					return nil
				}))
		}),
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
		"Readonly?": method0(func(this Value) Value {
			return False
		}),
	}
}
