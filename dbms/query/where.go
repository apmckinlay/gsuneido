// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"math"

	"slices"

	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/trace"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/tsc"
)

// NOTE: Where source and expr should NOT be modified,
// instead, construct a new one with NewWhere

type Where struct {
	t       QueryTran
	colSels map[string][]span // from NewWhere, result of perField
	// tbl will be set if the source is a Table, nil otherwise
	tbl *Table
	// idxSel is for the chosen index
	idxSel *idxSel
	expr   *ast.Nary // And
	// curPtrng is idxSel.ptrngs[idxSelPos] adjusted by Select (selOrg, selEnd)
	curPtrng pointRange
	// selOrg and selEnd are set by Select
	selOrg string
	selEnd string
	ctx    ast.Context

	selectCols []string
	selectVals []string

	idxSels []idxSel // from optInit, result of perIndex
	Query1
	// idxSelPos is the current index in idxSel.ptrngs
	idxSelPos int

	nIn  int
	nOut int

	// conflict is true if the expression conflicts and selects nothing
	conflict bool
	// singleton is true if we know the result is at most one record
	// because there is a single point select on a key
	singleton bool
	selSet    bool

	// exprMore is whether expr has more than idxSels
	exprMore bool
	optInited
	optimized bool
}

type optInited byte

const (
	optInitNo optInited = iota
	optInitInProgress
	optInitYes
)

type whereApproach struct {
	index []string
	cost  Cost
}

func NewWhere(src Query, expr ast.Expr, t QueryTran) *Where {
	if !set.Subset(src.Columns(), expr.Columns()) {
		panic("where: nonexistent columns: " + str.Join(", ",
			set.Difference(expr.Columns(), src.Columns())))
	}
	if nary, ok := expr.(*ast.Nary); !ok || nary.Tok != tok.And {
		expr = &ast.Nary{Tok: tok.And, Exprs: []ast.Expr{expr}}
	}
	w := &Where{Query1: Query1{source: src}, expr: expr.(*ast.Nary), t: t}
	w.header = src.Header()
	w.rowSiz.Set(src.rowSize())
	w.singleTbl.Set(src.SingleTable())
	w.lookCost.Set(src.lookupCost())
	w.calcFixed()
	if !w.conflict {
		w.conflict = w.exprFalse()
	}
	if !w.conflict {
		fields := w.source.Header().Physical()
		w.colSels, w.exprMore = perField(w.expr.Exprs, fields)
		// fmt.Println("colSels", w.colSels)
		w.conflict = (w.colSels == nil)
	}
	return w
}

func (w *Where) SetTran(t QueryTran) {
	w.t = t
	w.source.SetTran(t)
	w.ctx.Tran = nil
}

func (w *Where) String() string {
	s := "where"
	if w.conflict {
		s += " /*NOTHING*/"
		if w.optimized {
			return s
		}
	}
	if w.optimized && w.singleton {
		s += "*1"
	}
	if len(w.expr.Exprs) > 0 {
		s += " " + w.expr.Echo()
	}
	if w.slow() {
		s += fmt.Sprintf(" /*SLOW %d->%d*/", w.nIn, w.nOut)
	}
	return s
}

// calcFixed sets w.fixed and may set w.conflict
func (w *Where) calcFixed() {
	efixed, conflict := w.exprsToFixed()
	if conflict {
		w.conflict = true
		return
	}
	fixed, none := combineFixed(w.source.Fixed(), efixed)
	if none {
		w.conflict = true
		return
	}
	w.fixed = fixed
}

func (w *Where) exprsToFixed() (fixed []Fixed, conflict bool) {
	for _, e := range w.expr.Exprs {
		fixed, conflict = addFixed(fixed, e)
		if conflict {
			return nil, true
		}
	}
	return fixed, false
}

func addFixed(fixed []Fixed, e ast.Expr) ([]Fixed, bool) {
	// MAYBE: handle OR, could use colSels
	if b, ok := e.(*ast.Binary); ok && (b.Tok == tok.Is || b.Tok == tok.Lte) {
		if id, ok := b.Lhs.(*ast.Ident); ok {
			if c, ok := b.Rhs.(*ast.Constant); ok {
				if b.Tok == tok.Is || c.Val == EmptyStr {
					return fixedAnd(fixed, id.Name, c.Val)
				}
			}
		}
	} else if col, vals := inToFixed(e); col != "" {
		return fixedAnd(fixed, col, vals...)
	}
	return fixed, false
}

func inToFixed(e ast.Expr) (col string, vals []Value) {
	in, ok := e.(*ast.In)
	if !ok {
		return "", nil
	}
	id, ok := in.E.(*ast.Ident)
	if !ok {
		return "", nil
	}
	for _, e2 := range in.Exprs {
		if c, ok := e2.(*ast.Constant); ok {
			vals = append(vals, c.Val)
		} else {
			return "", nil
		}
	}
	return id.Name, vals
}

// fixedAnd adds col,vals to fixed, handling if col already exists
func fixedAnd(fixed []Fixed, col string, vals ...Value) ([]Fixed, bool) {
	vs := make([]string, len(vals))
	for i, v := range vals {
		vs[i] = Pack(v.(Packable))
	}
	for i, f := range fixed {
		if f.col == col {
			v := set.Intersect(f.values, vs)
			if len(v) == 0 {
				return nil, true // conflict
			}
			if len(v) == len(f.values) {
				return fixed, false // no change
			}
			fixed := slc.Clone(fixed)
			fixed[i].values = v
			return fixed, false
		}
	}
	// col not found
	return append(fixed, Fixed{col: col, values: vs}), false
}

func (w *Where) Keys() [][]string {
	if w.keys == nil {
		w.optInit()
		if w.singleton || w.conflict {
			return [][]string{{}} // intentionally {} not nil
		}
		//TODO treat unique indexes with a where != "" as keys
		w.keys = w.source.Keys()
		assert.That(w.keys != nil)
	}
	return w.keys
}

func (w *Where) fastSingle() bool {
	if w.fast1.NotSet() {
		if w.source.fastSingle() {
			w.fast1.Set(true)
		} else {
			w.optInit()
			w.fast1.Set(w.singleton || w.conflict)
		}
	}
	return w.fast1.Get()
}

func (w *Where) Indexes() [][]string {
	w.optInit() // sets indexes
	return w.indexes
}

func (w *Where) Nrows() (int, int) {
	w.optInit() // calls calcNrows
	return w.nNrows.Get(), w.pNrows.Get()
}

func (w *Where) calcNrows() (int, int) {
	assert.That(w.optInited == optInitInProgress)
	srcNrows, srcPop := w.source.Nrows()
	if w.conflict || srcPop == 0 {
		return 0, srcPop
	}
	if w.singleton {
		return 1, srcPop
	}
	if len(w.idxSels) == 0 {
		return srcNrows / 2, srcPop
	}
	est := math.MaxInt
	nsrc := float64(srcNrows)
	for i := range w.idxSels {
		is := &w.idxSels[i]
		n := int(math.Round(is.frac * nsrc))
		if n < est {
			est = n
		}
	}
	if w.exprMore {
		est /= 2 // ??? adjust for additional restrictions
	}
	return est, srcPop
}

func (w *Where) Transform() Query {
	if w.conflict {
		return NewNothing(w)
	}
	src := w.source.Transform()
	if len(w.expr.Exprs) == 0 {
		// remove empty where
		return src
	}
	switch q := src.(type) {
	case *Nothing:
		return NewNothing(w)
	case *Tables:
		return w.tablesLookup(q)
	case *Where:
		// combine consecutive where's
		exprs := slc.With(q.expr.Exprs, w.expr.Exprs...)
		e := &ast.Nary{Tok: tok.And, Exprs: exprs}
		return NewWhere(q.source, e, w.t).Transform()
	case *Project:
		// move where before project
		q = newProject(NewWhere(q.source, w.expr, w.t), q.columns)
		return q.Transform()
	case *Rename:
		// move where before rename
		newExpr := renameExpr(w.expr, q)
		src := NewWhere(q.source, newExpr, w.t)
		return NewRename(src, q.from, q.to).Transform()
	case *Extend:
		// move where before extend, unless it depends on rules
		var before, after []ast.Expr
		for _, e := range w.expr.Exprs {
			if q.needRule(e.Columns()) {
				after = append(after, e)
			} else {
				before = append(before, replaceExpr(e, q.cols, q.exprs))
			}
		}
		if before == nil { // no split
			return w.transform(src)
		}
		src := NewWhere(q.source,
			&ast.Nary{Tok: tok.And, Exprs: before}, w.t)
		q = NewExtend(src, q.cols, q.exprs)
		if after == nil {
			return q.Transform()
		}
		e := &ast.Nary{Tok: tok.And, Exprs: after}
		return NewWhere(q, e, w.t).Transform()
	case *Summarize:
		// split where before & after summarize
		cols1 := q.source.Columns()
		var before, after []ast.Expr
		for _, e := range w.expr.Exprs {
			if set.Subset(cols1, e.Columns()) {
				before = append(before, e)
			} else {
				after = append(after, e)
			}
		}
		if before == nil { // no split
			return w.transform(src)
		}
		src := NewWhere(q.source,
			&ast.Nary{Tok: tok.And, Exprs: before}, w.t)
		q = NewSummarize(src, q.hint, q.by, q.cols, q.ops, q.ons)
		if after == nil {
			return q.Transform()
		}
		e := &ast.Nary{Tok: tok.And, Exprs: after}
		return NewWhere(q, e, w.t).Transform()
	case *Intersect:
		// distribute where over intersect
		// no project because Intersect Columns are the intersection
		src1 := NewWhere(q.source1, w.expr, w.t)
		src2 := NewWhere(q.source2, w.expr, w.t)
		return NewIntersect(src1, src2).Transform()
	case *Minus:
		// distribute where over minus
		// need project because Minus Columns are just the left side's
		src1 := NewWhere(q.source1, w.expr, w.t)
		src2 := NewWhere(q.source2, w.project(q.source2), w.t)
		return NewMinus(src1, src2).Transform()
	case *Union:
		// distribute where over union
		// need project because Union Columns is the union
		src1 := NewWhere(q.source1, w.project(q.source1), w.t)
		src2 := NewWhere(q.source2, w.project(q.source2), w.t)
		return NewUnion(src1, src2).Transform()
	case *Times:
		// split where over times
		return w.split(q, func(src1, src2 Query) Query {
			return NewTimes(src1, src2).Transform()
		})
	case *Join:
		// split where over join
		return w.split(q, func(src1, src2 Query) Query {
			return q.With(src1, src2).Transform()
		})
	case *LeftJoin:
		if w.leftJoinToJoin(q) {
			return w.split(q, func(src1, src2 Query) Query {
				return NewJoin(src1, src2, q.by, w.t).Transform()
			})
		}
		// split where over leftjoin (left side only)
		cols1 := q.source1.Columns()
		var common, exprs1 []ast.Expr
		for _, e := range w.expr.Exprs {
			if set.Subset(cols1, e.Columns()) {
				exprs1 = append(exprs1, e)
			} else {
				common = append(common, e)
			}
		}
		if exprs1 == nil { // no split
			return w.transform(src)
		}
		src1 := NewWhere(q.source1,
			&ast.Nary{Tok: tok.And, Exprs: exprs1}, w.t)
		q2 := q.With(src1, q.source2).Transform()
		if common == nil {
			return q2
		}
		e := &ast.Nary{Tok: tok.And, Exprs: common}
		return NewWhere(q2, e, w.t).Transform()
	default:
		return w.transform(src)
	}
}

func (w *Where) transform(src Query) Query {
	if src != w.source {
		return NewWhere(src, w.expr, w.t)
	}
	return w
}

func (w *Where) tablesLookup(tables *Tables) Query {
	// Optimize: tables where table = <string>
	// This is to handle the speed issue from heavy use of TableExists?.
	// It could be more general.
	col, val := w.lookup1()
	if col != "table" {
		return w
	}
	s, ok := val.ToStr()
	if !ok {
		return NewNothing(w)
	}
	return NewTablesLookup(tables.tran, s)
}

func (w *Where) lookup1() (string, Value) {
	if len(w.expr.Exprs) == 1 {
		if b, ok := w.expr.Exprs[0].(*ast.Binary); ok && b.Tok == tok.Is {
			if id, ok := b.Lhs.(*ast.Ident); ok {
				if c, ok := b.Rhs.(*ast.Constant); ok {
					return id.Name, c.Val
				}
			}
		}
	}
	return "", nil
}

func (w *Where) leftJoinToJoin(lj *LeftJoin) bool {
	cols := lj.source2.Header().GetFields()
	cols = set.Difference(cols, lj.by)
	for _, e := range w.expr.Exprs {
		if set.Subset(cols, e.Columns()) && ast.CantBeEmpty(e, cols) {
			return true
		}
	}
	return false
}

func (w *Where) project(q Query) *ast.Nary {
	srcCols := q.Columns()
	exprCols := w.expr.Columns()
	missing := set.Difference(exprCols, srcCols)
	expr := replaceExpr(w.expr, missing, nEmpty(len(missing)))
	if nary, ok := expr.(*ast.Nary); !ok || nary.Tok != tok.And {
		expr = &ast.Nary{Tok: tok.And, Exprs: []ast.Expr{expr}}
	}
	return expr.(*ast.Nary)
}

var emptyConstant = ast.Constant{Val: EmptyStr}

func nEmpty(n int) []ast.Expr {
	list := make([]ast.Expr, n)
	for i := range list {
		list[i] = &emptyConstant
	}
	return list
}

func (w *Where) split(q2 Query, newQ2 func(Query, Query) Query) Query {
	src2 := q2.(q2i).Source2()
	src1 := q2.(q2i).Source()
	cols1 := src1.Columns()
	cols2 := src2.Columns()
	var common, exprs1, exprs2 []ast.Expr
	for _, e := range w.expr.Exprs {
		used := false
		if set.Subset(cols1, e.Columns()) {
			exprs1 = append(exprs1, e)
			used = true
		}
		if set.Subset(cols2, (e.Columns())) {
			if used {
				e = replaceExpr(e, nil, nil) // copy Binary/In CouldEvalRaw
			}
			exprs2 = append(exprs2, e)
			used = true
		}
		if !used {
			common = append(common, e)
		}
	}
	if exprs1 != nil {
		src1 = NewWhere(src1, &ast.Nary{Tok: tok.And, Exprs: exprs1}, w.t)
	}
	if exprs2 != nil {
		src2 = NewWhere(src2, &ast.Nary{Tok: tok.And, Exprs: exprs2}, w.t)
	}
	if exprs1 != nil || exprs2 != nil {
		q2 = newQ2(src1, src2).Transform()
	}
	if exprs1 == nil && exprs2 == nil {
		return w.transform(q2)
	}
	if common != nil {
		e := &ast.Nary{Tok: tok.And, Exprs: common}
		return NewWhere(q2, e, w.t)
	}
	return q2
}

// optimize ---------------------------------------------------------

func (w *Where) optimize(mode Mode, index []string, frac float64) (f Cost, v Cost, a any) {
	// defer func() { fmt.Println("Where opt", index, frac, "=", f, v, a) }()
	w.optInit()
	if w.conflict {
		// fmt.Println("Where opt CONFLICT")
		return 0, 0, nil
	}
	assert.That(!w.singleton || index == nil)
	// we always have the option of just filtering (no specific index use)
	filterFixCost, filterVarCost := Optimize(w.source, mode, index, frac)
	if w.tbl == nil || w.tbl.singleton {
		// fmt.Println("Where opt", index, frac, "= filter", filterFixCost, filterVarCost)
		return filterFixCost, filterVarCost, nil
	}
	// where on table
	cost, idx := w.bestIndex(index, frac)
	if cost >= impossible {
		// only use the filter if there are no possible idxSel
		// fmt.Println("Where opt", index, frac, "= filter", filterFixCost, filterVarCost)
		return filterFixCost, filterVarCost, nil
	}
	// fmt.Println("Where opt", index, frac, "= ", idx, cost)
	return 0, cost, whereApproach{index: idx, cost: cost}
}

// exprFalse checks if any expressions folded to false
func (w *Where) exprFalse() bool {
	//TODO also check for always true
	for _, expr := range w.expr.Exprs {
		if c, ok := expr.(*ast.Constant); ok && c.Val == False {
			return true
		}
	}
	return false
}

func (w *Where) optInit() {
	if w.optInited == optInitYes {
		return
	}
	assert.That(w.optInited == optInitNo)
	w.optInited = optInitInProgress
	w.tbl, _ = w.source.(*Table)
	if !w.conflict && w.tbl != nil {
		w.idxSels = w.perIndex(w.colSels)
		// fmt.Println("idxSels", w.idxSels)
		if !w.exprMore {
			// check if any colSels were not used by idxSels
			for _, idxSel := range w.idxSels {
				for i := 0; i < idxSel.nfields; i++ {
					delete(w.colSels, idxSel.index[i])
				}
			}
			w.exprMore = len(w.colSels) > 0
		}
	}
	w.setNrows(w.calcNrows())
	if !w.singleton {
		w.indexes = w.source.Indexes()
	}
	w.optInited = optInitYes
}

// bestIndex returns the best (lowest cost) index with an idxSel
// that satisfies the required order (or impossible)
func (w *Where) bestIndex(order []string, frac float64) (Cost, []string) {
	best := newBestIndex()
	for _, idx := range w.source.Indexes() {
		if ordered(idx, order, w.fixed) {
			if is := w.getIdxSel(idx); is != nil {
				varcost := w.source.lookupCost() * len(is.ptrngs)
				tblFixCost, tblVarCost, _ :=
					w.tbl.optimize(CursorMode, idx, frac*is.fracRange)
				assert.That(tblFixCost == 0)
				varcost += tblVarCost
				best.update(idx, 0, varcost)
			}
		}
	}
	return best.varcost, best.index
}

func (w *Where) getIdxSel(index []string) *idxSel {
	for i := range w.idxSels {
		is := &w.idxSels[i]
		if slices.Equal(index, is.index) {
			return is
		}
	}
	return nil
}

func (w *Where) setApproach(index []string, frac float64, app any, tran QueryTran) {
	w.optimized = true
	if w.conflict {
		return
	}
	if app != nil {
		wapp := app.(whereApproach)
		idx := wapp.index
		w.tbl.setIndex(idx)
		w.idxSel = w.getIdxSel(idx)
		w.tbl.cacheSetCost(frac*w.idxSel.frac, 0, wapp.cost)
		w.idxSelPos = -1
	} else { // filter
		w.source = SetApproach(w.source, index, frac, tran)
	}
	w.header = w.source.Header()
	w.ctx.Hdr = w.header
}

// execution --------------------------------------------------------

// MakeSuTran is injected by dbms to avoid import cycle
var MakeSuTran func(qt QueryTran) *SuTran

func (w *Where) Get(th *Thread, dir Dir) Row {
	defer func(t uint64) { w.tget += tsc.Read() - t }(tsc.Read())
	w.ngets++
	if w.selSet && w.selOrg == ixkey.Max && w.selEnd == "" {
		return nil // conflict from Select
	}
	for {
		row := w.get(th, dir)
		if w.filter(th, row) {
			w.nOut++
			if row == nil {
				w.slowQueries()
			}
			return row
		}
	}
}

func (w *Where) get(th *Thread, dir Dir) Row {
	if w.idxSel == nil {
		w.nIn++
		return w.source.Get(th, dir)
	}
	for {
		if w.idxSelPos != -1 && w.curPtrng.isRange() {
			if row := w.tbl.Get(th, dir); row != nil {
				return row
			}
		}
		if !w.advance(dir) {
			return nil // eof
		}
		if w.curPtrng.isRange() {
			w.tbl.SelectRaw(w.curPtrng.org, w.curPtrng.end)
		} else { // point
			if row := w.tbl.lookup(w.curPtrng.org); row != nil {
				w.nIn++
				return row
			}
		}
	}
}

func (w *Where) filter(th *Thread, row Row) bool {
	if row == nil {
		return true
	}
	if w.selectCols != nil &&
		!singletonFilter(w.header, row, w.selectCols, w.selectVals) {
		return false
	}
	if w.ctx.Tran == nil {
		w.ctx.Tran = MakeSuTran(w.t)
	}
	w.ctx.Th = th
	w.ctx.Row = row
	defer func() { w.ctx.Th, w.ctx.Row = nil, nil }()
	return w.expr.Eval(&w.ctx) == True
}

func (w *Where) advance(dir Dir) bool {
	if w.idxSelPos == -1 { // rewound
		if dir == Prev {
			w.idxSelPos = len(w.idxSel.ptrngs)
		}
	}
	for {
		if dir == Prev {
			w.idxSelPos--
		} else { // Next
			w.idxSelPos++
		}
		if w.idxSelPos < 0 || len(w.idxSel.ptrngs) <= w.idxSelPos {
			return false // eof
		}
		pr := w.idxSel.ptrngs[w.idxSelPos]
		if w.selSet {
			pr = pr.intersect(w.selOrg, w.selEnd)
			if pr.conflict() {
				continue
			}
		}
		w.curPtrng = pr
		return true
	}
}

func (w *Where) Rewind() {
	w.source.Rewind()
	w.idxSelPos = -1
}

func (w *Where) Select(cols, vals []string) {
	// fmt.Println("Where", w.tbl.name, "Select", cols, unpack(vals))
	w.nsels++
	w.Rewind()
	w.selOrg, w.selEnd = "", ""
	w.selSet = false
	w.selectCols = nil
	w.selectVals = nil
	if cols == nil && vals == nil { // clear select
		if w.idxSel == nil {
			w.source.Select(nil, nil)
		}
		return
	}
	// Note: conflict could come from any of expr, not just fixed.
	// But to evaluate that would require building a Row.
	// It should be rare.
	satisfied, conflict := selectFixed(cols, vals, w.Fixed())
	if conflict {
		w.selOrg, w.selEnd = ixkey.Max, ""
		w.selSet = true
		return
	}
	if satisfied {
		w.selOrg, w.selEnd = "", ""
		w.selSet = false
		return
	}

	if w.singleton {
		w.selectCols = cols
		w.selectVals = vals
		return
	}

	if w.idxSel == nil {
		// if this where is not using an index selection
		// then just pass the Select to the source
		w.source.Select(cols, vals)
		return
	}

	cols, vals = w.addFixed(cols, vals)
	w.selOrg, w.selEnd = selKeys(w.idxSel.encoded, w.idxSel.index, cols, vals)
	w.selSet = true
}

func (w *Where) addFixed(cols []string, vals []string) ([]string, []string) {
	cols = slices.Clip(cols)
	vals = slices.Clip(vals)
	for _, fix := range w.fixed {
		if fix.single() && !slices.Contains(cols, fix.col) {
			cols = append(cols, fix.col)
			vals = append(vals, fix.values[0])
		}
	}
	return cols, vals
}

func (w *Where) Lookup(th *Thread, cols, vals []string) Row {
	w.nlooks++
	if conflictFixed(cols, vals, w.Fixed()) {
		return nil
	}
	if w.singleton {
		// can't use source.Lookup because cols may not match source index
		w.Rewind()
		row := w.Get(th, Next)
		if row == nil || !singletonFilter(w.header, row, cols, vals) {
			return nil
		}
		return row
	}
	cols, vals = w.addFixed(cols, vals)
	row := w.source.Lookup(th, cols, vals)
	if !w.filter(th, row) {
		row = nil
	}
	return row
}

func singletonFilter(
	hdr *Header, row Row, cols []string, vals []string) bool {
	for i, col := range cols {
		if row.GetRaw(hdr, col) != vals[i] {
			return false
		}
	}
	return true
}

func (w *Where) slowQueries() {
	if w.slow() && trace.SlowQuery.On() {
		trace.SlowQuery.Println(w.nIn, "->", w.nOut)
		trace.Println(strategy(w, 1))
		w.nIn = 0
		w.nOut = 0
	}
}
func (w *Where) slow() bool {
	return w.nIn > 100 && w.nIn > w.nOut*100
}

func (w *Where) Simple(th *Thread) []Row {
	w.ctx.Hdr = w.header
	rows := w.source.Simple(th)
	dst := 0
	for _, row := range rows {
		if w.filter(nil, row) {
			rows[dst] = row
			dst++
		}
	}
	return rows[:dst]
}
