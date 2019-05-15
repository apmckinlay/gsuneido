package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin("Transaction(read=false, update=false, block=false)",
	func(t *Thread, args ...Value) Value {
		read := ToBool(args[0])
		update := ToBool(args[1])
		if read == update {
			panic("usage: Transaction(read:) or Transaction(update:)")
		}
		itran := t.Dbms().Transaction(update)
		st := NewSuTran(itran)
		if args[2] == False {
			return st
		}
		// block form
		defer func() {
			if e := recover(); e != nil && e != BlockReturn {
				st.Rollback()
			} else {
				st.Complete()
			}
		}()
		return t.CallWithArgs(args[2], st)
	})

func init() {
	TranMethods = Methods{
		"Complete": method0(func(this Value) Value {
			this.(*SuTran).Complete()
			return nil
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
