package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin("Transaction(read=false, update=false, block=false)",
	func(th *Thread, args []Value) Value {
		read := ToBool(args[0])
		update := ToBool(args[1])
		if read == true && update == true {
			panic("usage: Transaction(read:) or Transaction(update:)")
		}
		itran := th.Dbms().Transaction(update)
		st := NewSuTran(itran, update)
		if args[2] == False {
			return st
		}
		// block form
		defer func() {
			if e := recover(); e != nil && e != BlockReturn {
				st.Rollback()
				panic(e)
			} else {
				st.Complete()
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
		"Ended?": method0(func(this Value) Value {
			return SuBool(this.(*SuTran).Ended())
		}),
		"Query": methodRaw("(@args)",
			func(th *Thread, as *ArgSpec, this Value, args []Value) Value {
				query, args := extractQuery(th, queryBlockParams, as, args)
				q := this.(*SuTran).Query(query)
				if args[1] == False {
					return q
				}
				// block form
				defer func() {
					q.Close()
				}()
				return th.Call(args[1], q)
			}),
		"QueryDo": methodRaw("(@args)",
			func(th *Thread, as *ArgSpec, this Value, args []Value) Value {
				query, _ := extractQuery(th, queryParams, as, args)
				return IntVal(this.(*SuTran).Request(query))
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
	row, hdr := st.GetRow(query, dir)
	if row == nil {
		return False
	}
	return SuRecordFromRow(row, hdr, st)
}
