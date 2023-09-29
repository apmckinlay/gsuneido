// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
)

// DefConcat must be called to make Concat available
func DefConcat() {
	builtin(Concat, "(s, t)")
}

func Concat(x, y Value) Value {
	return NewSuConcat().Add(AsStr(x)).Add(AsStr(y))
}
