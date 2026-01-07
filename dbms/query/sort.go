// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"slices"
	"strings"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/tsc"
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

func NewSort(src Query, reverse bool, order []string) *Sort {
	if !set.Subset(src.Columns(), order) {
		panic("sort: nonexistent columns: " +
			str.Join(", ", set.Difference(order, src.Columns())))
	}
	sort := Sort{reverse: reverse}
	sort.source = src
	sort.order = order
	sort.header = src.Header()
	sort.keys = src.Keys()
	sort.fixed = src.Fixed()
	sort.setNrows(src.Nrows())
	sort.rowSiz.Set(src.rowSize())
	sort.fast1.Set(src.fastSingle())
	sort.singleTbl.Set(src.SingleTable())
	return &sort
}

func (sort *Sort) String() string {
	r := ""
	if sort.reverse {
		r = "reverse"
	}
	if sort.index != nil { // optimized
		return r
	}
	return "sort " + str.Opt(r, " ") + str.Join(", ", sort.order)
}

func (sort *Sort) Order() []string {
	return sort.order
}

func (*Sort) Indexes() [][]string {
	panic(assert.ShouldNotReachHere())
}

func (sort *Sort) knowExactNrows() bool {
	return sort.source.knowExactNrows()
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
	best := bestOrdered(src, sort.order, mode, frac, sort.fixed)
	if fixcost+varcost < best.fixcost+best.varcost {
		return fixcost, varcost, sortApproach{index: sort.order}
	}
	return best.fixcost, best.varcost, sortApproach{index: best.index}
}

// bestOrdered returns the best index that supplies the required order
// taking fixed into consideration.
func bestOrdered(q Query, order []string, mode Mode, frac float64, fixed []Fixed) bestIndex {
	best := newBestIndex()
	for _, ix := range q.Indexes() {
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

func (sort *Sort) Get(th *Thread, dir Dir) Row {
	defer func(t uint64) { sort.tget += tsc.Read() - t }(tsc.Read())
	if sort.reverse {
		dir = dir.Reverse()
	}
	row := sort.source.Get(th, dir)
	if row != nil {
		sort.ngets++
	}
	return row
}

func (sort *Sort) Select(cols, vals []string) {
	sort.nsels++
	sort.source.Select(cols, vals)
}

func (sort *Sort) Simple(th *Thread) []Row {
	rev := 1
	if sort.reverse {
		rev = -1
	}
	rows := sort.source.Simple(th)
	cmp := func(xrow, yrow Row) int {
		for _, col := range sort.order {
			x := xrow.GetRawVal(sort.header, col, th, nil)
			y := yrow.GetRawVal(sort.header, col, th, nil)
			if c := strings.Compare(x, y); c != 0 {
				return c * rev
			}
		}
		return 0
	}
	slices.SortFunc(rows, cmp)
	return rows
}
