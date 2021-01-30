// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import "github.com/apmckinlay/gsuneido/util/sset"

type Times struct {
	Query2
}

func (t *Times) String() string {
	return t.Query2.String("TIMES")
}

func (t *Times) Columns() []string {
	return sset.Union(t.source.Columns(), t.source2.Columns())
}

func (t *Times) Transform() Query {
	t.source = t.source.Transform()
	t.source2 = t.source2.Transform()
	return t
}
