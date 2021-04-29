// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/sset"
	"github.com/apmckinlay/gsuneido/util/ssset"
)

type Intersect struct {
	Compatible
}

type intersectApproach struct {
	keyIndex []string
	reverse  bool
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
	return (min + max) / 2 // estimate half way between
}

func (it *Intersect) rowSize() int {
	return (it.source.rowSize() + it.source2.rowSize()) / 2
}

func (it *Intersect) Transform() Query {
	it.source = it.source.Transform()
	it.source2 = it.source2.Transform()
	return it
}

func (it *Intersect) optimize(mode Mode, index []string) (Cost, interface{}) {
	it.keyIndex = it.source2.Keys()[0]
	cost1, key1 := it.cost(it.source, it.source2, mode, index)
	cost2, key2 := it.cost(it.source2, it.source, mode, index) // reversed
	cost2 += outOfOrder
	cost := ints.Min(cost1, cost2)
	app := intersectApproach{keyIndex: key1}
	if cost2 < cost1 {
		app = intersectApproach{keyIndex: key2, reverse: true}
	}
	return cost, &app
}

func (*Intersect) cost(source, source2 Query, mode Mode, index []string) (
	cost Cost, key []string) {
	key = bestKey(source2, mode)
	// iterate source and lookups on source2
	cost = Optimize(source, mode, index) +
		LookupCost(source2, mode, key, source.nrows())
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

func (it *Intersect) Select(cols, org []string) {
	it.source.Select(cols, org)
}

// COULD have a "merge" strategy (like Union)
