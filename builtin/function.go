// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	SuFuncMethods = Methods{
		"Disasm": method1("(source = false)", func(this, a Value) Value {
			fn := this.(*SuFunc)
			if a == False {
				return SuStr(DisasmOps(fn))
			}
			return SuStr(DisasmMixed(fn, ToStr(a)))
		}),
		"StartCoverage": method1("(count = false)", func(this, a Value) Value {
			fn := this.(*SuFunc)
			fn.StartCoverage(ToBool(a))
			return nil
		}),
		"StopCoverage": method0(func(this Value) Value {
			fn := this.(*SuFunc)
			fn.StopCoverage()
			return nil
			// cover := fn.StopCoverage()
			// ob := &SuObject{}
			// for _, c := range cover {
			// 	for i := 0; i < 16; i++ {
			// 		ob.Add(SuBool(c&(1<<i) != 0))
			// 	}
			// }
			// return ob
		}),
	}
}
