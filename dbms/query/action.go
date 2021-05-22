// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
	"github.com/apmckinlay/gsuneido/db19"
	. "github.com/apmckinlay/gsuneido/runtime"
)

func DoAction(ut *db19.UpdateTran, action string) int {
	a := ParseAction(action)
	return a.execute(ut)
}

type insertRecordAction struct {
	record *SuRecord
	query  Query
}

func (a *insertRecordAction) String() string {
	return "insert " + a.record.Show() + " into " + a.query.String()
}

func (a *insertRecordAction) execute(ut *db19.UpdateTran) int {
	a.query.SetTran(ut)
	var th Thread // ???
	rec := a.record.ToRecord(&th, a.query.Header())
	a.query.Output(rec)
	return 1
}

type insertQueryAction struct {
	query Query
	table string
}

func (a *insertQueryAction) String() string {
	return "insert " + a.query.String() + " into " + a.table
}

func (a *insertQueryAction) execute(ut *db19.UpdateTran) int {
	return 0 //TODO
}

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

func (a *updateAction) execute(ut *db19.UpdateTran) int {
	q := SetupKey(a.query, UpdateMode, ut)
	table := q.Updateable()
	if table == "" {
		panic("update: query not updateable")
	}
	hdr := q.Header()
	th := &Thread{}
	var prev uint64
	n := 0
	for row := q.Get(Next); row != nil; row = q.Get(Next) {
		// avoid getting stuck on the same record
		if row[0].Off == prev {
			continue
		}
		r := SuRecordFromRow(row, hdr, table, nil)
		context := &ast.Context{T: th, Rec: r}
		for i, col := range a.cols {
			r.Put(th, SuStr(col), a.exprs[i].Eval(context))
		}
		newrec := r.ToRecord(th, hdr)
		prev = ut.Update(table, row[0].Off, newrec)
		n++
	}
	return n
}

type deleteAction struct {
	query Query
}

func (a *deleteAction) String() string {
	return "delete " + a.query.String()
}

func (a *deleteAction) execute(ut *db19.UpdateTran) int {
	q, _ := Setup(a.query, UpdateMode, ut)
	table := q.Updateable()
	if table == "" {
		panic("delete: query not updateable")
	}
	n := 0
	for row := q.Get(Next); row != nil; row = q.Get(Next) {
		ut.Delete(table, row[0].Off)
		n++
	}
	return n
}
