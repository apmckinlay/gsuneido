// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
)

var _ = builtin(Cmp, "(x, y)")

func Cmp(x, y Value) Value {
	return IntVal(x.Compare(y))
}
