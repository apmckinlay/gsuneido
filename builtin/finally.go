// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin(Finally, "(main, final)")

func Finally(t *Thread, args []Value) Value {
	defer func() {
		e := recover()
		func() {
			defer func() {
				if e != nil {
					recover() // if main block panics, ignore finally panic
				}
			}()
			t.Call(args[1])
		}()
		if e != nil {
			panic(e)
		}
	}()
	return t.Call(args[0])
}
