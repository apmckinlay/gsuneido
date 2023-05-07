// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/ord"
)

type Minus struct {
	Compatible1
}

type minusApproach struct {
	keyIndex []string
}

func NewMinus(src1, src2 Query) *Minus {
	var m Minus
	m.Compatible = *newCompatible(src1, src2)
	m.fixed = m.fixed1
	m.header = m.source1.Header()
	return &m
}

func (m *Minus) String() string {
	return m.String2(m.stringOp())
}

func (m *Minus) stringOp() string {
	return m.Compatible.stringOp("MINUS", "")
}

func (m *Minus) Keys() [][]string {
	return m.source1.Keys()
}

func (m *Minus) fastSingle() bool {
	return m.source1.fastSingle()
}

func (m *Minus) Indexes() [][]string {
	return m.source1.Indexes()
}

func (m *Minus) Nrows() (int, int) {
	n1, p1 := m.source1.Nrows()
	n2, p2 := m.source2.Nrows()
	calc := func(n1, n2 int) int {
		min := ord.Max(0, n1-n2) // all common
		max := n1                // none common
		return (min + max) / 2
	}
	return calc(n1, n2), calc(p1, p2)
}

func (m *Minus) Transform() Query {
	if m.disjoint != "" {
		return m.source1
	}
	src1 := m.source1.Transform()
	if _, ok := src1.(*Nothing); ok {
		return NewNothing(m.Columns())
	}
	src2 := m.source2.Transform()
	if _, ok := src2.(*Nothing); ok {
		return src1
	}
	if src1 != m.source1 || src2 != m.source2 {
		return NewMinus(src1, src2)
	}
	return m
}

func (m *Minus) optimize(mode Mode, index []string, frac float64) (Cost, Cost, any) {
	assert.That(m.disjoint == "") // eliminated by Transform
	// iterate source and lookup on source2
	fixcost, varcost := Optimize(m.source1, mode, index, frac)
	nrows1, _ := m.source1.Nrows()
	best2 := bestKey2(m.source2, mode, int(float64(nrows1)*frac))
	return fixcost + best2.fixcost, varcost + best2.varcost,
		&minusApproach{keyIndex: best2.index}
}

func (m *Minus) setApproach(index []string, frac float64, approach any, tran QueryTran) {
	ap := approach.(*minusApproach)
	m.keyIndex = ap.keyIndex
	m.source1 = SetApproach(m.source1, index, frac, tran)
	m.source2 = SetApproach(m.source2, m.keyIndex, 0, tran)
	m.header = m.source1.Header()
}

func (m *Minus) Get(th *Thread, dir Dir) Row {
	for {
		row := m.source1.Get(th, dir)
		if row == nil || !m.source2Has(th, row) {
			return row
		}
	}
}

func (m *Minus) Lookup(th *Thread, cols, vals []string) Row {
	row := m.source1.Lookup(th, cols, vals)
	if row == nil || !m.source2Has(th, row) {
		return row
	}
	return nil
}

// COULD have a "merge" strategy (like Union)
