// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"slices"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/tsc"
)

type Intersect struct {
	Compatible1
	conflict bool
}

type intersectApproach struct {
	keyIndex []string
	reverse  bool
}

func NewIntersect(src1, src2 Query) *Intersect {
	it := Intersect{}
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
	return it.Compatible.String("intersect")
}

func (it *Intersect) getHeader() *Header {
	hdr := it.source1.Header()
	cols := set.Intersect(it.source1.Columns(), it.source2.Columns())
	return NewHeader(hdr.Fields, cols)
}

func (it *Intersect) getKeys() [][]string {
	k := set.IntersectFn(it.source1.Keys(), it.source2.Keys(), set.Equal[string])
	if len(k) == 0 {
		k = [][]string{it.Columns()}
	}
	return k
}

func (it *Intersect) getIndexes() [][]string {
	return set.UnionFn(it.source1.Indexes(), it.source2.Indexes(), slices.Equal)
}

func (it *Intersect) getFixed() []Fixed {
	// same as Join
	fixed, none := combineFixed(it.source1.Fixed(), it.source2.Fixed())
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
	if src1 != it.source1 || src2 != it.source2 {
		return NewIntersect(src1, src2)
	}
	return it
}

func (it *Intersect) optimize(mode Mode, index []string, frac float64) (Cost, Cost, any) {
	assert.That(it.disjoint == "") // eliminated by Transform
	fixcost1, varcost1, key1 := it.cost(it.source1, it.source2, mode, index, frac)
	fixcost2, varcost2, key2 := it.cost(it.source2, it.source1, mode, index, frac)
	fixcost2 += outOfOrder
	if fixcost1+varcost1 < fixcost2+varcost2 {
		return fixcost1, varcost1, &intersectApproach{keyIndex: key1}
	}
	return fixcost2, varcost2, &intersectApproach{keyIndex: key2, reverse: true}
}

func (*Intersect) cost(src1, src2 Query, mode Mode, index []string, frac float64) (
	Cost, Cost, []string) {
	// iterate source and lookup on source2
	fixcost1, varcost1 := Optimize(src1, mode, index, frac)
	nrows1, _ := src1.Nrows()
	best2 := bestLookupKey(src2, mode, int(float64(nrows1)*frac))
	return fixcost1 + best2.fixcost, varcost1 + best2.varcost, best2.index
}

func (it *Intersect) setApproach(index []string, frac float64, approach any,
	tran QueryTran) {
	ap := approach.(*intersectApproach)
	it.keyIndex = ap.keyIndex
	if ap.reverse {
		it.source1, it.source2 = it.source2, it.source1
	}
	it.source1 = SetApproach(it.source1, index, frac, tran)
	it.source2 = SetApproach(it.source2, it.keyIndex, 0, tran)
	it.header = it.getHeader()
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

func (it *Intersect) Lookup(th *Thread, cols, vals []string) Row {
	it.nlooks++
	row := it.source1.Lookup(th, cols, vals)
	if row == nil || it.source2Has(th, row) {
		return row
	}
	return nil
}

// COULD have a "merge" strategy (like Union)

func (it *Intersect) Simple(th *Thread) []Row {
	cols := it.Columns()
	rows1 := it.source1.Simple(th)
	rows2 := it.source2.Simple(th)
	dst := 0
	for _, row1 := range rows1 {
		for _, row2 := range rows2 {
			if EqualRows(it.source1.Header(), row1, it.source2.Header(), row2,
				cols, th, nil) {
				rows1[dst] = row1
				dst++
				break
			}
		}
	}
	return rows1[:dst]
}
