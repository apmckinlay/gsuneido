// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

// cache is a very simple cache for query costs
// It uses linear search so it is only suitable for small numbers of entries.
// It does not limit the size of the cache.
type cache struct {
	entries []cacheEntry
}

type cacheEntry struct {
	index []string
	cost  Cost
}

// cacheAdd adds an entry to the cache.
// It does *not* check if the item already exists
// because it assumes you previously tried cacheGet.
func (c *cache) cacheAdd(index []string, cost Cost) {
	assert.That(cost >= 0)
	c.entries = append(c.entries, cacheEntry{index: index, cost: cost})
}

// cacheGet returns the cost associated with an index
// or -1 if the index as not been added.
func (c *cache) cacheGet(index []string) Cost {
	for i := range c.entries {
		if str.Equal(index, c.entries[i].index) {
			return c.entries[i].cost
		}
	}
	return -1
}
