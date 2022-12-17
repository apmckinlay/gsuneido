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
	cost    int
}

type cacheEntry struct {
	mode     Mode
	index    []string
	cost     Cost
	approach any
}

// cacheAdd adds an entry to the cache.
// It does *not* check if the item already exists
// because it assumes you previously tried cacheGet.
func (c *cache) cacheAdd(mode Mode, index []string, cost Cost, approach any) {
	assert.Msg("cache cost < 0").That(cost >= 0)
	c.entries = append(c.entries,
		cacheEntry{mode: mode, index: index, cost: cost, approach: approach})
}

// cacheGet returns the cost and approach associated with an index
// or -1 if the index as not been added.
func (c *cache) cacheGet(mode Mode, index []string) (Mode, Cost, any) {
	for i := range c.entries {
		if mode == c.entries[i].mode &&
			slices.Equal(index, c.entries[i].index) {
			slc.Swap(c.entries, 0, i) // so chosen approach is first
			return c.entries[0].mode, c.entries[0].cost, c.entries[0].approach
		}
	}
	return -1, -1, nil
}

func (c *cache) cacheSetCost(cost int) {
	c.cost = cost
}

func (c *cache) cacheCost() int {
	return c.cost
}
