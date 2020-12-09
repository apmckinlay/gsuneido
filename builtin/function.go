// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin0("EnableCoverage()", func() Value {
	atomic.StoreInt64(&options.Coverage, 1)
	return nil
})

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
			if atomic.LoadInt64(&options.Coverage) == 0 {
				panic("coverage not enabled")
			}
			fn := this.(*SuFunc)
			fn.StartCoverage(ToBool(a))
			return nil
		}),
		"StopCoverage": method0(func(this Value) Value {
			fn := this.(*SuFunc)
			return fn.StopCoverage()
		}),
	}
}
