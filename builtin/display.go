// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtin(display, "(value, quotes=0)")

func display(t *Thread, args []Value) Value {
	defer func(q int) { t.Quote = q }(t.Quote)
	t.Quote = ToInt(args[1])
	return SuStr(Display(t, args[0]))
}
