package runtime

import "github.com/apmckinlay/gsuneido/runtime/types"

// SuTran is a database transaction
type SuTran struct {
	itran ITran
	state tstate
	CantConvert
}

func NewSuTran(itran ITran) *SuTran {
	return &SuTran{itran: itran}
}

type tstate byte

const (
	active tstate = iota
	committed
	commitFailed
	aborted
)

var _ Value = (*SuTran)(nil)

func (*SuTran) Get(*Thread, Value) Value {
	panic("transaction does not support get")
}

func (*SuTran) Put(*Thread, Value, Value) {
	panic("transaction does not support put")
}

func (*SuTran) RangeTo(int, int) Value {
	panic("transaction does not support range")
}

func (*SuTran) RangeLen(int, int) Value {
	panic("transaction does not support range")
}

func (*SuTran) Hash() uint32 {
	panic("transaction hash not implemented")
}

func (*SuTran) Hash2() uint32 {
	panic("transaction hash not implemented")
}

func (st *SuTran) Equal(other interface{}) bool {
	if t2, ok := other.(*SuTran); ok {
		return st == t2
	}
	return false
}

func (*SuTran) Compare(Value) int {
	panic("transaction compare not implemented")
}

func (*SuTran) Call(*Thread, *ArgSpec) Value {
	panic("can't call transaction")
}

func (*SuTran) Type() types.Type {
	return types.Transaction
}

func (st *SuTran) String() string {
	return st.itran.String()
}

// TranMethods is initialized by the builtin package
var TranMethods Methods

func (*SuTran) Lookup(_ *Thread, method string) Callable {
	return TranMethods[method]
}

func (st *SuTran) Complete() {
	if st.state == aborted || st.state == commitFailed {
		panic("can't Complete a transaction after failure or Rollback")
	}
	conflict := st.itran.Complete()
	if conflict == "" {
		st.state = committed
	} else {
		st.state = commitFailed
		panic("transaction.Complete failed: " + conflict)
	}
}

func (st *SuTran) Rollback() {
	if st.state == committed {
		panic("can't Rollback a transaction after Complete")
	}
	st.itran.Abort()
	st.state = aborted
}

func (st *SuTran) GetRow(query string, prev, single bool) (Row, *Header) {
	st.ckActive()
	return st.itran.Get(query, single, prev)
}

func (st *SuTran) Erase(adr int) {
	st.ckActive()
	st.itran.Erase(adr)
}

func (st *SuTran) Update(adr int, rec Record) {
	st.ckActive()
	st.itran.Update(adr, rec)
}

func (st *SuTran) Request(req string) int {
	st.ckActive()
	return st.itran.Request(req)
}

func (st *SuTran) Query(query string) *SuQuery {
	st.ckActive()
	iquery := st.itran.Query(query)
	return NewSuQuery(st, query, iquery)
}

func (st *SuTran) ckActive() {
	if st.state != active {
		panic("cannot use ended transaction")
	}
}
