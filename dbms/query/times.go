// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Times struct {
	Query2
	rewound bool
	row1    Row
}

func NewTimes(src, src2 Query) *Times {
	if !set.Disjoint(src.Columns(), src2.Columns()) {
		panic("times: common columns not allowed: " + str.Join(", ",
			set.Intersect(src.Columns(), src2.Columns())))
	}
	return &Times{Query2: Query2{source: src, source2: src2}, rewound: true}
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
	// no columns in common so no keys in common
	// so there won't be any duplicates in the result
	return t.keypairs()
}

func (t *Times) Indexes() [][]string {
	// no columns in common so no indexes in common
	return slc.With(t.source.Indexes(), t.source2.Indexes()...)
}

func (t *Times) Fixed() []Fixed {
	fixed, _ := combineFixed(t.source.Fixed(), t.source2.Fixed())
	return fixed
}

func (t *Times) rowSize() int {
	return t.source.rowSize() + t.source2.rowSize()
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

func (t *Times) optimize(mode Mode, index []string, frac float64) (Cost, Cost, any) {
	opt := func(src1, src2 Query) (Cost, Cost) {
		nrows1, _ := src1.Nrows()
		fixcost1, varcost1 := Optimize(src1, mode, index, frac)
		fixcost2, varcost2 := Optimize(src2, mode, nil, frac*float64(nrows1))
		return fixcost1 + fixcost2, varcost1 + varcost2
	}
	fixFwd, varFwd := opt(t.source, t.source2)
	fixRev, varRev := opt(t.source2, t.source)
	fixRev += outOfOrder
	if fixFwd+varFwd < fixRev+varRev {
		return fixFwd, varFwd, false
	}
	return fixRev, varRev, true
}

func (t *Times) setApproach(index []string, frac float64, approach any, tran QueryTran) {
	if approach.(bool) {
		t.source, t.source2 = t.source2, t.source
	}
	t.source = SetApproach(t.source, index, frac, tran)
	nrows1, _ := t.source.Nrows()
	t.source2 = SetApproach(t.source2, nil, frac*float64(nrows1), tran)
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

func (t *Times) Get(th *Thread, dir Dir) Row {
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
	return JoinRows(t.row1, row2)
}

func (t *Times) Select(cols, vals []string) {
	t.source.Select(cols, vals)
	t.source2.Rewind()
	t.rewound = true
}

func (t *Times) Lookup(th *Thread, cols, vals []string) Row {
	t.Select(cols, vals)
	row := t.Get(th, Next)
	t.Select(nil, nil) // clear select
	return row
}

func (t *Times) lookupCost() int {
	return t.source.lookupCost() * 2 // ???
}
