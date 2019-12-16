// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	SuExceptMethods = Methods{
		"Callstack": method0(func(this Value) Value {
			return this.(*SuExcept).Callstack
		}),
	}
}
