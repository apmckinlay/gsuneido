// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"strings"

	"github.com/apmckinlay/gsuneido/runtime/types"
)

// SuQueryCursor is the common base for SuQuery and SuCursor
type SuQueryCursor struct {
	CantConvert
	// which is either "Cursor" or "Query"
	which string
	query string
	iqc   IQueryCursor
	eof   Dir
}

func newQueryCursor(which string, query string, iqc IQueryCursor) *SuQueryCursor {
	return &SuQueryCursor{which: which, query: query, iqc: iqc}
}

func (qc *SuQueryCursor) Get(*Thread, Value) Value {
	panic(qc.which + " does not support get")
}

func (qc *SuQueryCursor) Put(*Thread, Value, Value) {
	panic(qc.which + " does not support put")
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
	return qc.iqc.Keys()
}

func (qc *SuQueryCursor) Order() Value {
	return qc.iqc.Order()
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

func NewSuQuery(tran *SuTran, query string, iquery IQuery) *SuQuery {
	return &SuQuery{*newQueryCursor("Query", query, iquery), tran}
}

var _ Value = (*SuQuery)(nil)

func (q *SuQuery) Equal(other interface{}) bool {
	if q2, ok := other.(*SuQuery); ok {
		return q == q2
	}
	return false
}

func (*SuQuery) Type() types.Type {
	return types.Query
}

// QueryMethods is initialized by the builtin package
var QueryMethods Methods

func (*SuQuery) Lookup(_ *Thread, method string) Callable {
	return QueryMethods[method]
}

func (q *SuQuery) GetRec(dir Dir) Value {
	if dir == q.eof {
		return False
	}
	if q.tran.Ended() {
		panic("cannot use a completed transaction")
	}
	row := q.iqc.(IQuery).Get(dir)
	if row == nil {
		q.eof = dir
		return False
	}
	q.eof = 0
	return SuRecordFromRow(row, q.iqc.Header(), q.tran)
}

func (q *SuQuery) Output(th *Thread, ob Container) {
	rec := ob.ToRecord(th, q.iqc.Header())
	q.iqc.(IQuery).Output(rec)
}

// ------------------------------------------------------------------

// SuCursor is a database cursor
type SuCursor struct {
	SuQueryCursor
}

func NewSuCursor(query string, icursor ICursor) *SuCursor {
	return &SuCursor{*newQueryCursor("Cursor", query, icursor)}
}

func (q *SuCursor) Equal(other interface{}) bool {
	if q2, ok := other.(*SuCursor); ok {
		return q == q2
	}
	return false
}

func (*SuCursor) Type() types.Type {
	return types.Cursor
}

// CursorMethods is initialized by the builtin package
var CursorMethods Methods

func (*SuCursor) Lookup(_ *Thread, method string) Callable {
	if f, ok := CursorMethods[method]; ok {
		return f
	}
	return QueryMethods[method]
}

func (q *SuCursor) GetRec(tran *SuTran, dir Dir) Value {
	if dir == q.eof {
		return False
	}
	row := q.iqc.(ICursor).Get(tran.itran, dir)
	if row == nil {
		q.eof = dir
		return False
	}
	q.eof = 0
	return SuRecordFromRow(row, q.iqc.Header(), tran)
}
