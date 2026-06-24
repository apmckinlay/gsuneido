// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"slices"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/slc"
)

// cache is a very simple cache for query costs
// It uses linear search so it is only suitable for a small numbers of entries.
// It does not limit the size of the cache. (no eviction)
type cache struct {
	entries []cacheEntry
}

type cacheEntry struct {
	approach any
	use      Use
	index    []string
	fixcost  Cost
	varcost  Cost
	frac     float64
}

// cacheAdd adds an entry to the cache.
// It does *not* check if the item already exists
// because it assumes you previously tried cacheGet.
func (c *cache) cacheAdd(use Use, index []string, frac float64,
	fixcost Cost, varcost Cost, approach any) {
	assert.That(fixcost >= 0)
	assert.That(varcost >= 0)
	c.entries = append(c.entries,
		cacheEntry{use: use, index: index, frac: frac,
			fixcost: fixcost, varcost: varcost, approach: approach})
}

// cacheGet returns the cost and approach associated with an index
// or -1 if the index has not been added.
func (c *cache) cacheGet(use Use, index []string, frac float64) (
	fixcost, varcost Cost, approach any) {
	for i := range c.entries {
		if use == c.entries[i].use &&
			frac == c.entries[i].frac &&
			slices.Equal(index, c.entries[i].index) {
			slc.Swap(c.entries, 0, i) // so chosen approach is first
			return c.entries[0].fixcost, c.entries[0].varcost, c.entries[0].approach
		}
	}
	return -1, -1, nil
}

func (c *cache) cacheClear() {
	c.entries = nil
}

// cache2 is a separate cache for v2 unified optimize.
// It uses Require as the key, which includes cols, frac, and nlookups.
// Linear search matching all Require fields.
type cache2 struct {
	entries []cacheEntry2
}

type cacheEntry2 struct {
	Require
	fixcost  Cost
	varcost  Cost
	approach any
}

func (c *cache2) cacheAdd2(req Require, fixcost, varcost Cost, approach any) {
	assert.That(fixcost >= 0)
	assert.That(varcost >= 0)
	c.entries = append(c.entries,
		cacheEntry2{Require: req, fixcost: fixcost, varcost: varcost, approach: approach})
}

func (c *cache2) cacheGet2(req Require) (fixcost, varcost Cost, approach any) {
	for i := range c.entries {
		e := &c.entries[i]
		if slices.Equal(e.cols, req.cols) &&
			e.frac == req.frac && e.nlookups == req.nlookups {
			slc.Swap(c.entries, 0, i)
			return c.entries[0].fixcost, c.entries[0].varcost, c.entries[0].approach
		}
	}
	return -1, -1, nil
}

func (c *cache2) cacheClear2() {
	c.entries = nil
}
