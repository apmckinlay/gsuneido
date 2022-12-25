// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/trace"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/ord"
	"github.com/apmckinlay/gsuneido/util/generic/set"
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

func NewJoin(src, src2 Query, by []string) *Join {
	b := set.Intersect(src.Columns(), src2.Columns())
	if len(b) == 0 {
		panic("join: common columns required")
	}
	if by == nil {
		by = b
	} else if !set.Equal(by, b) {
		panic("join: by does not match common columns")
	}
	jn := &Join{Query2: Query2{source: src, source2: src2}, by: by}
	k1 := containsKey(by, src.Keys())
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
	return parenQ2(jn.source) + " " + jn.stringOp() + " " + paren(jn.source2)
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
	return set.Union(jn.source.Columns(), jn.source2.Columns())
}

func (jn *Join) Indexes() [][]string {
	// can really only provide source.indexes() but optimize may swap.
	// optimize will return impossible for source2 indexes.
	return set.UnionFn(jn.source.Indexes(), jn.source2.Indexes(), slices.Equal[string])
}

func (jn *Join) Keys() [][]string {
	switch jn.joinType {
	case one_one:
		return set.UnionFn(jn.source.Keys(), jn.source2.Keys(), set.Equal[string])
	case one_n:
		return jn.source2.Keys()
	case n_one:
		return jn.source.Keys()
	case n_n:
		return jn.keypairs()
	default:
		panic("unknown join type")
	}
}

func (jn *Join) Fixed() []Fixed {
	fixed, none := combineFixed(jn.source.Fixed(), jn.source2.Fixed())
	if none {
		jn.conflict = true
	}
	return fixed
}

func (jn *Join) Transform() Query {
	if jn.Fixed(); jn.conflict {
		return NewNothing(jn.Columns())
	}
	jn.source = jn.source.Transform()
	jn.source2 = jn.source2.Transform()
	// propagate Nothing
	if _, ok := jn.source.(*Nothing); ok {
		return NewNothing(jn.Columns())
	}
	if _, ok := jn.source2.(*Nothing); ok {
		return NewNothing(jn.Columns())
	}
	return jn
}

func (jn *Join) optimize(mode Mode, index []string) (Cost, Cost, any) {
	fwd := joinopt(jn.source, jn.source2, jn.joinType, jn.Nrows,
		mode, index, jn.by)
	rev := joinopt(jn.source2, jn.source, jn.joinType.reverse(), jn.Nrows,
		mode, index, jn.by)
	rev.fixcost += outOfOrder
	if trace.JoinOpt.On() {
		trace.JoinOpt.Println(mode, index)
		trace.Println("    fwd", fwd)
		trace.Println("    rev", rev)
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

func joinopt(src1, src2 Query, joinType joinType, nrows func() (int, int),
	mode Mode, index, by []string) bestIndex {
	// always have to read all of source 1
	fixcost1, varcost1 := Optimize(src1, mode, index)
	if fixcost1+varcost1 >= impossible {
		return newBestIndex() // impossible
	}
	best2 := bestGrouped(src2, mode, nil, by)
	if best2.index == nil {
		return newBestIndex() // impossible
	}
	nrows1, _ := src1.Nrows()
	nrows2, _ := src2.Nrows()
	nrows2 = ord.Max(1, nrows2) // avoid division by zero
	nr, _ := nrows()
	// NOTE: lookupCost is not really correct because we use Select
	return bestIndex{
		index:   best2.index,
		fixcost: fixcost1 + best2.fixcost,
		varcost: varcost1 + (nrows1 * src2.lookupCost()) +
			(best2.varcost * nr / nrows2),
	}
}

func (jn *Join) setApproach(mode Mode, index []string, approach any, tran QueryTran) {
	ap := approach.(*joinApproach)
	if ap.reverse {
		jn.source, jn.source2 = jn.source2, jn.source
		jn.joinType = jn.joinType.reverse()
	}
	jn.source = SetApproach(jn.source, mode, index, tran)
	jn.source2 = SetApproach(jn.source2, mode, ap.index2, tran)
	jn.hdr1 = jn.source.Header()
}

func (jn *Join) Nrows() (int, int) {
	n1, p1 := jn.source.Nrows()
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
		p1 = ord.Max(1, p1) // avoid division by zero
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
	return jn.source.rowSize() + jn.source2.rowSize()
}

func (jn *Join) lookupCost() int {
	return jn.source.lookupCost() * 2 // ???
}

// execution

func (jn *Join) Rewind() {
	jn.source.Rewind()
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
	jn.row1 = jn.source.Get(th, dir)
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
	jn.source.Select(cols, vals)
	jn.row2 = nil
}

func (jn *Join) Lookup(th *Thread, cols, vals []string) Row {
	defer jn.Rewind()
	jn.row1 = jn.source.Lookup(th, cols, vals)
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

func NewLeftJoin(src, src2 Query, by []string) *LeftJoin {
	return &LeftJoin{Join: *NewJoin(src, src2, by)}
}

func (lj *LeftJoin) String() string {
	return parenQ2(lj.source) + " " + lj.stringOp() + " " + paren(lj.source2)
}

func (lj *LeftJoin) stringOp() string {
	return "LEFTJOIN" + lj.bystr()
}

func (lj *LeftJoin) Indexes() [][]string {
	return lj.source.Indexes()
}

func (lj *LeftJoin) Keys() [][]string {
	// can't use source2.Keys() like Join.Keys()
	// because multiple right sides can be missing/blank
	switch lj.joinType {
	case one_one, n_one:
		return lj.source.Keys()
	case one_n, n_n:
		return lj.keypairs()
	default:
		panic("unknown join type")
	}
}

func (lj *LeftJoin) Fixed() []Fixed {
	return lj.source.Fixed()
}

func (lj *LeftJoin) Transform() Query {
	if lj.Join.Fixed(); lj.conflict {
		return lj.source.Transform() // remove useless left join
	}
	lj.source = lj.source.Transform()
	lj.source2 = lj.source2.Transform()
	// propagate Nothing
	if _, ok := lj.source.(*Nothing); ok {
		return NewNothing(lj.Columns())
	}
	if _, ok := lj.source2.(*Nothing); ok {
		return lj.source
	}
	return lj
}

func (lj *LeftJoin) optimize(mode Mode, index []string) (Cost, Cost, any) {
	best := joinopt(lj.source, lj.source2, lj.joinType, lj.Nrows,
		mode, index, lj.by)
	return best.fixcost, best.varcost, &joinApproach{index2: best.index}
}

func (lj *LeftJoin) setApproach(mode Mode, index []string, approach any, tran QueryTran) {
	ap := approach.(*joinApproach)
	lj.source = SetApproach(lj.source, mode, index, tran)
	lj.source2 = SetApproach(lj.source2, mode, ap.index2, tran)
	lj.empty2 = make(Row, len(lj.source2.Header().Fields))
	lj.hdr1 = lj.source.Header()
}

func (lj *LeftJoin) Nrows() (int, int) {
	n1, p1 := lj.source.Nrows()
	n2, p2 := lj.source2.Nrows()
	return lj.nrows(n1, p1, n2, p2), lj.pop(p1, p2)
}

func (lj *LeftJoin) nrows(n1, p1, n2, p2 int) int {
	switch lj.joinType {
	case one_one, n_one:
		return n1
	case one_n:
		p1 = ord.Max(1, p1) // avoid division by zero
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
	lj.row1 = lj.source.Lookup(th, cols, vals)
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
