// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/types"
	"github.com/apmckinlay/gsuneido/util/dnum"
)

func NewSuQueryNode(q Query) Value {
	return SuQueryNode{q: q}
}

type SuQueryNode struct {
	ValueBase[SuQueryNode]
	q Query
}

func (SuQueryNode) Type() types.Type {
	return types.QueryNode
}

func (SuQueryNode) Equal(any) bool {
	return false
}

func (SuQueryNode) SetConcurrent() {
	// read-only so nothing to do
}

func (n SuQueryNode) Get(_ *Thread, key Value) Value {
	return n.q.ValueGet(key)
}

func qryBase(q Query, key Value) Value {
	switch StrOrEmpty(key) {
	case "string":
		return SuStr(q.String())
	case "String":
		return SuStr(String(q))
	case "nrows":
		n, _ := q.Nrows()
		return IntVal(n)
	case "pop":
		_, p := q.Nrows()
		return IntVal(p)
	case "fast1":
		return SuBool(q.fastSingle())
	case "nchild":
		return Zero // overridden by Query1 and Query2
	case "frac":
		return SuDnum{Dnum: dnum.FromFloat(q.Metrics().frac)}
	case "fixcost":
		return IntVal(q.Metrics().fixcost)
	case "varcost":
		return IntVal(q.Metrics().varcost)
	case "cost":
		m := q.Metrics()
		return IntVal(m.fixcost + m.varcost)
	case "costself":
		return IntVal(q.Metrics().costself)
	case "tget":
		return IntVal(q.Metrics().tget)
	case "tgetself":
		return IntVal(q.Metrics().tgetself)
	case "ngets":
		return IntVal(q.Metrics().ngets)
	case "nsels":
		return IntVal(q.Metrics().nsels)
	case "nlooks":
		return IntVal(q.Metrics().nlooks)
	}
	return nil
}

func (tbl *Table) ValueGet(key Value) Value {
	switch StrOrEmpty(key) {
	case "type":
		return SuStr("table")
	case "name":
		return SuStr(tbl.name)
	}
	return qryBase(tbl, key)
}

func (ts *Tables) ValueGet(key Value) Value {
	switch StrOrEmpty(key) {
	case "type":
		return SuStr("table")
	case "name":
		return SuStr("tables")
	}
	return qryBase(ts, key)
}

func (cs *Columns) ValueGet(key Value) Value {
	switch StrOrEmpty(key) {
	case "type":
		return SuStr("table")
	case "name":
		return SuStr("columns")
	}
	return qryBase(cs, key)
}

func (is *Indexes) ValueGet(key Value) Value {
	switch StrOrEmpty(key) {
	case "type":
		return SuStr("table")
	case "name":
		return SuStr("indexes")
	}
	return qryBase(is, key)
}

func (vs *Views) ValueGet(key Value) Value {
	switch StrOrEmpty(key) {
	case "type":
		return SuStr("table")
	case "name":
		return SuStr("views")
	}
	return qryBase(vs, key)
}

func (his *History) ValueGet(key Value) Value {
	switch StrOrEmpty(key) {
	case "type":
		return SuStr("table")
	case "name":
		return SuStr("history")
	}
	return qryBase(his, key)
}

func (st *StatsTable) ValueGet(key Value) Value {
	switch StrOrEmpty(key) {
	case "type":
		return SuStr("table")
	case "name":
		return SuStr("dbstats")
	}
	return qryBase(st, key)
}

func (no *Nothing) ValueGet(key Value) Value {
	switch StrOrEmpty(key) {
	case "type":
		return SuStr("nothing")
	}
	return qryBase(no, key)
}

func (pn *ProjectNone) ValueGet(key Value) Value {
	switch StrOrEmpty(key) {
	case "type":
		return SuStr("projectNone")
	}
	return qryBase(pn, key)
}

//-------------------------------------------------------------------

func query1(q Query, key Value) Value {
	switch StrOrEmpty(key) {
	case "source":
		return NewSuQueryNode(q.(q1i).Source())
	case "nchild":
		return One
	}
	return qryBase(q, key)
}

func (e *Extend) ValueGet(key Value) Value {
	switch StrOrEmpty(key) {
	case "type":
		return SuStr("extend")
	}
	return query1(e, key)
}

func (p *Project) ValueGet(key Value) Value {
	switch StrOrEmpty(key) {
	case "type":
		return SuStr("project")
	}
	return query1(p, key)
}

func (r *Rename) ValueGet(key Value) Value {
	switch StrOrEmpty(key) {
	case "type":
		return SuStr("rename")
	}
	return query1(r, key)
}

func (sort *Sort) ValueGet(key Value) Value {
	switch StrOrEmpty(key) {
	case "type":
		return SuStr("sort")
	}
	return query1(sort, key)
}

func (su *Summarize) ValueGet(key Value) Value {
	switch StrOrEmpty(key) {
	case "type":
		return SuStr("summarize")
	}
	return query1(su, key)
}

func (ti *TempIndex) ValueGet(key Value) Value {
	switch StrOrEmpty(key) {
	case "type":
		return SuStr("tempindex")
	}
	return query1(ti, key)
}

func (w *Where) ValueGet(key Value) Value {
	switch StrOrEmpty(key) {
	case "type":
		return SuStr("where")
	case "expr":
		return w.expr
	}
	return query1(w, key)
}

func (v *View) ValueGet(key Value) Value {
	switch StrOrEmpty(key) {
	case "type":
		return SuStr("view")
	case "name":
		return SuStr(v.name)
	}
	return query1(v, key)
}

//-------------------------------------------------------------------

var Two = IntVal(2)

func query2(q Query, key Value) Value {
	switch StrOrEmpty(key) {
	case "source1":
		return NewSuQueryNode(q.(q2i).Source())
	case "source2":
		return NewSuQueryNode(q.(q2i).Source2())
	case "nchild":
		return Two
	}
	return qryBase(q, key)
}

func (u *Union) ValueGet(key Value) Value {
	switch StrOrEmpty(key) {
	case "type":
		return SuStr("union")
	}
	return query2(u, key)
}

func (it *Intersect) ValueGet(key Value) Value {
	switch StrOrEmpty(key) {
	case "type":
		return SuStr("intersect")
	}
	return query2(it, key)
}

func (m *Minus) ValueGet(key Value) Value {
	switch StrOrEmpty(key) {
	case "type":
		return SuStr("minus")
	}
	return query2(m, key)
}

func (t *Times) ValueGet(key Value) Value {
	switch StrOrEmpty(key) {
	case "type":
		return SuStr("times")
	}
	return query2(t, key)
}

func (jn *Join) ValueGet(key Value) Value {
	switch StrOrEmpty(key) {
	case "type":
		return SuStr("join")
	}
	return query2(jn, key)
}

func (lj *LeftJoin) ValueGet(key Value) Value {
	switch StrOrEmpty(key) {
	case "type":
		return SuStr("leftjoin")
	}
	return query2(lj, key)
}

func (sj *SemiJoin) ValueGet(key Value) Value {
	switch StrOrEmpty(key) {
	case "type":
		return SuStr("semijoin")
	}
	return query2(sj, key)
}
