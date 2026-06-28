// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"slices"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/set"
	"github.com/apmckinlay/gsuneido/util/slc"
	"github.com/apmckinlay/gsuneido/util/tsc"
)

type Intersect struct {
	Compatible1
	conflict   bool
	qt         QueryTran
	prevFixed1 Fixed
	prevFixed2 Fixed
}

type intersectApproach struct {
	keyIndex   []string
	reverse    bool
	req1, req2 Require
}

func NewIntersect(src1, src2 Query, t QueryTran) *Intersect {
	return newIntersect(src1, src2, t, nil, nil)
}

func newIntersect(src1, src2 Query, t QueryTran, prevFixed1, prevFixed2 Fixed) *Intersect {
	it := Intersect{qt: t, prevFixed1: prevFixed1, prevFixed2: prevFixed2}
	it.Compatible = *newCompatible(src1, src2)
	it.header = it.getHeader()
	it.keys = it.getKeys()
	it.indexes = it.getIndexes()
	it.fixed = it.getFixed()
	it.setNrows(it.getNrows())
	it.rowSiz.Set(src1.rowSize())
	it.fast1.Set(src1.fastSingle() && src2.fastSingle())
	it.lookCost.Set(it.getLookupCost())
	return &it
}

func (it *Intersect) String() string {
	return "intersect"
}

func (it *Intersect) getHeader() *Header {
	hdr := it.source1.Header()
	cols := set.Intersect(it.source1.Columns(), it.source2.Columns())
	return NewHeader(hdr.Fields, cols)
}

func (it *Intersect) getKeys() [][]string {
	k := slc.With(it.source1.Keys(), it.source2.Keys()...)
	return minimizeKeys(k)
}

func (it *Intersect) getIndexes() [][]string {
	idx1 := it.source1.Indexes()
	idx2 := it.source2.Indexes()
	if isEmptyKey(idx1) {
		return idx2
	} else if isEmptyKey(idx2) {
		return idx1
	}
	return set.UnionFn(idx1, idx2, slices.Equal)
}

func (it *Intersect) getFixed() Fixed {
	// same as Join
	fixed, none := it.source1.Fixed().Combine(it.source2.Fixed())
	if none {
		it.conflict = true
	}
	return fixed
}

func (it *Intersect) getNrows() (int, int) {
	if it.disjoint != "" {
		return 0, 0
	}
	nrows1, pop1 := it.source1.Nrows()
	nrows2, pop2 := it.source2.Nrows()
	maxNrows := min(nrows1, nrows2)
	maxPop := min(pop1, pop2)
	return maxNrows / 2, maxPop / 2 // estimate half
}

func (it *Intersect) Transform() Query {
	if it.disjoint != "" {
		return NewNothing(it)
	}
	if it.Fixed(); it.conflict {
		return NewNothing(it)
	}
	src1 := it.source1.Transform()
	if _, ok := src1.(*Nothing); ok {
		return NewNothing(it)
	}
	src2 := it.source2.Transform()
	if _, ok := src2.(*Nothing); ok {
		return NewNothing(it)
	}
	fix1, fix2 := src1.Fixed(), src2.Fixed()
	if !fix1.Equal(it.prevFixed1) || !fix2.Equal(it.prevFixed2) {
		src1 = compatCopyFixed(fix2, fix1, src1, it.qt)
		if src1 == nil {
			return NewNothing(it)
		}
		src2 = compatCopyFixed(fix1, fix2, src2, it.qt)
		if src2 == nil {
			return NewNothing(it)
		}
		it.prevFixed1, it.prevFixed2 = fix1, fix2
	}
	if src1 != it.source1 || src2 != it.source2 {
		return newIntersect(src1, src2, it.qt, it.prevFixed1, it.prevFixed2).Transform()
	}
	return it
}

// compatCopyFixed wraps copyFixed (join) with a conflict check:
// if fromFixed has non-"" values for columns not in to (treated as ""),
// it can never match and nil is returned.
func compatCopyFixed(fromFixed, toFixed Fixed, to Query, t QueryTran) Query {
	cols := to.Columns()
	for _, f := range fromFixed {
		if !slices.Contains(cols, f.col) && !slices.Contains(f.values, "") {
			return nil
		}
	}
	return copyFixed(fromFixed, toFixed, to, cols, t)
}

func (it *Intersect) optimize(mode Mode, req Require) (Cost, Cost, any) {
	assert.That(it.disjoint == "") // eliminated by Transform
	fixcost1, varcost1, ap1 := it.cost2(mode, req, false)
	fixcost2, varcost2, ap2 := it.cost2(mode, req, true)
	fixcost2 += outOfOrder
	if fixcost1+varcost1 < fixcost2+varcost2 {
		return fixcost1, varcost1, ap1
	}
	return fixcost2, varcost2, ap2
}

func (it *Intersect) cost2(mode Mode, req Require, reverse bool) (Cost, Cost, *intersectApproach) {
	// iterate source and lookup on source2
	src1, src2 := it.source1, it.source2
	if reverse {
		src1, src2 = src2, src1
	}
	fixcost1, varcost1 := Optimize(src1, mode, req)
	nrows1, _ := src1.Nrows()
	nlookups := req.LookupCount(nrows1)
	req2 := LookupReq(src2.Columns(), nlookups)
	fc2, vc2 := Optimize(src2, mode, req2)
	if fc2+vc2 >= impossible {
		return impossible, impossible, nil
	}
	return fixcost1 + fc2, varcost1 + vc2,
		&intersectApproach{keyIndex: req2.cols, req1: req, req2: req2, reverse: reverse}
}

func (it *Intersect) setApproach(req Require, approach any, tran QueryTran) {
	ap := approach.(*intersectApproach)
	it.keyIndex = ap.keyIndex
	if ap.reverse {
		it.source1, it.source2 = it.source2, it.source1
	}
	it.source1 = SetApproach(it.source1, ap.req1, tran)
	it.source2 = SetApproach(it.source2, ap.req2, tran)
	it.header = it.getHeader()
	it.src1Only = set.Difference(it.source1.Columns(), it.source2.Columns())
}

func (it *Intersect) Get(th *Thread, dir Dir) Row {
	defer func(t uint64) { it.tget += tsc.Read() - t }(tsc.Read())
	for {
		row := it.source1.Get(th, dir)
		if row == nil {
			return nil
		}
		if it.source2Has(th, row) {
			it.ngets++
			return row
		}
	}
}

func (it *Intersect) Lookup(th *Thread, sels Sels) Row {
	it.nlooks++
	row := it.source1.Lookup(th, sels)
	if row == nil || it.source2Has(th, row) {
		return row
	}
	return nil
}

// COULD have a "merge" strategy (like Union)

func (it *Intersect) Simple(th *Thread) []Row {
	rows1 := it.source1.Simple(th)
	rows2 := it.source2.Simple(th)
	dst := 0
	for _, row1 := range rows1 {
		for _, row2 := range rows2 {
			if it.equal(th, row1, row2) {
				rows1[dst] = row1
				dst++
				break
			}
		}
	}
	return rows1[:dst]
}
