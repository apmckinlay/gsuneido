// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	. "github.com/apmckinlay/gsuneido/runtime"
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
}

type joinApproach struct {
	reverse bool
	index2  []string
}

type joinType int

const (
	one_one joinType = iota + 1
	one_n
	n_one
	n_n
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
	jn := &Join{Query2: Query2{Query1: Query1{source: src}, source2: src2}, by: by}
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
	return jn.string("JOIN")
}

func (jn *Join) string(op string) string {
	by := ""
	if len(jn.by) > 0 {
		by = "by" + str.Join("(,)", jn.by) + " "
	}
	return parenQ2(jn.source) + " " + op + " " +
		str.Opt(jn.joinType.String(), " ") + by + paren(jn.source2)
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

func (jn *Join) optimize(mode Mode, index []string) (Cost, any) {
	defer be(gin("Join", jn, index))
	fwd := jn.opt(jn.source, jn.source2, jn.joinType, mode, index)
	rev := jn.opt(jn.source2, jn.source, jn.joinType.reverse(), mode, index)
	rev.cost += outOfOrder
	trace("forward", fwd, "reverse", rev)
	approach := &joinApproach{}
	if rev.cost < fwd.cost {
		fwd = rev
		approach.reverse = true
	}
	if fwd.index == nil {
		return impossible, nil
	}
	approach.index2 = fwd.index
	return fwd.cost, approach
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

func (jn *Join) opt(src1, src2 Query, joinType joinType,
	mode Mode, index []string) bestIndex {
	trace("OPT", paren(src1), "JOIN", joinType, paren(src2))
	// always have to read all of source 1
	cost1 := Optimize(src1, mode, index)
	if cost1 >= impossible {
		return newBestIndex()
	}
	best := bestGrouped(src2, mode, nil, jn.by)
	if best.index == nil {
		return best
	}
	nrows1 := src1.Nrows()
	// should only be taking a portion of the variable cost2,
	// not the fixed temp index cost2 (so 2/3 instead of 1/2)
	cost := cost1 + (nrows1 * src2.lookupCost()) + (best.cost * 2 / 3)
	trace("join opt", cost1, "+ (", nrows1, "*", src2.lookupCost(), ") + (",
		best.cost, "* 2/3 ) =", cost)
	best.cost = cost
	return best
}

func (jn *Join) setApproach(index []string, approach any, tran QueryTran) {
	ap := approach.(*joinApproach)
	if ap.reverse {
		jn.source, jn.source2 = jn.source2, jn.source
		jn.joinType = jn.joinType.reverse()
	}
	jn.source = SetApproach(jn.source, index, tran)
	jn.source2 = SetApproach(jn.source2, ap.index2, tran)
	jn.hdr1 = jn.source.Header()
}

func (jn *Join) Nrows() int {
	// n_one and one_n assume records will have matching counterparts
	nrows1 := jn.source.Nrows()
	nrows2 := jn.source2.Nrows()
	var nrows int
	switch jn.joinType {
	case one_one:
		nrows = ord.Min(nrows1, nrows2)
	case n_one:
		nrows = nrows1
	case one_n:
		nrows = nrows2
	case n_n:
		nrows = nrows1 * nrows2
	default:
		panic("shouldn't reach here")
	}
	return nrows / 2 // actual will be between 0 and nrows so estimate halfway
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
	jn.source2.Select(jn.by, jn.projectRow(jn.row1))
	return true
}

func (jn *Join) projectRow(row Row) []string {
	key := make([]string, len(jn.by))
	for i, col := range jn.by {
		key[i] = row.GetRaw(jn.hdr1, col)
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
	jn.source2.Select(jn.by, jn.projectRow(jn.row1))
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
	return lj.string("LEFTJOIN")
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

func (lj *LeftJoin) optimize(mode Mode, index []string) (Cost, any) {
	best := lj.opt(lj.source, lj.source2, lj.joinType, mode, index)
	return best.cost, &joinApproach{index2: best.index}
}

func (lj *LeftJoin) setApproach(index []string, approach any, tran QueryTran) {
	ap := approach.(*joinApproach)
	lj.source = SetApproach(lj.source, index, tran)
	lj.source2 = SetApproach(lj.source2, ap.index2, tran)
	lj.empty2 = make(Row, len(lj.source2.Header().Fields))
	lj.hdr1 = lj.source.Header()
}

func (lj *LeftJoin) Nrows() int {
	return lj.source.Nrows()
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
	lj.source2.Select(lj.by, lj.projectRow(lj.row1))
	row2 := lj.source2.Get(th, Next)
	if lj.shouldOutput(row2) {
		if row2 == nil {
			return JoinRows(lj.row1, lj.empty2)
		}
		return JoinRows(lj.row1, row2)
	}
	return nil
}
