// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"strings"

	"github.com/apmckinlay/gsuneido/core/types"
)

// SuQueryCursor ----------------------------------------------------

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

// SuQueryCursor is the common base for SuQuery and SuCursor
type SuQueryCursor struct {
	ckActive func()
	iqc      IQueryCursor
	query    string
	eof      Dir
	closed   bool
}

var _ ISuQueryCursor = (*SuQueryCursor)(nil)

func (qc *SuQueryCursor) Columns() Value {
	qc.ckActive()
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
	qc.ckActive()
	return SuObjectOfStrs(qc.iqc.Keys())
}

func (qc *SuQueryCursor) Order() Value {
	qc.ckActive()
	return SuObjectOfStrs(qc.iqc.Order())
}

func (qc *SuQueryCursor) Rewind() {
	qc.ckActive()
	qc.iqc.Rewind()
	qc.eof = 0
}

func (qc *SuQueryCursor) RuleColumns() Value {
	qc.ckActive()
	hdr := qc.iqc.Header()
	ob := &SuObject{}
	for _, col := range hdr.Rules() {
		ob.Add(SuStr(col))
	}
	return ob
}

func (qc *SuQueryCursor) Strategy(formatted bool) Value {
	qc.ckActive()
	return SuStr(qc.iqc.Strategy(formatted))
}

func (qc *SuQueryCursor) Close() {
	qc.ckActive()
	qc.closed = true
	qc.iqc.Close()
}

// SuQuery ------------------------------------------------------------

type SuQuery struct {
	ValueBase[SuQuery]
	tran *SuTran
	SuQueryCursor
}

func NewSuQuery(th *Thread, tran *SuTran, query string, iquery IQuery) *SuQuery {
	q := &SuQuery{tran: tran,
		SuQueryCursor: SuQueryCursor{query: query, iqc: iquery}}
	q.SuQueryCursor.ckActive = q.ckActive
	return q
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
	q.ckActive()
	if dir == q.eof {
		return False
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
	q.ckActive()
	rec := ob.ToRecord(th, q.iqc.Header())
	q.iqc.(IQuery).Output(th, rec)
}

func (q *SuQuery) ckActive() {
	q.tran.ckActive()
	if q.closed {
		panic("can't use closed query")
	}
}

// SuCursor ---------------------------------------------------------

type SuCursor struct {
	ValueBase[SuCursor]
	SuQueryCursor
}

func NewSuCursor(th *Thread, query string, icursor ICursor) *SuCursor {
	q := &SuCursor{SuQueryCursor: SuQueryCursor{query: query, iqc: icursor}}
	q.SuQueryCursor.ckActive = q.ckActive
	return q
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
	tran.ckActive()
	q.ckActive()
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

func (q *SuCursor) ckActive() {
	if q.closed {
		panic("can't use closed cursor")
	}
}
