// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/setord"
	"github.com/apmckinlay/gsuneido/util/sset"
	"github.com/apmckinlay/gsuneido/util/strs"
)

type Times struct {
	Query2
	rewound bool
	row1    runtime.Row
}

func (t *Times) Init() {
	t.Query2.Init()
	t.rewound = true
	if !sset.Disjoint(t.source.Columns(), t.source2.Columns()) {
		panic("times: common columns not allowed: " + strs.Join(", ",
			sset.Intersect(t.source.Columns(), t.source2.Columns())))
	}
}

func (t *Times) String() string {
	return t.Query2.String2("TIMES")
}

func (t *Times) Columns() []string {
	return sset.Union(t.source.Columns(), t.source2.Columns())
}

func (t *Times) Keys() [][]string {
	// there are no columns in common so no keys in common
	// so there won't be any duplicates in the result
	return t.keypairs()
}

func (t *Times) Indexes() [][]string {
	return setord.Union(t.source.Indexes(), t.source2.Indexes())
}

func (t *Times) Transform() Query {
	t.source = t.source.Transform()
	t.source2 = t.source2.Transform()
	return t
}

func (t *Times) optimize(mode Mode, index []string) (Cost, interface{}) {
	cost := Optimize(t.source, mode, index) +
		t.source.Nrows()*Optimize(t.source2, mode, nil)
	costRev := Optimize(t.source2, mode, index) +
		t.source2.Nrows()*Optimize(t.source, mode, nil) + outOfOrder
	if cost < costRev {
		return cost, false
	}
	return costRev, true
}

func (t *Times) setApproach(index []string, approach interface{}, tran QueryTran) {
	if approach.(bool) {
		t.source, t.source2 = t.source2, t.source
	}
	t.source = SetApproach(t.source, index, tran)
	t.source2 = SetApproach(t.source2, nil, tran)
}

// execution --------------------------------------------------------

func (t *Times) Rewind() {
	t.rewound = true
	t.source.Rewind()
	t.source2.Rewind()
}

func (t *Times) Get(dir runtime.Dir) runtime.Row {
	row2 := t.source2.Get(dir)
	if t.rewound {
		t.rewound = false
		t.row1 = t.source.Get(dir)
		if t.row1 == nil || row2 == nil {
			return nil
		}
	}
	if row2 == nil {
		t.row1 = t.source.Get(dir)
		if t.row1 == nil {
			return nil
		}
		t.source2.Rewind()
		row2 = t.source2.Get(dir)
	}
	return runtime.JoinRows(t.row1, row2)
}

func (t *Times) Select(cols, vals []string) {
	t.source.Select(cols, vals)
	t.source2.Rewind()
	t.rewound = true
}
