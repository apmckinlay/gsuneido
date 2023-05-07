// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Sort struct {
	order []string
	sortApproach
	Query1
	reverse bool
}

type sortApproach struct {
	index []string
}

func NewSort(src Query, reverse bool, cols []string) *Sort {
	if !set.Subset(src.Columns(), cols) {
		panic("sort: nonexistent columns: " +
			str.Join(", ", set.Difference(cols, src.Columns())))
	}
	sort := Sort{reverse: reverse, order: cols}
	sort.source = src
	sort.header = src.Header()
	return &sort
}

func (sort *Sort) String() string {
	return sort.source.String() + str.Opt(" ", sort.stringOp())
}

func (sort *Sort) stringOp() string {
	r := ""
	if sort.reverse {
		r = "reverse"
	}
	if sort.index != nil {
		return r
	}
	return "SORT " + str.Opt(r, " ") + str.Join(", ", sort.order)
}

func (sort *Sort) Transform() Query {
	src := sort.source.Transform()
	if _, ok := src.(*Nothing); ok {
		return src
	}
	if src != sort.source {
		return NewSort(src, sort.reverse, sort.order)
	}
	return sort
}

func (sort *Sort) optimize(mode Mode, index []string, frac float64) (Cost, Cost, any) {
	assert.That(index == nil)
	src := sort.source
	fixcost, varcost := Optimize(src, mode, sort.order, frac) // adds temp index if needed
	best := bestOrdered(src, src.Indexes(), sort.order, mode, frac)
	if fixcost+varcost < best.fixcost+best.varcost {
		return fixcost, varcost, sortApproach{index: sort.order}
	}
	return best.fixcost, best.varcost, sortApproach{index: best.index}
}

// bestOrdered returns the best index that supplies the required order
// taking fixed into consideration.
func bestOrdered(q Query, indexes [][]string, order []string,
	mode Mode, frac float64) bestIndex {
	best := newBestIndex()
	fixed := q.Fixed()
	for _, ix := range indexes {
		if ordered(ix, order, fixed) {
			fixcost, varcost := Optimize(q, mode, ix, frac)
			best.update(ix, fixcost, varcost)
		}
	}
	return best
}

func (sort *Sort) setApproach(_ []string, frac float64, approach any, tran QueryTran) {
	sort.sortApproach = approach.(sortApproach)
	sort.source = SetApproach(sort.source, sort.index, frac, tran)
	sort.header = sort.source.Header()
}

// execution --------------------------------------------------------

// Only implements reverse.
// The actual sorting is done with a TempIndex

func (sort *Sort) Get(th *runtime.Thread, dir runtime.Dir) runtime.Row {
	if sort.reverse {
		dir = dir.Reverse()
	}
	return sort.source.Get(th, dir)
}

func (sort *Sort) Select(cols, vals []string) {
	sort.source.Select(cols, vals)
}

func (sort *Sort) Ordering() []string {
	return sort.order
}
