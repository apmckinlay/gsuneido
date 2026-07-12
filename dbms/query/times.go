// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/set"
	"github.com/apmckinlay/gsuneido/util/slc"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/tsc"
)

type Times struct {
	joinLike
	row1 Row
	state
}

func NewTimes(src1, src2 Query) *Times {
	if !set.Disjoint(src1.Columns(), src2.Columns()) {
		panic("times: common columns not allowed: " + str.Join(", ",
			set.Intersect(src1.Columns(), src2.Columns())))
	}
	t := &Times{joinLike: newJoinLike(src1, src2)}
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
	idx1 := t.source1.Indexes()
	idx2 := t.source2.Indexes()
	if isEmptyKey(idx1) {
		return idx2
	} else if isEmptyKey(idx2) {
		return idx1
	}
	// no columns in common so no indexes in common
	return slc.With(idx1, idx2...)
}

func (t *Times) getFixed() Fixed {
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

func (t *Times) optimize(mode Mode, req Require) (Cost, Cost, any) {
	opt := func(src1, src2 Query) (Cost, Cost) {
		req1, req2 := reqsForTimes(src1, src2, req)
		fixcost1, varcost1 := Optimize(src1, mode, req1)
		fixcost2, varcost2 := Optimize(src2, mode, req2)
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

func (t *Times) setApproach(req Require, approach any, tran QueryTran) {
	if approach.(bool) {
		t.source1, t.source2 = t.source2, t.source1
	}
	req1, req2 := reqsForTimes(t.source1, t.source2, req)
	t.source1 = SetApproach(t.source1, req1, tran)
	t.source2 = SetApproach(t.source2, req2, tran)
	t.header = t.getHeader()
}

func reqsForTimes(src1, src2 Query, req Require) (Require, Require) {
	nrows1, _ := src1.Nrows()
	nrows2, _ := src2.Nrows()
	var frac2 float32
	if nrows1 == 0 {
		frac2 = float32(1) / float32(max(1, nrows2))
	} else {
		frac2 = req.frac * float32(nrows1)
	}
	req1 := req
	if nrows2 == 0 {
		req1.frac = float32(1) / float32(max(1, nrows1))
	}
	return req1, NoneReq(frac2)
}

func (t *Times) getNrows() (int, int) {
	n1, p1 := t.source1.Nrows()
	n2, p2 := t.source2.Nrows()
	return n1 * n2, p1 * p2
}

// execution --------------------------------------------------------

func (t *Times) Rewind() {
	t.state = rewound
	t.source1.Rewind()
	t.source2.Rewind()
}

func (t *Times) Get(th *Thread, dir Dir) Row {
	defer func(t0 uint64) { t.tget += tsc.Read() - t0 }(tsc.Read())
	if t.state == eof {
		return nil
	}
	row2 := t.source2.Get(th, dir)
	if t.state == rewound {
		t.state = within
		t.row1 = t.source1.Get(th, dir)
		if t.row1 == nil || row2 == nil {
			t.state = eof
			return nil
		}
	}
	if row2 != nil && t.row1 == nil {
		t.row1 = t.source1.Get(th, dir)
		if t.row1 == nil {
			t.state = eof
			return nil
		}
	}
	if row2 == nil {
		t.row1 = t.source1.Get(th, dir)
		if t.row1 == nil {
			t.state = eof
			return nil
		}
		t.source2.Rewind()
		row2 = t.source2.Get(th, dir)
		if row2 == nil {
			t.state = eof
			return nil
		}
	}
	t.ngets++
	return JoinRows(t.row1, row2)
}

func (t *Times) Select(sels Sels) {
	t.nsels++
	t.Rewind()
	sel1, sel2 := t.splitSelect(sels)
	t.source1.Select(sel1)
	t.source2.Select(sel2)
}

func (t *Times) Lookup(th *Thread, sels Sels) Row {
	t.nlooks++
	return lookupViaSelectGet(t, th, sels)
}

func (t *Times) Simple(th *Thread) []Row {
	rows1 := t.source1.Simple(th)
	rows2 := t.source2.Simple(th)
	assert.That(len(rows1)*len(rows2) < maxSimple)
	rows := make([]Row, 0, len(rows1)*len(rows2))
	for _, row1 := range rows1 {
		for _, row2 := range rows2 {
			rows = append(rows, JoinRows(row1, row2))
		}
	}
	return rows
}
