// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = exportMethods(&SuExceptMethods)

var _ = method(except_Callstack, "()")

func except_Callstack(this Value) Value {
	return this.(*SuExcept).Callstack
}
