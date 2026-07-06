// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// View is a pass-through wrapper around a Query
// it is optionally added by parsequery
// so that format doesn't expand views.
// It should not be used when doing Transform and Optimize
// because it blocks identification of child queries.
type View struct {
	Query1
	name string
}

var _ Query = (*View)(nil)

func NewView(name string, src Query) *View {
	v := &View{name: name}
	v.source = src
	v.header = src.Header()
	v.keys = src.Keys()
	v.fixed = src.Fixed()
	v.setNrows(src.Nrows())
	v.rowSiz.Set(src.rowSize())
	v.fast1.Set(src.fastSingle())
	v.singleTbl.Set(src.SingleTable())
	v.lookCost.Set(src.lookupCost())
	return v
}

func (v *View) knowExactNrows() bool {
	return v.source.knowExactNrows()
}

func (v *View) Transform() Query {
	panic(assert.ShouldNotReachHere())
}

func (v *View) optimize(mode Mode, req Require) (
	Cost, Cost, any) {
	panic(assert.ShouldNotReachHere())
}

func (*Query1) setApproach(Require, any, QueryTran) {
	panic(assert.ShouldNotReachHere())
}

func (v *View) Get(th *Thread, dir Dir) Row {
	panic(assert.ShouldNotReachHere())
}

func (v *View) Select(Sels) {
	panic(assert.ShouldNotReachHere())
}

func (v *View) Lookup(*Thread, Sels) Row {
	panic(assert.ShouldNotReachHere())
}

func (v *View) String() string {
	return v.name + " /*view*/"
}

func (*View) Simple(*Thread) []Row {
	panic(assert.ShouldNotReachHere())
}
