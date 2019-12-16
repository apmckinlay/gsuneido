// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

// Concat must be called to be available
func Concat() {
	builtin2("Concat(s, t)", func(x, y Value) Value {
		return NewSuConcat().Add(AsStr(x)).Add(AsStr(y))
	})
}
