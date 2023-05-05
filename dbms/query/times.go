// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Times struct {
	row1 Row
	joinLike
	rewound bool
}

func NewTimes(src1, src2 Query) *Times {
	if !set.Disjoint(src1.Columns(), src2.Columns()) {
		panic("times: common columns not allowed: " + str.Join(", ",
			set.Intersect(src1.Columns(), src2.Columns())))
	}
	t := Times{rewound: true}
	t.source1, t.source2 = src1, src2
	return &t
}

func (t *Times) String() string {
	return t.Query2.String2("TIMES")
}

func (t *Times) stringOp() string {
	return "TIMES"
}

func (t *Times) Columns() []string {
	return set.Union(t.source1.Columns(), t.source2.Columns())
}

func (t *Times) Keys() [][]string {
	// no columns in common so no keys in common
	// so there won't be any duplicates in the result
	return t.keypairs()
}

func (t *Times) Indexes() [][]string {
	// no columns in common so no indexes in common
	return slc.With(t.source1.Indexes(), t.source2.Indexes()...)
}

func (t *Times) Fixed() []Fixed {
	t.ensureFixed()
	fixed, conflict := combineFixed(t.fixed1, t.fixed2)
	assert.That(!conflict) // because no common columns
	return fixed
}

func (t *Times) rowSize() int {
	return t.source1.rowSize() + t.source2.rowSize()
}

func (t *Times) Transform() Query {
	src1 := t.source1.Transform()
	if _, ok := src1.(*Nothing); ok {
		return NewNothing(t.Columns())
	}
	src2 := t.source2.Transform()
	if _, ok := src2.(*Nothing); ok {
		return NewNothing(t.Columns())
	}
	if src1 != t.source1 || src2 != t.source2 {
		return NewTimes(src1, src2)
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
	fixFwd, varFwd := opt(t.source1, t.source2)
	fixRev, varRev := opt(t.source2, t.source1)
	fixRev += outOfOrder
	if fixFwd+varFwd < fixRev+varRev {
		return fixFwd, varFwd, false
	}
	return fixRev, varRev, true
}

func (t *Times) setApproach(index []string, frac float64, approach any, tran QueryTran) {
	if approach.(bool) {
		t.source1, t.source2 = t.source2, t.source1
	}
	t.source1 = SetApproach(t.source1, index, frac, tran)
	t.saIndex = index
	nrows1, _ := t.source1.Nrows()
	t.source2 = SetApproach(t.source2, nil, frac*float64(nrows1), tran)
}

func (t *Times) Nrows() (int, int) {
	n1, p1 := t.source1.Nrows()
	n2, p2 := t.source2.Nrows()
	return n1 * n2, p1 * p2
}

// execution --------------------------------------------------------

func (t *Times) Rewind() {
	t.rewound = true
	t.source1.Rewind()
	t.source2.Rewind()
}

func (t *Times) Get(th *Thread, dir Dir) Row {
	if t.conflict1 || t.conflict2 {
		return nil
	}
	row2 := t.source2.Get(th, dir)
	if t.rewound {
		t.rewound = false
		t.row1 = t.source1.Get(th, dir)
		if t.row1 == nil || row2 == nil {
			return nil
		}
	}
	if row2 == nil {
		t.row1 = t.source1.Get(th, dir)
		if t.row1 == nil {
			return nil
		}
		t.source2.Rewind()
		row2 = t.source2.Get(th, dir)
	}
	return JoinRows(t.row1, row2)
}

func (t *Times) Select(cols, vals []string) {
	t.Rewind()
	if cols == nil { // clear
		t.conflict1, t.conflict2 = false, false
		t.source1.Select(nil, nil)
		t.source2.Select(nil, nil)
		return
	}
	if t.fastSingle() {
		t.source2.Select(t.selectByCols(cols, vals))
		return
	}
	sel1cols, sel1vals := t.splitSelect(cols, vals)
	if t.conflict1 || t.conflict2 {
		return
	}
	t.source1.Select(sel1cols, sel1vals)
}

func (t *Times) Lookup(th *Thread, cols, vals []string) Row {
	if t.fastSingle() {
		t.source2.Select(t.selectByCols(cols, vals))
	} else {
		sel1cols, sel1vals := t.splitSelect(cols, vals)
		if t.conflict1 || t.conflict2 {
			return nil
		}
		t.source1.Select(sel1cols, sel1vals)
	}
	row := t.Get(th, Next)
	t.Select(nil, nil) // clear select
	return row
}

func (t *Times) lookupCost() int {
	return t.source1.lookupCost() * 2 // ???
}

func (t *Times) fastSingle() bool {
	return t.source1.fastSingle() && t.source2.fastSingle()
}
