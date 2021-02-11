// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

type Minus struct {
	Compatible
}

func (m *Minus) String() string {
	return m.Query2.String("MINUS")
}

func (m *Minus) Keys() [][]string {
	return m.source.Keys()
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
