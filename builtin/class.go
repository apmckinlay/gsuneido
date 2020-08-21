// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

func init() {
	ClassMethods = Methods{
		"*new*": methodRaw("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args []Value) Value {
				return this.(*SuClass).New(t, as)
			}),
		"Base?": method("(class)", func(t *Thread, this Value, args []Value) Value {
			return nilToFalse(this.(*SuClass).Finder(t,
				func(v Value, _ *MemBase) Value {
					if v == args[0] {
						return True
					}
					return nil
				}))
		}),
		"Method?": method("(string)",
			func(t *Thread, this Value, args []Value) Value {
				m := ToStr(args[0])
				return nilToFalse(this.(Findable).Finder(t, func(c Value, mb *MemBase) Value {
					if _, ok := c.(*SuClass); ok {
						if x, ok := mb.Data[m]; ok && x.Type() == types.Function {
							return True
						}
					}
					return nil
				}))
			}),
		"MethodClass": method("(string)",
			func(t *Thread, this Value, args []Value) Value {
				m := ToStr(args[0])
				return nilToFalse(this.(Findable).Finder(t, func(c Value, mb *MemBase) Value {
					if _, ok := c.(*SuClass); ok {
						if x, ok := mb.Data[m]; ok && x.Type() == types.Function {
							return c
						}
					}
					return nil
				}))
			}),
		"Readonly?": method0(func(this Value) Value {
			return True
		}),
	}
}
