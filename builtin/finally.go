// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin("Finally(main_block, final_block)",
	func(t *Thread, args []Value) Value {
		for i := 0; i < 1; i++ { // workaround for 1.14 bug
			defer func() {
				e := recover()
				func() {
					for i := 0; i < 1; i++ { // workaround for 1.14 bug
						defer func() {
							if e != nil {
								recover() // if main block panics, ignore finally panic
							}
						}()
					}
					t.Call(args[1])
				}()
				if e != nil {
					panic(e)
				}
			}()
		}
		return t.Call(args[0])
	})
