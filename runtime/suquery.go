// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"strings"

	"github.com/apmckinlay/gsuneido/runtime/types"
)

// SuQueryCursor is the common base for SuQuery and SuCursor
type SuQueryCursor struct {
	query string
	iqc   IQueryCursor
	eof   Dir
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
	Strategy(formatted bool) Value
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

func (qc *SuQueryCursor) Strategy(formatted bool) Value {
	return SuStr(qc.iqc.Strategy(formatted))
}

// ------------------------------------------------------------------

// SuQuery is a database query
type SuQuery struct {
	ValueBase[SuQuery]
	SuQueryCursor
	tran *SuTran
}

func NewSuQuery(th *Thread, tran *SuTran, query string, iquery IQuery) *SuQuery {
	return &SuQuery{tran: tran,
		SuQueryCursor: SuQueryCursor{query: query, iqc: iquery}}
}

var _ Value = (*SuQuery)(nil)

func (q *SuQuery) Equal(other any) bool {
	return q == other
}

func (*SuQuery) Type() types.Type {
	return types.Query
}

func (q *SuQuery) String() string {
	return "Query('" + q.query + "')"
}

func (*SuQuery) SetConcurrent() {
	// FIXME
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
	ValueBase[SuCursor]
	SuQueryCursor
}

func NewSuCursor(th *Thread, query string, icursor ICursor) *SuCursor {
	return &SuCursor{SuQueryCursor: SuQueryCursor{query: query, iqc: icursor}}
}

func (q *SuCursor) Equal(other any) bool {
	return q == other
}

func (*SuCursor) Type() types.Type {
	return types.Cursor
}

func (q *SuCursor) String() string {
	return "Cursor('" + q.query + "')"
}

func (*SuCursor) SetConcurrent() {
	// FIXME
}

// CursorMethods is initialized by the builtin package
var CursorMethods Methods

func (q *SuCursor) Lookup(_ *Thread, method string) Callable {
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
