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
	data      *SuObject
	updatable bool
}

func NewSuTran(itran ITran, updatable bool) *SuTran {
	return &SuTran{itran: itran, updatable: updatable}
}

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

func (st *SuTran) SetConcurrent() {
	//FIXME not thread safe
}

// TranMethods is initialized by the builtin package
var TranMethods Methods

var gnTrans = Global.Num("Transactions")

func (st *SuTran) Lookup(t *Thread, method string) Callable {
	return Lookup(t, TranMethods, gnTrans, method)
}

//-------------------------------------------------------------------

func (st *SuTran) Complete() {
	if conflict := st.itran.Complete(); conflict != "" {
		panic("transaction.Complete failed: " + conflict)
	}
}

func (st *SuTran) Conflict() string {
	return st.itran.Conflict()
}

func (st *SuTran) Ended() bool {
	return st.itran.Ended()
}

func (st *SuTran) Delete(table string, off uint64) {
	st.ckActive()
	st.itran.Delete(table, off)
}

func (st *SuTran) GetRow(query string, dir Dir) (Row, *Header, string) {
	st.ckActive()
	return st.itran.Get(query, dir)
}

func (st *SuTran) Query(th *Thread, query string) *SuQuery {
	st.ckActive()
	iquery := st.itran.Query(query)
	return NewSuQuery(th, st, query, iquery)
}

func (st *SuTran) ReadCount() int {
	return st.itran.ReadCount()
}

func (st *SuTran) Action(action string) int {
	st.ckActive()
	return st.itran.Action(action)
}

func (st *SuTran) Rollback() {
	if err := st.itran.Abort(); err != "" {
		panic("transaction Rollback failed: " + err)
	}
}

func (st *SuTran) Updatable() bool {
	return st.updatable
}

func (st *SuTran) Update(table string, off uint64, rec Record) uint64 {
	st.ckActive()
	return st.itran.Update(table, off, rec)
}

func (st *SuTran) WriteCount() int {
	return st.itran.WriteCount()
}

func (st *SuTran) ckActive() {
	if st.itran.Ended() {
		panic("can't use ended transaction")
	}
}

func (st *SuTran) Data() *SuObject {
	if st.data == nil {
		st.data = &SuObject{}
	}
	return st.data
}
