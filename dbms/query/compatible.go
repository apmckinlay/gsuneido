// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"slices"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/str"
)

// Compatible is shared by Intersect, Minus, and Union
type Compatible struct {
	st       *SuTran
	disjoint string
	allCols  []string
	keyIndex []string
	Query2
}

// newCompatible sets disjoint
func newCompatible(src1, src2 Query) *Compatible {
	c := &Compatible{}
	c.source1, c.source2 = src1, src2
	c.allCols = set.Union(c.source1.Columns(), c.source2.Columns())
	fixed1 := src1.Fixed()
	fixed2 := src2.Fixed()
	cols1 := src1.Columns()
	cols2 := src2.Columns()
	for _, f1 := range fixed1 {
		for _, f2 := range fixed2 {
			if f1.col == f2.col && set.Disjoint(f1.values, f2.values) {
				c.disjoint = f1.col
				goto done
			}
		}
	}
	for _, f1 := range fixed1 {
		if !slices.Contains(cols2, f1.col) && !slices.Contains(f1.values, "") {
			c.disjoint = f1.col
			goto done
		}
	}
	for _, f2 := range fixed2 {
		if !slices.Contains(cols1, f2.col) && !slices.Contains(f2.values, "") {
			c.disjoint = f2.col
			goto done
		}
	}
done:
	return c
}

func (c *Compatible) String(s string) string {
	if c.keyIndex != nil {
		s += str.Join("(,)", c.keyIndex)
	}
	return s
}

func (c *Compatible) SetTran(t QueryTran) {
	c.st = MakeSuTran(t)
}

// source2Has returns true if a row from source exists in source2.
// It does Lookup on source2.
func (c *Compatible) source2Has(th *Thread, row Row) bool {
	if c.disjoint != "" {
		return false
	}
	vals := make([]string, len(c.keyIndex))
	for i, col := range c.keyIndex {
		vals[i] = row.GetRawVal(c.source1.Header(), col, th, c.st)
	}
	row2 := c.source2.Lookup(th, c.keyIndex, vals)
	return row2 != nil && c.equal(th, row, row2)
}

func (c *Compatible) equal(th *Thread, row1, row2 Row) bool {
	if c.disjoint != "" {
		return false
	}
	return EqualRows(c.source1.Header(), row1, c.source2.Header(), row2,
		c.allCols, th, c.st)
}

// bestLookupKey is used by Intersect and Minus
func bestLookupKey(q Query, mode Mode, nrows int) bestIndex {
	//TODO possibly exclude fixed
	best := newBestIndex()
	for _, key := range q.Keys() {
		fixcost, varcost := LookupCost(q, mode, key, nrows)
		best.update(key, fixcost, varcost)
	}
	return best
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
	c1.nsels++
	c1.source1.Select(cols, vals)
}

func (c1 *Compatible1) getLookupCost() int {
	cost := c1.source1.lookupCost()
	if c1.disjoint == "" {
		cost += c1.source2.lookupCost()
	}
	return cost
}
