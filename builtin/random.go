// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"math/rand"

	. "github.com/apmckinlay/gsuneido/core"
)

var _ = builtin(Random, "(limit)")

func Random(arg Value) Value {
	limit := IfInt(arg)
	return IntVal(rand.Intn(limit))
}
