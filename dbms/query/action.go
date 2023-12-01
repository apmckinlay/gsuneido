// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"github.com/apmckinlay/gsuneido/compile/ast"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
)

func DoAction(th *Thread, ut *db19.UpdateTran, action string) int {
	a := ParseAction(action, ut, th.Sviews())
	return a.execute(th, ut)
}

//-------------------------------------------------------------------

type insertRecordAction struct {
	record *SuRecord
	query  Query
}

func (a *insertRecordAction) String() string {
	return "insert " + a.record.Show() + " into " + a.query.String()
}

func (a *insertRecordAction) execute(th *Thread, ut *db19.UpdateTran) int {
	a.query.SetTran(ut)
	rec := a.record.ToRecord(th, a.query.Header())
	a.query.Output(th, rec)
	return 1
}

//-------------------------------------------------------------------

// NOTE: doesn't execute rules or output _deps

type insertQueryAction struct {
	query Query
	table string
}

func (a *insertQueryAction) String() string {
	return "insert " + a.query.String() + " into " + a.table
}

func (a *insertQueryAction) execute(th *Thread, ut *db19.UpdateTran) int {
	qr, _, _ := Setup(a.query, ReadMode, ut)
	hdr := qr.Header()
	fields := ut.GetSchema(a.table).Columns
	n := 0
	for row := qr.Get(th, Next); row != nil; row = qr.Get(th, Next) {
		rb := RecordBuilder{}
		var tsField string
		for _, f := range fields {
			if f == "-" || strings.HasSuffix(f, "_deps") {
				rb.AddRaw("")
			} else if strings.HasSuffix(f, "_TS") {
				if tsField != "" {
					panic("multiple _TS fields not supported")
				}
				rb.Add(th.Timestamp())
			} else {
				rb.AddRaw(row.GetRaw(hdr, f))
			}
		}
		rec := rb.Trim().Build()
		ut.Output(th, a.table, rec)
		n++
	}
	return n
}

//-------------------------------------------------------------------

type updateAction struct {
	query Query
	cols  []string
	exprs []ast.Expr
}

func (a *updateAction) String() string {
	s := "update " + a.query.String() + " set "
	sep := ""
	for i := range a.cols {
		s += sep + a.cols[i] + " = " + a.exprs[i].String()
		sep = ", "
	}
	return s
}

func (a *updateAction) execute(th *Thread, ut *db19.UpdateTran) int {
	q := SetupKey(a.query, UpdateMode, ut)
	table := q.Updateable()
	if table == "" {
		panic("update: query not updateable")
	}
	hdr := q.Header()
	tran := MakeSuTran(ut)
	ctx := ast.Context{Th: th, Tran: tran, Hdr: hdr}
	n := 0
	prev := uint64(0)
	for row := q.Get(th, Next); row != nil; row = q.Get(th, Next) {
		// avoid getting stuck on the same record
		if row[0].Off == prev {
			continue
		}
		ctx.Row = row
		r := SuRecordFromRow(row, hdr, table, tran)
		for i, col := range a.cols {
			r.Put(th, SuStr(col), a.exprs[i].Eval(&ctx))
		}
		newrec := r.ToRecord(th, hdr)
		prev = ut.Update(th, table, row[0].Off, newrec)
		n++
	}
	return n
}

//-------------------------------------------------------------------

type deleteAction struct {
	query Query
}

func (a *deleteAction) String() string {
	return "delete " + a.query.String()
}

func (a *deleteAction) execute(th *Thread, ut *db19.UpdateTran) int {
	//TODO optimize deleting all records of table (but still check foreign keys)
	q, _, _ := Setup(a.query, UpdateMode, ut)
	table := q.Updateable()
	if table == "" {
		panic("delete: query not updateable")
	}
	n := 0
	prev := uint64(0)
	for row := q.Get(th, Next); row != nil; row = q.Get(th, Next) {
		if row[0].Off == prev {
			continue
		}
		prev = row[0].Off
		ut.Delete(th, table, row[0].Off)
		n++
	}
	return n

	// offs := []uint64{}
	// for row := q.Get(th, Next); row != nil; row = q.Get(th, Next) {
	// 	offs = append(offs, row[0].Off)
	// }
	// for _, off := range offs {
	// 	ut.Delete(table, off)
	// }
	// return len(offs)
}
