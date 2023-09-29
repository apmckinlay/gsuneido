// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"github.com/apmckinlay/gsuneido/core/types"
)

//lint:ignore U1000 incorrect
type suTransaction struct{}

// SuTran is a database transaction
type SuTran struct {
	ValueBase[*suTransaction]
	itran     ITran
	data      *SuObject
	updatable bool
}

func NewSuTran(itran ITran, updatable bool) *SuTran {
	return &SuTran{itran: itran, updatable: updatable}
}

var _ Value = (*SuTran)(nil)

func (st *SuTran) Equal(other any) bool {
	return st == other
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

func (st *SuTran) Lookup(th *Thread, method string) Callable {
	return Lookup(th, TranMethods, gnTrans, method)
}

//-------------------------------------------------------------------

func (st *SuTran) Asof(val Value) Value {
	var asof int64
	switch val {
	case False:
		asof = 0
	case One:
		asof = +1
	case MinusOne:
		asof = -1
	default:
		asof = val.(SuDate).UnixMilli()
	}
	asof = st.itran.Asof(asof)
	if asof == 0 {
		return False
	}
	return SuDateFromUnixMilli(asof)
}

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

func (st *SuTran) Delete(th *Thread, table string, off uint64) {
	st.ckActive()
	st.itran.Delete(th, table, off)
}

func (st *SuTran) GetRow(th *Thread, query string, dir Dir) (Row, *Header, string) {
	st.ckActive()
	return st.itran.Get(th, query, dir, nil)
}

func (st *SuTran) Query(th *Thread, query string) *SuQuery {
	st.ckActive()
	iquery := st.itran.Query(query, nil)
	return NewSuQuery(th, st, query, iquery)
}

func (st *SuTran) ReadCount() int {
	return st.itran.ReadCount()
}

func (st *SuTran) Action(th *Thread, action string) int {
	st.ckActive()
	return st.itran.Action(th, action, nil)
}

func (st *SuTran) Rollback() {
	if err := st.itran.Abort(); err != "" {
		panic("transaction Rollback failed: " + err)
	}
}

func (st *SuTran) Updatable() bool {
	return st.updatable
}

func (st *SuTran) Update(th *Thread, table string, off uint64, rec Record) uint64 {
	st.ckActive()
	return st.itran.Update(th, table, off, rec)
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
