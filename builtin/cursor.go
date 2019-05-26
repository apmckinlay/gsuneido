package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtinRaw("Cursor(@args)",
	func(th *Thread, as *ArgSpec, args []Value) Value {
		query, args := extractQuery(th, queryBlockParams, as, args)
		icursor := th.Dbms().Cursor(query)
		c := NewSuCursor(query, icursor)
		if args[1] == False {
			return c
		}
		// block form
		defer func() {
			c.Close()
		}()
		return th.CallWithArgs(args[1], c)
	})

// see also QueryMethods
func init() {
	CursorMethods = Methods{
		"Next": method1("(transaction)", func(this, arg Value) Value {
			return this.(*SuCursor).GetRec(arg.(*SuTran), Next)
		}),
		"Prev": method1("(transaction)", func(this, arg Value) Value {
			return this.(*SuCursor).GetRec(arg.(*SuTran), Prev)
		}),
		"Output": method("(transaction, record)",
			func(_ *Thread, _ Value, _ []Value) Value {
				panic("cursor.Output is not supported")
			}),
	}
}
