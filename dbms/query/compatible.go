// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"slices"
	"sync/atomic"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/dbg"
	"github.com/apmckinlay/gsuneido/util/set"
)

// Compatible is shared by Intersect, Minus, and Union
type Compatible struct {
	Query2
	st          *SuTran
	disjoint    string
	allCols     []string
	src1Only    []string
	lookupCols  []string
	lookupCache lookupCache
}

// newCompatible sets disjoint
var (
	compatCacheProbes atomic.Int64
	compatCacheMisses atomic.Int64
)

var _ = AddInfo("query.compatible.cacheProbes", &compatCacheProbes)
var _ = AddInfo("query.compatible.cacheMisses", &compatCacheMisses)

func newCompatible(src1, src2 Query) *Compatible {
	c := &Compatible{}
	c.source1, c.source2 = src1, src2
	cols1 := src1.Columns()
	cols2 := src2.Columns()
	c.allCols = set.Union(cols1, cols2)
	fixed1 := src1.Fixed()
	fixed2 := src2.Fixed()
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
	c.lookupCache.SetCounters(&compatCacheProbes, &compatCacheMisses)
	return c
}

func (c *Compatible) SetTran(t QueryTran) {
	c.st = MakeSuTran(t)
	c.lookupCache.Reset()
	c.Query2.SetTran(t)
}

// source2Has returns true if a row from source exists in source2.
// It does Lookup on source2 using all source2 columns.
func (c *Compatible) source2Has(th *Thread, row1 Row) bool {
	dbg.Assert(c.disjoint == "")
	hdr1 := c.source1.Header()
	for _, col := range c.src1Only {
		if x := row1.GetRaw(hdr1, col); x != "" && x[0] != PackForward {
			return false
		}
	}
	dbg.Assert(set.Equal(c.lookupCols, c.source2.Columns()))
	sels := makeSels(hdr1, row1, c.lookupCols, th, c.st)
	row2 := c.lookupCache.Lookup(c.source2, sels, th, c.st)
	return row2 != nil &&
		EqualRows(hdr1, row1, c.source2.Header(), row2, c.src1Only, th, c.st)
}

func makeSels(hdr *Header, row Row, cols []string, th *Thread, st *SuTran) Sels {
	sels := make(Sels, len(cols))
	for i, col := range cols {
		sels[i] = Sel{col: col, val: row.GetRawVal(hdr, col, th, st)}
	}
	return sels
}

func (c *Compatible) equal(th *Thread, row1, row2 Row) bool {
	if c.disjoint != "" {
		return false
	}
	return EqualRows(c.source1.Header(), row1, c.source2.Header(), row2,
		c.allCols, th, c.st)
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

func (c1 *Compatible1) Select(sels Sels) {
	c1.nsels++
	c1.source1.Select(sels)
}
