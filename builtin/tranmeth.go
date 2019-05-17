package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin("Transaction(read=false, update=false, block=false)",
	func(th *Thread, args ...Value) Value {
		read := ToBool(args[0])
		update := ToBool(args[1])
		if read == update {
			panic("usage: Transaction(read:) or Transaction(update:)")
		}
		itran := th.Dbms().Transaction(update)
		st := NewSuTran(itran)
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
		return th.CallWithArgs(args[2], st)
	})

var queryParams = params("(query, block = false)")

func init() {
	TranMethods = Methods{
		"Complete": method0(func(this Value) Value {
			this.(*SuTran).Complete()
			return nil
		}),
		"Query": methodRaw("(@args)",
			func(th *Thread, as *ArgSpec, this Value, args ...Value) Value {
				query := buildQuery("Query", as, args)
				q := this.(*SuTran).Query(query)
				args = th.Args(queryParams, as)
				if args[1] == False {
					return q
				}
				// block form
				defer func() {
					q.Close()
				}()
				return th.CallWithArgs(args[1], q)
			}),
		"QueryDo": methodRaw("(@args)",
			func(th *Thread, as *ArgSpec, this Value, args ...Value) Value {
				query := buildQuery("QueryDo", as, args)
				return IntVal(this.(*SuTran).Request(query))
			}),
		"Query1": methodRaw("(@args)",
			func(th *Thread, as *ArgSpec, this Value, args ...Value) Value {
				return tranQueryOne(this.(*SuTran), "Query1", false, true, as, args...)
			}),
		"QueryFirst": methodRaw("(@args)",
			func(th *Thread, as *ArgSpec, this Value, args ...Value) Value {
				return tranQueryOne(this.(*SuTran), "QueryFirst", false, false, as, args...)
			}),
		"QueryLast": methodRaw("(@args)",
			func(th *Thread, as *ArgSpec, this Value, args ...Value) Value {
				return tranQueryOne(this.(*SuTran), "QueryLast", true, false, as, args...)
			}),
		"Rollback": method0(func(this Value) Value {
			this.(*SuTran).Rollback()
			return nil
		}),
	}
}

func tranQueryOne(st *SuTran, which string, prev, single bool,
	as *ArgSpec, args ...Value) Value {
	query := buildQuery(which, as, args)
	row, hdr := st.GetRow(query, prev, single)
	return SuRecordFromRow(row, hdr, st)
}
