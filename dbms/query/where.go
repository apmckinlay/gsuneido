// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"sort"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/compile/ast"
	"github.com/apmckinlay/gsuneido/compile/tokens"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/sset"
	"github.com/apmckinlay/gsuneido/util/strs"
	"github.com/apmckinlay/gsuneido/util/strss"
)

type Where struct {
	Query1
	expr  *ast.Nary // And
	fixed []Fixed
	t     QueryTran
	// tbl will be set if the source is a Table, nil otherwise
	tbl       *Table
	optInited bool
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
	sel     string
	selSet  bool
	context *ast.Context
}

type whereApproach struct {
	index []string
}

func (w *Where) Init() {
	w.Query1.Init()
	if !sset.Subset(w.source.Columns(), w.expr.Columns()) {
		panic("where: nonexistent columns: " + strs.Join(", ",
			sset.Difference(w.expr.Columns(), w.source.Columns())))
	}
	w.tbl, _ = w.source.(*Table)
}

func (w *Where) SetTran(t QueryTran) {
	w.t = t
	w.source.SetTran(t)
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
				w.fixed = append(w.fixed, NewFixed(id.Name, c.Val))
			}
		}
	}
}

func (w *Where) Nrows() int {
	if w.conflict {
		return 0
	}
	if len(w.idxSels) == 0 {
		return w.source.Nrows() / 2
	}
	var n int
	nmin := ints.MaxInt
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
				q.source.Init()
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
				q.source.Init()
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
				q.source.Init()
			}
			if rest != nil {
				w.expr = &ast.Nary{Tok: tok.And, Exprs: rest}
			} else {
				moved = true
			}
		case *Intersect:
			// distribute where over intersect
			q.source = &Where{Query1: Query1{source: q.source}, expr: w.expr}
			q.source.Init()
			q.source2 = &Where{Query1: Query1{source: q.source2}, expr: w.expr}
			q.source2.Init()
			moved = true
		case *Minus:
			// distribute where over minus
			q.source = &Where{Query1: Query1{source: q.source}, expr: w.expr}
			q.source.Init()
			q.source2 = &Where{Query1: Query1{source: q.source2},
				expr: w.project(q.source2)}
			q.source2.Init()
			moved = true
		case *Union:
			// distribute where over union
			q.source = &Where{Query1: Query1{source: q.source},
				expr: w.project(q.source)}
			q.source.Init()
			q.source2 = &Where{Query1: Query1{source: q.source2},
				expr: w.project(q.source2)}
			q.source.Init()
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
				q.source.Init()
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
		cols := lj.source2.Header().GetFields()
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

// optimize ---------------------------------------------------------

func (w *Where) optimize(mode Mode, index []string) (Cost, interface{}) {
	w.conflictCheck()
	if w.conflict {
		return 0, nil
	}
	// we always have the option of just filtering (no specific index use)
	filterCost := Optimize(w.source, mode, index)
	if w.tbl == nil || w.tbl.singleton {
		return filterCost, nil
	}
	if !w.optInited {
		w.optInit()
	}
	cost, index := w.bestIndex(index)
	if cost >= filterCost {
		return filterCost, nil
	}
	return cost, whereApproach{index: index}
}

func (w *Where) conflictCheck() {
	for _, expr := range w.expr.Exprs {
		if c, ok := expr.(*ast.Constant); ok && c.Val == runtime.False {
			w.conflict = true
			return
		}
	}
}

func (w *Where) optInit() {
	w.Fixed()                   // calc before altering expression
	cmps := w.extractCompares() // NOTE: modifies expr
	if w.conflict {
		return
	}
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
			if bin, ok := expr.(*ast.Binary); ok && bin.Tok != tokens.Isnt {
				cmp := cmpExpr{
					col: bin.Lhs.(*ast.Ident).Name,
					op:  bin.Tok,
					val: bin.Rhs.(*ast.Constant).Packed,
				}
				cmps = append(cmps, cmp)
			} else if in, ok := expr.(*ast.In); ok {
				cmp := cmpExpr{
					col:  in.E.(*ast.Ident).Name,
					op:   tokens.In,
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

// bestIndex returns the best (lowest cost) index
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
					tblCost, _ := w.tbl.optimize(0, idx)
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
		if strs.Equal(index, is.index) {
			return is
		}
	}
	return nil
}

func (w *Where) setApproach(index []string, app interface{}, tran QueryTran) {
	if w.conflict {
		return
	}
	if app != nil {
		idx := app.(whereApproach).index
		assert.Msg("index", index, "idx", idx).
			That(index == nil || w.singleton || strs.HasPrefix(idx, index))
		w.tbl.setIndex(idx)
		w.idxSel = w.getIdxSel(idx)
		w.idxSelPos = -1
	} else {
		w.source = SetApproach(w.source, index, tran)
	}
}

// cmpExpr is <field> <op> <constant> or <field> in (<constants>)
// which can be evaluated packed
type cmpExpr struct {
	col  string
	op   tokens.Token
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
	if cmp.op == tokens.In {
		return filter{vals: cmp.vals}
	}
	// else binary
	switch cmp.op {
	case tokens.Is:
		return filter{vals: []string{cmp.val}}
	case tokens.Lt:
		return filter{end: limit{val: cmp.val}}
	case tokens.Lte:
		return filter{end: limit{val: cmp.val, inc: 1}}
	case tokens.Gt:
		return filter{org: limit{val: cmp.val, inc: 1}, end: limitMax}
	case tokens.Gte:
		return filter{org: limit{val: cmp.val}, end: limitMax}
	default:
		panic("shouldn't reach here")
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
	return len(f.vals) == 0 && cmp(f.org, f.end) >= 0
}

func (f *filter) andWith(f2 filter) {
	if f.isRange() && f2.isRange() {
		if cmp(f2.org, f.org) > 0 {
			f.org = f2.org
		}
		if cmp(f2.end, f.end) < 0 {
			f.end = f2.end
		}
	} else if !f.isRange() && !f2.isRange() {
		f.vals = sset.Intersect(f.vals, f2.vals)
		f.org, f.end = limit{}, limit{}
	} else { // set & range => set
		if f.isRange() {
			f.vals = f2.vals
			f2.org, f2.end = f.org, f.end
		}
		vals := make([]string, 0, len(f.vals)/2)
		for _, v := range f.vals {
			lim := limit{val: v}
			if cmp(f2.org, lim) <= 0 && cmp(lim, f2.end) < 0 {
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

func cmp(x, y limit) int {
	cmp := strings.Compare(x.val, y.val)
	if cmp == 0 {
		cmp = ints.Compare(x.inc, y.inc)
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
	s := strs.Join(",", is.index)
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
						enc.Add("")
					}
					enc2.Add(f.end.val)
					if f.end.inc == 1 {
						enc2.Add("")
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
	iIndex := strss.Index(w.tbl.indexes, idx)
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

func (w *Where) Header() *runtime.Header {
	return w.source.Header()
}

// MakeSuTran is injected by dbms to avoid import cycle
var MakeSuTran func(qt QueryTran) *runtime.SuTran

func (w *Where) Get(dir runtime.Dir) runtime.Row {
	if w.conflict {
		return nil
	}
	for {
		row := w.get(dir)
		if w.filter(row) {
			return row
		}
	}
}

func (w *Where) get(dir runtime.Dir) runtime.Row {
	if w.idxSel == nil {
		return w.source.Get(dir)
	}
	if w.idxSel.isRanges() {
		return w.getRange(dir)
	}
	return w.getPoint(dir)
}

func (w *Where) filter(row runtime.Row) bool {
	if row == nil {
		return true
	}
	if w.context == nil {
		w.context = &ast.Context{T: &runtime.Thread{}}
	}
	if w.hdr == nil {
		w.hdr = w.source.Header()
	}
	w.context.Rec = runtime.SuRecordFromRow(row, w.hdr, "", MakeSuTran(w.t))
	return w.expr.Eval(w.context) == runtime.True
}

func (w *Where) getRange(dir runtime.Dir) runtime.Row {
	for {
		if w.idxSelPos != -1 {
			if row := w.tbl.Get(dir); row != nil {
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
	if !w.advance(dir) {
		return nil
	}
	return w.tbl.lookup(w.curPtrng.org)
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
	cols = strs.Cow(cols)
	vals = strs.Cow(vals)
	for _, fix := range w.Fixed() {
		if len(fix.values) == 1 {
			cols = append(cols, fix.col)
			vals = append(vals, fix.values[0])
		}
	}
	if w.idxSel == nil {
		w.source.Select(cols, vals)
		return
	}
	w.sel = selEncode(w.idxSel.encoded, w.idxSel.index, cols, vals)
	w.Rewind()
	w.selSet = true
}

func (w *Where) Lookup(cols, vals []string) runtime.Row {
	row := w.source.Lookup(cols, vals)
	if !w.filter(row) {
		row = nil
	}
	return row
}
