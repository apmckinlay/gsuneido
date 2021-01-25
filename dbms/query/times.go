// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

type Times struct {
	Query2
}

func (t *Times) String() string {
	return t.Query2.String("times")
}

func (t *Times) Transform() Query {
	t.source = t.source.Transform()
	t.source2 = t.source2.Transform()
	return t
}
