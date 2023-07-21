// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/regex"
)

var _ = builtin(Transaction, "(read=nil, update=nil, block=false)")

func Transaction(th *Thread, args []Value) Value {
	if (args[0] == nil) == (args[1] == nil) {
		panic("usage: Transaction(read:) or Transaction(update:)")
	}
	var update bool
	if args[0] == nil {
		update = ToBool(args[1])
	} else {
		update = !ToBool(args[0])
	}
	itran := th.Dbms().Transaction(update)
	st := NewSuTran(itran, update)
	if args[2] == False {
		return st
	}
	// block form
	defer func() {
		if !st.Ended() {
			e := recover()
			if e != nil && e != BlockReturn {
				st.Rollback()
			} else {
				st.Complete()
			}
			if e != nil {
				panic(e)
			}
		}
	}()
	return th.Call(args[2], st)
}

var queryBlockParams = params("(query, block = false)")

var _ = exportMethods(&TranMethods)

var _ = method(tran_Complete, "()")

func tran_Complete(this Value) Value {
	this.(*SuTran).Complete()
	return nil
}

var _ = method(tran_Conflict, "()")

func tran_Conflict(this Value) Value {
	return SuStr(this.(*SuTran).Conflict())
}

var _ = method(tran_Data, "()")

func tran_Data(this Value) Value {
	return this.(*SuTran).Data()
}

var _ = method(tran_EndedQ, "()")

func tran_EndedQ(this Value) Value {
	return SuBool(this.(*SuTran).Ended())
}

var _ = method(tran_Query, "(@args)")

func tran_Query(th *Thread, as *ArgSpec, this Value, args []Value) Value {
	query, args := extractQuery(th, &queryBlockParams, as, args)
	mustNotBeAction(query)
	q := this.(*SuTran).Query(th, query)
	if args[1] == False {
		return q
	}
	// block form
	defer func() {
		if !this.(*SuTran).Ended() {
			q.Close()
		}
	}()
	return th.Call(args[1], q)
}

var _ = method(tran_QueryDo, "(@args)")

func tran_QueryDo(th *Thread, as *ArgSpec, this Value, args []Value) Value {
	query, _ := extractQuery(th, &queryParams, as, args)
	return IntVal(this.(*SuTran).Action(th, query))
}

var _ = method(tran_Query1, "(@args)")

func tran_Query1(th *Thread, as *ArgSpec, this Value, args []Value) Value {
	return tranQueryOne(th, this.(*SuTran), as, args, Only)
}

var _ = method(tran_QueryFirst, "(@args)")

func tran_QueryFirst(th *Thread, as *ArgSpec, this Value, args []Value) Value {
	return tranQueryOne(th, this.(*SuTran), as, args, Next)
}

var _ = method(tran_QueryLast, "(@args)")

func tran_QueryLast(th *Thread, as *ArgSpec, this Value, args []Value) Value {
	return tranQueryOne(th, this.(*SuTran), as, args, Prev)
}

var _ = method(tran_ReadCount, "()")

func tran_ReadCount(this Value) Value {
	return IntVal(this.(*SuTran).ReadCount())
}

var _ = method(tran_Rollback, "()")

func tran_Rollback(this Value) Value {
	this.(*SuTran).Rollback()
	return nil
}

var _ = method(tran_UpdateQ, "()")

func tran_UpdateQ(this Value) Value {
	return SuBool(this.(*SuTran).Updatable())
}

var _ = method(tran_WriteCount, "()")

func tran_WriteCount(this Value) Value {
	return IntVal(this.(*SuTran).WriteCount())
}

var _ = method(tran_Asof, "(asof = false)")

func tran_Asof(this, arg Value) Value {
	return this.(*SuTran).Asof(arg)
}

func tranQueryOne(th *Thread, st *SuTran, as *ArgSpec, args []Value, dir Dir) Value {
	query, _ := extractQuery(th, &queryParams, as, args)
	row, hdr, table := st.GetRow(th, query, dir)
	if row == nil {
		return False
	}
	return SuRecordFromRow(row, hdr, table, st)
}

var requestRegex = regex.Compile(`(?i)\A(insert|delete|update)\>`)

func mustNotBeAction(query string) {
	if requestRegex.Matches(query) {
		panic("transaction.Query: use QueryDo for insert, delete, or update requests")
	}
}
