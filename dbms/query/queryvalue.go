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
	switch key {
	case SuStr("string"):
		return SuStr(q.String())
	case SuStr("String"):
		return SuStr(String(q))
	case SuStr("nrows"):
		n, _ := q.Nrows()
		return IntVal(n)
	case SuStr("pop"):
		_, p := q.Nrows()
		return IntVal(p)
	case SuStr("fast1"):
		return SuBool(q.fastSingle())
	case SuStr("nchild"):
		return Zero // overridden by Query1 and Query2
	case SuStr("frac"):
		return SuDnum{Dnum: dnum.FromFloat(q.Metrics().frac)}
	case SuStr("fixcost"):
		return IntVal(q.Metrics().fixcost)
	case SuStr("varcost"):
		return IntVal(q.Metrics().varcost)
	case SuStr("cost"):
		m := q.Metrics()
		return IntVal(m.fixcost + m.varcost)
	case SuStr("costself"):
		return IntVal(int(q.Metrics().costself))
	case SuStr("tget"):
		return Int64Val(int64(q.Metrics().tget))
	case SuStr("tgetself"):
		return Int64Val(int64(q.Metrics().tgetself))
	case SuStr("ngets"):
		return IntVal(int(q.Metrics().ngets))
	case SuStr("nsels"):
		return IntVal(int(q.Metrics().nsels))
	case SuStr("nlooks"):
		return IntVal(int(q.Metrics().nlooks))
	}
	return nil
}

func (tbl *Table) ValueGet(key Value) Value {
	switch key {
	case SuStr("type"):
		return SuStr("table")
	case SuStr("name"):
		return SuStr(tbl.name)
	}
	return qryBase(tbl, key)
}

func (ts *Tables) ValueGet(key Value) Value {
	switch key {
	case SuStr("type"):
		return SuStr("table")
	case SuStr("name"):
		return SuStr("tables")
	}
	return qryBase(ts, key)
}

func (tl *TablesLookup) ValueGet(key Value) Value {
	switch key {
	case SuStr("type"):
		return SuStr("tablelookup")
	case SuStr("name"):
		return SuStr(tl.table)
	}
	return qryBase(tl, key)
}

func (cs *Columns) ValueGet(key Value) Value {
	switch key {
	case SuStr("type"):
		return SuStr("table")
	case SuStr("name"):
		return SuStr("columns")
	}
	return qryBase(cs, key)
}

func (is *Indexes) ValueGet(key Value) Value {
	switch key {
	case SuStr("type"):
		return SuStr("table")
	case SuStr("name"):
		return SuStr("indexes")
	}
	return qryBase(is, key)
}

func (vs *Views) ValueGet(key Value) Value {
	switch key {
	case SuStr("type"):
		return SuStr("table")
	case SuStr("name"):
		return SuStr("views")
	}
	return qryBase(vs, key)
}

func (his *History) ValueGet(key Value) Value {
	switch key {
	case SuStr("type"):
		return SuStr("table")
	case SuStr("name"):
		return SuStr("history")
	}
	return qryBase(his, key)
}

func (no *Nothing) ValueGet(key Value) Value {
	switch key {
	case SuStr("type"):
		return SuStr("nothing")
	}
	return qryBase(no, key)
}

func (pn *ProjectNone) ValueGet(key Value) Value {
	switch key {
	case SuStr("type"):
		return SuStr("projectNone")
	}
	return qryBase(pn, key)
}

//-------------------------------------------------------------------

func query1(q Query, key Value) Value {
	switch key {
	case SuStr("source"):
		return NewSuQueryNode(q.(q1i).Source())
	case SuStr("nchild"):
		return One
	}
	return qryBase(q, key)
}

func (e *Extend) ValueGet(key Value) Value {
	switch key {
	case SuStr("type"):
		return SuStr("extend")
	}
	return query1(e, key)
}

func (p *Project) ValueGet(key Value) Value {
	switch key {
	case SuStr("type"):
		return SuStr("project")
	}
	return query1(p, key)
}

func (r *Rename) ValueGet(key Value) Value {
	switch key {
	case SuStr("type"):
		return SuStr("rename")
	}
	return query1(r, key)
}

func (sort *Sort) ValueGet(key Value) Value {
	switch key {
	case SuStr("type"):
		return SuStr("sort")
	}
	return query1(sort, key)
}

func (su *Summarize) ValueGet(key Value) Value {
	switch key {
	case SuStr("type"):
		return SuStr("summarize")
	}
	return query1(su, key)
}

func (ti *TempIndex) ValueGet(key Value) Value {
	switch key {
	case SuStr("type"):
		return SuStr("tempindex")
	}
	return query1(ti, key)
}

func (w *Where) ValueGet(key Value) Value {
	switch key {
	case SuStr("type"):
		return SuStr("where")
	case SuStr("expr"):
		return w.expr
	}
	return query1(w, key)
}

func (v *View) ValueGet(key Value) Value {
	switch key {
	case SuStr("type"):
		return SuStr("view")
	case SuStr("name"):
		return SuStr(v.name)
	}
	return query1(v, key)
}

//-------------------------------------------------------------------

var Two = IntVal(2)

func query2(q Query, key Value) Value {
	switch key {
	case SuStr("source1"):
		return NewSuQueryNode(q.(q2i).Source())
	case SuStr("source2"):
		return NewSuQueryNode(q.(q2i).Source2())
	case SuStr("nchild"):
		return Two
	}
	return qryBase(q, key)
}

func (u *Union) ValueGet(key Value) Value {
	switch key {
	case SuStr("type"):
		return SuStr("union")
	}
	return query2(u, key)
}

func (it *Intersect) ValueGet(key Value) Value {
	switch key {
	case SuStr("type"):
		return SuStr("intersect")
	}
	return query2(it, key)
}

func (m *Minus) ValueGet(key Value) Value {
	switch key {
	case SuStr("type"):
		return SuStr("minus")
	}
	return query2(m, key)
}

func (t *Times) ValueGet(key Value) Value {
	switch key {
	case SuStr("type"):
		return SuStr("times")
	}
	return query2(t, key)
}

func (jn *Join) ValueGet(key Value) Value {
	switch key {
	case SuStr("type"):
		return SuStr("join")
	}
	return query2(jn, key)
}

func (lj *LeftJoin) ValueGet(key Value) Value {
	switch key {
	case SuStr("type"):
		return SuStr("leftjoin")
	}
	return query2(lj, key)
}
