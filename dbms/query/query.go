// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"github.com/apmckinlay/gsuneido/util/str"
)

type Query interface {
	//TODO
	String() string
}

type Sort struct {
	Query1
	reverse bool
	columns []string
}

func (sort *Sort) String() string {
	s := sort.Query1.String() + " sort "
	if sort.reverse {
		s += "reverse "
	}
	return s + str.Join(", ", sort.columns...)
}

type Table struct {
	name string
}

func (tbl *Table) String() string {
	return tbl.name
}

type Query1 struct {
	source Query
}

func (q1 *Query1) String() string {
	return q1.source.String()
}

type Project struct {
	Query1
	columns []string
}

func (p *Project) String() string {
	return p.Query1.String() + " project " + str.Join(", ", p.columns...)
}

type Remove struct {
	Query1
	columns []string
}

func (r *Remove) String() string {
	return r.Query1.String() + " remove " + str.Join(", ", r.columns...)
}

type Query2 struct {
	Query1
	source2 Query
}

func (q2 *Query2) String(op string) string {
	return q2.Query1.String() + " " + op + " " + q2.source2.String()
}

type Intersect struct {
	Query2
}

func (i *Intersect) String() string {
	return i.Query2.String("intersect")
}

type Minus struct {
	Query2
}

func (m *Minus) String() string {
	return m.Query2.String("minus")
}

type Times struct {
	Query2
}

func (t *Times) String() string {
	return t.Query2.String("times")
}

type Union struct {
	Query2
}

func (u *Union) String() string {
	return u.Query2.String("union")
}

type Join struct {
	Query2
	by []string
}

func (j *Join) String() string {
	return j.string("join")
}

func (j *Join) string(op string) string {
	by := ""
	if len(j.by) > 0 {
		by = "by" + str.Join("(,)", j.by...) + " "
	}
	return j.Query1.String() + " " + op + " " + by + j.source2.String()
}

type LeftJoin struct {
	Join
}

func (lj *LeftJoin) String() string {
	return lj.string("leftjoin")
}

type Rename struct {
	Query1
	renames []Rename1
}

func (r *Rename) String() string {
	sep := ""
	var sb strings.Builder
	for _, ren := range r.renames {
		sb.WriteString(sep)
		sb.WriteString(ren.String())
		sep = ", "
	}
	return r.Query1.String() + " rename " + sb.String()
}

type Rename1 struct {
	from string
	to   string
}

func (rn Rename1) String() string {
	return rn.from + " to " + rn.to
}

type Summarize struct {
	Query1
	by   []string
	cols []string
	ops  []string
	ons  []string
}

func (su *Summarize) String() string {
	s := su.Query1.String() + " summarize "
	if len(su.by) > 0 {
		s += str.Join(", ", su.by...) + ", "
	}
	sep := ""
	for i := range su.cols {
		s += sep
		sep = ", "
		if su.cols[i] != "" {
			s += su.cols[i] + " = "
		}
		s += su.ops[i]
		if su.ops[i] != "count" {
			s += " " + su.ons[i]
		}
	}
	return s
}

type Where struct {
	Query1
	expr Expr
}

func (w *Where) String() string {
	return w.Query1.String() + " where " + w.expr.String()
}

type Extend struct {
	Query1
	cols  []string
	exprs []Expr
}

func (e *Extend) String() string {
	s := e.Query1.String() + " extend "
	sep := ""
	for i, c := range e.cols {
		s += sep + c
		sep = ", "
		if e.exprs[i] != nil {
			s += " = " + e.exprs[i].String()
		}
	}
	return s
}
