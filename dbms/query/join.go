// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"slices"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/trace"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/hash"
	"github.com/apmckinlay/gsuneido/util/lrucache"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/tsc"
)

/*
	joinLike
		Times - symmetric
		joinBase
			Join - symmetric
			LeftJoin - asymmetric
*/

// joinLike is common stuff for Join, LeftJoin, and Times
type joinLike struct {
	// sel2cols/vals are from an incoming Select and used by Get
	sel2cols, sel2vals []string
	Query2
}

// joinBase is common stuff for Join and LeftJoin
type joinBase struct {
	qt         QueryTran
	st         *SuTran
	lookup     *lookupInfo
	cache      *lrucache.Cache[lookupKey, Row]
	by         []string
	prevFixed1 []Fixed
	prevFixed2 []Fixed
	row1       Row
	row2       Row // nil when we need a new row1
	lookupRow  Row
	joinLike
	joinType
	optimized bool
}

type Join struct {
	joinBase
	conflict bool
}

type lookupInfo struct {
	keys1    [][]string
	fixed1   []Fixed
	fallback bool
}

type joinApproach struct {
	index1  []string
	index2  []string
	frac2   float64
	reverse bool
}

type joinType int

const (
	one_one joinType = iota + 1 //lint:ignore ST1003 for clarity
	one_n                       //lint:ignore ST1003 for clarity
	n_one                       //lint:ignore ST1003 for clarity
	n_n                         //lint:ignore ST1003 for clarity
)

func (jt joinType) String() string {
	switch jt {
	case 0:
		return ""
	case one_one:
		return "1:1"
	case one_n:
		return "1:n"
	case n_one:
		return "n:1"
	case n_n:
		return "n:n"
	default:
		panic("bad joinType")
	}
}

func NewJoin(src1, src2 Query, by []string, t QueryTran) Query {
	return newJoin(src1, src2, by, t, nil, nil)
}

func newJoin(src1, src2 Query, by []string, t QueryTran,
	prevFixed1, prevFixed2 []Fixed) *Join {
	jn := &Join{joinBase: newJoinBase(src1, src2, by, t,
		prevFixed1, prevFixed2)}
	jn.keys = jn.getKeys()
	jn.indexes = jn.getIndexes()
	fixed, none := combineFixed(src1.Fixed(), src2.Fixed())
	jn.conflict = none
	jn.fixed = fixed
	jn.setNrows(jn.getNrows())
	jn.fast1.Set(src1.fastSingle() && src2.fastSingle())
	return jn
}

func (jn *Join) With(src1, src2 Query) *Join {
	return newJoin(src1, src2, jn.by, jn.qt, jn.prevFixed1, jn.prevFixed2)
}

func newJoinBase(src1, src2 Query, by []string, t QueryTran,
	prevFixed1, prevFixed2 []Fixed) joinBase {
	b := set.Intersect(src1.Columns(), src2.Columns())
	if len(b) == 0 {
		panic("join: common columns required")
	}
	if by == nil {
		by = b
	} else if !set.Equal(by, b) {
		panic("join: by does not match common columns")
	}
	jb := joinBase{qt: t, st: MakeSuTran(t), joinLike: newJoinLike(src1, src2), prevFixed1: prevFixed1, prevFixed2: prevFixed2}
	jb.by = by
	k1 := hasKey(by, src1.Keys(), src1.Fixed())
	k2 := hasKey(by, src2.Keys(), src2.Fixed())
	if k1 && k2 {
		jb.joinType = one_one
	} else if k1 {
		jb.joinType = one_n
	} else if k2 {
		jb.joinType = n_one
	} else {
		jb.joinType = n_n
	}
	return jb
}

func fixedToExpr(col string, values []string) ast.Expr {
	es := make([]ast.Expr, len(values))
	for i, v := range values {
		es[i] = &ast.Constant{Val: Unpack(v)}
	}
	var folder ast.Folder
	return folder.In(&ast.Ident{Name: col}, es)
}

func newJoinLike(src1, src2 Query) joinLike {
	jl := joinLike{}
	jl.source1, jl.source2 = src1, src2
	jl.header = jl.getHeader()
	jl.rowSiz.Set(jl.source1.rowSize() + jl.source2.rowSize())
	jl.lookCost.Set(src1.lookupCost() * 2) // ???
	return jl
}

func (jn *Join) String() string {
	return jn.joinBase.String("join")
}

func (jb *joinBase) String(op string) string {
	if jb.optimized {
		op += " " + jb.joinType.String()
	} else if jb.joinType == n_n {
		op += " /*MANY TO MANY*/"
	}
	return op + " by" + str.Join("(,)", jb.by)
}

func (jb *joinBase) SetTran(qt QueryTran) {
	jb.qt = qt
	jb.st = MakeSuTran(qt)
	jb.cache.Reset()
}

func (jl *joinLike) getHeader() *Header {
	return JoinHeaders(jl.source1.Header(), jl.source2.Header())
}

func (jn *Join) getIndexes() [][]string {
	// can really only provide source.indexes() but optimize may swap.
	// optimize will return impossible for source2 indexes.
	return set.UnionFn(jn.source1.Indexes(), jn.source2.Indexes(), slices.Equal)
}

func (jn *Join) getKeys() [][]string {
	switch jn.joinType {
	case one_one:
		return set.UnionFn(jn.source1.Keys(), jn.source2.Keys(), set.Equal[string])
	case one_n:
		return jn.source2.Keys()
	case n_one:
		return jn.source1.Keys()
	case n_n:
		return jn.keypairs()
	default:
		panic("unknown join type")
	}
}

func (jn *Join) Transform() Query {
	jn.optimized = true
	if jn.conflict {
		return NewNothing(jn)
	}
	src1 := jn.source1.Transform()
	if _, ok := src1.(*Nothing); ok {
		return NewNothing(jn)
	}
	src2 := jn.source2.Transform()
	if _, ok := src2.(*Nothing); ok {
		return NewNothing(jn)
	}
	fix1, fix2 := src1.Fixed(), src2.Fixed()
	if !equalFixed(fix1, jn.prevFixed1) || !equalFixed(fix2, jn.prevFixed2) {
		src1 = copyFixed(fix2, fix1, src1, jn.by, jn.qt)
		src2 = copyFixed(fix1, fix2, src2, jn.by, jn.qt)
		jn.prevFixed1, jn.prevFixed2 = fix1, fix2
	}
	if src1 != jn.source1 || src2 != jn.source2 {
		return jn.With(src1, src2).Transform()
	}
	return jn
}

// copyFixed adds a Where to `to` for fixed from `from` for common by columns.
// It doesn't check if the where is conflicting, that should be handled by Where
func copyFixed(fromFixed, toFixed []Fixed, to Query, by []string, t QueryTran) Query {
	var exprs []ast.Expr
	for _, col := range by {
		if frf := getFixed(fromFixed, col); frf != nil {
			if tof := getFixed(toFixed, col); !set.Equal(frf, tof) {
				exprs = append(exprs, fixedToExpr(col, frf))
			}
		}
	}
	if len(exprs) == 0 {
		return to
	}
	expr := &ast.Nary{Tok: tok.And, Exprs: exprs}
	return NewWhere(to, expr, t)
}

var joinRev = 0 // tests can set to impossible to prevent reverse

func (jn *Join) optimize(mode Mode, index []string, frac float64) (Cost, Cost, any) {
	fwd := joinopt(jn.source1, jn.source2, jn.Nrows,
		mode, index, frac, jn.by, jn.fixed)
	rev := joinopt(jn.source2, jn.source1, jn.Nrows,
		mode, index, frac, jn.by, jn.fixed)
	rev.fixcost += outOfOrder + joinRev
	if trace.JoinOpt.On() {
		trace.JoinOpt.Println(mode, index, frac)
		trace.Println("    fwd index1", fwd.index1, "index2", fwd.index2,
			"=", fwd.fixcost, fwd.varcost)
		trace.Println("    rev index1", rev.index1, "index2", rev.index2,
			"=", rev.fixcost, rev.varcost)
		trace.Println(strategy(jn, 1))
	}
	approach := &joinApproach{}
	if rev.fixcost+rev.varcost < fwd.fixcost+fwd.varcost {
		fwd = rev
		approach.reverse = true
	}
	if fwd.fixcost == impossible {
		return impossible, impossible, nil
	}
	approach.index1 = fwd.index1
	approach.index2 = fwd.index2
	approach.frac2 = fwd.frac2
	return fwd.fixcost, fwd.varcost, approach
}

func (jt joinType) reverse() joinType {
	switch jt {
	case one_n:
		return n_one
	case n_one:
		return one_n
	}
	return jt
}

type joinCost struct {
	index1  []string
	index2  []string
	frac2   float64
	fixcost Cost
	varcost Cost
}

func joinopt(src1, src2 Query, nrows func() (int, int),
	mode Mode, index []string, frac float64, by []string, fixed []Fixed) joinCost {
	// always have to read all of source 1
	fixcost1, varcost1, index := optOrdered(src1, mode, index, frac, fixed)
	if fixcost1+varcost1 >= impossible {
		return joinCost{fixcost: impossible}
	}
	nrows1, _ := src1.Nrows()
	nrows2, _ := src2.Nrows()
	read2, _ := nrows()
	frac2 := float64(read2) * frac / float64(max(1, nrows2))
	best2 := bestGrouped(src2, mode, nil, frac2, by)
	if best2.index == nil {
		return joinCost{fixcost: impossible}
	}
	varcost2 := Cost(frac * float64(nrows1*src2.lookupCost()))
	// trace.Println("joinopt", joinType, "frac", frac)
	// trace.Println("   ", nrows1, joinType, nrows2, "=> read2", read2, "=> frac2", frac2)
	// trace.Println("    best2", best2.index, "=", best2.fixcost, best2.varcost)
	// trace.Println("    nrows1", nrows1, "lookups", nrows1 * src2.lookupCost())
	return joinCost{index1: index, index2: best2.index, frac2: frac2,
		fixcost: fixcost1 + best2.fixcost,
		varcost: varcost1 + varcost2 + best2.varcost,
	}
}

func optOrdered(q Query, mode Mode, index []string, frac float64, fixed []Fixed) (Cost, Cost, []string) {
	//TODO singleton ?
	//TODO all cols fixed ?
	if len(index) > 0 && len(fixed) > 0 {
		best := bestOrdered(q, index, mode, frac, fixed)
		if best.index != nil {
			// fmt.Println("best", best.index, "index", index, "indexes", q.Indexes(),
			// 	"fixed", fixedStr(fixed))
			index = best.index
		}
	}
	fixcost, varcost := Optimize(q, mode, index, frac)
	return fixcost, varcost, index
}

func (jn *Join) setApproach(index []string, frac float64, approach any, tran QueryTran) {
	ap := approach.(*joinApproach)
	if ap.reverse {
		jn.source1, jn.source2 = jn.source2, jn.source1
		jn.joinType = jn.joinType.reverse()
	}
	jn.source1 = SetApproach(jn.source1, ap.index1, frac, tran)
	jn.source2 = SetApproach(jn.source2, ap.index2, ap.frac2, tran)
	jn.header = jn.getHeader()
}

func (jn *Join) getNrows() (int, int) {
	n1, p1 := jn.source1.Nrows()
	n2, p2 := jn.source2.Nrows()
	return jn.nrows(n1, p1, n2, p2), jn.pop(p1, p2)
}

func (jn *Join) nrows(n1, p1, n2, p2 int) int {
	switch jn.joinType {
	case one_one:
		return min(n1, n2)
	case n_one:
		n1, p1, n2, p2 = n2, p2, n1, p1
		fallthrough
	case one_n:
		p1 = max(1, p1) // avoid divide by zero
		p2 = max(1, p2)
		if n1 <= p1*n2/p2 { // rearranged n1/p1 <= n2/p2 (for integer math)
			return n1 * p2 / p1
		}
		return n2
	case n_n:
		return (n1 * n2) / 2 // estimate half
	default:
		panic(assert.ShouldNotReachHere())
	}
}

func (jn *Join) pop(p1, p2 int) int {
	switch jn.joinType {
	case one_one:
		return min(p1, p2)
	case n_one:
		return p1
	case one_n:
		return p2
	case n_n:
		return (p1 * p2) / 2 // estimate half
	default:
		panic(assert.ShouldNotReachHere())
	}
}

// execution --------------------------------------------------------

func (jb *joinBase) Rewind() {
	jb.source1.Rewind()
	jb.source2.Rewind()
	jb.rewind()
}
func (jb *joinBase) rewind() {
	jb.row1 = nil
	jb.row2 = nil
	jb.lookupRow = nil
}

func (jn *Join) Get(th *Thread, dir Dir) Row {
	defer func(t uint64) { jn.tget += tsc.Read() - t }(tsc.Read())
	for {
		if jn.row2 == nil && !jn.nextRow1(th, dir) {
			return nil
		}
		if jn.joinType == one_one || jn.joinType == n_one {
			jn.row2 = jn.lookupRow
			jn.lookupRow = nil
		} else {
			jn.row2 = jn.source2.Get(th, dir)
		}
		if jn.row2 != nil {
			// assert.That(jn.equalBy(th, jn.st, jn.row1, jn.row2))
			jn.ngets++
			return JoinRows(jn.row1, jn.row2)
		}
	}
}

func (jn *Join) nextRow1(th *Thread, dir Dir) bool {
	jn.row1 = jn.source1.Get(th, dir)
	if jn.row1 == nil {
		return false
	}
	// fmt.Println("Join row1", jn.row1)
	// assert.That(set.Disjoint(jn.by, jn.sel2cols))
	sel2cols := append(jn.sel2cols, jn.by...)
	sel2vals := append(jn.sel2vals, jn.projectRow1(th, jn.row1)...)
	if jn.joinType == one_one || jn.joinType == n_one {
		jn.lookupRow = jn.cachedLookup(th, sel2cols, sel2vals)
	} else {
		jn.source2.Select(sel2cols, sel2vals)
	}
	return true
}

func (jb *joinBase) cachedLookup(th *Thread, cols, vals []string) Row {
	if jb.st == nil || !jb.st.Updatable() {
		// Read-only transaction: use cache
		if jb.cache == nil {
			jb.cache = lrucache.New[lookupKey, Row](100)
		}
		key := lookupKey{cols: cols, vals: vals}
		joinCacheProbes.Add(1)
		return jb.cache.GetPut(key, func(k lookupKey) Row {
			joinCacheMisses.Add(1)
			return jb.source2.Lookup(th, k.cols, k.vals)
		})
	}
	// Updatable transaction: bypass cache
	return jb.source2.Lookup(th, cols, vals)
}

func (jb *joinBase) projectRow1(th *Thread, row Row) []string {
	key := make([]string, len(jb.by))
	for i, col := range jb.by {
		key[i] = row.GetRawVal(jb.source1.Header(), col, th, jb.st)
	}
	return key
}

func (jn *Join) Select(cols, vals []string) {
	// fmt.Println(jn.strategy(), "Select", cols, unpack(vals))
	jn.nsels++
	jn.rewind()
	jn.select1(cols, vals)
}

// select1 splits the select, calls source1 Select, and sets sel2cols/vals.
// It is used by Join, LeftJoin, and Times.
func (jl *joinLike) select1(cols, vals []string) {
	if cols == nil { // clear
		jl.source1.Select(nil, nil)
		jl.source2.Select(nil, nil)
		jl.sel2cols, jl.sel2vals = nil, nil
		return
	}
	sel1cols, sel1vals, sel2cols, sel2vals := jl.splitSelect(cols, vals)
	jl.source1.Select(sel1cols, sel1vals)
	jl.sel2cols, jl.sel2vals = sel2cols, sel2vals
}

func (jl *joinLike) splitSelect(cols, vals []string) (
	sel1cols, sel1vals, sel2cols, sel2vals []string) {
	columns1 := jl.source1.Columns()
	columns2 := jl.source2.Columns()
	for i, col := range cols {
		if slices.Contains(columns1, col) { // includes common
			sel1cols = append(sel1cols, col)
			sel1vals = append(sel1vals, vals[i])
		} else if slices.Contains(columns2, col) {
			sel2cols = append(sel2cols, col)
			sel2vals = append(sel2vals, vals[i])
		}
	}
	return
}

func (jn *Join) Lookup(th *Thread, cols, vals []string) Row {
	// fmt.Println(jn.strategy(), "Lookup", cols, unpack(vals))
	jn.nlooks++
	sel1cols, sel1vals, sel2cols, sel2vals := jn.splitSelect(cols, vals)
	if jn.lookupFallback(sel1cols) {
		// log.Println("INFO Join Lookup fallback to Select & Get")
		jn.rewind()
		jn.source1.Select(sel1cols, sel1vals)
		defer jn.Select(nil, nil)
		jn.sel2cols, jn.sel2vals = sel2cols, sel2vals
		x := jn.Get(th, Next)
		// if x != nil { // verify unique
		// 	assert.That(jn.Get(th, Next) == nil)
		// }
		return x
	}
	row1 := jn.source1.Lookup(th, sel1cols, sel1vals)
	if row1 == nil {
		return nil
	}
	var row2 Row
	sel2cols = append(sel2cols, jn.by...)
	sel2vals = append(sel2vals, jn.projectRow1(th, row1)...)
	if jn.joinType == one_one || jn.joinType == n_one {
		row2 = jn.cachedLookup(th, sel2cols, sel2vals)
	} else {
		jn.source2.Select(sel2cols, sel2vals)
		defer jn.Select(nil, nil)
		row2 = jn.source2.Get(th, Next)
	}
	if row2 == nil {
		return nil
	}
	// assert.That(jn.equalBy(th, jn.st, row1, row2))
	return JoinRows(row1, row2)
}

func (jb *joinBase) lookupFallback(sel1cols []string) bool {
	if jb.lookup == nil { // memoize
		jb.lookup = &lookupInfo{
			keys1:  jb.source1.Keys(),
			fixed1: jb.source1.Fixed(),
		}
	}
	if !hasKey(sel1cols, jb.lookup.keys1, jb.lookup.fixed1) {
		// can't do lookup on source1
		// this can happen (rarely) because there's no way to tell Optimize
		// that we want to do lookups with the index
		if !jb.lookup.fallback {
			jb.lookup.fallback = true
			// log.Println("INFO query", which, "Lookup fallback to Select & Get")
			// fmt.Println("sel1cols", sel1cols, "keys1", jb.lookup.keys1, "fixed1", jb.lookup.fixed1)
		}
		return true
	}
	return false
}

func (jn *Join) Simple(th *Thread) []Row {
	st := MakeSuTran(jn.qt)
	rows1 := jn.source1.Simple(th)
	rows2 := jn.source2.Simple(th)
	rows := make([]Row, 0, len(rows1))
	for _, row1 := range rows1 {
		for _, row2 := range rows2 {
			if jn.equalBy(th, st, row1, row2) {
				rows = append(rows, JoinRows(row1, row2))
			}
		}
	}
	return rows
}

func (jb *joinBase) equalBy(th *Thread, st *SuTran, row1, row2 Row) bool {
	for _, f := range jb.by {
		if row1.GetRawVal(jb.source1.Header(), f, th, st) !=
			row2.GetRawVal(jb.source2.Header(), f, th, st) {
			return false
		}
	}
	return true
}

// LeftJoin ---------------------------------------------------------

type LeftJoin struct {
	empty2 Row
	joinBase
}

func NewLeftJoin(src1, src2 Query, by []string, t QueryTran) *LeftJoin {
	return newLeftJoin(src1, src2, by, t, nil, nil)
}

func newLeftJoin(src1, src2 Query, by []string, t QueryTran,
	prevFixed1, prevFixed2 []Fixed) *LeftJoin {
	lj := &LeftJoin{joinBase: newJoinBase(src1, src2, by, t,
		prevFixed1, prevFixed2)}
	lj.keys = lj.getKeys()
	lj.indexes = lj.source1.Indexes()
	lj.fixed = lj.getFixed()
	lj.setNrows(lj.getNrows())
	lj.fast1.Set(src1.fastSingle() &&
		(lj.joinType == one_one || lj.joinType == n_one || src2.fastSingle()))
	return lj
}

func (lj *LeftJoin) With(src1, src2 Query) *LeftJoin {
	return newLeftJoin(src1, src2, lj.by, lj.qt, lj.prevFixed1, lj.prevFixed2)
}

func (lj *LeftJoin) String() string {
	return lj.joinBase.String("leftjoin")
}

func (lj *LeftJoin) getKeys() [][]string {
	// can't use source2.Keys() like Join.Keys()
	// because multiple right sides can be missing/blank
	switch lj.joinType {
	case one_one, n_one:
		return lj.source1.Keys()
	case one_n, n_n:
		return lj.keypairs()
	default:
		panic("unknown join type")
	}
}

func (lj *LeftJoin) getFixed() []Fixed {
	fixed1 := lj.source1.Fixed()
	fixed2 := lj.source2.Fixed()
	if len(fixed2) == 0 {
		return fixed1
	}
	result := make([]Fixed, 0, len(fixed1)+len(fixed2))
	// all of fixed1
	result = append(result, fixed1...)
	// fixed2 that are not common with source1
	for _, f2 := range fixed2 {
		if !slices.Contains(lj.by, f2.col) {
			// add "" because source2 row can be empty
			result = append(result, fixedWith(f2, ""))
		}
	}
	return result
}

func (lj *LeftJoin) Transform() Query {
	lj.optimized = true
	src1 := lj.source1.Transform()
	if _, ok := src1.(*Nothing); ok {
		return NewNothing(lj)
	}
	src2 := lj.source2.Transform()
	_, src2Nothing := src2.(*Nothing)
	fix1, fix2 := src1.Fixed(), src2.Fixed()
	if src2Nothing || fixedConflict(fix1, fix2) {
		// remove useless left join
		return keepCols(src1, src2, lj.Header())
	}
	if !equalFixed(fix1, lj.prevFixed1) || !equalFixed(fix2, lj.prevFixed2) {
		// for leftjoin, we can only apply fixed from src1 to src2
		// can't do the reverse because src2 is "optional"
		src2 = copyFixed(fix1, fix2, src2, lj.by, lj.qt)
		lj.prevFixed1, lj.prevFixed2 = fix1, fix2
	}
	if src1 != lj.source1 || src2 != lj.source2 {
		return lj.With(src1, src2).Transform()
	}
	return lj
}

func fixedConflict(fixed1, fixed2 []Fixed) bool {
	for _, f2 := range fixed2 {
		if src1vals := getFixed(fixed1, f2.col); src1vals != nil {
			// field is in both
			if set.Disjoint(src1vals, f2.values) {
				return true // can't match anything
			}
		}
	}
	return false
}

func (lj *LeftJoin) optimize(mode Mode, index []string, frac float64) (Cost, Cost, any) {
	jc := joinopt(lj.source1, lj.source2, lj.Nrows,
		mode, index, frac, lj.by, lj.fixed)
	return jc.fixcost, jc.varcost,
		&joinApproach{index1: jc.index1, index2: jc.index2, frac2: jc.frac2}
}

func (lj *LeftJoin) setApproach(index []string, frac float64, approach any, tran QueryTran) {
	ap := approach.(*joinApproach)
	lj.source1 = SetApproach(lj.source1, ap.index1, frac, tran)
	lj.source2 = SetApproach(lj.source2, ap.index2, ap.frac2, tran)
	lj.empty2 = make(Row, len(lj.source2.Header().Fields))
	lj.header = lj.getHeader()
}

func (lj *LeftJoin) getNrows() (int, int) {
	n1, p1 := lj.source1.Nrows()
	n2, p2 := lj.source2.Nrows()
	return lj.nrows(n1, p1, n2, p2), lj.pop(p1, p2)
}

func (lj *LeftJoin) nrows(n1, p1, n2, p2 int) int {
	switch lj.joinType {
	case one_one, n_one:
		return n1
	case one_n:
		p1 = max(1, p1) // avoid divide by zero
		p2 = max(1, p2)
		if n1 <= p1*n2/p2 { // rearranged n1/p1 <= n2/p2 (for integer math)
			return n1 * p2 / p1
		}
		return n2
	case n_n:
		return max(n1, (n1*n2)/2) // estimate half
	default:
		panic(assert.ShouldNotReachHere())
	}
}

func (lj *LeftJoin) pop(n1, n2 int) int {
	switch lj.joinType {
	case one_one, n_one:
		return n1
	case one_n:
		return n2
	case n_n:
		return max(n1, (n1*n2)/2) // estimate half
	default:
		panic(assert.ShouldNotReachHere())
	}
}

// execution

func (lj *LeftJoin) Get(th *Thread, dir Dir) (r Row) {
	defer func(t uint64) { lj.tget += tsc.Read() - t }(tsc.Read())
	row1out := true
	for {
		if lj.row2 == nil {
			lj.row1 = lj.source1.Get(th, dir)
			if lj.row1 == nil {
				return nil
			}
			lj.source2.Select(lj.by, lj.projectRow1(th, lj.row1))
			row1out = false
		}
		lj.row2 = lj.source2.Get(th, dir)
		// fmt.Println(lj.strategy(), "row2", lj.row2)
		if !row1out || lj.row2 != nil {
			row1out = true // regardless of filter
			row2 := lj.row2
			if row2 == nil {
				row2 = lj.empty2
			} else {
				// assert.That(lj.equalBy(th, lj.st, lj.row1, row2))
			}
			if lj.filter2(row2) {
				lj.ngets++
				return JoinRows(lj.row1, row2)
			}
		}
	}
}

func (lj *LeftJoin) filter2(row2 Row) bool {
	// fmt.Println(lj.strategy(), "filter", lj.sel2cols, unpack(lj.sel2vals))
	for i, col := range lj.sel2cols {
		x := row2.GetRaw(lj.source2.Header(), col)
		assert.That(len(x) == 0 || x[0] != PackForward)
		if x != lj.sel2vals[i] {
			return false
		}
	}
	return true
}

func (lj *LeftJoin) Select(cols, vals []string) {
	// fmt.Println(lj.strategy(), "Select", cols, unpack(vals))
	lj.nsels++
	lj.rewind()
	lj.select1(cols, vals)
}

func (lj *LeftJoin) Lookup(th *Thread, cols, vals []string) Row {
	lj.nlooks++
	defer lj.Select(nil, nil)
	sel1cols, sel1vals, sel2cols, sel2vals := lj.splitSelect(cols, vals)
	lj.sel2cols, lj.sel2vals = sel2cols, sel2vals
	if lj.lookupFallback(sel1cols) {
		// log.Println("INFO LeftJoin Lookup fallback to Select & Get")
		lj.rewind()
		lj.source1.Select(sel1cols, sel1vals)
		x := lj.Get(th, Next)
		// if x != nil { // verify unique
		// 	assert.That(lj.Get(th, Next) == nil)
		// }
		return x
	}
	row1 := lj.source1.Lookup(th, sel1cols, sel1vals)
	if row1 == nil {
		return nil
	}
	lj.source2.Select(lj.by, lj.projectRow1(th, row1))
	row2 := lj.source2.Get(th, Next)
	if row2 == nil {
		row2 = lj.empty2
	} else {
		// assert.That(lj.equalBy(th, lj.st, row1, row2))
	}
	if !lj.filter2(row2) {
		return nil
	}
	return JoinRows(row1, row2)
}

func (lj *LeftJoin) Simple(th *Thread) []Row {
	empty2 := make(Row, len(lj.source2.Header().Fields))
	rows1 := lj.source1.Simple(th)
	rows2 := lj.source2.Simple(th)
	rows := make([]Row, 0, len(rows1))
	for i1 := range len(rows1) {
		row1out := false
		for i2 := range len(rows2) {
			if lj.equalBy(th, lj.st, rows1[i1], rows2[i2]) {
				rows = append(rows, JoinRows(rows1[i1], rows2[i2]))
				row1out = true
			}
		}
		if !row1out {
			rows = append(rows, JoinRows(rows1[i1], empty2))
		}
	}
	return rows
}

// joinCache --------------------------------------------------------

type lookupKey struct {
	cols []string
	vals []string
}

func (lk lookupKey) Hash() uint64 {
	h := uint64(0)
	for _, col := range lk.cols {
		h = h*131 + hash.String(col)
	}
	for _, val := range lk.vals {
		h = h*131 + hash.String(val)
	}
	return h
}

func (lk lookupKey) Equal(other any) bool {
	if o, ok := other.(lookupKey); ok {
		return slices.Equal(lk.cols, o.cols) && slices.Equal(lk.vals, o.vals)
	}
	return false
}

var (
	joinCacheProbes atomic.Int64
	joinCacheMisses atomic.Int64
)

var _ = AddInfo("query.join.cacheProbes", &joinCacheProbes)
var _ = AddInfo("query.join.cacheMisses", &joinCacheMisses)
