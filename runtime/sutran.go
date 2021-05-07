// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"github.com/apmckinlay/gsuneido/runtime/types"
)

// SuTran is a database transaction
type SuTran struct {
	CantConvert
	itran     ITran
	conflict  string
	data      *SuObject
	updatable bool
	state     tstate
}

func NewSuTran(itran ITran, updatable bool) *SuTran {
	return &SuTran{itran: itran, updatable: updatable}
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

func (*SuTran) GetPut(*Thread, Value, Value, func(x, y Value) Value, bool) Value {
	panic("transaction does not support update")
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
	st2, ok := other.(*SuTran)
	return ok && st == st2
}

func (*SuTran) Compare(Value) int {
	panic("transaction compare not implemented")
}

func (*SuTran) Call(*Thread, Value, *ArgSpec) Value {
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

var gnTrans = Global.Num("Transactions")

func (*SuTran) Lookup(t *Thread, method string) Callable {
	return Lookup(t, TranMethods, gnTrans, method)
}

//-------------------------------------------------------------------

func (st *SuTran) Complete() {
	if st.state == aborted || st.state == commitFailed {
		panic("can't Complete a transaction after failure or Rollback")
	}
	st.conflict = st.itran.Complete()
	if st.conflict == "" {
		st.state = committed
	} else {
		st.state = commitFailed
		panic("transaction.Complete failed: " + st.conflict)
	}
}

func (st *SuTran) Conflict() string {
	return st.conflict
}

func (st *SuTran) Ended() bool {
	return st.state != active
}

func (st *SuTran) Erase(off uint64) {
	st.ckActive()
	st.itran.Erase(off)
}

func (st *SuTran) GetRow(query string, dir Dir) (Row, *Header) {
	st.ckActive()
	return st.itran.Get(query, dir)
}

func (st *SuTran) Query(query string) *SuQuery {
	st.ckActive()
	iquery := st.itran.Query(query)
	return NewSuQuery(st, query, iquery)
}

func (st *SuTran) ReadCount() int {
	return st.itran.ReadCount()
}

func (st *SuTran) Action(action string) int {
	st.ckActive()
	return st.itran.Action(action)
}

func (st *SuTran) Rollback() {
	if st.state == committed {
		panic("can't Rollback transaction after Complete")
	}
	if st.state != aborted {
		st.itran.Abort()
		st.state = aborted
	}
}

func (st *SuTran) Updatable() bool {
	return st.updatable
}

func (st *SuTran) Update(off uint64, rec Record) uint64 {
	st.ckActive()
	return st.itran.Update(off, rec)
}

func (st *SuTran) WriteCount() int {
	return st.itran.WriteCount()
}

func (st *SuTran) ckActive() {
	if st.state != active {
		panic("can't use ended transaction")
	}
}

func (st *SuTran) Data() *SuObject {
	if st.data == nil {
		st.data = &SuObject{}
	}
	return st.data
}
