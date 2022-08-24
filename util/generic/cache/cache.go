// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package cache

import "sync"

func New[K comparable, E any](getter func(key K) E) *Cache[K, E] {
	return &Cache[K, E]{getter: getter}
}

const cacheSize = 8 // preferably a power of 2

// Cache is a simple cache.
// Zero value is not usable, use New.
type Cache[K comparable, E any] struct {
	key    [cacheSize]K
	val    [cacheSize]E
	used   [cacheSize]bool
	i      int
	getter func(key K) E
}

func (c *Cache[K, E]) Get(key K) E {
	i := c.i
	for {
		if c.used[c.i] && c.key[c.i] == key {
			return c.val[c.i]
		}
		c.i = (c.i + 1) % cacheSize
		if c.i == i { // not found
			// current slot is recently used, so advance to & reuse next slot
			c.i = (c.i + 1) % cacheSize
			val := c.getter(key)
			c.key[c.i] = key
			c.val[c.i] = val
			c.used[c.i] = true
			return val
		}
	}
}

func NewConc[K comparable, E any](getter func(key K) E) *ConcCache[K, E] {
	return &ConcCache[K, E]{Cache: Cache[K, E]{getter: getter}}
}

type ConcCache[K comparable, E any] struct {
	Cache[K, E]
	lock sync.Mutex
}

func (c *ConcCache[K, E]) Get(key K) E {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.Cache.Get(key)
}
