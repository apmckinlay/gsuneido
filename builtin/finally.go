// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
)

var _ = builtin(Finally, "(main, final)")

func Finally(th *Thread, args []Value) Value {
	defer func() {
		e := recover()
		func() {
			defer func() {
				if e != nil && e != BlockReturn {
					recover() // if main block panics, ignore finally panic
				}
			}()
			returnThrow := th.ReturnThrow
			th.ReturnThrow = false
			result := th.Call(args[1])
			// the following should match interp op.Call*
			if th.ReturnThrow {
				th.ReturnThrow = false
				if result != EmptyStr && result != True {
					if s, ok := result.ToStr(); ok {
						panic(s)
					}
					panic("return value not checked")
				}
			}
			th.ReturnThrow = returnThrow
		}()
		if e != nil {
			panic(e)
		}
	}()
	return th.Call(args[0])
}
