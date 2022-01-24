// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/setord"
	"github.com/apmckinlay/gsuneido/util/setset"
	"github.com/apmckinlay/gsuneido/util/sset"
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
	return it.String2("INTERSECT", "")
}

func (it *Intersect) Columns() []string {
	return sset.Intersect(it.source.Columns(), it.source2.Columns())
}

func (it *Intersect) Keys() [][]string {
	k := setset.Intersect(it.source.Keys(), it.source2.Keys())
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
	return setord.Union(it.source.Indexes(), it.source2.Indexes())
}

func (it *Intersect) Nrows() int {
	if it.disjoint != "" {
		return 0
	}
	min := 0
	max := ints.Min(it.source.Nrows(), it.source2.Nrows())
	return (min + max) / 2 // estimate half way between
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

func (it *Intersect) optimize(mode Mode, index []string) (Cost, interface{}) {
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
	cost = Optimize(source, mode, index) +
		LookupCost(source2, mode, key, source.Nrows())
	return cost, key
}

func (it *Intersect) setApproach(index []string, approach interface{},
	tran QueryTran) {
	ap := approach.(*intersectApproach)
	it.keyIndex = ap.keyIndex
	if ap.reverse {
		it.source, it.source2 = it.source2, it.source
	}
	it.source = SetApproach(it.source, index, tran)
	it.source2 = SetApproach(it.source2, it.keyIndex, tran)
}

func (it *Intersect) Header() *runtime.Header {
	hdr := it.source.Header()
	return runtime.NewHeader(hdr.Fields, it.Columns())
}

func (it *Intersect) Get(dir runtime.Dir) runtime.Row {
	if it.disjoint != "" {
		return nil
	}
	for {
		row := it.source.Get(dir)
		if row == nil || it.source2Has(row) {
			return row
		}
	}
}

func (it *Intersect) Select(cols, vals []string) {
	it.source.Select(cols, vals)
}

// COULD have a "merge" strategy (like Union)
