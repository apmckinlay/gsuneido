// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/ints"
)

type Minus struct {
	Compatible
}

type minusApproach struct {
	keyIndex []string
}

func NewMinus(src, src2 Query) *Minus {
	m := &Minus{Compatible: Compatible{
		Query2: Query2{Query1: Query1{source: src}, source2: src2}}}
	m.init()
	return m
}

func (m *Minus) String() string {
	return m.String2("MINUS", "")
}

func (m *Minus) Keys() [][]string {
	return m.source.Keys()
}

func (m *Minus) Indexes() [][]string {
	return m.source.Indexes()
}

func (m *Minus) Nrows() int {
	n1 := m.source.Nrows()
	min := ints.Max(0, n1-m.source2.Nrows())
	max := n1
	return (min + max) / 2
}

func (m *Minus) Transform() Query {
	// remove if disjoint
	if m.disjoint != "" {
		return m.source.Transform()
	}
	m.source = m.source.Transform()
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

func (m *Minus) optimize(mode Mode, index []string) (Cost, interface{}) {
	// iterate source and lookups on source2
	cost := Optimize(m.source, mode, index) +
		(m.source.Nrows() * m.source2.lookupCost())
	keyIndex := bestKey(m.source2, mode)
	if keyIndex == nil {
		return impossible, nil
	}
	approach := &minusApproach{keyIndex: keyIndex}
	return cost, approach
}

func (m *Minus) setApproach(index []string, approach interface{}, tran QueryTran) {
	ap := approach.(*minusApproach)
	m.keyIndex = ap.keyIndex
	m.source = SetApproach(m.source, index, tran)
	m.source2 = SetApproach(m.source2, m.keyIndex, tran)
}

func (m *Minus) Header() *runtime.Header {
	return m.source.Header()
}

func (m *Minus) Get(dir runtime.Dir) runtime.Row {
	if m.disjoint != "" {
		return m.source.Get(dir)
	}
	for {
		row := m.source.Get(dir)
		if row == nil || !m.source2Has(row) {
			return row
		}
	}
}

func (m *Minus) Select(cols, vals []string) {
	m.source.Select(cols, vals)
}

// COULD have a "merge" strategy (like Union)
