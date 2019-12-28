// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin("Trace(value, block = false)",
	func(t *Thread, args []Value) Value {
		if s, ok := args[0].ToStr(); ok {
			if args[1] != False {
				panic("usage: Trace(string) or Trace(flags, block)")
			}
			Trace(s)
		} else {
			oldFlags := options.Trace
			options.Trace = ToInt(args[0])
			if 0 == (options.Trace & (options.TraceConsole | options.TraceLogFile)) {
				options.Trace |= options.TraceConsole | options.TraceLogFile
			}
			if args[1] != False {
				defer func() {
					options.Trace = oldFlags
				}()
				return t.Call(args[1])
			}
		}
		return nil
	})
