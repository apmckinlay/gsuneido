// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/regex"
)

var _ = builtin("Transaction(read=nil, update=nil, block=false)",
	func(th *Thread, args []Value) Value {
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
		if itran == nil {
			panic("too many active transactions")
		}
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
	})

var queryBlockParams = params("(query, block = false)")

func init() {
	TranMethods = Methods{
		"Complete": method0(func(this Value) Value {
			this.(*SuTran).Complete()
			return nil
		}),
		"Conflict": method0(func(this Value) Value {
			return SuStr(this.(*SuTran).Conflict())
		}),
		"Data": method0(func(this Value) Value {
			return this.(*SuTran).Data()
		}),
		"Ended?": method0(func(this Value) Value {
			return SuBool(this.(*SuTran).Ended())
		}),
		"Query": methodRaw("(@args)",
			func(th *Thread, as *ArgSpec, this Value, args []Value) Value {
				query, args := extractQuery(th, queryBlockParams, as, args)
				mustNotBeAction(query)
				q := this.(*SuTran).Query(query)
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
			}),
		"QueryDo": methodRaw("(@args)",
			func(th *Thread, as *ArgSpec, this Value, args []Value) Value {
				query, _ := extractQuery(th, queryParams, as, args)
				return IntVal(this.(*SuTran).Action(query))
			}),
		"Query1": methodRaw("(@args)",
			func(th *Thread, as *ArgSpec, this Value, args []Value) Value {
				return tranQueryOne(th, this.(*SuTran), as, args, Only)
			}),
		"QueryFirst": methodRaw("(@args)",
			func(th *Thread, as *ArgSpec, this Value, args []Value) Value {
				return tranQueryOne(th, this.(*SuTran), as, args, Next)
			}),
		"QueryLast": methodRaw("(@args)",
			func(th *Thread, as *ArgSpec, this Value, args []Value) Value {
				return tranQueryOne(th, this.(*SuTran), as, args, Prev)
			}),
		"ReadCount": method0(func(this Value) Value {
			return IntVal(this.(*SuTran).ReadCount())
		}),
		"Rollback": method0(func(this Value) Value {
			this.(*SuTran).Rollback()
			return nil
		}),
		"Update?": method0(func(this Value) Value {
			return SuBool(this.(*SuTran).Updatable())
		}),
		"WriteCount": method0(func(this Value) Value {
			return IntVal(this.(*SuTran).WriteCount())
		}),
	}
}

func tranQueryOne(th *Thread, st *SuTran, as *ArgSpec, args []Value, dir Dir) Value {
	query, _ := extractQuery(th, queryParams, as, args)
	row, hdr, table := st.GetRow(query, dir)
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
