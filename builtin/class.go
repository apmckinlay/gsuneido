// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/options"
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
		"StartCoverage": method1("(count = false)", func(this, a Value) Value {
			if atomic.LoadInt64(&options.Coverage) == 0 {
				panic("coverage not enabled")
			}
			c := this.(*SuClass)
			c.StartCoverage(ToBool(a))
			return nil
		}),
		"StopCoverage": method0(func(this Value) Value {
			c := this.(*SuClass)
			return c.StopCoverage()
		}),
	}
}
