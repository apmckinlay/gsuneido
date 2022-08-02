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
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/ord"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/str"
	"golang.org/x/exp/slices"
)

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
	// singleton is true if we know the result is at most one record
	// because there is a single point select on a key
	singleton bool
	idxSels   []idxSel
	// conflict is true if the expression conflicts and selects nothing
	conflict bool
	// idxSel is for the chosen index
	idxSel *idxSel
	// idxSelPos is the current index in idxSel.ptrngs
	idxSelPos int
	// curPtrng is idxSel.ptrngs[idxSelPos] adjusted by Select (selOrg, selEnd)
	curPtrng pointRange
	hdr      *runtime.Header
	// sel is set by Select
	sel    string
	selSet bool

	selectCols []string
	selectVals []string
}

type whereApproach struct {
	index []string
}

func NewWhere(src Query, expr ast.Expr, t QueryTran) *Where {
	if !set.Subset(src.Columns(), expr.Columns()) {
		panic("where: nonexistent columns: " + str.Join(", ",
			set.Difference(expr.Columns(), src.Columns())))
	}
	if nary, ok := expr.(*ast.Nary); !ok || nary.Tok != tok.And {
		expr = &ast.Nary{Tok: tok.And, Exprs: []ast.Expr{expr}}
	}
	return &Where{Query1: Query1{source: src}, expr: expr.(*ast.Nary), t: t}
}

func (w *Where) SetTran(t QueryTran) {
	w.t = t
	w.source.SetTran(t)
	w.ctx.Tran = nil
}

func (w *Where) String() string {
	s := parenQ2(w.source) + " WHERE"
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
	if w.fixed != nil {
		return w.fixed
	}
	whereFixed := w.exprsToFixed()
	fixed, none := combineFixed(w.source.Fixed(), whereFixed)
	if none {
		w.conflict = true
	}
	return fixed
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

func (w *Where) Nrows() int {
	if w.conflict {
		return 0
	}
	if len(w.idxSels) == 0 {
		return w.source.Nrows() / 2
	}
	var n int
	nmin := math.MaxInt
	nsrc := float64(w.source.Nrows())
	for i := range w.idxSels {
		ix := &w.idxSels[i]
		if ix.isRanges() {
			n = int(ix.frac * nsrc)
		} else { // points
			n = len(ix.ptrngs)
		}
		if n < nmin {
			nmin = n
		}
	}
	if w.exprMore {
		nmin /= 2 // ??? adjust for additional restrictions
	}
	return nmin
}

func (w *Where) Transform() Query {
	if len(w.expr.Exprs) == 0 {
		// remove empty where
		return w.source.Transform()
	}
	if w.Fixed(); w.conflict || w.exprConflict() {
		return NewNothing(w.Columns())
	}
	if tl := w.tablesLookup(); tl != nil {
		return tl
	}

	if lj := w.leftJoinToJoin(); lj != nil {
		// convert leftjoin to join
		w.source = NewJoin(lj.source, lj.source2, lj.by)
	}
	moved := false
	switch q := w.source.(type) {
	case *Where:
		// combine consecutive where's
		for {
			n := len(q.expr.Exprs)
			exprs := append(q.expr.Exprs[:n:n], w.expr.Exprs...) // copy on write
			w.expr = &ast.Nary{Tok: tok.And, Exprs: exprs}
			w.source = q.source
			if q, _ = w.source.(*Where); q == nil {
				break
			}
		}
		return w.Transform() // RECURSE
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
			q.source = NewWhere(q.source, newExpr.(*ast.Nary), w.t)
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
			q.source = NewWhere(q.source,
				&ast.Nary{Tok: tok.And, Exprs: src1}, w.t)
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
			if set.Subset(cols1, e.Columns()) {
				src1 = append(src1, e)
			} else {
				rest = append(rest, e)
			}
		}
		if src1 != nil {
			q.source = NewWhere(q.source,
				&ast.Nary{Tok: tok.And, Exprs: src1}, w.t)
		}
		if rest != nil {
			w.expr = &ast.Nary{Tok: tok.And, Exprs: rest}
		} else {
			moved = true
		}
	case *Intersect:
		// distribute where over intersect
		// no project because Intersect Columns are the intersection
		q.source = NewWhere(q.source, w.expr, w.t)
		q.source2 = NewWhere(q.source2, w.expr, w.t)
		moved = true
	case *Minus:
		// distribute where over minus
		// need project because Minus Columns are just the left side's
		q.source = NewWhere(q.source, w.expr, w.t)
		q.source2 = NewWhere(q.source2, w.project(q.source2), w.t)
		moved = true
	case *Union:
		// distribute where over union
		// need project because Union Columns is the union
		q.source = NewWhere(q.source, w.project(q.source), w.t)
		q.source2 = NewWhere(q.source2, w.project(q.source2), w.t)
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
			if set.Subset(cols1, e.Columns()) {
				src1 = append(src1, e)
			} else {
				common = append(common, e)
			}
		}
		if src1 != nil {
			q.source = NewWhere(q.source,
				&ast.Nary{Tok: tok.And, Exprs: src1}, w.t)
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
	// propagate Nothing
	if _, ok := w.source.(*Nothing); ok {
		return w.source
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

func (w *Where) split(q2 *Query2) bool {
	cols1 := q2.source.Columns()
	cols2 := q2.source2.Columns()
	var common, src1, src2 []ast.Expr
	for _, e := range w.expr.Exprs {
		used := false
		if set.Subset(cols1, e.Columns()) {
			src1 = append(src1, e)
			used = true
		}
		if set.Subset(cols2, (e.Columns())) {
			if used {
				e = replaceExpr(e, nil, nil) // copy Binary/In CouldEvalRaw
			}
			src2 = append(src2, e)
			used = true
		}
		if !used {
			common = append(common, e)
		}
	}
	if src1 != nil {
		q2.source = NewWhere(q2.source, &ast.Nary{Tok: tok.And, Exprs: src1}, w.t)
	}
	if src2 != nil {
		q2.source2 = NewWhere(q2.source2, &ast.Nary{Tok: tok.And, Exprs: src2}, w.t)
	}
	if common != nil {
		w.expr = &ast.Nary{Tok: tok.And, Exprs: common}
	} else {
		return true
	}
	return false
}

// optimize ---------------------------------------------------------

func (w *Where) optimize(mode Mode, index []string) (Cost, any) {
	if !w.optInited {
		w.optInit()
		if w.conflict {
			return 0, nil
		}
	}
	// we always have the option of just filtering (no specific index use)
	filterCost := Optimize(w.source, mode, index)
	if w.tbl == nil || w.tbl.singleton {
		return filterCost, nil
	}
	cost, index := w.bestIndex(index)
	if cost >= impossible {
		// only use the filter if there are no possible idxSel
		return filterCost, nil
	}
	return cost, whereApproach{index: index}
}

func (w *Where) exprConflict() bool {
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
	if w.tbl, _ = w.source.(*Table); w.tbl == nil {
		return
	}
	w.whereFixed = w.exprsToFixed()
	w.fixed = w.Fixed() // cache
	cmps := w.extractCompares()
	w.exprMore = len(cmps) < len(w.expr.Exprs)
	colSels := w.comparesToFilters(cmps)
	if w.conflict {
		return
	}
	w.idxSels = w.colSelsToIdxSels(colSels)
	w.exprMore = w.exprMore || len(w.idxSels) > 1
	if !w.exprMore {
		// check if any colSels were not used by idxSels
		for _, idxSel := range w.idxSels {
			for _, col := range idxSel.index {
				delete(colSels, col)
			}
		}
		w.exprMore = w.exprMore || len(colSels) > 0
	}
}

// extractCompares finds sub-expressions like <field> <op> <constant>
func (w *Where) extractCompares() []cmpExpr {
	cols := w.tbl.schema.Columns
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
func (w *Where) bestIndex(order []string) (Cost, []string) {
	if w.singleton {
		order = nil
	}
	best := newBestIndex()
	for _, idx := range w.source.Indexes() {
		if ordered(idx, order, w.fixed) {
			if is := w.getIdxSel(idx); is != nil {
				cost := w.source.lookupCost() * len(is.ptrngs)
				if is.isRanges() {
					tblCost, _ := w.tbl.optimize(CursorMode, idx)
					cost += int(is.frac * float64(tblCost))
				}
				best.update(idx, cost)
			}
		}
	}
	return best.cost, best.index
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

func (w *Where) setApproach(index []string, app any, tran QueryTran) {
	if w.conflict {
		return
	}
	if app != nil {
		idx := app.(whereApproach).index
		w.tbl.setIndex(idx)
		w.idxSel = w.getIdxSel(idx)
		w.idxSelPos = -1
	} else {
		w.source = SetApproach(w.source, index, tran)
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
		assert.ShouldNotReachHere()
		return filter{}
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

func (pr pointRange) intersect(sel string, encode bool) pointRange {
	if pr.isPoint() {
		if pr.org == sel {
			return pr
		}
	} else { // range
		if pr.org <= sel && sel < pr.end {
			pr.org = sel
			if !encode {
				pr.end = sel + "\x00"
			} else {
				pr.end = sel + ixkey.Sep + ixkey.Max
			}
			return pr
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
	nrows := float64(w.tbl.Nrows())
	for _, pr := range ptrngs {
		if pr.end == "" { // lookup
			frac += 1 / nrows
		} else { // range
			frac += float64(w.t.RangeFrac(w.tbl.name, iIndex, pr.org, pr.end))
		}
	}
	return frac
}

// execution --------------------------------------------------------

// MakeSuTran is injected by dbms to avoid import cycle
var MakeSuTran func(qt QueryTran) *runtime.SuTran

func (w *Where) Get(th *runtime.Thread, dir runtime.Dir) runtime.Row {
	if w.conflict {
		return nil
	}
	for {
		row := w.get(th, dir)
		if w.filter(th, row) {
			return row
		}
	}
}

func (w *Where) get(th *runtime.Thread, dir runtime.Dir) runtime.Row {
	if w.idxSel == nil {
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
		!w.singletonFilter(row, w.selectCols, w.selectVals) {
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
			pr = pr.intersect(w.sel, w.idxSel.encoded)
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
	if w.conflict {
		return
	}

	w.Rewind()
	if cols == nil && vals == nil { // clear select
		w.sel = ""
		w.selSet = false
		w.selectCols = nil
		w.selectVals = nil
		return
	}
	satisfied, conflict := w.selectFixed(cols, vals)
	if conflict {
		w.sel = ixkey.Max
		w.selSet = true
		return
	}
	if satisfied {
		w.sel = ""
		w.selSet = false
		return
	}

	if w.singleton {
		w.selectCols = cols
		w.selectVals = vals
		return
	}

	cols = slices.Clip(cols)
	vals = slices.Clip(vals)
	for _, fix := range w.Fixed() {
		if len(fix.values) == 1 {
			cols = append(cols, fix.col)
			vals = append(vals, fix.values[0])
		}
	}
	if w.idxSel == nil {
		// if this where is not using an index selection
		// then just pass the Select to the source
		w.source.Select(cols, vals)
		return
	}
	w.sel = selOrg(w.idxSel.encoded, w.idxSel.index, cols, vals, false)
	w.selSet = true
}

func (w *Where) selectFixed(cols, vals []string) (satisfied, conflict bool) {
	// Note: conflict could come from any of expr, not just fixed.
	// But to evaluate that would require building a Row.
	// It should be rare.
	w.Fixed()
	satisfied = true
	for i, col := range cols {
		if fv := getFixed(w.whereFixed, col); len(fv) == 1 {
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
	if w.conflict {
		return nil
	}
	if w.singleton {
		// can't use source.Lookup because cols may not match source index
		w.Rewind()
		row := w.Get(th, runtime.Next)
		if row == nil || !w.singletonFilter(row, cols, vals) {
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

func (w *Where) singletonFilter(
	row runtime.Row, cols []string, vals []string) bool {
	for i, col := range cols {
		if row.GetRaw(w.hdr, col) != vals[i] {
			return false
		}
	}
	return true
}
