// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/util/assert"
	"golang.org/x/exp/slices"
)

// cache is a very simple cache for query costs
// It uses linear search so it is only suitable for a small numbers of entries.
// It does not limit the size of the cache. (no eviction)
type cache struct {
	entries []cacheEntry
}

type cacheEntry struct {
	index    []string
	cost     Cost
	approach interface{}
}

// cacheAdd adds an entry to the cache.
// It does *not* check if the item already exists
// because it assumes you previously tried cacheGet.
func (c *cache) cacheAdd(index []string, cost Cost, approach interface{}) {
	assert.Msg("cache cost < 0").That(cost >= 0)
	c.entries = append(c.entries,
		cacheEntry{index: index, cost: cost, approach: approach})
}

// cacheGet returns the cost and approach associated with an index
// or -1 if the index as not been added.
func (c *cache) cacheGet(index []string) (Cost, interface{}) {
	for i := range c.entries {
		if slices.Equal(index, c.entries[i].index) {
			return c.entries[i].cost, c.entries[i].approach
		}
	}
	return -1, nil
}
