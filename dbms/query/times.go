// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/util/sset"
	"github.com/apmckinlay/gsuneido/util/ssset"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Times struct {
	Query2
}

func (t *Times) Init() {
	t.Query2.Init()
	if !sset.Disjoint(t.source.Columns(), t.source2.Columns()) {
		panic("times: common columns not allowed: " + str.Join(", ",
			sset.Intersect(t.source.Columns(), t.source2.Columns())))
	}
}

func (t *Times) String() string {
	return t.Query2.String("TIMES")
}

func (t *Times) Columns() []string {
	return sset.Union(t.source.Columns(), t.source2.Columns())
}

func (t *Times) Keys() [][]string {
	// there are no columns in common so no keys in common
	// so there won't be any duplicates in the result
	return t.keypairs()
}

func (t *Times) Indexes() [][]string {
	return ssset.Union(t.source.Indexes(), t.source2.Indexes())
}

func (t *Times) Transform() Query {
	t.source = t.source.Transform()
	t.source2 = t.source2.Transform()
	return t
}
