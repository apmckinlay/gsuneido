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
	return t.Query2.String2("TIMES")
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

func (t *Times) optimize(mode Mode, index []string) (Cost, interface{}) {
	cost := Optimize(t.source, mode, index) +
		t.source.nrows()*Optimize(t.source2, mode, nil)
	costRev := Optimize(t.source2, mode, index) +
		t.source2.nrows()*Optimize(t.source, mode, nil) + outOfOrder
	if cost < costRev {
		return cost, false
	}
	return costRev, true
}

func (t *Times) setApproach(index []string, approach interface{}, tran QueryTran) {
	if approach.(bool) {
		t.source, t.source2 = t.source2, t.source
	}
	t.source = SetApproach(t.source, index, tran)
	t.source2 = SetApproach(t.source2, nil, tran)
}
