// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/sset"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Extend struct {
	Query1
	cols     []string   // modified by Project.transform
	exprs    []ast.Expr // modified by Project.transform
	exprCols []string
	fixed    []Fixed
	hdr      *runtime.Header
}

func (e *Extend) Init() {
	e.Query1.Init()
	e.checkDependencies()
	e.init()
}

func (e *Extend) checkDependencies() {
	avail := sset.Copy(e.source.Columns())
	for i := range e.cols {
		if e.exprs[i] != nil {
			ecols := e.exprs[i].Columns()
			if !sset.Subset(avail, ecols) {
				panic("extend: invalid column(s) in expressions: " +
					str.Join(", ", sset.Difference(ecols, avail)))
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
	i := str.List(e.cols).Index(col)
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
				e.fixed = append(e.fixed,
					Fixed{col: e.cols[i], values: []runtime.Value{c.Val}})
			}
		}
	}
	e.fixed = combineFixed(e.fixed, e.source.Fixed())
	return e.fixed
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

func (e *Extend) Header() *runtime.Header {
	hdr := e.source.Header()
	cols := sset.Union(hdr.Columns, e.cols)
	flds := append(hdr.Fields, e.cols)
	return runtime.NewHeader(flds, cols)
}

func (e *Extend) Get(dir runtime.Dir) runtime.Row {
	row := e.source.Get(dir)
	if row == nil {
		return nil // eof
	}
	if e.hdr == nil {
		e.hdr = e.Header()
	}
	var th runtime.Thread // ???
	context := ast.Context{T: &th,
		Rec: runtime.SuRecordFromRow(row, e.hdr, nil)}
	var rb runtime.RecordBuilder
	for i, col := range e.cols {
		if e := e.exprs[i]; e != nil {
			val := e.Eval(&context)
			rb.Add(val.(runtime.Packable))
			context.Rec.PreSet(runtime.SuStr(col), val)
		}
	}
	return append(row, runtime.DbRec{Record: rb.Build()})
}

func (e *Extend) Output(rec runtime.Record) {
	e.source.Output(rec)
}
