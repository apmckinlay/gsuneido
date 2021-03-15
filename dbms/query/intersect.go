// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/sset"
	"github.com/apmckinlay/gsuneido/util/ssset"
)

type Intersect struct {
	Compatible
}

func (it *Intersect) String() string {
	return it.String2("INTERSECT")
}

func (it *Intersect) Columns() []string {
	return sset.Intersect(it.source.Columns(), it.source2.Columns())
}

func (it *Intersect) Keys() [][]string {
	k := ssset.Intersect(it.source.Keys(), it.source2.Keys())
	if len(k) == 0 {
		k = [][]string{it.Columns()}
	}
	return k
}

func (it *Intersect) Indexes() [][]string {
	return ssset.Union(it.source.Indexes(), it.source2.Indexes())
}

func (it *Intersect) nrows() int {
	if it.disjoint != "" {
		return 0
	}
	min := 0
	max := ints.Min(it.source.nrows(), it.source2.nrows())
	return (min + max) / 2  // guess half way between
}

func (it *Intersect) rowSize() int {
	return (it.source.rowSize() + it.source2.rowSize()) / 2
}

func (it *Intersect) Transform() Query {
	it.source = it.source.Transform()
	it.source2 = it.source2.Transform()
	return it
}

func (it *Intersect) optimize(mode Mode, index []string, act action) Cost {
	it.keyIndex = it.source2.Keys()[0]
	cost1, key1 := it.cost(it.source, it.source2, mode, index)
	cost2, key2 := it.cost(it.source2, it.source, mode, index) // reversed
	cost2 += outOfOrder
	cost := ints.Min(cost1, cost2)
	if act == freeze {
		if cost2 < cost1 {
			it.source, it.source2, key1 = it.source2, it.source, key2
		}
		it.keyIndex = key1
		Optimize(it.source, mode, index, freeze)
		Optimize(it.source2, mode, it.keyIndex, freeze)
	}
	return cost
}

func (*Intersect) cost(source, source2 Query, mode Mode, index []string) (
	cost Cost, key []string) {
	key = bestKey(source2, mode)
	// iterate source and lookups on source2
	cost = Optimize(source, mode, index, assess) +
		(source.nrows() * source2.lookupCost())
	return
}
