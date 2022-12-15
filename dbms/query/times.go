// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/str"
	"golang.org/x/exp/slices"
)

type Times struct {
	Query2
	rewound bool
	row1    runtime.Row
}

func NewTimes(src, src2 Query) *Times {
	if !set.Disjoint(src.Columns(), src2.Columns()) {
		panic("times: common columns not allowed: " + str.Join(", ",
			set.Intersect(src.Columns(), src2.Columns())))
	}
	return &Times{Query2: Query2{Query1: Query1{source: src}, source2: src2},
		rewound: true}
}

func (t *Times) String() string {
	return t.Query2.String2("TIMES")
}

func (t *Times) stringOp() string {
	return "TIMES"
}

func (t *Times) Columns() []string {
	return set.Union(t.source.Columns(), t.source2.Columns())
}

func (t *Times) Keys() [][]string {
	// there are no columns in common so no keys in common
	// so there won't be any duplicates in the result
	return t.keypairs()
}

func (t *Times) Indexes() [][]string {
	return set.UnionFn(t.source.Indexes(), t.source2.Indexes(), slices.Equal[string])
}

func (t *Times) Transform() Query {
	t.source = t.source.Transform()
	t.source2 = t.source2.Transform()
	// propagate Nothing
	if _, ok := t.source.(*Nothing); ok {
		return NewNothing(t.Columns())
	}
	if _, ok := t.source2.(*Nothing); ok {
		return NewNothing(t.Columns())
	}
	return t
}

func (t *Times) optimize(mode Mode, index []string) (Cost, any) {
	nrows1, _ := t.source.Nrows()
	cost := Optimize(t.source, mode, index) +
		nrows1*Optimize(t.source2, mode, nil)
	nrows2, _ := t.source2.Nrows()
	costRev := Optimize(t.source2, mode, index) +
		nrows2*Optimize(t.source, mode, nil) + outOfOrder
	if cost < costRev {
		return cost, false
	}
	return costRev, true
}

func (t *Times) setApproach(mode Mode, index []string, approach any, tran QueryTran) {
	if approach.(bool) {
		t.source, t.source2 = t.source2, t.source
	}
	t.source = SetApproach(t.source, mode, index, tran)
	t.source2 = SetApproach(t.source2, mode, nil, tran)
}

func (t *Times) Nrows() (int, int) {
	n1, p1 := t.source.Nrows()
	n2, p2 := t.source2.Nrows()
	return n1 * n2, p1 * p2
}

// execution --------------------------------------------------------

func (t *Times) Rewind() {
	t.rewound = true
	t.source.Rewind()
	t.source2.Rewind()
}

func (t *Times) Get(th *runtime.Thread, dir runtime.Dir) runtime.Row {
	row2 := t.source2.Get(th, dir)
	if t.rewound {
		t.rewound = false
		t.row1 = t.source.Get(th, dir)
		if t.row1 == nil || row2 == nil {
			return nil
		}
	}
	if row2 == nil {
		t.row1 = t.source.Get(th, dir)
		if t.row1 == nil {
			return nil
		}
		t.source2.Rewind()
		row2 = t.source2.Get(th, dir)
	}
	return runtime.JoinRows(t.row1, row2)
}

func (t *Times) Select(cols, vals []string) {
	t.source.Select(cols, vals)
	t.source2.Rewind()
	t.rewound = true
}
