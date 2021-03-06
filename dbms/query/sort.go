// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/sset"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Sort struct {
	Query1
	reverse bool
	columns []string
	frozen  bool
}

func (sort *Sort) Init() {
	sort.Query1.Init()
	if !sset.Subset(sort.source.Columns(), sort.columns) {
		panic("sort: nonexistent columns: " +
			str.Join(", ", sset.Difference(sort.columns, sort.source.Columns())))
	}
}

func (sort *Sort) String() string {
	s := sort.Query1.String()
	r := ""
	if sort.reverse {
		r = " reverse"
	}
	if sort.frozen {
		return s + r
	}
	return s + " SORT" + r + " " + str.Join(", ", sort.columns)
}

func (sort *Sort) Transform() Query {
	sort.source = sort.source.Transform()
	return sort
}

func (sort *Sort) optimize(mode Mode, index []string, act action) Cost {
	assert.That(index == nil)
	src := sort.source
	best := sort.bestPrefixed(src.Indexes(), sort.columns, mode)
	cost := Optimize(src, mode, sort.columns, false)
	if !act {
		return ints.Min(cost, best.cost)
	}
	sort.frozen = true
	if cost <= best.cost {
		return Optimize(src, mode, sort.columns, true)
	}
	// optimize1 to avoid tempindex
	return optimize1(src, mode, best.index, true)
}
