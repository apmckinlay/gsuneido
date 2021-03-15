// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import "github.com/apmckinlay/gsuneido/util/ints"

type Minus struct {
	Compatible
}

func (m *Minus) String() string {
	return m.String2("MINUS")
}

func (m *Minus) Keys() [][]string {
	return m.source.Keys()
}

func (m *Minus) Indexes() [][]string {
	return m.source.Indexes()
}

func (m *Minus) nrows() int {
	n1 := m.source.nrows()
	min := ints.Max(0, n1-m.source2.nrows())
	max := n1
	return (min + max) / 2
}

func (m *Minus) Transform() Query {
	if m.disjoint == "" {
		m.source = m.source.Transform()
		m.source2 = m.source2.Transform()
		return m
	}
	// remove if disjoint
	return m.source.Transform()
}

func (m *Minus) optimize(mode Mode, index []string, act action) Cost {
	m.keyIndex = bestKey(m.source2, mode)
	// iterate source and lookups on source2
	cost := Optimize(m.source, mode, index, act) +
		(m.source.nrows() * m.source2.lookupCost())
	if act == freeze {
		Optimize(m.source2, mode, m.keyIndex, freeze)
	}
	return cost
}
