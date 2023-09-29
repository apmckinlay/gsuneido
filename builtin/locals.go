// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import . "github.com/apmckinlay/gsuneido/core"

var _ = builtin(Locals, "(i)")

func Locals(th *Thread, args []Value) Value {
	return th.Locals(ToInt(args[0]))
}
