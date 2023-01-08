// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/trace"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/ord"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/str"
	"golang.org/x/exp/slices"
)

type Join struct {
	Query2
	by []string
	joinType
	conflict bool
	hdr1     *Header
	row1     Row
	row2     Row // nil when we need a new row1
	st       *SuTran
}

type joinApproach struct {
	reverse bool
	index2  []string
	frac2   float64
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
	b := set.Intersect(src1.Columns(), src2.Columns())
	if len(b) == 0 {
		panic("join: common columns required")
	}
	if by == nil {
		by = b
	} else if !set.Equal(by, b) {
		panic("join: by does not match common columns")
	}
	jn := &Join{Query2: Query2{source1: src1, source2: src2}, by: by}
	k1 := containsKey(by, src1.Keys())
	k2 := containsKey(by, src2.Keys())
	if k1 && k2 {
		jn.joinType = one_one
	} else if k1 {
		jn.joinType = one_n
	} else if k2 {
		jn.joinType = n_one
	} else {
		jn.joinType = n_n
	}
	return jn
}

func (jn *Join) String() string {
	return parenQ2(jn.source1) + " " + jn.stringOp() + " " + paren(jn.source2)
}

func (jn *Join) stringOp() string {
	return "JOIN" + jn.bystr()
}

func (jn *Join) bystr() string {
	if len(jn.by) == 0 {
		return ""
	}
	return " " + str.Opt(jn.joinType.String(), " ") + "by" + str.Join("(,)", jn.by)
}

func (jn *Join) SetTran(t QueryTran) {
	jn.st = MakeSuTran(t)
}

func (jn *Join) Columns() []string {
	return set.Union(jn.source1.Columns(), jn.source2.Columns())
}

func (jn *Join) Indexes() [][]string {
	// can really only provide source.indexes() but optimize may swap.
	// optimize will return impossible for source2 indexes.
	return set.UnionFn(jn.source1.Indexes(), jn.source2.Indexes(), slices.Equal[string])
}

func (jn *Join) Keys() [][]string {
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

func (*Join) fastSingle() bool {
	return false
}

func (jn *Join) Fixed() []Fixed {
	fixed, none := combineFixed(jn.source1.Fixed(), jn.source2.Fixed())
	if none {
		jn.conflict = true
	}
	return fixed
}

func (jn *Join) Transform() Query {
	if jn.Fixed(); jn.conflict {
		return NewNothing(jn.Columns())
	}
	jn.source1 = jn.source1.Transform()
	jn.source2 = jn.source2.Transform()
	// propagate Nothing
	if _, ok := jn.source1.(*Nothing); ok {
		return NewNothing(jn.Columns())
	}
	if _, ok := jn.source2.(*Nothing); ok {
		return NewNothing(jn.Columns())
	}
	return jn
}

func (jn *Join) optimize(mode Mode, index []string, frac float64) (Cost, Cost, any) {
	fwd := joinopt(jn.source1, jn.source2, jn.joinType, jn.Nrows,
		mode, index, frac, jn.by)
	rev := joinopt(jn.source2, jn.source1, jn.joinType.reverse(), jn.Nrows,
		mode, index, frac, jn.by)
	rev.fixcost += outOfOrder
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
	jn.hdr1 = jn.source1.Header()
}

func (jn *Join) Nrows() (int, int) {
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

func (jn *Join) rowSize() int {
	return jn.source1.rowSize() + jn.source2.rowSize()
}

func (jn *Join) lookupCost() int {
	return jn.source1.lookupCost() * 2 // ???
}

// execution

func (jn *Join) Rewind() {
	jn.source1.Rewind()
	jn.row1 = nil
	jn.row2 = nil
}

func (jn *Join) Get(th *Thread, dir Dir) Row {
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
	jn.source2.Select(jn.by, jn.projectRow(th, jn.row1))
	return true
}

func (jn *Join) projectRow(th *Thread, row Row) []string {
	key := make([]string, len(jn.by))
	for i, col := range jn.by {
		key[i] = row.GetRawVal(jn.hdr1, col, th, jn.st)
	}
	return key
}

func (jn *Join) Select(cols, vals []string) {
	jn.source1.Select(cols, vals)
	jn.row2 = nil
}

func (jn *Join) Lookup(th *Thread, cols, vals []string) Row {
	defer jn.Rewind()
	jn.row1 = jn.source1.Lookup(th, cols, vals)
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

// LeftJoin ---------------------------------------------------------

type LeftJoin struct {
	Join
	row1out bool
	empty2  Row
}

func NewLeftJoin(src1, src2 Query, by []string) *LeftJoin {
	return &LeftJoin{Join: *NewJoin(src1, src2, by)}
}

func (lj *LeftJoin) String() string {
	return parenQ2(lj.source1) + " " + lj.stringOp() + " " + paren(lj.source2)
}

func (lj *LeftJoin) stringOp() string {
	return "LEFTJOIN" + lj.bystr()
}

func (lj *LeftJoin) Indexes() [][]string {
	return lj.source1.Indexes()
}

func (lj *LeftJoin) Keys() [][]string {
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

func (lj *LeftJoin) Fixed() []Fixed {
	return lj.source1.Fixed()
}

func (lj *LeftJoin) Transform() Query {
	if lj.Join.Fixed(); lj.conflict {
		return lj.source1.Transform() // remove useless left join
	}
	lj.source1 = lj.source1.Transform()
	lj.source2 = lj.source2.Transform()
	// propagate Nothing
	if _, ok := lj.source1.(*Nothing); ok {
		return NewNothing(lj.Columns())
	}
	if _, ok := lj.source2.(*Nothing); ok {
		return lj.source1
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
	lj.hdr1 = lj.source1.Header()
}

func (lj *LeftJoin) Nrows() (int, int) {
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
	for {
		if lj.row2 == nil && !lj.nextRow1(th, dir) {
			return nil
		}
		lj.row2 = lj.source2.Get(th, dir)
		if lj.shouldOutput(lj.row2) {
			if lj.row2 == nil {
				return JoinRows(lj.row1, lj.empty2)
			}
			return JoinRows(lj.row1, lj.row2)
		}
	}
}

func (lj *LeftJoin) nextRow1(th *Thread, dir Dir) bool {
	lj.row1out = false
	return lj.Join.nextRow1(th, dir)
}

func (lj *LeftJoin) shouldOutput(row Row) bool {
	if !lj.row1out {
		lj.row1out = true
		return true
	}
	return row != nil
}

func (lj *LeftJoin) Lookup(th *Thread, cols, vals []string) Row {
	defer lj.Rewind()
	lj.row1 = lj.source1.Lookup(th, cols, vals)
	if lj.row1 == nil {
		return nil
	}
	lj.row1out = false
	lj.source2.Select(lj.by, lj.projectRow(th, lj.row1))
	row2 := lj.source2.Get(th, Next)
	if lj.shouldOutput(row2) {
		if row2 == nil {
			return JoinRows(lj.row1, lj.empty2)
		}
		return JoinRows(lj.row1, row2)
	}
	return nil
}
