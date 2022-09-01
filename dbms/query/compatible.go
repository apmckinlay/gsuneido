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
	Query2
	allCols  []string
	disjoint string
	keyIndex []string
	hdr1     *Header
	hdr2     *Header
	st       *SuTran
}

func (c *Compatible) init() {
	c.allCols = set.Union(c.source.Columns(), c.source2.Columns())
	fixed1 := c.source.Fixed()
	fixed2 := c.source2.Fixed()
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
	cols1 := c.source.Columns()
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

// source2Has returns true if a row from source exists in source2
func (c *Compatible) source2Has(th *Thread, row Row) bool {
	if c.disjoint != "" {
		return false
	}
	if c.hdr1 == nil {
		c.hdr1 = c.source.Header()
		c.hdr2 = c.source2.Header()
	}
	fixed := c.Fixed()
	cols := slices.Clip(c.keyIndex)
	vals := make([]string, 0, len(c.keyIndex)+len(fixed))
	for _, col := range c.keyIndex {
		vals = append(vals, row.GetRawVal(c.hdr1, col, th, c.st))
	}
	for _, fix := range fixed {
		if len(fix.values) == 1 {
			cols = append(cols, fix.col)
			vals = append(vals, fix.values[0])
		}
	}
	row2 := c.source2.Lookup(th, cols, vals)
	return row2 != nil && c.equal(row, row2, th)
}

func (c *Compatible) equal(row1, row2 Row, th *Thread) bool {
	if c.disjoint != "" {
		return false
	}
	return EqualRows(c.hdr1, row1, c.hdr2, row2, c.allCols, th, c.st)
}

func bestKey(q Query, mode Mode) []string {
	var best []string
	bestCost := impossible
	for _, key := range q.Keys() {
		cost := Optimize(q, mode, key)
		cost += (len(key) - 1) * cost / 20 // ??? prefer shorter keys
		if cost < bestCost {
			best = key
		}
	}
	return best
}
