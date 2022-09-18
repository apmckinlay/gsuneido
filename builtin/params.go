// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

type paramsable interface {
	Params() string
}

var _ = exportMethods(&ParamsMethods)

var _ = method(fn_Params, "()")

func fn_Params(this Value) Value {
	fn := this.(paramsable)
	return SuStr(fn.Params())
}
