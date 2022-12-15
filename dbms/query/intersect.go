// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/generic/ord"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"golang.org/x/exp/slices"
)

type Intersect struct {
	Compatible
	conflict bool
}

type intersectApproach struct {
	keyIndex []string
	reverse  bool
}

func NewIntersect(src, src2 Query) *Intersect {
	it := &Intersect{Compatible: Compatible{
		Query2: Query2{Query1: Query1{source: src}, source2: src2}}}
	it.init()
	return it
}

func (it *Intersect) String() string {
	return it.String2(it.stringOp())
}

func (it *Intersect) stringOp() string {
	return it.Compatible.stringOp("INTERSECT", "")
}

func (it *Intersect) Columns() []string {
	return set.Intersect(it.source.Columns(), it.source2.Columns())
}

func (it *Intersect) Keys() [][]string {
	k := set.IntersectFn(it.source.Keys(), it.source2.Keys(), set.Equal[string])
	if len(k) == 0 {
		k = [][]string{it.Columns()}
	}
	return k
}

func (it *Intersect) Fixed() []Fixed {
	fixed, none := FixedIntersect(it.source.Fixed(), it.source2.Fixed())
	if none {
		it.conflict = true
	}
	return fixed
}

func (it *Intersect) Indexes() [][]string {
	return set.UnionFn(it.source.Indexes(), it.source2.Indexes(), slices.Equal[string])
}

func (it *Intersect) Nrows() (int, int) {
	if it.disjoint != "" {
		return 0, 0
	}
	nrows1, pop1 := it.source.Nrows()
	nrows2, pop2 := it.source2.Nrows()
	maxNrows := ord.Min(nrows1, nrows2)
	maxPop := ord.Min(pop1, pop2)
	return maxNrows / 2, maxPop / 2 // estimate half
}

func (it *Intersect) rowSize() int {
	return (it.source.rowSize() + it.source2.rowSize()) / 2
}

func (it *Intersect) Transform() Query {
	if it.Fixed(); it.conflict {
		return NewNothing(it.Columns())
	}
	it.source = it.source.Transform()
	it.source2 = it.source2.Transform()
	// propagate Nothing
	if _, ok := it.source.(*Nothing); ok {
		return NewNothing(it.Columns())
	}
	if _, ok := it.source2.(*Nothing); ok {
		return NewNothing(it.Columns())
	}
	return it
}

func (it *Intersect) optimize(mode Mode, index []string) (Cost, any) {
	cost1, key1 := it.cost(it.source, it.source2, mode, index)
	cost2, key2 := it.cost(it.source2, it.source, mode, index) // reversed
	cost2 += outOfOrder
	if cost1 < cost2 {
		return cost1, &intersectApproach{keyIndex: key1}
	}
	return cost2, &intersectApproach{keyIndex: key2, reverse: true}
}

func (*Intersect) cost(source, source2 Query, mode Mode, index []string) (
	cost Cost, key []string) {
	key = bestKey(source2, mode)
	// iterate source and lookups on source2
	cost = Optimize(source, mode, index)
	nrows1, _ := source.Nrows()
	cost += LookupCost(source2, mode, key, nrows1)
	return cost, key
}

func (it *Intersect) setApproach(mode Mode, index []string, approach any,
	tran QueryTran) {
	ap := approach.(*intersectApproach)
	it.keyIndex = ap.keyIndex
	if ap.reverse {
		it.source, it.source2 = it.source2, it.source
	}
	it.source = SetApproach(it.source, mode, index, tran)
	it.source2 = SetApproach(it.source2, mode, it.keyIndex, tran)
}

func (it *Intersect) Header() *runtime.Header {
	hdr := it.source.Header()
	return runtime.NewHeader(hdr.Fields, it.Columns())
}

func (it *Intersect) Get(th *runtime.Thread, dir runtime.Dir) runtime.Row {
	if it.disjoint != "" {
		return nil
	}
	for {
		row := it.source.Get(th, dir)
		if row == nil || it.source2Has(th, row) {
			return row
		}
	}
}

func (it *Intersect) Select(cols, vals []string) {
	it.source.Select(cols, vals)
}

// COULD have a "merge" strategy (like Union)
