// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/sset"
)

type Compatible struct {
	Query2
	allCols  []string
	disjoint string
	keyIndex []string
}

func (c *Compatible) Init() {
	c.Query2.Init()
	c.allCols = sset.Union(c.source.Columns(), c.source2.Columns())
	fixed1 := c.source.Fixed()
	fixed2 := c.source2.Fixed()
	for _, f1 := range fixed1 {
		for _, f2 := range fixed2 {
			if f1.col == f2.col && vDisjoint(f1.values, f2.values) {
				c.disjoint = f1.col
				return
			}
		}
	}
	cols2 := c.source2.Columns()
	for _, f1 := range fixed1 {
		if !sset.Contains(cols2, f1.col) && !vContains(f1.values, EmptyStr) {
			c.disjoint = f1.col
			return
		}
	}
	cols1 := c.source.Columns()
	for _, f2 := range fixed2 {
		if !sset.Contains(cols1, f2.col) && !vContains(f2.values, EmptyStr) {
			c.disjoint = f2.col
			return
		}
	}
}

func (c *Compatible) String2(op string) string {
	if c.disjoint != "" {
		op += "-DISJOINT(" + c.disjoint + ")"
	}
	return c.Query2.String2(op)
}
