// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/tsc"
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
	t := &Times{joinLike: newJoinLike(src1, src2), rewound: true}
	t.keys = t.getKeys()
	t.indexes = t.getIndexes()
	t.fixed = t.getFixed()
	t.setNrows(t.getNrows())
	t.fast1.Set(src1.fastSingle() && src2.fastSingle())
	return t
}

func (t *Times) String() string {
	return "times"
}

func (t *Times) getKeys() [][]string {
	// no columns in common so no keys in common
	// so there won't be any duplicates in the result
	return t.keypairs()
}

func (t *Times) getIndexes() [][]string {
	// no columns in common so no indexes in common
	return slc.With(t.source1.Indexes(), t.source2.Indexes()...)
}

func (t *Times) getFixed() []Fixed {
	// no common columns so no overlap
	return slc.With(t.source1.Fixed(), t.source2.Fixed()...)
}

func (t *Times) Transform() Query {
	src1 := t.source1.Transform()
	if _, ok := src1.(*Nothing); ok {
		return NewNothing(t)
	}
	src2 := t.source2.Transform()
	if _, ok := src2.(*Nothing); ok {
		return NewNothing(t)
	}
	if src1 != t.source1 || src2 != t.source2 {
		return NewTimes(src1, src2)
	}
	return t
}

func (t *Times) optimize(mode Mode, index []string, frac float64) (Cost, Cost, any) {
	//TODO index could be split, first part on source1, second part on source2
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
	nrows1, _ := t.source1.Nrows()
	t.source2 = SetApproach(t.source2, nil, frac*float64(nrows1), tran)
	t.header = t.getHeader()
}

func (t *Times) getNrows() (int, int) {
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
	defer func(t0 uint64) { t.tget += tsc.Read() - t0 }(tsc.Read())
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
	t.select1(cols, vals)
	if len(t.sel2cols) > 0 {
		t.source2.Select(t.sel2cols, t.sel2vals)
	}
}

func (t *Times) Lookup(th *Thread, cols, vals []string) Row {
	// could use source1.Lookup like (Left)Join
	// but Times isn't used much
	t.select1(cols, vals)
	if len(t.sel2cols) > 0 {
		t.source2.Select(t.sel2cols, t.sel2vals)
	}
	return t.Get(th, Next)
}

func (t *Times) Simple(th *Thread) []Row {
	rows1 := t.source1.Simple(th)
	rows2 := t.source2.Simple(th)
	assert.Msg("Times Simple too big").That(len(rows1)*len(rows2) < maxSimple)
	rows := make([]Row, 0, len(rows1)*len(rows2))
	for _, row1 := range rows1 {
		for _, row2 := range rows2 {
			rows = append(rows, JoinRows(row1, row2))
		}
	}
	return rows
}
