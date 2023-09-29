// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
)

var _ = builtin(Cursor, "(@args)")

func Cursor(th *Thread, as *ArgSpec, args []Value) Value {
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
}

// see also QueryMethods

var _ = exportMethods(&CursorMethods)

var _ = method(cursor_Next, "(transaction)")

func cursor_Next(th *Thread, this Value, args []Value) Value {
	return this.(*SuCursor).GetRec(th, args[0].(*SuTran), Next)
}

var _ = method(cursor_Prev, "(transaction)")

func cursor_Prev(th *Thread, this Value, args []Value) Value {
	return this.(*SuCursor).GetRec(th, args[0].(*SuTran), Prev)
}

var _ = method(cursor_Output, "(transaction, record)")

func cursor_Output(*Thread, Value, []Value) Value {
	panic("cursor.Output is not supported")
}
