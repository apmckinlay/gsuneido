// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/trace"
)

var _ = builtin("Trace(value, block = false)",
	func(t *Thread, args []Value) Value {
		if s, ok := args[0].ToStr(); ok {
			if args[1] != False {
				panic("usage: Trace(string) or Trace(flags, block)")
			}
			trace.Print(s + "\n")
		} else {
			prev := trace.Set(ToInt(args[0]))
			if args[1] != False {
				defer func() {
					trace.Set(prev)
				}()
				return t.Call(args[1])
			}
		}
		return nil
	})
