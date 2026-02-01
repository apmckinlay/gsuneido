// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"slices"

	"github.com/apmckinlay/gsuneido/compile/ast"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/set"
	"github.com/apmckinlay/gsuneido/util/slc"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/tsc"
)

type Extend struct {
	t        QueryTran
	ctx      ast.RowContext
	cols     []string
	exprs    []ast.Expr
	exprCols []string // columns used in exprs
	physical []string // cols with exprs
	selCols  []string
	selVals  []string
	srcFlds  []string
	fwd      map[int]string
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
	e.srcFlds = src.Header().Physical()
	e.setNrows(src.Nrows())
	e.rowSiz.Set(e.getRowSize())
	e.fast1.Set(src.fastSingle())
	e.singleTbl.Set(!e.hasExprs && src.SingleTable())
	e.lookCost.Set(src.lookupCost())

	for i, expr := range e.exprs {
		if c, ok := expr.(*ast.Constant); ok {
			c.Packed = Pack(c.Val.(Packable))
		}
		if id, ok := expr.(*ast.Ident); ok && !e.header.HasField(id.Name) {
			assert.That(id.Name != e.cols[i])
			if e.fwd == nil {
				e.fwd = make(map[int]string)
			}
			e.fwd[i] = string(rune(PackForward)) + id.Name
			// fmt.Println("Extend: forward", e.cols[i], "=", id)
		}
	}
	return e
}

func (e *Extend) checkDependencies() {
	avail := slc.Clone(e.source.Columns())
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
	s := "extend "
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
		for i := range len(e.cols) {
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
	if e.source.fastSingle() {
		index = e.filterSourceIndex(index)
	} else if !set.Disjoint(index, e.cols) {
		return impossible, impossible, nil
	}
	fixcost, varcost := Optimize(e.source, mode, index, frac)
	return fixcost, varcost, nil
}

func (e *Extend) setApproach(index []string, frac float64, _ any, tran QueryTran) {
	if e.source.fastSingle() {
		index = e.filterSourceIndex(index)
	}
	e.source = SetApproach(e.source, index, frac, tran)
	e.header = e.getHeader()
	e.ctx.Hdr = e.header
}

// filterSourceIndex filters out extended columns from the index for the source query
// This is needed when the source is fastSingle and accepts any columns as index
func (e *Extend) filterSourceIndex(index []string) []string {
	filteredIndex := make([]string, 0, len(index))
	for _, col := range index {
		if !slices.Contains(e.cols, col) {
			filteredIndex = append(filteredIndex, col)
		}
	}
	return filteredIndex
}

// execution --------------------------------------------------------

func (e *Extend) Get(th *Thread, dir Dir) Row {
	defer func(t uint64) { e.tget += tsc.Read() - t }(tsc.Read())
	if e.conflict {
		return nil
	}
	for {
		row := e.source.Get(th, dir)
		if row == nil {
			return nil
		}
		if !e.hasExprs {
			e.ngets++
			return row
		}
		rec := e.extendRow(th, row)
		if e.filter(rec, th, row) {
			e.ngets++
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
	for i, expr := range e.exprs {
		if expr != nil {
			if c, ok := expr.(*ast.Constant); ok {
				rb.AddRaw(c.Packed)
			} else if f, ok := ast.IsField(expr, e.srcFlds); ok {
				fld := f
				rb.AddRaw(row.GetRawVal(e.header, fld, e.ctx.Th, e.ctx.Tran))
			} else if f, ok := e.fwd[i]; ok {
				rb.AddRaw(f)
			} else {
				// incrementally build record so extends can see previous ones
				e.ctx.Row = append(row, DbRec{Record: rb.Build()})
				val := expr.Eval(&e.ctx)
				rb.Add(val.(Packable))
			}
		}
	}
	return rb.Trim().Build()
}

func (e *Extend) filter(rec Record, th *Thread, row Row) bool {
	var extrow Row
	for i, col := range e.selCols {
		j := slices.Index(e.physical, col)
		var x string
		if j >= 0 {
			x = rec.GetRaw(j)
		}
		if j == -1 || (len(x) > 0 && x[0] == PackForward) {
			if extrow == nil {
				extrow = append(row, DbRec{Record: rec})
			}
			x = extrow.GetRawVal(e.header, col, th, e.ctx.Tran)
		}
		if x != e.selVals[i] {
			return false
		}
	}
	return true
}

func (e *Extend) Select(cols, vals []string) {
	// fmt.Println("Extend Select", cols, unpack(vals))
	e.nsels++
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
	e.nlooks++
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
	if !e.filter(rec, th, row) {
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

func (e *Extend) Simple(th *Thread) []Row {
	e.header = e.getHeader()
	e.ctx.Hdr = e.header
	rows := e.source.Simple(th)
	if e.hasExprs {
		for i, row := range rows {
			rows[i] = append(row, DbRec{Record: e.extendRow(th, row)})
		}
	}
	return rows
}

func (e *Extend) knowExactNrows() bool {
	return e.source.knowExactNrows()
}
