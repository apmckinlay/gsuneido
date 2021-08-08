// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/sset"
	"github.com/apmckinlay/gsuneido/util/strs"
)

type Extend struct {
	Query1
	cols     []string   // modified by Project.transform
	exprs    []ast.Expr // modified by Project.transform
	exprCols []string
	fixed    []Fixed
	hdr      *Header
}

func NewExtend(src Query, cols []string, exprs []ast.Expr) *Extend {
	e := &Extend{Query1: Query1{source: src}, cols: cols, exprs: exprs}
	e.checkDependencies()
	e.init()
	return e
}

func (e *Extend) checkDependencies() {
	avail := sset.Copy(e.source.Columns())
	for i := range e.cols {
		if e.exprs[i] != nil {
			ecols := e.exprs[i].Columns()
			if !sset.Subset(avail, ecols) {
				panic("extend: invalid column(s) in expressions: " +
					strs.Join(", ", sset.Difference(ecols, avail)))
			}
		}
		avail = append(avail, e.cols[i])
	}
}

func (e *Extend) init() {
	srcCols := e.source.Columns()
	if !sset.Disjoint(e.cols, srcCols) {
		panic("extend: column(s) already exist")
	}
	var cols []string
	for _, x := range e.exprs {
		if x != nil {
			cols = sset.Union(cols, x.Columns())
		}
	}
	e.exprCols = cols
}

func (e *Extend) String() string {
	s := parenQ2(e.source) + " EXTEND "
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
	return sset.Union(e.source.Columns(), e.cols)
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
	// combine Renames
	for e2, ok := e.source.(*Extend); ok; e2, ok = e.source.(*Extend) {
		e.cols = append(e2.cols, e.cols...)
		e.exprs = append(e2.exprs, e.exprs...)
		e.source = e2.source
		e.init()
	}
	e.source = e.source.Transform()
	return e
}

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
	i := strs.Index(e.cols, col)
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
	if e.fixed != nil { // once only
		return e.fixed
	}
	for i := 0; i < len(e.cols); i++ {
		if expr := e.exprs[i]; expr != nil {
			if c, ok := expr.(*ast.Constant); ok {
				e.fixed = append(e.fixed, NewFixed(e.cols[i], c.Val))
			}
		}
	}
	e.fixed = combineFixed(e.fixed, e.source.Fixed())
	return e.fixed
}

func (e *Extend) SingleTable() bool {
	return false
}

func (e *Extend) optimize(mode Mode, index []string) (Cost, interface{}) {
	if !sset.Disjoint(index, e.cols) {
		return impossible, nil
	}
	return Optimize(e.source, mode, index), nil
}

func (e *Extend) setApproach(index []string, _ interface{}, tran QueryTran) {
	e.source = SetApproach(e.source, index, tran)
}

// execution --------------------------------------------------------

func (e *Extend) Header() *Header {
	hdr := e.source.Header()
	cols := sset.Union(hdr.Columns, e.cols)
	flds := append(hdr.Fields, e.cols)
	return NewHeader(flds, cols)
}

func (e *Extend) Get(dir Dir) Row {
	row := e.source.Get(dir)
	return e.extendRow(row)
}

func (e *Extend) extendRow(row Row) Row {
	if row == nil {
		return nil // eof
	}
	if e.hdr == nil {
		e.hdr = e.Header()
	}
	var th Thread // ???
	context := ast.Context{T: &th,
		Rec: SuRecordFromRow(row, e.hdr, "", nil)}
	var rb RecordBuilder
	for i, col := range e.cols {
		if e := e.exprs[i]; e != nil {
			val := e.Eval(&context)
			rb.Add(val.(Packable))
			context.Rec.PreSet(SuStr(col), val)
		} else {
			rb.Add(EmptyStr.(Packable))
		}
	}
	return append(row, DbRec{Record: rb.Build()})
}

func (e *Extend) Select(cols, vals []string) {
	e.source.Select(cols, vals)
}

func (e *Extend) Lookup(cols, vals []string) Row {
	row := e.source.Lookup(cols, vals)
	return e.extendRow(row)
}
