// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/str"
	"golang.org/x/exp/slices"
)

type Extend struct {
	Query1
	cols     []string   // modified by Project.transform
	exprs    []ast.Expr // modified by Project.transform
	exprCols []string
	fixed    []Fixed
	hasExprs bool
	hdr      *Header
	t        QueryTran
	ctx      ast.Context
	conflict bool
}

func NewExtend(src Query, cols []string, exprs []ast.Expr) *Extend {
	e := &Extend{Query1: Query1{source: src}, cols: cols, exprs: exprs}
	e.checkDependencies()
	e.init()
	return e
}

func (e *Extend) checkDependencies() {
	avail := slices.Clone(e.source.Columns())
	for i := range e.cols {
		if e.exprs[i] != nil {
			ecols := e.exprs[i].Columns()
			if !set.Subset(avail, ecols) {
				panic("extend: invalid column(s) in expressions: " +
					str.Join(", ", set.Difference(ecols, avail)))
			}
		}
		avail = append(avail, e.cols[i])
	}
}

func (e *Extend) init() {
	srcCols := e.source.Columns()
	if !set.Disjoint(e.cols, srcCols) {
		panic("extend: column(s) already exist")
	}
	var cols []string
	for _, expr := range e.exprs {
		if expr != nil {
			e.hasExprs = true
			cols = set.Union(cols, expr.Columns())
		}
	}
	e.exprCols = cols
}

func (e *Extend) SetTran(t QueryTran) {
	e.t = t
	e.source.SetTran(t)
	e.ctx.Tran = nil
}

func (e *Extend) String() string {
	return parenQ2(e.source) + " " + e.stringOp()
}

func (e *Extend) stringOp() string {
	s := "EXTEND "
	sep := ""
	for i, c := range e.cols {
		s += sep + c
		sep = ", "
		if e.exprs[i] != nil {
			s += " = " + e.exprs[i].Echo()
		}
	}
	return s
}

func (e *Extend) Columns() []string {
	return set.Union(e.source.Columns(), e.cols)
}

func (e *Extend) rowSize() int {
	nsc := len(e.source.Columns())
	nc := len(e.Columns())
	return e.source.rowSize() * nc / nsc
}

func (e *Extend) Transform() Query {
	// remove empty Extends
	if len(e.cols) == 0 {
		return e.source.Transform()
	}
	// combine Extends
	for e2, ok := e.source.(*Extend); ok; e2, ok = e.source.(*Extend) {
		e.cols = append(e2.cols, e.cols...)
		e.exprs = append(e2.exprs, e.exprs...)
		e.source = e2.source
		e.init()
	}
	e.source = e.source.Transform()
	// propagate Nothing
	if _, ok := e.source.(*Nothing); ok {
		return NewNothing(e.Columns())
	}
	return e
}

// hasRules is used by Project transformExtend
func (e *Extend) hasRules() bool {
	for _, e := range e.exprs {
		if e == nil {
			return true
		}
	}
	return false
}

func (e *Extend) needRule(cols []string) bool {
	for _, col := range cols {
		if e.needRule2(col) {
			return true
		}
	}
	return false
}

func (e *Extend) needRule2(col string) bool {
	i := slices.Index(e.cols, col)
	if i == -1 {
		return false // fld is not a result of extend
	}
	if e.exprs[i] == nil {
		return true // direct dependency
	}
	exprdeps := e.exprs[i].Columns()
	return e.needRule(exprdeps) // recursive
}

func (e *Extend) Fixed() []Fixed {
	if e.fixed != nil {
		return e.fixed
	}
	fixed := append([]Fixed{}, e.source.Fixed()...) // copy
	for i := 0; i < len(e.cols); i++ {
		if expr := e.exprs[i]; expr != nil {
			if c, ok := expr.(*ast.Constant); ok {
				fixed = append(fixed, NewFixed(e.cols[i], c.Val))
			}
		}
	}
	return fixed
}

func (e *Extend) SingleTable() bool {
	if e.hasExprs {
		return false
	}
	return e.source.SingleTable()
}

func (e *Extend) optimize(mode Mode, index []string) (Cost, Cost, any) {
	if !set.Disjoint(index, e.cols) {
		return impossible, impossible, nil
	}
	fixcost, varcost := Optimize(e.source, mode, index)
	return fixcost, varcost, nil
}

func (e *Extend) setApproach(mode Mode, index []string, _ any, tran QueryTran) {
	e.source = SetApproach(e.source, mode, index, tran)
	e.hdr = e.Header() // cache for Get
	e.ctx.Hdr = e.hdr
	e.fixed = e.Fixed() // cache
}

// execution --------------------------------------------------------

func (e *Extend) Header() *Header {
	if e.hdr != nil {
		return e.hdr
	}
	hdr := e.source.Header()
	cols := append(hdr.Columns, e.cols...)
	flds := hdr.Fields
	if e.hasExprs {
		physical := make([]string, 0, len(cols))
		for i, col := range e.cols {
			if e.exprs[i] != nil {
				physical = append(physical, col)
			}
		}
		flds = append(hdr.Fields, physical)
	}
	return NewHeader(flds, cols)
}

func (e *Extend) Get(th *Thread, dir Dir) Row {
	if e.conflict {
		return nil
	}
	row := e.source.Get(th, dir)
	return e.extendRow(th, row)
}

func (e *Extend) extendRow(th *Thread, row Row) Row {
	if row == nil || !e.hasExprs {
		return row // eof
	}
	e.ctx.Th = th
	defer func() { e.ctx.Th = nil }()
	if e.ctx.Tran == nil {
		e.ctx.Tran = MakeSuTran(e.t)
	}
	var rb RecordBuilder
	for _, expr := range e.exprs {
		if expr != nil {
			// incrementally build record so extends can see previous ones
			e.ctx.Row = append(row, DbRec{Record: rb.Build()})
			val := expr.Eval(&e.ctx)
			rb.Add(val.(Packable))
		}
	}
	return append(row, DbRec{Record: rb.Build()})
}

func (e *Extend) Select(cols, vals []string) {
	fixed := e.Fixed()
	satisfied := true
	e.conflict = false
	for i, col := range cols {
		if fv := getFixed(fixed, col); len(fv) == 1 {
			if fv[0] != vals[i] {
				e.conflict = true
				break
			}
		} else {
			satisfied = false
		}
	}
	if e.conflict {
		return
	} else if satisfied {
		e.source.Select(nil, nil) // clear select
	} else {
		e.source.Select(cols, vals)
	}
}

func (e *Extend) Lookup(th *Thread, cols, vals []string) Row {
	row := e.source.Lookup(th, cols, vals)
	return e.extendRow(th, row)
}
