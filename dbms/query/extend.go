// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"slices"

	"github.com/apmckinlay/gsuneido/compile/ast"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Extend struct {
	t        QueryTran
	ctx      ast.Context
	cols     []string
	exprs    []ast.Expr
	exprCols []string // columns used in exprs
	physical []string // cols with exprs
	selCols  []string
	selVals  []string
	Query1
	hasExprs bool
	conflict bool
}

func NewExtend(src Query, cols []string, exprs []ast.Expr) *Extend {
	e := &Extend{Query1: Query1{source: src}, cols: cols, exprs: exprs}
	e.checkDependencies()
	srcCols := e.source.Columns()
	if !set.Disjoint(e.cols, srcCols) {
		panic("extend: column(s) already exist")
	}
	var exprCols []string
	for _, expr := range e.exprs {
		if expr != nil {
			e.hasExprs = true
			exprCols = set.Union(exprCols, expr.Columns())
		}
	}
	e.exprCols = exprCols
	e.physical = e.getPhysical()
	e.header = e.getHeader()
	e.keys = src.Keys()
	e.indexes = src.Indexes()
	e.setNrows(src.Nrows())
	e.rowSiz.Set(e.getRowSize())
	e.fast1.Set(src.fastSingle())
	e.singleTbl.Set(!e.hasExprs && src.SingleTable())
	e.lookCost.Set(src.lookupCost())
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

func (e *Extend) getPhysical() []string {
	if !e.hasExprs {
		return nil
	}
	physical := make([]string, 0, len(e.cols))
	for i, col := range e.cols {
		if e.exprs[i] != nil {
			physical = append(physical, col)
		}
	}
	return physical
}

func (e *Extend) getHeader() *Header {
	srchdr := e.source.Header()
	cols := append(srchdr.Columns, e.cols...)
	flds := srchdr.Fields
	if e.physical != nil {
		flds = append(flds, e.physical)
	}
	return NewHeader(flds, cols)
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

func (e *Extend) getRowSize() int {
	return e.source.rowSize() + len(e.cols)*16 // ???
}

func (e *Extend) Transform() Query {
	src := e.source.Transform()

	// remove empty Extends
	if len(e.cols) == 0 {
		return src
	}
	cols := e.cols
	exprs := e.exprs
	// combine Extends
	if e2, ok := src.(*Extend); ok {
		src = e2.source
		cols = append(e2.cols, cols...)
		exprs = append(e2.exprs, exprs...)
	}
	if _, ok := src.(*Nothing); ok {
		return NewNothing(e)
	}
	if src != e.source {
		return NewExtend(src, cols, exprs)
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
	if e.fixed == nil {
		e.fixed = append([]Fixed{}, e.source.Fixed()...) // non-nil copy
		for i := 0; i < len(e.cols); i++ {
			if expr := e.exprs[i]; expr != nil {
				switch expr := expr.(type) {
				case *ast.Constant: // col = <Constant>
					e.fixed = append(e.fixed, NewFixed(e.cols[i], expr.Val))
				case *ast.Ident: // col = <Ident>
					if v := getFixed(e.fixed, expr.Name); v != nil {
						e.fixed = append(e.fixed, Fixed{col: e.cols[i], values: v})
					}
				}
			}
		}
		assert.That(e.fixed != nil)
	}
	return e.fixed
}

func (e *Extend) optimize(mode Mode, index []string, frac float64) (
	Cost, Cost, any) {
	if !set.Disjoint(index, e.cols) {
		return impossible, impossible, nil
	}
	fixcost, varcost := Optimize(e.source, mode, index, frac)
	return fixcost, varcost, nil
}

func (e *Extend) setApproach(index []string, frac float64, _ any, tran QueryTran) {
	e.source = SetApproach(e.source, index, frac, tran)
	e.header = e.getHeader()
	e.ctx.Hdr = e.header
}

// execution --------------------------------------------------------

func (e *Extend) Get(th *Thread, dir Dir) Row {
	if e.conflict {
		return nil
	}
	for {
		row := e.source.Get(th, dir)
		if row == nil {
			return nil
		}
		if !e.hasExprs {
			return row
		}
		rec := e.extendRow(th, row)
		if e.filter(rec) {
			return append(row, DbRec{Record: rec})
		}
	}
}

func (e *Extend) extendRow(th *Thread, row Row) Record {
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
	return rb.Trim().Build()
}

func (e *Extend) filter(rec Record) bool {
	for i, col := range e.selCols {
		j := slices.Index(e.physical, col)
		x := rec.GetRaw(j)
		if x != e.selVals[i] {
			return false
		}
	}
	return true
}

func (e *Extend) Select(cols, vals []string) {
	// fmt.Println("Extend Select", cols, unpack(vals))
	e.conflict = false
	e.selCols, e.selVals = nil, nil
	if cols == nil && vals == nil {
		e.source.Select(nil, nil) // clear select
		return
	}
	satisfied, conflict := selectFixed(cols, vals, e.Fixed())
	if conflict {
		e.conflict = true
	} else if satisfied {
		e.source.Select(nil, nil) // clear select
	} else {
		e.source.Select(e.splitSelect(cols, vals))
	}
}

func (e *Extend) Lookup(th *Thread, cols, vals []string) Row {
	if conflictFixed(cols, vals, e.Fixed()) {
		return nil
	}
	defer func() {
		e.selCols, e.selVals = nil, nil
	}()
	srccols, srcvals := e.splitSelect(cols, vals)
	row := e.source.Lookup(th, srccols, srcvals)
	if row == nil {
		return nil
	}
	if !e.hasExprs {
		return row
	}
	rec := e.extendRow(th, row)
	if !e.filter(rec) {
		return nil
	}
	return append(row, DbRec{Record: rec})
}

func (e *Extend) splitSelect(cols, vals []string) ([]string, []string) {
	var ecols, evals, srccols, srcvals []string
	for i, col := range cols {
		if slices.Contains(e.cols, col) {
			ecols = append(ecols, col)
			evals = append(evals, vals[i])
		} else {
			srccols = append(srccols, col)
			srcvals = append(srcvals, vals[i])
		}
	}
	e.selCols, e.selVals = ecols, evals
	return srccols, srcvals
}
