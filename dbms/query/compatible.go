// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"golang.org/x/exp/slices"
)

// Compatible is shared by Intersect, Minus, and Union
type Compatible struct {
	hdr1     *Header
	hdr2     *Header
	st       *SuTran
	disjoint string
	fixed    []Fixed // set by operations
	allCols  []string
	keyIndex []string
	Query2
}

// init sets disjoint
func (c *Compatible) init(calcFixed func(fixed1, fixed2 []Fixed) []Fixed) {
	c.allCols = set.Union(c.source1.Columns(), c.source2.Columns())
	fixed1 := c.source1.Fixed()
	fixed2 := c.source2.Fixed()
	c.fixed = calcFixed(fixed1, fixed2)
	for _, f1 := range fixed1 {
		for _, f2 := range fixed2 {
			if f1.col == f2.col && set.Disjoint(f1.values, f2.values) {
				c.disjoint = f1.col
				return
			}
		}
	}
	cols2 := c.source2.Columns()
	for _, f1 := range fixed1 {
		if !slices.Contains(cols2, f1.col) && !slices.Contains(f1.values, "") {
			c.disjoint = f1.col
			return
		}
	}
	cols1 := c.source1.Columns()
	for _, f2 := range fixed2 {
		if !slices.Contains(cols1, f2.col) && !slices.Contains(f2.values, "") {
			c.disjoint = f2.col
			return
		}
	}
}

func (c *Compatible) String2(op string) string {
	return c.Query2.String2(op)
}

func (c *Compatible) stringOp(op, strategy string) string {
	if c.disjoint != "" {
		op += "-DISJOINT(" + c.disjoint + ")"
	}
	return op + strategy
}

func (c *Compatible) SetTran(t QueryTran) {
	c.st = MakeSuTran(t)
}

func (c *Compatible) Fixed() []Fixed {
	return c.fixed // set by init
}

// source2Has returns true if a row from source exists in source2
func (c *Compatible) source2Has(th *Thread, row Row) bool {
	if c.disjoint != "" {
		return false
	}
	if c.hdr1 == nil { // once only
		c.hdr1 = c.source1.Header()
		c.hdr2 = c.source2.Header()
	}
	vals := make([]string, len(c.keyIndex))
	for i, col := range c.keyIndex {
		vals[i] = row.GetRawVal(c.hdr1, col, th, c.st)
	}
	row2 := c.source2.Lookup(th, c.keyIndex, vals)
	return row2 != nil && c.equal(th, row, row2)
}

func (c *Compatible) equal(th *Thread, row1, row2 Row) bool {
	if c.disjoint != "" {
		return false
	}
	return EqualRows(c.hdr1, row1, c.hdr2, row2, c.allCols, th, c.st)
}

func bestKey2(src2 Query, mode Mode, nrows int) bestIndex {
	best := newBestIndex()
	for _, key := range src2.Keys() {
		fixcost, varcost := LookupCost(src2, mode, key, nrows)
		best.update(key, fixcost, varcost)
	}
	return best
}

func (c *Compatible) lookupCost() int {
	cost := c.source1.lookupCost()
	if c.disjoint == "" {
		cost += c.source2.lookupCost()
	}
	return cost
}

//-------------------------------------------------------------------

// Compatible1 is embedded by Intersect and Minus
// (that return a subset of source1 records)
type Compatible1 struct {
	Compatible
}

func (c1 *Compatible1) Rewind() {
	c1.source1.Rewind()
}

func (c1 *Compatible1) Select(cols, vals []string) {
	c1.source1.Select(cols, vals)
}

func (c1 *Compatible1) rowSize() int {
	return c1.source1.rowSize()
}
