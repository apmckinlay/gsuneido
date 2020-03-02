// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtin("Display(value, quotes=0)", // raw to get thread
	func(t *Thread, args []Value) Value {
		t.Quote = ToInt(args[1])
		defer func() { t.Quote = 0 }()
		return SuStr(Display(t, args[0]))
	})
