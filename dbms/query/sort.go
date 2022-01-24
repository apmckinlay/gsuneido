// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/sset"
	"github.com/apmckinlay/gsuneido/util/strs"
)

type Sort struct {
	Query1
	reverse bool
	columns []string
	sortApproach
}

type sortApproach struct {
	index []string
}

func NewSort(src Query, reverse bool, cols []string) *Sort {
	if !sset.Subset(src.Columns(), cols) {
		panic("sort: nonexistent columns: " +
			strs.Join(", ", sset.Difference(cols, src.Columns())))
	}
	return &Sort{Query1: Query1{source: src}, reverse: reverse, columns: cols}
}

func (sort *Sort) String() string {
	s := sort.source.String()
	r := ""
	if sort.reverse {
		r = " reverse"
	}
	if sort.index != nil {
		return s + r
	}
	return s + " SORT" + r + " " + strs.Join(", ", sort.columns)
}

func (sort *Sort) Transform() Query {
	sort.source = sort.source.Transform()
	// propagate Nothing
	if _, ok := sort.source.(*Nothing); ok {
		return NewNothing(sort.Columns())
	}
	return sort
}

func (sort *Sort) optimize(mode Mode, index []string) (Cost, interface{}) {
	assert.That(index == nil)
	src := sort.source
	cost := Optimize(src, mode, sort.columns) // adds temp index if needed
	best := sort.bestOrdered(src.Indexes(), sort.columns, mode)
	trace("SORT", "cost", cost, "best", best)
	if cost < best.cost {
		return cost, sortApproach{index: sort.columns}
	}
	return best.cost, sortApproach{index: best.index}
}

// bestOrdered returns the best index that supplies the required order
// taking fixed into consideration.
func (q1 *Query1) bestOrdered(indexes [][]string, order []string,
	mode Mode) bestIndex {
	best := newBestIndex()
	fixed := q1.source.Fixed()
	for _, ix := range indexes {
		if ordered(ix, order, fixed) {
			cost := Optimize(q1.source, mode, ix)
			best.update(ix, cost)
		}
	}
	return best
}

func (sort *Sort) setApproach(_ []string, approach interface{}, tran QueryTran) {
	sort.sortApproach = approach.(sortApproach)
	sort.source = SetApproach(sort.source, sort.index, tran)
}

// execution --------------------------------------------------------

// Only implements reverse.
// The actual sorting is done with a TempIndex

func (sort *Sort) Get(dir runtime.Dir) runtime.Row {
	if sort.reverse {
		dir = dir.Reverse()
	}
	return sort.source.Get(dir)
}

func (sort *Sort) Select(cols, vals []string) {
	sort.source.Select(cols, vals)
}

func (sort *Sort) Ordering() []string {
	return sort.columns
}
