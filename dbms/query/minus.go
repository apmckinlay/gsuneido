// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/set"
	"github.com/apmckinlay/gsuneido/util/tsc"
)

type Minus struct {
	Compatible1
	qt         QueryTran
	prevFixed1 Fixed
	prevFixed2 Fixed
}

type minusApproach struct {
	keyIndex   []string
	req1, req2 Require
}

func NewMinus(src1, src2 Query, t QueryTran) *Minus {
	return newMinus(src1, src2, t, nil, nil)
}

func newMinus(src1, src2 Query, t QueryTran, prevFixed1, prevFixed2 Fixed) *Minus {
	m := Minus{qt: t, prevFixed1: prevFixed1, prevFixed2: prevFixed2}
	m.Compatible = *newCompatible(src1, src2)
	m.header = src1.Header()
	m.keys = src1.Keys()
	m.indexes = src1.Indexes()
	m.fixed = src1.Fixed()
	m.setNrows(m.getNrows())
	m.rowSiz.Set(src1.rowSize())
	m.fast1.Set(src1.fastSingle())
	m.lookCost.Set(m.getLookupCost())
	return &m
}

func (m *Minus) String() string {
	return "minus"
}

func (m *Minus) getNrows() (int, int) {
	n1, p1 := m.source1.Nrows()
	n2, p2 := m.source2.Nrows()
	calc := func(n1, n2 int) int {
		min := max(0, n1-n2) // all common
		max := n1            // none common
		return (min + max) / 2
	}
	return calc(n1, n2), calc(p1, p2)
}

func (m *Minus) Transform() Query {
	src1 := m.source1.Transform()
	if m.disjoint != "" {
		return src1
	}
	if _, ok := src1.(*Nothing); ok {
		return NewNothing(m)
	}
	src2 := m.source2.Transform()
	if _, ok := src2.(*Nothing); ok {
		return src1
	}
	fix1, fix2 := src1.Fixed(), src2.Fixed()
	if !fix1.Equal(m.prevFixed1) || !fix2.Equal(m.prevFixed2) {
		src2 = compatCopyFixed(fix1, fix2, src2, m.qt)
		if src2 == nil {
			return src1
		}
		m.prevFixed1, m.prevFixed2 = fix1, fix2
	}
	if src1 != m.source1 || src2 != m.source2 {
		return newMinus(src1, src2, m.qt, m.prevFixed1, m.prevFixed2).Transform()
	}
	return m
}

func (m *Minus) optimize(mode Mode, req Require) (Cost, Cost, any) {
	assert.That(m.disjoint == "")
	fixcost1, varcost1 := Optimize(m.source1, mode, req)
	nrows1, _ := m.source1.Nrows()
	nlookups := req.LookupCount(nrows1)
	req2 := LookupReq(m.source2.Columns(), nlookups)
	fc2, vc2 := Optimize(m.source2, mode, req2)
	if fc2+vc2 >= impossible {
		return impossible, impossible, nil
	}
	return fixcost1 + fc2, varcost1 + vc2,
		&minusApproach{keyIndex: req2.cols, req1: req, req2: req2}
}

func (m *Minus) setApproach(req Require, approach any, tran QueryTran) {
	ap := approach.(*minusApproach)
	m.keyIndex = ap.keyIndex
	m.source1 = SetApproach(m.source1, ap.req1, tran)
	m.source2 = SetApproach(m.source2, ap.req2, tran)
	m.header = m.source1.Header()
	m.src1Only = set.Difference(m.source1.Columns(), m.source2.Columns())
}

func (m *Minus) Get(th *Thread, dir Dir) Row {
	defer func(t uint64) { m.tget += tsc.Read() - t }(tsc.Read())
	for {
		row := m.source1.Get(th, dir)
		if row == nil {
			return nil
		}
		if !m.source2Has(th, row) {
			m.ngets++
			return row
		}
	}
}

func (m *Minus) Lookup(th *Thread, sels Sels) Row {
	m.nlooks++
	row := m.source1.Lookup(th, sels)
	if row == nil || !m.source2Has(th, row) {
		return row
	}
	return nil
}

// COULD have a "merge" strategy (like Union)

func (m *Minus) Simple(th *Thread) []Row {
	rows1 := m.source1.Simple(th)
	rows2 := m.source2.Simple(th)
	dst := 0
outer:
	for _, row1 := range rows1 {
		for _, row2 := range rows2 {
			if m.equal(th, row1, row2) {
				continue outer
			}
		}
		rows1[dst] = row1
		dst++
	}
	return rows1[:dst]
}
