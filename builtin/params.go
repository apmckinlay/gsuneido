// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
)

type paramsable interface {
	Params() string
}

var _ = exportMethods(&ParamsMethods, "params")

var _ = method(params_Params, "()")

func params_Params(this Value) Value {
	fn := this.(paramsable)
	return SuStr(fn.Params())
}
