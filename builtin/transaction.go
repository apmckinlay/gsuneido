package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

var _ = builtin("Transaction(read=false, update=false, block=false)",
	func(t *Thread, args ...Value) Value {
		read := ToBool(args[0])
		update := ToBool(args[1])
		if read == update {
			panic("usage: Transaction(read:) or Transaction(update:)")
		}
		it := t.Dbms().Transaction(update)
		st := SuTran{it: it}
		if args[2] == False {
			return st
		}
		// block form
		defer func() {
			if e := recover(); e != nil && e != BlockReturn {
				st.rollback()
			} else {
				st.complete()
			}
		}()
		return t.CallWithArgs(args[2], st)
	})

// SuTran references a database transaction
type SuTran struct {
	it    ITran
	state tstate
	CantConvert
}

type tstate byte

const (
	active tstate = iota
	committed
	commitFailed
	aborted
)

var _ Value = (*SuTran)(nil)

func (SuTran) Get(*Thread, Value) Value {
	panic("transaction does not support get")
}

func (SuTran) Put(*Thread, Value, Value) {
	panic("transaction does not support put")
}

func (SuTran) RangeTo(int, int) Value {
	panic("transaction does not support range")
}

func (SuTran) RangeLen(int, int) Value {
	panic("transaction does not support range")
}

func (SuTran) Hash() uint32 {
	panic("transaction hash not implemented")
}

func (SuTran) Hash2() uint32 {
	panic("transaction hash not implemented")
}

func (st SuTran) Equal(other interface{}) bool {
	if t2, ok := other.(SuTran); ok {
		return st == t2
	}
	return false
}

func (SuTran) Compare(Value) int {
	panic("transaction compare not implemented")
}

func (SuTran) Call(*Thread, *ArgSpec) Value {
	panic("can't call transaction")
}

func (SuTran) Type() types.Type {
	return types.Transaction
}

func (st SuTran) String() string {
	return st.it.String()
}

var tranMethods = Methods{
	"Complete": method0(func(this Value) Value {
		this.(SuTran).complete()
		return nil
	}),
	"Query1": methodRaw("(@args)",
		func(th *Thread, as *ArgSpec, this Value, args ...Value) Value {
			return this.(SuTran).queryOne("Query1", false, true, as, args...)
		}),
	"QueryFirst": methodRaw("(@args)",
		func(th *Thread, as *ArgSpec, this Value, args ...Value) Value {
			return this.(SuTran).queryOne("QueryFirst", false, false, as, args...)
		}),
	"QueryLast": methodRaw("(@args)",
		func(th *Thread, as *ArgSpec, this Value, args ...Value) Value {
			return this.(SuTran).queryOne("QueryLast", true, false, as, args...)
		}),
	"Rollback": method0(func(this Value) Value {
		this.(SuTran).rollback()
		return nil
	}),
}

func (SuTran) Lookup(_ *Thread, method string) Callable {
	return tranMethods[method]
}

func (st SuTran) complete() {
	if st.state == aborted || st.state == commitFailed {
		panic("can't Complete a transaction after failure or Rollback")
	}
	conflict := st.it.Complete()
	if conflict == "" {
		st.state = committed
	} else {
		st.state = commitFailed
		panic("transaction.Complete failed: " + conflict)
	}
}

func (st SuTran) rollback() {
	if st.state == committed {
		panic("can't Rollback a transaction after Complete")
	}
	st.it.Abort()
	st.state = aborted
}

func (st SuTran) queryOne(which string, prev, single bool,
	as *ArgSpec, args ...Value) Value {
	query := buildQuery(which, as, args)
	row, hdr := st.it.Get(query, prev, single)
	// fmt.Println(hdr)
	// fmt.Println(row)
	return SuRecordFromRow(row, hdr)
}
