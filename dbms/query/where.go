// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"math"
	"sync/atomic"

	"slices"

	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/trace"
	"github.com/apmckinlay/gsuneido/db19/index/iface"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/set"
	"github.com/apmckinlay/gsuneido/util/slc"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/tsc"
)

var whereSingletonCount atomic.Int64
var _ = AddInfo("query.where.singleton", &whereSingletonCount)

// NOTE: Where source and expr should NOT be modified,
// instead, construct a new one with NewWhere

type Where struct {
	Query1
	t         QueryTran
	colSels   map[string][]span // from NewWhere, result of perField
	mergedBuf map[string][]span // reusable buffer for mergedPerCol
	// tbl will be set if the source is a table, nil otherwise
	tbl *Table
	// idxSel is for the chosen index
	idxSelBase   *idxSel
	idxSelActive *idxSel
	selConflict  bool      // conflict from Select
	expr         *ast.Nary // And

	// curPtrng is idxSelActive.prefixRanges[idxSelPos]
	curPtrng pointRange

	rowCtx ast.RowContext

	// singleSels are set by Select when singleton
	singleSels Sels

	idxSels []*idxSel // from optInit, result of perIndex
	// idxSelPos is the current index in idxSel.ptrngs
	idxSelPos int
	// eof is set when advance returns false (past end or past beginning)
	// to enforce "plain stick" cursor behavior: once past eof, stay past eof
	eof bool

	nIn  int
	nOut int

	// conflict is true if the expression conflicts and selects nothing
	conflict bool
	// singleton is true if we know the result is at most one record
	// because there is a single point select on a key
	// e.g. key(a,b) where a = 1 and b = 2
	// NOTE: this does NOT mean source.fastSingle
	singleton bool

	optInited      // used by optinit
	optimized bool // set by setApproach, used by String

	ixCtx  ixContext
	ixExpr ast.Expr

	srcIndex []string // set by setApproach, used by Lookup
	wFixed   Fixed

	wfrac float64
}

type optInited byte

const (
	optInitNo optInited = iota
	optInitInProgress
	optInitYes
)

type whereApproach struct {
	index []string
	*idxSel
	cost Cost
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
		w.expr.CanEvalRaw(fields)
		w.colSels = perField(w.expr.Exprs, fields)
		// fmt.Println("colSels", w.colSels)
		w.conflict = (w.colSels == nil)
	}
	return w
}

func (w *Where) SetTran(t QueryTran) {
	w.t = t
	w.source.SetTran(t)
	w.rowCtx.Tran = nil
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
	w.wFixed = efixed
	fixed, none := w.source.Fixed().Combine(efixed)
	if none {
		w.conflict = true
		return
	}
	w.fixed = fixed
}

func (w *Where) exprsToFixed() (fixed Fixed, conflict bool) {
	for _, e := range w.expr.Exprs {
		fixed, conflict = addFixed(fixed, e)
		if conflict {
			return nil, true
		}
	}
	return fixed, false
}

func addFixed(fixed Fixed, e ast.Expr) (Fixed, bool) {
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
func fixedAnd(fixed Fixed, col string, vals ...Value) (Fixed, bool) {
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
			fixed2 := slc.Clone(fixed)
			fixed2[i].values = v
			return fixed2, false
		}
	}
	// col not found
	return append(fixed, Fix{col: col, values: vs}), false
}

func (w *Where) Keys() [][]string {
	if w.keys == nil {
		w.optInit()
		if w.singleton || w.conflict {
			w.keys = emptyKey // intentionally {} not nil
		} else {
			//TODO treat unique indexes with a where != "" as keys
			w.keys = w.source.Keys()
			// logically, we should remove fixed from the keys
			// but it causes problems for operations choosing a key index
		}
		assert.That(w.keys != nil) // once only
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
	// Note: the estimated row count should be consistent
	// with bestIndex and WhereCost
	// since cost is primarily driven by number of rows
	srcNrows, srcPop := w.source.Nrows()
	if w.conflict || srcPop == 0 {
		return 0, srcPop
	}
	if w.singleton {
		return 1, srcPop
	}
	if len(w.idxSels) == 0 {
		return int(math.Round(float64(srcNrows) * unknownFrac)), srcPop
	}
	est := int(math.Round(w.wfrac * float64(srcNrows)))
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
		return newProject(NewWhere(q.source, w.expr, w.t), q.columns).Transform()
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
				before = append(before, replaceExpr(e, q.cols, q.exprs, false))
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
		return NewIntersect(src1, src2, w.t).Transform()
	case *Minus:
		// distribute where over minus
		// need project because Minus Columns are just the left side's
		src1 := NewWhere(q.source1, w.expr, w.t)
		src2 := NewWhere(q.source2, w.project(q.source2), w.t)
		return NewMinus(src1, src2, w.t).Transform()
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
	case *SemiJoin:
		// move where before semijoin
		src1 := NewWhere(q.source1, w.expr, w.t)
		return q.With(src1, q.source2).Transform()
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
	flds := lj.source2.Header().GetFields()
	flds = set.Difference(flds, lj.by)
	for _, e := range w.expr.Exprs {
		if set.Subset(flds, e.Columns()) && !ast.CanBeEmpty(e) {
			return true
		}
	}
	return false
}

func (w *Where) project(q Query) *ast.Nary {
	srcCols := q.Columns()
	exprCols := w.expr.Columns()
	missing := set.Difference(exprCols, srcCols)
	expr := replaceExpr(w.expr, missing, nEmpty(len(missing)), false)
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
				e = replaceExpr(e, nil, nil, true) // clone
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

func (w *Where) optimize(mode Mode, req Require) (Cost, Cost, any) {
	w.optInit()
	if w.conflict {
		return 0, 0, nil
	}
	if w.tbl == nil || w.tbl.isSingleton() {
		fixcost, varcost := Optimize(w.source, mode, req)
		return fixcost, varcost, nil
	}
	if req.use == ReqUnique {
		return w.optWhereLookup(req)
	}
	return w.optWhereIdx(req)
}

func (w *Where) optWhereIdx(req Require) (Cost, Cost, any) {
	if w.singleton {
		// here singleton == fastSingle
		// because source is a Table and Table keys are a subset of indexes
		cost := w.source.lookupCost()
		isel := w.idxSels[0]
		return 0, cost, &whereApproach{index: isel.index, cost: cost, idxSel: isel}
	}
	type bestIdx struct {
		index  []string
		idxSel *idxSel
	}
	best := newBest[bestIdx]()
	for _, idx := range w.source.Indexes() {
		if !req.SatisfiedByWithFixed(idx, w.fixed) {
			continue
		}
		_, varCost, _ := w.tbl.optimize(0, OrderReq(idx, 1.0))
		irFrac := 1.0
		ifFrac := 1.0
		dfFrac := w.wfrac
		isel := w.getIdxSel(idx)
		if isel != nil {
			irFrac = isel.prefixFrac * isel.skipFrac
			ifFrac = isel.indexFilterFrac
			dfFrac = isel.dataFilterFrac
		}
		cost := WhereCost(float64(varCost), float64(req.frac), irFrac, ifFrac, dfFrac)
		cost += Cost(req.nseeks) * w.source.lookupCost()
		best.update(0, cost, bestIdx{index: idx, idxSel: isel})
	}
	if best.none() {
		return impossible, impossible, nil
	}
	return 0, best.varcost, &whereApproach{index: best.data.index,
		cost: best.varcost, idxSel: best.data.idxSel}
}

func (w *Where) optWhereLookup(req Require) (Cost, Cost, any) {
	if w.singleton {
		isel := w.idxSels[0]
		cost := w.tbl.lookupCostFor(w.tbl.indexi(isel.index))
		return 0, cost, &whereApproach{index: isel.index, cost: cost, idxSel: isel}
	}
	best := newBest[[]string]()
	for idxi, idx := range w.tbl.indexes {
		if indexCovered(idx, req.cols, w.fixed) {
			varcost := Cost(req.nseeks) * w.tbl.lookupCostFor(idxi)
			best.update(0, varcost, idx)
		}
	}
	if best.none() {
		return impossible, impossible, nil
	}
	return 0, best.varcost, &whereApproach{index: best.data, cost: best.varcost}
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
	}
	// detect singleton when fixed covers a key (for non-Table sources).
	// Required so bestLookupIndex doesn't pick an index with extra columns
	// that sels can't cover at Lookup time
	// (lookupIndexEligible allows any index when nColsUnfixed == 0).
	if !w.singleton && !w.conflict && w.tbl == nil {
		if slices.ContainsFunc(w.source.Keys(), w.fixed.All) {
			w.singleton = true
		}
	}
	w.setNrows(w.calcNrows())
	if w.singleton {
		w.indexes = emptyKey
	} else {
		w.indexes = w.source.Indexes()
	}
	w.optInited = optInitYes
}

func WhereCost(cost, inFrac, irFrac, ifFrac, dfFrac float64) Cost {
	// Model the cost of reading via a particular index.
	// We have 3 potential filter stages:
	// 1. Index Range (irFrac)
	// 2. Index Filter (ifFrac)
	// 3. Data Filter (dataFilter)

	// this is necessary because some places e.g. LookupCost pass 0 frac
	if inFrac == 0 {
		return 0
	}

	// index range
	cost *= irFrac

	// index filter
	indexCost := .2 * cost
	dataCost := .8 * cost * ifFrac
	cost = indexCost + dataCost

	// data filter (dfFrac) does NOT affect cost (other than pessimism)
	// because you still have to read everything

	// apply inFrac
	if ifFrac < 1 || dfFrac < 1 {
		inFrac = .25 + (.75 * inFrac) // pessimistic guard
	}
	cost *= inFrac

	return int(math.Round(cost))
}

func (w *Where) getIdxSel(index []string) *idxSel {
	for _, isel := range w.idxSels {
		if slices.Equal(index, isel.index) {
			return isel
		}
	}
	return nil
}

func (w *Where) setApproach(req Require, approach any, tran QueryTran) {
	w.optimized = true
	if w.singleton {
		whereSingletonCount.Add(1)
	}
	if w.conflict {
		return
	}
	if approach == nil {
		w.source = SetApproach(w.source, req, tran)
		w.srcIndex = req.cols
		w.tbl = nil
	} else {
		app := approach.(*whereApproach)
		w.tbl.SetIndex(app.index)
		w.srcIndex = app.index
		if app.idxSel != nil {
			w.ixCtx.cols = app.index
			w.ixCtx.encodes = w.tbl.IndexEncodes(app.index)
			w.ixExpr = w.exprsFor(w.ixCtx.cols)
			w.idxSelBase = app.idxSel
			w.idxSelActive = w.idxSelBase
			w.tbl.setCost(float64(req.frac)*app.idxSel.prefixFrac*app.idxSel.skipFrac,
				0, app.cost)
			w.idxSelPos = -1
		} else {
			w.tbl.setCost(float64(req.frac), 0, app.cost)
		}
	}
	w.header = w.source.Header()
	w.rowCtx.Hdr = w.header
}

func (w *Where) exprsFor(cols []string) ast.Expr {
	var exprs []ast.Expr
	for _, e := range w.expr.Exprs {
		if set.Subset(cols, e.Columns()) {
			exprs = append(exprs, e)
		}
	}
	if len(exprs) == 0 {
		return nil
	} else if len(exprs) == 1 {
		return exprs[0]
	}
	return &ast.Nary{Tok: tok.And, Exprs: exprs}
}

// execution --------------------------------------------------------

// MakeSuTran is injected by dbms to avoid import cycle
var MakeSuTran func(qt QueryTran) *SuTran

func (w *Where) Get(th *Thread, dir Dir) Row {
	defer func(t uint64) { w.tget += tsc.Read() - t }(tsc.Read())
	if w.selConflict || w.eof {
		return nil
	}
	// apply the non-indexed filtering
	for {
		row := w.get(th, dir)
		if w.filter(th, row) {
			w.nOut++
			if row == nil {
				w.eof = true
				w.slowQueries()
				return nil
			}
			w.ngets++
			return row
		}
	}
}

func (w *Where) get(th *Thread, dir Dir) Row {
	if w.idxSelActive == nil {
		w.nIn++
		if w.tbl == nil {
			return w.source.Get(th, dir)
		}
		return w.getFilter(th, dir)
	}
	// loop over the prefix index ranges/points
	for {
		if w.idxSelPos != -1 && w.idxSelPos < len(w.idxSelActive.prefixRanges) &&
			w.curPtrng.isRange() {
			if row := w.getFilter(th, dir); row != nil {
				w.nIn++
				return row
			}
		}
		if !w.advance(dir) {
			return nil // eof
		}
		if w.curPtrng.isRange() {
			if w.idxSelActive.skipLen > 0 {
				w.tbl.SelectSkipScan(iface.Range(w.curPtrng),
					iface.Range(w.idxSelActive.skipRange),
					w.idxSelActive.skipStart)
			} else {
				w.tbl.SelectRaw(w.curPtrng.Org, w.curPtrng.End)
			}
		} else { // point
			if row := w.tbl.LookupRaw(w.curPtrng.Org); row != nil {
				w.nIn++
				return row
			}
		}
	}
}

func (w *Where) getFilter(th *Thread, dir Dir) Row {
	var filterFunc func(string) bool
	if w.ixExpr != nil {
		w.ixCtx.th = th
		filterFunc = func(key string) bool {
			w.ixCtx.key = key
			return w.ixExpr.Eval(&w.ixCtx) == True
		}
	}
	return w.tbl.GetFilter(dir, filterFunc)
}

// filter applies the entire where expression
// and also selectSelCols/Vals singletonFilter
func (w *Where) filter(th *Thread, row Row) bool {
	if row == nil {
		return true
	}
	if w.singleSels != nil &&
		!singletonFilter(w.header, row, w.singleSels) {
		return false
	}
	if w.rowCtx.Tran == nil {
		w.rowCtx.Tran = MakeSuTran(w.t)
	}
	w.rowCtx.Th = th
	w.rowCtx.Row = row
	defer func() { w.rowCtx.Th, w.rowCtx.Row = nil, nil }()
	return w.expr.Eval(&w.rowCtx) == True
}

// advance moves to the next prefix index range
func (w *Where) advance(dir Dir) bool {
	if w.idxSelPos == -1 { // rewound
		if dir == Prev {
			w.idxSelPos = len(w.idxSelActive.prefixRanges)
		}
	}
	if dir == Prev {
		w.idxSelPos--
	} else { // Next
		w.idxSelPos++
	}
	if w.idxSelPos < 0 || len(w.idxSelActive.prefixRanges) <= w.idxSelPos {
		return false // eof
	}
	w.curPtrng = w.idxSelActive.prefixRanges[w.idxSelPos]
	return true
}

func (w *Where) Rewind() {
	w.source.Rewind()
	w.idxSelPos = -1
	w.eof = false
}

func (w *Where) Select(sels Sels) {
	// fmt.Println("Where Select", cols, unpack(vals))
	w.nsels++
	w.Rewind()
	w.idxSelActive = w.idxSelBase
	w.selConflict = false
	w.singleSels = nil
	if sels == nil { // clear select
		if w.idxSelBase == nil {
			w.source.Select(nil)
		}
		return
	}
	// Note: conflict could come from any of expr, not just fixed.
	// But to evaluate that would require building a Row.
	// It should be rare.
	satisfied, conflict := w.Fixed().Match(sels)
	if conflict {
		w.selConflict = true
		return
	}
	if satisfied {
		return
	}

	if w.singleton {
		w.singleSels = sels
		return
	}

	if w.idxSelBase == nil {
		// add applicable fixed to sels for the source
		// (same as Lookup) since source index may include fixed columns
		sels = slices.Clip(sels)
		for _, fix := range w.fixed {
			if fix.Single() && slices.Contains(w.srcIndex, fix.col) &&
				!sels.HasCol(fix.col) {
				sels = append(sels, Sel{fix.col, fix.values[0]})
			}
		}
		w.source.Select(sels)
		return
	}

	var isel *idxSel
	isel, w.selConflict = w.recalcIdxSel(w.idxSelBase.index, w.idxSelBase.mode, sels)
	if w.selConflict {
		return
	}
	w.idxSelActive = isel
}

func (w *Where) Lookup(th *Thread, sels Sels) Row {
	// sels (plus fixed) specify a single source row
	w.nlooks++
	if w.Fixed().Conflicts(sels) {
		return nil
	}
	if w.fastSingle() || w.srcIndex == nil {
		// can't use source.Lookup because cols may not cover the source index.
		// fastSingle: source is a literal singleton (empty key).
		// srcIndex == nil: fixed covers a key (singleton detected in optInit),
		// so Optimize passed index=nil and setApproach left srcIndex nil.
		w.Rewind()
		return GetNext1(w, th, sels)
	}
	cloned := false
	sels = slices.Clip(sels)
	for _, fix := range w.fixed {
		if fix.Single() && slices.Contains(w.srcIndex, fix.col) &&
			!sels.HasCol(fix.col) {
			sels = append(sels, Sel{fix.col, fix.values[0]})
			cloned = true // because they're clipped, append will realloc
		}
	}
	isels, osels := Split(cloned, sels, w.srcIndex)
	var residual Sels
	for _, sel := range osels {
		if !w.fixed.Has(sel.col) {
			residual = append(residual, sel)
		}
	}

	row := w.source.Lookup(th, slc.With(isels, residual...))
	if !w.filter(th, row) {
		row = nil
	}
	return row
}

// Split partitions flds and vals, returning sub-slices.
// It clones the slices only if modifications are needed.
func Split(cloned bool, sels Sels, index []string) (isels, osels Sels) {
	pivot := func(i int) bool {
		return slices.Contains(index, sels[i].col)
	}
	swap := func(i, j int) {
		if !cloned {
			sels = slc.Clone(sels)
			cloned = true
		}
		sels[i], sels[j] = sels[j], sels[i]
	}
	i := slc.Partition(len(sels), pivot, swap)
	isels = sels[:i]
	osels = sels[i:]
	if len(isels) == 0 {
		isels = nil
	}
	if len(osels) == 0 {
		osels = nil
	}
	return
}

func singletonFilter(hdr *Header, row Row, sels Sels) bool {
	for _, sel := range sels {
		x := row.GetRaw(hdr, sel.col)
		assert.That(len(x) == 0 || x[0] != PackForward)
		if x != sel.val {
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

func (w *Where) InCount() int {
	return w.nIn
}

func (w *Where) Simple(th *Thread) []Row {
	ast.Unraw(w.expr)
	w.rowCtx.Hdr = w.header
	rows := w.source.Simple(th)
	dst := 0
	for _, row := range rows {
		if w.filter(th, row) {
			rows[dst] = row
			dst++
		}
	}
	return rows[:dst]
}

//-------------------------------------------------------------------

type ixContext struct {
	th      *Thread
	key     string
	cols    []string
	encodes bool
}

func (c *ixContext) GetVal(id string) Value {
	if ascii.IsUpper(id[0]) {
		return Global.GetName(c.Thread(), id)
	}
	return Unpack(c.GetRaw(id))
}

func (c *ixContext) GetRaw(id string) string {
	if !c.encodes {
		return c.key
	}
	return ixkey.Decode1(c.key, slices.Index(c.cols, id))
}

func (c *ixContext) Thread() *Thread {
	return c.th
}
