// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/sset"
	"github.com/apmckinlay/gsuneido/util/str"
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

func (sort *Sort) Init() {
	sort.Query1.Init()
	if !sset.Subset(sort.source.Columns(), sort.columns) {
		panic("sort: nonexistent columns: " +
			str.Join(", ", sset.Difference(sort.columns, sort.source.Columns())))
	}
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
	return s + " SORT" + r + " " + str.Join(", ", sort.columns)
}

func (sort *Sort) Transform() Query {
	sort.source = sort.source.Transform()
	return sort
}

func (sort *Sort) optimize(mode Mode, index []string) (Cost, interface{}) {
	assert.That(index == nil)
	src := sort.source
	cost := Optimize(src, mode, sort.columns)
	best := sort.bestPrefixed(src.Indexes(), sort.columns, mode)
	trace("SORT", "cost", cost, "best", best.cost)
	if cost <= best.cost {
		return cost, sortApproach{index: sort.columns}
	}
	return best.cost, sortApproach{index: best.index}
}

func (sort *Sort) setApproach(_ []string, approach interface{}, tran QueryTran) {
	sort.sortApproach = approach.(sortApproach)
	sort.source = SetApproach(sort.source, sort.index, tran)
}

// execution --------------------------------------------------------

func (sort *Sort) Header() *runtime.Header {
	return sort.source.Header()
}

func (sort *Sort) Get(dir runtime.Dir) runtime.Row {
	if sort.reverse {
		dir = dir.Reverse()
	}
	return sort.source.Get(dir)
}

func (sort *Sort) Select(cols, org []string) {
	sort.source.Select(cols, org)
}
