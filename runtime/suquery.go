// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"strings"

	"github.com/apmckinlay/gsuneido/runtime/types"
)

// SuQueryCursor is the common base for SuQuery and SuCursor
type SuQueryCursor struct {
	owner *Thread
	CantConvert
	// which is either "Cursor" or "Query"
	which string
	query string
	iqc   IQueryCursor
	eof   Dir
}

func (qc *SuQueryCursor) Get(*Thread, Value) Value {
	panic(qc.which + " does not support get")
}

func (qc *SuQueryCursor) Put(*Thread, Value, Value) {
	panic(qc.which + " does not support put")
}

func (qc *SuQueryCursor) GetPut(*Thread, Value, Value, func(x, y Value) Value, bool) Value {
	panic(qc.which + " does not support update")
}

func (qc *SuQueryCursor) RangeTo(int, int) Value {
	panic(qc.which + " does not support range")
}

func (qc *SuQueryCursor) RangeLen(int, int) Value {
	panic(qc.which + " does not support range")
}

func (qc *SuQueryCursor) Hash() uint32 {
	panic(qc.which + " hash not implemented")
}

func (qc *SuQueryCursor) Hash2() uint32 {
	panic(qc.which + " hash not implemented")
}

func (qc *SuQueryCursor) Compare(Value) int {
	panic(qc.which + " compare not implemented")
}

func (qc *SuQueryCursor) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call " + qc.which)
}

func (qc *SuQueryCursor) String() string {
	return qc.which + "('" + qc.query + "')"
}

func (*SuQueryCursor) SetConcurrent() {
	// allows multiple threads to reference but only owner can use
}

//-------------------------------------------------------------------

// ISuQueryCursor is the common interface to SuQuery and SuCursor
type ISuQueryCursor interface {
	Close()
	Columns() Value
	Keys() Value
	Order() Value
	Rewind()
	RuleColumns() Value
	Strategy() Value
}

var _ ISuQueryCursor = (*SuQueryCursor)(nil)

func (qc *SuQueryCursor) Close() {
	qc.iqc.Close()
}

func (qc *SuQueryCursor) Columns() Value {
	hdr := qc.iqc.Header()
	ob := &SuObject{}
	for _, col := range hdr.Columns {
		if !strings.HasSuffix(col, "_deps") {
			ob.Add(SuStr(col))
		}
	}
	return ob
}

func (qc *SuQueryCursor) Keys() Value {
	return SuObjectOfStrs(qc.iqc.Keys())
}

func (qc *SuQueryCursor) Order() Value {
	return SuObjectOfStrs(qc.iqc.Order())
}

func (qc *SuQueryCursor) Rewind() {
	qc.iqc.Rewind()
	qc.eof = 0
}

func (qc *SuQueryCursor) RuleColumns() Value {
	hdr := qc.iqc.Header()
	ob := &SuObject{}
	for _, col := range hdr.Rules() {
		ob.Add(SuStr(col))
	}
	return ob
}

func (qc *SuQueryCursor) Strategy() Value {
	return SuStr(qc.iqc.Strategy())
}

// ------------------------------------------------------------------

// SuQuery is a database query
type SuQuery struct {
	SuQueryCursor
	tran *SuTran
}

func NewSuQuery(th *Thread, tran *SuTran, query string, iquery IQuery) *SuQuery {
	return &SuQuery{tran: tran, SuQueryCursor: SuQueryCursor{
		owner: th, which: "Query", query: query, iqc: iquery}}
}

var _ Value = (*SuQuery)(nil)

func (q *SuQuery) Equal(other interface{}) bool {
	q2, ok := other.(*SuQuery)
	return ok && q == q2
}

func (*SuQuery) Type() types.Type {
	return types.Query
}

// QueryMethods is initialized by the builtin package
var QueryMethods Methods

func (q *SuQuery) Lookup(_ *Thread, method string) Callable {
	//FIXME concurrency
	// if q.owner != th {
	// 	panic("can't use a query from a different thread")
	// }
	return QueryMethods[method]
}

func (q *SuQuery) GetRec(th *Thread, dir Dir) Value {
	if dir == q.eof {
		return False
	}
	if q.tran.Ended() {
		panic("can't use ended transaction")
	}
	row, table := q.iqc.(IQuery).Get(th, dir)
	if row == nil {
		q.eof = dir
		return False
	}
	q.eof = 0
	return SuRecordFromRow(row, q.iqc.Header(), table, q.tran)
}

func (q *SuQuery) Output(th *Thread, ob Container) {
	rec := ob.ToRecord(th, q.iqc.Header())
	q.iqc.(IQuery).Output(th, rec)
}

// ------------------------------------------------------------------

// SuCursor is a database cursor
type SuCursor struct {
	SuQueryCursor
}

func NewSuCursor(th *Thread, query string, icursor ICursor) *SuCursor {
	return &SuCursor{SuQueryCursor: SuQueryCursor{
		owner: th, which: "Cursor", query: query, iqc: icursor}}
}

func (q *SuCursor) Equal(other interface{}) bool {
	q2, ok := other.(*SuCursor)
	return ok && q == q2
}

func (*SuCursor) Type() types.Type {
	return types.Cursor
}

// CursorMethods is initialized by the builtin package
var CursorMethods Methods

func (q *SuCursor) Lookup(_ *Thread, method string) Callable {
	//FIXME concurrency
	// if q.owner != th {
	// 	panic("can't use a query from a different thread")
	// }
	if f, ok := CursorMethods[method]; ok {
		return f
	}
	return QueryMethods[method]
}

func (q *SuCursor) GetRec(th *Thread, tran *SuTran, dir Dir) Value {
	if dir == q.eof {
		return False
	}
	row, table := q.iqc.(ICursor).Get(th, tran.itran, dir)
	if row == nil {
		q.eof = dir
		return False
	}
	q.eof = 0
	return SuRecordFromRow(row, q.iqc.Header(), table, tran)
}
