// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/options"
)

var _ = builtin(CoverageEnable, "(enable)")

func CoverageEnable(a Value) Value {
	options.Coverage.Store(ToBool(a))
	return nil
}

var _ = exportMethods(&SuFuncMethods)

var _ = method(func_Disasm, "(source = false)")

func func_Disasm(this, a Value) Value {
	fn := this.(*SuFunc)
	if a == False {
		return SuStr(DisasmOps(fn))
	}
	return SuStr(DisasmMixed(fn, ToStr(a)))
}

var _ = method(func_StartCoverage, "(count = false)")

func func_StartCoverage(this, a Value) Value {
	if !options.Coverage.Load() {
		panic("coverage not enabled")
	}
	fn := this.(*SuFunc)
	fn.StartCoverage(ToBool(a))
	return nil
}

var _ = method(func_StopCoverage, "()")

func func_StopCoverage(this Value) Value {
	fn := this.(*SuFunc)
	return fn.StopCoverage()
}
