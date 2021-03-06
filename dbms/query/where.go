// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/sset"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Where struct {
	Query1
	expr  *ast.Nary // And
	fixed []Fixed
	t     QueryTran
}

func (w *Where) Init() {
	w.Query1.Init()
	if !sset.Subset(w.source.Columns(), w.expr.Columns()) {
		panic("select: nonexistent columns: " + str.Join(", ",
			sset.Difference(w.expr.Columns(), w.source.Columns())))
	}
}

func (w *Where) SetTran(t QueryTran) {
	w.t = t
	w.source.SetTran(t)
}

func (w *Where) String() string {
	return paren(w.source) + " WHERE " + w.expr.Echo()
}

func (w *Where) Fixed() []Fixed {
	if w.fixed != nil { // once only
		return w.fixed
	}
	for _, e := range w.expr.Exprs {
		w.addFixed(e)
	}
	w.fixed = combineFixed(w.fixed, w.source.Fixed())
	return w.fixed
}

func (w *Where) addFixed(e ast.Expr) {
	// MAYBE: handle IN
	if b, ok := e.(*ast.Binary); ok && b.Tok == tok.Is {
		if id, ok := b.Lhs.(*ast.Ident); ok {
			if c, ok := b.Rhs.(*ast.Constant); ok {
				w.fixed = append(w.fixed,
					Fixed{col: id.Name, values: []runtime.Value{c.Val}})
			}
		}
	}
}

func (w *Where) nrows() int {
	return w.source.nrows() / 2 //TODO
}

func (w *Where) dataSize() int {
	return w.source.dataSize() / 2 //TODO
}

func (w *Where) Transform() Query {
	if len(w.expr.Exprs) == 0 {
		// remove empty where
		return w.source.Transform()
	}
	if lj := w.leftJoinToJoin(); lj != nil {
		// convert leftjoin to join
		w.source = &Join{Query2: Query2{Query1: Query1{source: lj.source},
			source2: lj.source2}}
	}
	moved := false
	for {
		switch q := w.source.(type) {
		case *Where:
			// combine where's
			w.expr.Exprs = append(q.expr.Exprs, w.expr.Exprs...)
			w.source = q.source
			continue
		case *Project:
			// move where before project
			w.source = q.source
			q.source = w
			return q.Transform()
		case *Rename:
			// move where before rename
			newExpr := renameExpr(w.expr, q.to, q.from)
			w.source = q.source
			if newExpr == w.expr {
				q.source = w
			} else {
				q.source = &Where{Query1: Query1{source: q.source},
					expr: newExpr.(*ast.Nary)}
			}
			return q.Transform()
		case *Extend:
			// move where before extend, unless it depends on rules
			var src1, rest []ast.Expr
			for _, e := range w.expr.Exprs {
				if q.needRule(e.Columns()) {
					rest = append(rest, e)
				} else {
					src1 = append(src1, replaceExpr(e, q.cols, q.exprs))
				}
			}
			if src1 != nil {
				q.source = &Where{Query1: Query1{source: q.source},
					expr: &ast.Nary{Tok: tok.And, Exprs: src1}}
			}
			if rest != nil {
				w.expr = &ast.Nary{Tok: tok.And, Exprs: rest}
			} else {
				moved = true
			}
		case *Summarize:
			// split where before & after summarize
			cols1 := q.source.Columns()
			var src1, rest []ast.Expr
			for _, e := range w.expr.Exprs {
				if sset.Subset(cols1, e.Columns()) {
					src1 = append(src1, e)
				} else {
					rest = append(rest, e)
				}
			}
			if src1 != nil {
				q.source = &Where{Query1: Query1{source: q.source},
					expr: &ast.Nary{Tok: tok.And, Exprs: src1}}
			}
			if rest != nil {
				w.expr = &ast.Nary{Tok: tok.And, Exprs: rest}
			} else {
				moved = true
			}
		case *Intersect:
			// distribute where over intersect
			q.source = &Where{Query1: Query1{source: q.source}, expr: w.expr}
			q.source2 = &Where{Query1: Query1{source: q.source2}, expr: w.expr}
			moved = true
		case *Minus:
			// distribute where over minus
			q.source = &Where{Query1: Query1{source: q.source}, expr: w.expr}
			q.source2 = &Where{Query1: Query1{source: q.source2},
				expr: w.project(q.source2)}
			moved = true
		case *Union:
			// distribute where over union
			q.source = &Where{Query1: Query1{source: q.source},
				expr: w.project(q.source)}
			q.source2 = &Where{Query1: Query1{source: q.source2},
				expr: w.project(q.source2)}
			moved = true
		case *Times:
			// split where over product
			moved = w.split(&q.Query2)
		case *Join:
			// split where over join
			moved = w.split(&q.Query2)
		case *LeftJoin:
			// split where over leftjoin (left side only)
			cols1 := q.source.Columns()
			var common, src1 []ast.Expr
			for _, e := range w.expr.Exprs {
				if sset.Subset(cols1, e.Columns()) {
					src1 = append(src1, e)
				} else {
					common = append(common, e)
				}
			}
			if src1 != nil {
				q.source = &Where{Query1: Query1{source: q.source},
					expr: &ast.Nary{Tok: tok.And, Exprs: src1}}
			}
			if common != nil {
				w.expr = &ast.Nary{Tok: tok.And, Exprs: common}
			} else {
				moved = true
			}
		}
		w.source = w.source.Transform()
		if moved {
			return w.source
		}
		return w
	}
}

func (w *Where) leftJoinToJoin() *LeftJoin {
	if lj, ok := w.source.(*LeftJoin); ok {
		cols := lj.source2.Columns() //TODO .Header().Columns()
		cols = sset.Difference(cols, lj.by)
		for _, e := range w.expr.Exprs {
			if sset.Subset(cols, e.Columns()) && ast.CantBeEmpty(e, cols) {
				return lj
			}
		}
	}
	return nil
}

func (w *Where) project(q Query) *ast.Nary {
	srcCols := q.Columns()
	exprCols := w.expr.Columns()
	missing := sset.Difference(exprCols, srcCols)
	return replaceExpr(w.expr, missing, nEmpty(len(missing))).(*ast.Nary)
}

var emptyConstant = ast.Constant{Val: runtime.EmptyStr}

func nEmpty(n int) []ast.Expr {
	list := make([]ast.Expr, n)
	for i := range list {
		list[i] = &emptyConstant
	}
	return list
}

func (w *Where) split(q2 *Query2) bool {
	cols1 := q2.source.Columns()
	cols2 := q2.source2.Columns()
	var common, src1, src2 []ast.Expr
	for _, e := range w.expr.Exprs {
		used := false
		if sset.Subset(cols1, e.Columns()) {
			src1 = append(src1, e)
			used = true
		}
		if sset.Subset(cols2, (e.Columns())) {
			src2 = append(src2, e)
			used = true
		}
		if !used {
			common = append(common, e)
		}
	}
	if src1 != nil {
		q2.source = &Where{Query1: Query1{source: q2.source},
			expr: &ast.Nary{Tok: tok.And, Exprs: src1}}
	}
	if src2 != nil {
		q2.source2 = &Where{Query1: Query1{source: q2.source2},
			expr: &ast.Nary{Tok: tok.And, Exprs: src2}}
	}
	if common != nil {
		w.expr = &ast.Nary{Tok: tok.And, Exprs: common}
	} else {
		return true
	}
	return false
}
