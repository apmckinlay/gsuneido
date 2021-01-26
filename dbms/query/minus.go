// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

type Minus struct {
	Compatible
}

func (m *Minus) String() string {
	return m.Query2.String("minus")
}

func (m *Minus) Transform() Query {
	m.source = m.source.Transform()
	m.source2 = m.source2.Transform()
	return m
}
