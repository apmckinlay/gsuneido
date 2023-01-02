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

func NewMinus(src, src2 Query) *Minus {
	var m Minus
	m.source, m.source2 = src, src2
	m.init(m.calcFixed)
	return &m
}

func (m *Minus) String() string {
	return m.String2(m.stringOp())
}

func (m *Minus) stringOp() string {
	return m.Compatible.stringOp("MINUS", "")
}

func (m *Minus) Columns() []string {
	return m.source.Columns()
}

func (m *Minus) Keys() [][]string {
	return m.source.Keys()
}

func (m *Minus) Indexes() [][]string {
	return m.source.Indexes()
}

func (m *Minus) calcFixed(fixed1, fixed2 []Fixed) []Fixed {
	return fixed1
}

func (m *Minus) Nrows() (int, int) {
	n1, p1 := m.source.Nrows()
	n2, p2 := m.source2.Nrows()
	calc := func(n1, n2 int) int {
		min := ord.Max(0, n1-n2) // all common
		max := n1                // none common
		return (min + max) / 2
	}
	return calc(n1, n2), calc(p1, p2)
}

func (m *Minus) Transform() Query {
	m.source = m.source.Transform()
	// remove if disjoint
	if m.disjoint != "" {
		return m.source
	}
	m.source2 = m.source2.Transform()
	// propagate Nothing
	if _, ok := m.source.(*Nothing); ok {
		return NewNothing(m.Columns())
	}
	if _, ok := m.source2.(*Nothing); ok {
		return m.source
	}
	return m
}

func (m *Minus) optimize(mode Mode, index []string, frac float64) (Cost, Cost, any) {
	assert.That(m.disjoint == "") // eliminated by Transform
	// iterate source and lookup on source2
	fixcost, varcost := Optimize(m.source, mode, index, frac)
	nrows1, _ := m.source.Nrows()
	best2 := bestKey2(m.source2, mode, int(float64(nrows1)*frac))
	return fixcost + best2.fixcost, varcost + best2.varcost,
		&minusApproach{keyIndex: best2.index}
}

func (m *Minus) setApproach(index []string, frac float64, approach any, tran QueryTran) {
	ap := approach.(*minusApproach)
	m.keyIndex = ap.keyIndex
	m.source = SetApproach(m.source, index, frac, tran)
	m.source2 = SetApproach(m.source2, m.keyIndex, 0, tran)
}

func (m *Minus) Header() *Header {
	return m.source.Header()
}

func (m *Minus) Get(th *Thread, dir Dir) Row {
	for {
		row := m.source.Get(th, dir)
		if row == nil || !m.source2Has(th, row) {
			return row
		}
	}
}

func (m *Minus) Lookup(th *Thread, cols, vals []string) Row {
	row := m.source.Lookup(th, cols, vals)
	if row == nil || !m.source2Has(th, row) {
		return row
	}
	return nil
}

// COULD have a "merge" strategy (like Union)
