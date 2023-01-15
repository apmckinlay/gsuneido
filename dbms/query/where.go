// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/trace"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/ord"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/str"
	"golang.org/x/exp/slices"
)

// NOTE: Where source and expr should NOT be modified,
// instead, construct a new one with NewWhere

type Where struct {
	Query1
	expr       *ast.Nary // And
	whereFixed []Fixed
	fixed      []Fixed
	t          QueryTran
	// tbl will be set if the source is a Table, nil otherwise
	tbl       *Table
	optInited bool
	ctx       ast.Context
	// exprMore is whether expr has more than the minimum idxSel
	exprMore bool
	colSels  map[string]filter
	// singleton is true if we know the result is at most one record
	// because there is a single point select on a key
	singleton bool
	// nrows is the estimated number of result rows
	nrows int
	// srcpop is the source Nrows pop
	srcpop  int
	idxSels []idxSel
	// conflict is true if the expression conflicts and selects nothing
	conflict bool
	// idxSel is for the chosen index
	idxSel *idxSel
	// idxSelPos is the current index in idxSel.ptrngs
	idxSelPos int
	// curPtrng is idxSel.ptrngs[idxSelPos] adjusted by Select (selOrg, selEnd)
	curPtrng pointRange
	hdr      *runtime.Header
	// selOrg and selEnd are set by Select
	selOrg string
	selEnd string
	selSet bool

	selectCols []string
	selectVals []string

	nIn  int
	nOut int
}

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
	w.calcFixed()
	cmps := w.extractCompares()
	w.exprMore = len(cmps) < len(w.expr.Exprs)
	w.colSels = w.comparesToFilters(cmps)
	w.conflict = w.conflict || w.exprFalse()
	return w
}

func (w *Where) SetTran(t QueryTran) {
	w.t = t
	w.source.SetTran(t)
	w.ctx.Tran = nil
}

func (w *Where) String() string {
	return parenQ2(w.source) + " " + w.stringOp()
}

func (w *Where) stringOp() string {
	s := "WHERE"
	if w.conflict {
		return s + " nothing"
	}
	if w.singleton {
		s += "*1"
	}
	if len(w.expr.Exprs) > 0 {
		s += " " + w.expr.Echo()
	}
	return s
}

func (w *Where) Fixed() []Fixed {
	return w.fixed
}

// calcFixed sets w.whereFixed and w.fixed and may set w.conflict
func (w *Where) calcFixed() {
	w.whereFixed = w.exprsToFixed()
	fixed, none := combineFixed(w.source.Fixed(), w.whereFixed)
	if none {
		w.conflict = true
	}
	w.fixed = fixed
}

func (w *Where) exprsToFixed() []Fixed {
	var fixed []Fixed
	for _, e := range w.expr.Exprs {
		fixed = addFixed(fixed, e)
	}
	return fixed
}

func addFixed(fixed []Fixed, e ast.Expr) []Fixed {
	// MAYBE: handle IN and OR
	if b, ok := e.(*ast.Binary); ok && b.Tok == tok.Is {
		if id, ok := b.Lhs.(*ast.Ident); ok {
			if c, ok := b.Rhs.(*ast.Constant); ok {
				fixed = append(fixed, NewFixed(id.Name, c.Val))
			}
		}
	}
	return fixed
}

func (w *Where) Keys() [][]string {
	if !w.optInited {
		w.optInit()
	}
	if w.singleton || w.conflict {
		return [][]string{{}} // intentionally {} not nil
	}
	return w.source.Keys()
}

func (w *Where) fastSingle() bool {
	if w.source.fastSingle() {
		return true
	}
	if !w.optInited {
		w.optInit()
	}
	return w.singleton || w.conflict
}

func (w *Where) Indexes() [][]string {
	if !w.optInited {
		w.optInit()
	}
	if w.singleton {
		return [][]string{{}} // intentionally {} not nil
	}
	return w.source.Indexes()
}

func (w *Where) Nrows() (int, int) {
	if !w.optInited {
		w.optInit()
	}
	return w.nrows, w.srcpop
}

func (w *Where) calcNrows() (int, int) {
	assert.That(w.optInited)
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
		ix := &w.idxSels[i]
		var n int
		if ix.isRanges() {
			n = int(math.Round(ix.frac * nsrc))
		} else { // points
			n = len(ix.ptrngs)
		}
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
	if len(w.expr.Exprs) == 0 {
		// remove empty where
		return w.source.Transform()
	}
	if w.conflict {
		return NewNothing(w.Columns())
	}
	if tl := w.tablesLookup(); tl != nil {
		return tl
	}
	if lj := w.leftJoinToJoin(); lj != nil {
		// convert leftjoin to join
		src := NewJoin(lj.source1, lj.source2, lj.by)
		return NewWhere(src, w.expr, w.t).Transform()
	}
	switch q := w.source.(type) {
	case *Where:
		// combine consecutive where's
		exprs := w.expr.Exprs
		for {
			exprs = slc.With(q.expr.Exprs, exprs...)
			if src, ok := q.source.(*Where); ok {
				q = src
			} else {
				break
			}
		}
		e := &ast.Nary{Tok: tok.And, Exprs: exprs}
		return NewWhere(q.source, e, w.t).Transform()
	case *Project:
		// move where before project
		q = NewProject(NewWhere(q.source, w.expr, w.t), q.columns)
		return q.Transform()
	case *Rename:
		// move where before rename
		newExpr := renameExpr(w.expr, q.to, q.from)
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
			return w.transform()
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
			return w.transform()
		}
		src := NewWhere(q.source,
			&ast.Nary{Tok: tok.And, Exprs: before}, w.t)
		q = NewSummarize(src, q.by, q.cols, q.ops, q.ons)
		if after == nil {
			return q
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
		// split where over product
		return w.split(q, func(src1, src2 Query) Query {
			return NewTimes(src1, src2)
		})
	case *Join:
		// split where over join
		return w.split(q, func(src1, src2 Query) Query {
			return NewJoin(src1, src2, q.by)
		})
	case *LeftJoin:
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
			return w.transform()
		}
		src1 := NewWhere(q.source1,
			&ast.Nary{Tok: tok.And, Exprs: exprs1}, w.t)
		q2 := NewLeftJoin(src1, q.source2, q.by).Transform()
		if common == nil {
			return q2
		}
		e := &ast.Nary{Tok: tok.And, Exprs: common}
		return NewWhere(q2, e, w.t)
	default:
		return w.transform()
	}
}

func (w *Where) transform() Query {
	src := w.source.Transform()
	if _,ok := src.(*Nothing); ok {
		return NewNothing(w.Columns())
	}
	if src != w.source {
		return NewWhere(src, w.expr, w.t)
	}
	return w
}

func (w *Where) tablesLookup() Query {
	// Optimize: tables where table|tablename = <string>
	// This is to handle the speed issue from heavy use of TableExists?.
	// It could be more general.
	tables, ok := w.source.(*Tables)
	if !ok {
		return nil
	}
	col, val := w.lookup1()
	if col != "table" && col != "tablename" {
		return nil
	}
	s, ok := val.ToStr()
	if !ok {
		return NewNothing(w.Columns())
	}
	return NewTablesLookup(tables.tran, s)
}

func (w *Where) lookup1() (string, runtime.Value) {
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

func (w *Where) leftJoinToJoin() *LeftJoin {
	if lj, ok := w.source.(*LeftJoin); ok {
		cols := lj.source2.Header().GetFields()
		cols = set.Difference(cols, lj.by)
		for _, e := range w.expr.Exprs {
			if set.Subset(cols, e.Columns()) && ast.CantBeEmpty(e, cols) {
				return lj
			}
		}
	}
	return nil
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

var emptyConstant = ast.Constant{Val: runtime.EmptyStr}

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
	if common != nil {
		e := &ast.Nary{Tok: tok.And, Exprs: common}
		return NewWhere(q2, e, w.t)
	}
	return q2
}

// optimize ---------------------------------------------------------

func (w *Where) optimize(mode Mode, index []string, frac float64) (f Cost, v Cost, a any) {
	// defer fmt.Println("Where opt", index, frac, "=", f, v, a)
	if !w.optInited {
		w.optInit()
	}
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
		if c, ok := expr.(*ast.Constant); ok && c.Val == runtime.False {
			return true
		}
	}
	return false
}

func (w *Where) optInit() {
	w.optInited = true
	w.tbl, _ = w.source.(*Table)
	if !w.conflict && w.tbl != nil {
		w.idxSels = w.colSelsToIdxSels(w.colSels)
		w.exprMore = w.exprMore || len(w.idxSels) > 1
		if !w.exprMore {
			// check if any colSels were not used by idxSels
			for _, idxSel := range w.idxSels {
				for _, col := range idxSel.index {
					delete(w.colSels, col)
				}
			}
			w.exprMore = w.exprMore || len(w.colSels) > 0
		}
	}
	w.nrows, w.srcpop = w.calcNrows()
}

// extractCompares finds sub-expressions like <field> <op> <constant>
func (w *Where) extractCompares() []cmpExpr {
	cols := w.source.Header().Physical()
	cmps := make([]cmpExpr, 0, 4)
	for _, expr := range w.expr.Exprs {
		if expr.CanEvalRaw(cols) {
			if bin, ok := expr.(*ast.Binary); ok && bin.Tok != tok.Isnt {
				cmp := cmpExpr{
					col: bin.Lhs.(*ast.Ident).Name,
					op:  bin.Tok,
					val: bin.Rhs.(*ast.Constant).Packed,
				}
				cmps = append(cmps, cmp)
			} else if in, ok := expr.(*ast.In); ok {
				cmp := cmpExpr{
					col:  in.E.(*ast.Ident).Name,
					op:   tok.In,
					vals: in.Packed,
				}
				cmps = append(cmps, cmp)
			}
		}
	}
	return cmps
}

// comparesToFilters combines multiple compares for each field.
// It also sets w.conflict if any filter cannot match anything.
func (w *Where) comparesToFilters(cmps []cmpExpr) map[string]filter {
	// sort to group columns
	sort.Slice(cmps, func(i, j int) bool { return cmps[i].col < cmps[j].col })
	filters := make(map[string]filter)
	f := filterAll
	for i := range cmps {
		f.andWith(cmps[i].toFilter())
		if i+1 >= len(cmps) || cmps[i].col != cmps[i+1].col {
			// end of a group of one column
			if f.none() {
				w.conflict = true
				return nil
			}
			if !f.all() {
				filters[cmps[i].col] = f
			}
			f = filterAll
		}
	}
	return filters
}

// bestIndex returns the best (lowest cost) index with an idxSel
// that satisfies the required order (or impossible)
func (w *Where) bestIndex(order []string, frac float64) (Cost, []string) {
	best := newBestIndex()
	for _, idx := range w.source.Indexes() {
		if ordered(idx, order, w.fixed) {
			if is := w.getIdxSel(idx); is != nil {
				varcost := w.source.lookupCost() * len(is.ptrngs)
				if is.isRanges() {
					tblFixCost, tblVarCost, _ :=
						w.tbl.optimize(CursorMode, idx, frac*is.frac)
					assert.That(tblFixCost == 0)
					varcost += tblVarCost
				}
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
	w.hdr = w.source.Header()
	w.ctx.Hdr = w.hdr
}

// cmpExpr is <field> <op> <constant> or <field> in (<constants>)
// which can be evaluated packed
type cmpExpr struct {
	col  string
	op   tok.Token
	val  string   // packed (binary)
	vals []string // packed (in)
}

func (cmp cmpExpr) String() string {
	s := cmp.col + " " + cmp.op.String() + " "
	if cmp.vals == nil {
		s += packToStr(cmp.val)
	} else {
		sep := "("
		for _, val := range cmp.vals {
			s += sep + packToStr(val)
			sep = ", "
		}
		s += ")"
	}
	return s
}

func (cmp *cmpExpr) toFilter() filter {
	if cmp.op == tok.In {
		return filter{vals: cmp.vals}
	}
	// else binary
	switch cmp.op {
	case tok.Is:
		return filter{vals: []string{cmp.val}}
	case tok.Lt:
		return filter{end: limit{val: cmp.val}}
	case tok.Lte:
		return filter{end: limit{val: cmp.val, inc: 1}}
	case tok.Gt:
		return filter{org: limit{val: cmp.val, inc: 1}, end: limitMax}
	case tok.Gte:
		return filter{org: limit{val: cmp.val}, end: limitMax}
	default:
		panic(assert.ShouldNotReachHere())
	}
}

//-------------------------------------------------------------------

// filter is either a range or a list of values
type filter struct {
	// org is inclusive (>=)
	org limit
	// end is exclusive (<)
	end limit
	// vals is nil for a range (org and end)
	vals []string
}

var filterAll = filter{org: limit{val: ixkey.Min}, end: limit{val: ixkey.Max}}

func (f filter) String() string {
	if f.all() {
		return "<all>"
	}
	if f.none() {
		return "<none>"
	}
	if f.vals == nil {
		return "(" + f.org.String() + ".." + f.end.String() + ")"
	}
	if len(f.vals) == 1 {
		return packToStr(f.vals[0])
	}
	s := "["
	sep := ""
	for _, v := range f.vals {
		s += sep + packToStr(v)
		sep = ","
	}
	return s + "]"
}

func (f *filter) isRange() bool {
	return len(f.vals) == 0
}

func (f *filter) isPoint() bool {
	return len(f.vals) == 1
}

func (f *filter) all() bool {
	return len(f.vals) == 0 && f.org == limitMin && f.end == limitMax
}

func (f *filter) none() bool {
	return len(f.vals) == 0 && compare(f.org, f.end) >= 0
}

func (f *filter) andWith(f2 filter) {
	if f.isRange() && f2.isRange() {
		if compare(f2.org, f.org) > 0 {
			f.org = f2.org
		}
		if compare(f2.end, f.end) < 0 {
			f.end = f2.end
		}
	} else if !f.isRange() && !f2.isRange() {
		f.vals = set.Intersect(f.vals, f2.vals)
		f.org, f.end = limit{}, limit{}
	} else { // set & range => set
		if f.isRange() {
			f.vals = f2.vals
			f2.org, f2.end = f.org, f.end
		}
		vals := make([]string, 0, len(f.vals)/2)
		for _, v := range f.vals {
			lim := limit{val: v}
			if compare(f2.org, lim) <= 0 && compare(lim, f2.end) < 0 {
				vals = append(vals, v)
			}
		}
		f.vals = vals
		f.org, f.end = limit{}, limit{}
	}
}

type limit struct {
	val string
	inc int
}

var limitMin = limit{}
var limitMax = limit{val: ixkey.Max}

func compare(x, y limit) int {
	cmp := strings.Compare(x.val, y.val)
	if cmp == 0 {
		cmp = ord.Compare(x.inc, y.inc)
	}
	return cmp
}

// valRaw is for non-encoded (single field keys)
func (lim limit) valRaw() string {
	if lim.inc == 0 {
		return lim.val
	}
	return lim.val + "\x00"
}

func (lim limit) String() string {
	s := packToStr(lim.val)
	if lim.inc > 0 {
		s += "+"
	}
	return s
}

func packToStr(s string) string {
	if s == "" {
		return "''"
	}
	if s[0] == 0xff {
		return "<max>"
	}
	return runtime.Unpack(s).String()
}

//-------------------------------------------------------------------

type idxSel struct {
	index   []string
	encoded bool
	ptrngs  []pointRange
	frac    float64
}

func (is idxSel) String() string {
	s := str.Join(",", is.index)
	sep := ": "
	for _, pr := range is.ptrngs {
		s += sep + showKey(is.encoded, pr.org)
		sep = " | "
		if pr.isRange() {
			s += ".." + showKey(is.encoded, pr.end)
		}
	}
	if is.frac != 0 {
		s += " = " + strconv.FormatFloat(is.frac, 'g', 4, 64)
	}
	return s
}

func showKey(encode bool, key string) string {
	if !encode {
		return packToStr(key)
	}
	s := ""
	sep := ""
	for _, t := range ixkey.Decode(key) {
		s += sep + packToStr(t)
		sep = ","
	}
	return s
}

func (is idxSel) isRanges() bool {
	return is.ptrngs[0].isRange()
}

// singleton returns true if we know the result is at most one record
// because there is a single point select on a key
func (is idxSel) singleton() bool {
	return len(is.ptrngs) == 1 && is.ptrngs[0].end == ""
}

// colSelsToIdxSels takes filters on individual columns
// and returns filters on indexes (possibly multi-column).
// We can use an index if leading columns have vals
// and optionally a final column has a range.
// This does not need to be all the index columns, it can be a prefix of them.
// It also sets w.singleton
func (w *Where) colSelsToIdxSels(colSels map[string]filter) []idxSel {
	indexes := w.tbl.Indexes()
	idxSels := make([]idxSel, 0, len(indexes)/2)
	for i := range w.tbl.schema.Indexes {
		schix := &w.tbl.schema.Indexes[i]
		idx := schix.Columns
		key := schix.Mode == 'k'
		uniq := schix.Mode == 'u'
		encode := !key || len(idx) > 1
		if filters := colSelsToIdxFilters(colSels, idx); len(filters) > 0 {
			exploded := explodeFilters(filters, [][]filter{nil})
			comp := compositePtrngs(encode, exploded)
			for i := range comp {
				c := &comp[i]
				if c.isPoint() {
					lookup := len(exploded[i]) == len(idx) &&
						(key || (uniq && c.org != ""))
					if !lookup {
						// convert point to range
						if !encode {
							c.end = c.org + "\x00"
						} else {
							c.end = c.org + ixkey.Sep + ixkey.Max
						}
					}
				}
			}
			frac := w.idxFrac(idx, comp)
			idxSel := idxSel{index: idx, ptrngs: comp, frac: frac, encoded: encode}
			w.singleton = w.singleton || idxSel.singleton()
			idxSels = append(idxSels, idxSel)
		}
	}
	return idxSels
}

// colSelsToIdxFilters returns filters for the columns of an index
// or else nil if not possible.
func colSelsToIdxFilters(colSels map[string]filter, idx []string) []filter {
	const maxCount = 8 * 1024 // ???
	filters := make([]filter, 0, len(idx))
	count := 1
	for _, col := range idx {
		f, ok := colSels[col]
		if !ok {
			break
		}
		if len(f.vals) > 1 {
			count *= len(f.vals)
			if count > maxCount {
				return nil
			}
		}
		filters = append(filters, f)
		if f.isRange() {
			break // can't have anything after range
		}
	}
	return filters
}

// pointRange holds either a range or a single key (in org with end = "")
type pointRange struct {
	org string
	end string
}

func (pr pointRange) isPoint() bool {
	return pr.end == ""
}

func (pr pointRange) isRange() bool {
	return pr.end != ""
}

func (pr pointRange) String() string {
	if pr.conflict() {
		return "<empty>"
	}
	// WARNING: does NOT decode, intended for explode output
	// use idxSel.String for compositePtrngs output
	s := packToStr(pr.org)
	if pr.isRange() {
		s += ".." + packToStr(pr.end)
	}
	return s
}

func (pr pointRange) intersect(selOrg, selEnd string) pointRange {
	if pr.isPoint() {
		if pr.org == selOrg {
			return pr
		}
	} else { // range
		if pr.org <= selOrg && selOrg < pr.end {
			return pointRange{org: selOrg, end: selEnd}
		}
	}
	return pointRange{org: "z", end: "a"} // conflict
}

func (pr pointRange) conflict() bool {
	return pr.end != "" && pr.end < pr.org
}

// explodeFilters converts multi vals filters to multi single vals filters.
// e.g. (1|2)+(3|4) => 1+3, 1+4, 2+3, 2+4
// Each recursion processes one element from remaining, adding to prefixes.
func explodeFilters(remaining []filter, prefixes [][]filter) [][]filter {
	f := remaining[0]
	if len(f.vals) <= 1 { // single value or final range
		for i := range prefixes {
			prefixes[i] = append(prefixes[i], f)
		}
	} else { // len(f.vals) > 1
		newpre := make([][]filter, 0, len(f.vals)*len(prefixes))
		for i := range prefixes {
			pre := prefixes[i]
			// set cap to len so each append will make a new copy (COW)
			pre = pre[:len(pre):len(pre)]
			for _, v := range f.vals {
				p := append(pre, filter{vals: []string{v}})
				newpre = append(newpre, p)
			}
		}
		prefixes = newpre
	}
	if len(remaining) > 1 {
		return explodeFilters(remaining[1:], prefixes) // RECURSE
	}
	return prefixes
}

func compositePtrngs(encode bool, filters [][]filter) []pointRange {
	result := make([]pointRange, len(filters))
outer:
	for i, fs := range filters {
		if !encode {
			assert.That(len(fs) == 1)
			f := fs[0]
			if f.isPoint() {
				result[i] = pointRange{org: f.vals[0]}
			} else { // range
				result[i] = pointRange{org: f.org.valRaw(), end: f.end.valRaw()}
			}
		} else {
			var enc ixkey.Encoder
			for _, f := range fs {
				if f.isPoint() {
					enc.Add(f.vals[0])
				} else { // final range
					enc2 := enc.Dup()
					enc.Add(f.org.val)
					if f.org.inc == 1 {
						enc.Add(ixkey.Max)
					}
					enc2.Add(f.end.val)
					if f.end.inc == 1 {
						enc2.Add(ixkey.Max)
					}
					result[i] = pointRange{org: enc.String(), end: enc2.String()}
					continue outer
				}
			}
			result[i] = pointRange{org: enc.String()}
		}
	}
	return result
}

func (w *Where) idxFrac(idx []string, ptrngs []pointRange) float64 {
	iIndex := slc.IndexFn(w.tbl.indexes, idx, slices.Equal[string])
	if iIndex < 0 {
		panic("index not found")
	}
	var frac float64
	nrows1, _ := w.tbl.Nrows()
	for _, pr := range ptrngs {
		if pr.end == "" { // lookup
			if nrows1 > 0 {
				frac += 1 / float64(nrows1)
			}
		} else { // range
			frac += float64(w.t.RangeFrac(w.tbl.name, iIndex, pr.org, pr.end))
		}
	}
	assert.That(!math.IsNaN(frac) && !math.IsInf(frac, 0))
	return frac
}

// execution --------------------------------------------------------

// MakeSuTran is injected by dbms to avoid import cycle
var MakeSuTran func(qt QueryTran) *runtime.SuTran

func (w *Where) Get(th *runtime.Thread, dir runtime.Dir) runtime.Row {
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

func (w *Where) get(th *runtime.Thread, dir runtime.Dir) runtime.Row {
	if w.idxSel == nil {
		w.nIn++
		return w.source.Get(th, dir)
	}
	if w.idxSel.isRanges() {
		return w.getRange(th, dir)
	}
	return w.getPoint(dir)
}

func (w *Where) filter(th *runtime.Thread, row runtime.Row) bool {
	if row == nil {
		return true
	}
	if w.selectCols != nil &&
		!singletonFilter(w.hdr, row, w.selectCols, w.selectVals) {
		return false
	}
	if w.ctx.Tran == nil {
		w.ctx.Tran = MakeSuTran(w.t)
	}
	w.ctx.Th = th
	w.ctx.Row = row
	defer func() { w.ctx.Th, w.ctx.Row = nil, nil }()
	return w.expr.Eval(&w.ctx) == runtime.True
}

func (w *Where) getRange(th *runtime.Thread, dir runtime.Dir) runtime.Row {
	for {
		if w.idxSelPos != -1 {
			if row := w.tbl.Get(th, dir); row != nil {
				return row
			}
		}
		if !w.advance(dir) {
			return nil // eof
		}
		w.tbl.SelectRaw(w.curPtrng.org, w.curPtrng.end)
	}
}

func (w *Where) getPoint(dir runtime.Dir) runtime.Row {
	for {
		if !w.advance(dir) {
			return nil
		}
		if row := w.tbl.lookup(w.curPtrng.org); row != nil {
			w.nIn++
			return row
		}
	}
}

func (w *Where) advance(dir runtime.Dir) bool {
	if w.idxSelPos == -1 { // rewound
		if dir == runtime.Prev {
			w.idxSelPos = len(w.idxSel.ptrngs)
		}
	}
	for {
		if dir == runtime.Prev {
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
	w.Rewind()
	if cols == nil && vals == nil { // clear select
		w.selOrg, w.selEnd = "", ""
		w.selSet = false
		w.selectCols = nil
		w.selectVals = nil
		return
	}
	satisfied, conflict := selectFixed(cols, vals, w.whereFixed)
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

	cols = slices.Clip(cols)
	vals = slices.Clip(vals)
	for _, fix := range w.fixed {
		if len(fix.values) == 1 && !slices.Contains(cols, fix.col) {
			cols = append(cols, fix.col)
			vals = append(vals, fix.values[0])
		}
	}
	w.selOrg, w.selEnd = selKeys(w.idxSel.encoded, w.idxSel.index, cols, vals)
	w.selSet = true
}

func selectFixed(cols, vals []string, fixed []Fixed) (satisfied, conflict bool) {
	// Note: conflict could come from any of expr, not just fixed.
	// But to evaluate that would require building a Row.
	// It should be rare.
	satisfied = true
	for i, col := range cols {
		if fv := getFixed(fixed, col); len(fv) == 1 {
			if fv[0] != vals[i] {
				return false, true // conflict
			}
		} else {
			satisfied = false
		}
	}
	return satisfied, false
}

func (w *Where) Lookup(th *runtime.Thread, cols, vals []string) runtime.Row {
	if w.singleton {
		// can't use source.Lookup because cols may not match source index
		w.Rewind()
		row := w.Get(th, runtime.Next)
		if row == nil || !singletonFilter(w.hdr, row, cols, vals) {
			return nil
		}
		return row
	}
	row := w.source.Lookup(th, cols, vals)
	if !w.filter(th, row) {
		row = nil
	}
	return row
}

func singletonFilter(
	hdr *runtime.Header, row runtime.Row, cols []string, vals []string) bool {
	for i, col := range cols {
		if row.GetRaw(hdr, col) != vals[i] {
			return false
		}
	}
	return true
}

func (w *Where) slowQueries() {
	if w.nIn > 100 && w.nIn > w.nOut*100 && trace.SlowQuery.On() {
		trace.SlowQuery.Println(w.nIn, "->", w.nOut)
		trace.Println(format(w, 1))
		w.nIn = 0
		w.nOut = 0
	}
}
