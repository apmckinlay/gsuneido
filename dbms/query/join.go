// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"log"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/trace"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/ord"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/str"
	"golang.org/x/exp/slices"
)

/*
	joinLike
		Times
		joinBase
			Join
			LeftJoin
*/

// joinLike is common stuff for Join, LeftJoin, and Times
type joinLike struct {
	saIndex  []string
	sel2cols []string
	sel2vals []string
	Query2
	conflict1 bool
	conflict2 bool
}

// joinBase is common stuff for Join and LeftJoin
type joinBase struct {
	st     *SuTran
	lookup *lookupInfo
	by     []string
	row1   Row
	row2   Row // nil when we need a new row1
	joinLike
	joinType
}

type Join struct {
	joinBase
}

type lookupInfo struct {
	keys1    [][]string
	fixed1   []Fixed
	fallback bool
}

type joinApproach struct {
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

func NewJoin(src1, src2 Query, by []string) *Join {
	jn := &Join{joinBase: newJoinBase(src1, src2, by)}
	jn.keys = jn.getKeys()
	jn.indexes = jn.getIndexes()
	jn.fixed = jn.getFixed()
	jn.setNrows(jn.getNrows())
	jn.fast1.Set(src1.fastSingle() && src2.fastSingle())
	return jn
}

func newJoinBase(src1, src2 Query, by []string) joinBase {
	b := set.Intersect(src1.Columns(), src2.Columns())
	if len(b) == 0 {
		panic("join: common columns required")
	}
	if by == nil {
		by = b
	} else if !set.Equal(by, b) {
		panic("join: by does not match common columns")
	}
	jb := joinBase{joinLike: newJoinLike(src1, src2)}
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

func newJoinLike(src1, src2 Query) joinLike {
	jl := joinLike{}
	jl.source1, jl.source2 = src1, src2
	jl.header = jl.getHeader()
	jl.rowSiz.Set(jl.source1.rowSize() + jl.source2.rowSize())
	jl.lookCost.Set(src1.lookupCost() * 2) // ???
	return jl
}

func (jn *Join) String() string {
	return parenQ2(jn.source1) + " " + jn.stringOp() + " " + paren(jn.source2)
}

func (jn *Join) stringOp() string {
	return "JOIN" + jn.bystr()
}

func (jb *joinBase) bystr() string {
	if len(jb.by) == 0 {
		return ""
	}
	return " " + str.Opt(jb.joinType.String(), " ") + "by" + str.Join("(,)", jb.by)
}

func (jb *joinBase) SetTran(t QueryTran) {
	jb.st = MakeSuTran(t)
}

func (jl *joinLike) getHeader() *Header {
	return JoinHeaders(jl.source1.Header(), jl.source2.Header())
}

func (jn *Join) getIndexes() [][]string {
	// can really only provide source.indexes() but optimize may swap.
	// optimize will return impossible for source2 indexes.
	return set.UnionFn(jn.source1.Indexes(), jn.source2.Indexes(), slices.Equal[string])
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

func (jn *Join) getFixed() []Fixed {
	fixed, none := combineFixed(jn.source1.Fixed(), jn.source2.Fixed())
	if none {
		jn.conflict1 = true
	}
	return fixed
}

func (jn *Join) Transform() Query {
	cols := jn.Columns()
	if jn.Fixed(); jn.conflict1 {
		return NewNothing(cols)
	}
	src1 := jn.source1.Transform()
	if _, ok := src1.(*Nothing); ok {
		return NewNothing(cols)
	}
	src2 := jn.source2.Transform()
	if _, ok := src2.(*Nothing); ok {
		return NewNothing(cols)
	}
	if src1 != jn.source1 || src2 != jn.source2 {
		return NewJoin(src1, src2, jn.by)
	}
	return jn
}

var joinRev = 0 // tests can set to impossible to prevent reverse

func (jn *Join) optimize(mode Mode, index []string, frac float64) (Cost, Cost, any) {
	fwd := joinopt(jn.source1, jn.source2, jn.joinType, jn.Nrows,
		mode, index, frac, jn.by)
	rev := joinopt(jn.source2, jn.source1, jn.joinType.reverse(), jn.Nrows,
		mode, index, frac, jn.by)
	rev.fixcost += outOfOrder + joinRev
	if trace.JoinOpt.On() {
		trace.JoinOpt.Println(mode, index, frac)
		trace.Println("    fwd", fwd.index, "=", fwd.fixcost, fwd.varcost)
		trace.Println("    rev", rev.index, "=", rev.fixcost, rev.varcost)
		trace.Println(format(jn, 1))
	}
	approach := &joinApproach{}
	if rev.fixcost+rev.varcost < fwd.fixcost+fwd.varcost {
		fwd = rev
		approach.reverse = true
	}
	if fwd.index == nil {
		return impossible, impossible, nil
	}
	approach.index2 = fwd.index
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

type bestJoin struct {
	bestIndex
	frac2 float64
}

func joinopt(src1, src2 Query, joinType joinType, nrows func() (int, int),
	mode Mode, index []string, frac float64, by []string) bestJoin {
	if slc.Empty(index) && !src1.fastSingle() {
		return bestJoin{bestIndex: newBestIndex()} // impossible
	}
	// always have to read all of source 1
	fixcost1, varcost1 := Optimize(src1, mode, index, frac)
	if fixcost1+varcost1 >= impossible {
		return bestJoin{bestIndex: newBestIndex()} // impossible
	}
	nrows1, _ := src1.Nrows()
	nrows2, _ := src2.Nrows()
	read2, _ := nrows()
	frac2 := float64(read2) * frac / float64(ord.Max(1, nrows2))
	best2 := bestGrouped(src2, mode, nil, frac2, by)
	if best2.index == nil {
		return bestJoin{bestIndex: newBestIndex()} // impossible
	}
	varcost2 := Cost(frac * float64(nrows1*src2.lookupCost()))
	// trace.Println("joinopt", joinType, "frac", frac)
	// trace.Println("   ", nrows1, joinType, nrows2, "=> read2", read2, "=> frac2", frac2)
	// trace.Println("    best2", best2.index, "=", best2.fixcost, best2.varcost)
	// trace.Println("    nrows1", nrows1, "lookups", nrows1 * src2.lookupCost())
	return bestJoin{frac2: frac2, bestIndex: bestIndex{
		index:   best2.index,
		fixcost: fixcost1 + best2.fixcost,
		varcost: varcost1 + varcost2 + best2.varcost,
	}}
}

func (jn *Join) setApproach(index []string, frac float64, approach any, tran QueryTran) {
	ap := approach.(*joinApproach)
	if ap.reverse {
		jn.source1, jn.source2 = jn.source2, jn.source1
		jn.joinType = jn.joinType.reverse()
	}
	jn.source1 = SetApproach(jn.source1, index, frac, tran)
	jn.source2 = SetApproach(jn.source2, ap.index2, ap.frac2, tran)
	jn.header = jn.getHeader()
	jn.saIndex = index
}

func (jn *Join) getNrows() (int, int) {
	n1, p1 := jn.source1.Nrows()
	n2, p2 := jn.source2.Nrows()
	return jn.nrows(n1, p1, n2, p2), jn.pop(p1, p2)
}

func (jn *Join) nrows(n1, p1, n2, p2 int) int {
	switch jn.joinType {
	case one_one:
		return ord.Min(n1, n2)
	case n_one:
		n1, p1, n2, p2 = n2, p2, n1, p1
		fallthrough
	case one_n:
		p1 = ord.Max(1, p1) // avoid divide by zero
		p2 = ord.Max(1, p2)
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
		return ord.Min(p1, p2)
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

// execution

func (jb *joinBase) Rewind() {
	jb.source1.Rewind()
	jb.source2.Rewind()
	jb.row1 = nil
	jb.row2 = nil
}

func (jn *Join) Get(th *Thread, dir Dir) Row {
	if jn.conflict1 || jn.conflict2 {
		return nil
	}
	for {
		if jn.row2 == nil && !jn.nextRow1(th, dir) {
			return nil
		}
		jn.row2 = jn.source2.Get(th, dir)
		if jn.row2 != nil {
			return JoinRows(jn.row1, jn.row2)
		}
	}
}

func (jn *Join) nextRow1(th *Thread, dir Dir) bool {
	jn.row1 = jn.source1.Get(th, dir)
	if jn.row1 == nil {
		return false
	}
	sel2cols := append(jn.sel2cols, jn.by...)
	sel2vals := append(jn.sel2vals, jn.projectRow(th, jn.row1)...)
	jn.source2.Select(sel2cols, sel2vals)
	return true
}

func (jb *joinBase) projectRow(th *Thread, row Row) []string {
	key := make([]string, len(jb.by))
	for i, col := range jb.by {
		key[i] = row.GetRawVal(jb.source1.Header(), col, th, jb.st)
	}
	return key
}

func (jn *Join) Select(cols, vals []string) {
	// fmt.Println(jn.stringOp(), "Select", cols, unpack(vals))
	jn.select2(cols, vals, jn.fastSingle())
}

func (jb *joinBase) select2(cols, vals []string, fastSingle bool) {
	jb.Rewind()
	if cols == nil { // clear
		jb.conflict1, jb.conflict2 = false, false
		jb.source1.Select(nil, nil)
		jb.sel2cols, jb.sel2vals = nil, nil
		return
	}
	if fastSingle {
		jb.sel2cols, jb.sel2vals = jb.selectByCols(cols, vals)
		return
	}
	sel1cols, sel1vals := jb.splitSelect(cols, vals)
	if jb.conflict1 { // not conflict2 because of LeftJoin
		return
	}
	jb.source1.Select(sel1cols, sel1vals)
}

func (jn *Join) Lookup(th *Thread, cols, vals []string) Row {
	defer jn.Select(nil, nil) // clear select
	if jn.fastSingle() {
		jn.sel2cols, jn.sel2vals = jn.selectByCols(cols, vals)
		return jn.Get(th, Next)
	}
	sel1cols, sel1vals := jn.splitSelect(cols, vals)
	if jn.conflict1 || jn.conflict2 {
		return nil
	}
	if jn.lookupFallback(sel1cols) {
		jn.Select(cols, vals)
		return jn.Get(th, Next)
	}
	jn.row1 = jn.source1.Lookup(th, sel1cols, sel1vals)
	if jn.row1 == nil {
		return nil
	}
	jn.source2.Select(jn.by, jn.projectRow(th, jn.row1))
	row2 := jn.source2.Get(th, Next)
	if row2 == nil {
		return nil
	}
	return JoinRows(jn.row1, row2)
}

func (jl *joinLike) splitSelect(cols, vals []string) (
	sel1cols, sel1vals []string) {
	fixed1 := jl.source1.Fixed()
	fixed2 := jl.source2.Fixed()
	jl.conflict1, jl.conflict2 = false, false
	for i, col := range cols {
		if slices.Contains(jl.saIndex, col) {
			sel1cols = append(sel1cols, col)
			sel1vals = append(sel1vals, vals[i])
			continue
		}
		fixVals1 := getFixed(fixed1, col)
		if len(fixVals1) == 1 && fixVals1[0] != vals[i] {
			jl.conflict1 = true
		}
		fixVals2 := getFixed(fixed2, col)
		if len(fixVals2) == 1 && fixVals2[0] != vals[i] {
			jl.conflict2 = true
		}
		// extra Select cols should be fixed
		assert.That(len(fixVals1) == 1 || len(fixVals2) == 1)
	}
	return
}

func (jb *joinBase) lookupFallback(sel1cols []string) bool {
	if jb.lookup == nil {
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
			log.Println("INFO query join lookup fallback to slower select")
		}
		return true
	}
	return false
}

// LeftJoin ---------------------------------------------------------

type LeftJoin struct {
	empty2 Row
	joinBase
	row1out bool
}

func NewLeftJoin(src1, src2 Query, by []string) *LeftJoin {
	lj := &LeftJoin{joinBase: newJoinBase(src1, src2, by)}
	lj.keys = lj.getKeys()
	lj.indexes = lj.source1.Indexes()
	lj.fixed = lj.getFixed()
	lj.setNrows(lj.getNrows())
	lj.fast1.Set(src1.fastSingle() &&
		(lj.joinType == one_one || lj.joinType == n_one || src2.fastSingle()))
	return lj
}

func (lj *LeftJoin) String() string {
	return parenQ2(lj.source1) + " " + lj.stringOp() + " " + paren(lj.source2)
}

func (lj *LeftJoin) stringOp() string {
	return "LEFTJOIN" + lj.bystr()
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
	src1 := lj.source1.Transform()
	if _, ok := src1.(*Nothing); ok {
		return NewNothing(lj.Columns())
	}
	src2 := lj.source2.Transform()
	_, src2Nothing := src2.(*Nothing)
	_, none := combineFixed(src1.Fixed(), src2.Fixed())
	if none || src2Nothing {
		// remove useless left join
		return keepCols(src1, src2, lj.Header())
	}
	if src1 != lj.source1 || src2 != lj.source2 {
		return NewLeftJoin(src1, src2, lj.by)
	}
	return lj
}

func (lj *LeftJoin) optimize(mode Mode, index []string, frac float64) (Cost, Cost, any) {
	best := joinopt(lj.source1, lj.source2, lj.joinType, lj.Nrows,
		mode, index, frac, lj.by)
	return best.fixcost, best.varcost,
		&joinApproach{index2: best.index, frac2: best.frac2}
}

func (lj *LeftJoin) setApproach(index []string, frac float64, approach any, tran QueryTran) {
	ap := approach.(*joinApproach)
	lj.source1 = SetApproach(lj.source1, index, frac, tran)
	lj.source2 = SetApproach(lj.source2, ap.index2, ap.frac2, tran)
	lj.empty2 = make(Row, len(lj.source2.Header().Fields))
	lj.header = lj.getHeader()
	lj.saIndex = index
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
		p1 = ord.Max(1, p1) // avoid divide by zero
		p2 = ord.Max(1, p2)
		if n1 <= p1*n2/p2 { // rearranged n1/p1 <= n2/p2 (for integer math)
			return n1 * p2 / p1
		}
		return n2
	case n_n:
		return ord.Max(n1, (n1*n2)/2) // estimate half
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
		return ord.Max(n1, (n1*n2)/2) // estimate half
	default:
		panic(assert.ShouldNotReachHere())
	}
}

// execution

func (lj *LeftJoin) Get(th *Thread, dir Dir) Row {
	if lj.conflict1 {
		return nil
	}
	for {
		if lj.row2 == nil && !lj.nextRow1(th, dir) {
			return nil
		}
		if lj.conflict2 {
			return lj.filter(lj.row1, lj.empty2)
		}
		lj.row2 = lj.source2.Get(th, dir)
		if lj.shouldOutput(lj.row2) {
			if lj.row2 == nil {
				return lj.filter(lj.row1, lj.empty2)
			}
			return lj.filter(lj.row1, lj.row2)
		}
	}
}

func (lj *LeftJoin) nextRow1(th *Thread, dir Dir) bool {
	lj.row1out = false
	lj.row1 = lj.source1.Get(th, dir)
	if lj.row1 == nil {
		return false
	}
	lj.source2.Select(lj.by, lj.projectRow(th, lj.row1))
	return true
}

func (lj *LeftJoin) shouldOutput(row Row) bool {
	if !lj.row1out {
		lj.row1out = true
		return true
	}
	return row != nil
}

func (lj *LeftJoin) filter(row1, row2 Row) Row {
	if lj.fastSingle() {
		for i, col := range lj.sel2cols {
			if row2.GetRaw(lj.source2.Header(), col) != lj.sel2vals[i] {
				return nil
			}
		}
	}
	return JoinRows(row1, row2)
}

func (lj *LeftJoin) Select(cols, vals []string) {
	// fmt.Println(lj.stringOp(), "Select", cols, unpack(vals))
	lj.select2(cols, vals, lj.fastSingle())
}

func (lj *LeftJoin) Lookup(th *Thread, cols, vals []string) Row {
	defer lj.Select(nil, nil) // clear select
	if lj.fastSingle() {
		lj.sel2cols, lj.sel2vals = lj.selectByCols(cols, vals)
		return lj.Get(th, Next)
	}
	sel1cols, sel1vals := lj.splitSelect(cols, vals)
	if lj.conflict1 {
		return nil
	}
	if lj.lookupFallback(sel1cols) {
		lj.Select(cols, vals)
		return lj.Get(th, Next)
	}
	row1 := lj.source1.Lookup(th, sel1cols, sel1vals)
	if row1 == nil {
		return nil
	}
	if lj.conflict2 {
		return JoinRows(row1, lj.empty2)
	}
	lj.source2.Select(lj.by, lj.projectRow(th, row1))
	row2 := lj.source2.Get(th, Next)
	if row2 == nil {
		return JoinRows(row1, lj.empty2)
	}
	return JoinRows(row1, row2)
}
