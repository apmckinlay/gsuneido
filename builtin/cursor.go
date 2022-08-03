// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtinRaw("Cursor(@args)",
	func(th *Thread, as *ArgSpec, args []Value) Value {
		query, args := extractQuery(th, &queryBlockParams, as, args)
		icursor := th.Dbms().Cursor(query, nil)
		c := NewSuCursor(th, query, icursor)
		if args[1] == False {
			return c
		}
		// block form
		defer func() {
			c.Close()
		}()
		return th.Call(args[1], c)
	})

// see also QueryMethods
func init() {
	CursorMethods = Methods{
		"Next": method("(transaction)",
			func(th *Thread, this Value, args []Value) Value {
				return this.(*SuCursor).GetRec(th, args[0].(*SuTran), Next)
			}),
		"Prev": method("(transaction)",
			func(th *Thread, this Value, args []Value) Value {
				return this.(*SuCursor).GetRec(th, args[0].(*SuTran), Prev)
			}),
		"Output": method("(transaction, record)",
			func(*Thread, Value, []Value) Value {
				panic("cursor.Output is not supported")
			}),
	}
}
