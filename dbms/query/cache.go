// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"golang.org/x/exp/slices"
)

// cache is a very simple cache for query costs
// It uses linear search so it is only suitable for a small numbers of entries.
// It does not limit the size of the cache. (no eviction)
type cache struct {
	entries []cacheEntry
	frac    float64
	fixcost Cost
	varcost Cost
}

type cacheEntry struct {
	index    []string
	frac     float32 // deliberately less precision
	fixcost  Cost
	varcost  Cost
	approach any
}

// cacheAdd adds an entry to the cache.
// It does *not* check if the item already exists
// because it assumes you previously tried cacheGet.
func (c *cache) cacheAdd(index []string, frac float64,
	fixcost Cost, varcost Cost, approach any) {
	assert.Msg("cache fixcost < 0").That(fixcost >= 0)
	assert.Msg("cache varcost < 0").That(varcost >= 0)
	c.entries = append(c.entries,
		cacheEntry{index: index, frac: float32(frac),
			fixcost: fixcost, varcost: varcost, approach: approach})
}

// cacheGet returns the cost and approach associated with an index
// or -1 if the index as not been added.
func (c *cache) cacheGet(index []string, frac float64) (
	fixcost, varcost Cost, approach any) {
	for i := range c.entries {
		if float32(frac) == c.entries[i].frac &&
			slices.Equal(index, c.entries[i].index) {
			slc.Swap(c.entries, 0, i) // so chosen approach is first
			return c.entries[0].fixcost, c.entries[0].varcost, c.entries[0].approach
		}
	}
	return -1, -1, nil
}

func (c *cache) cacheSetCost(frac float64, fixcost, varcost Cost) {
	c.frac, c.fixcost, c.varcost = frac, fixcost, varcost
}

func (c *cache) cacheCost() (float64, Cost, Cost) {
	return c.frac, c.fixcost, c.varcost
}
