// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex

const cacheSize = 8 // should be power of 2
// var CacheGet = 0
// var CacheHit = 0

// Cache is a simple cache for compiled regex.
// Zero value is ready to go
type Cache struct {
	slots [cacheSize]slot
	i     int
}

type slot struct {
	rx  string
	pat Pattern
}

func (c *Cache) Get(rx string) Pattern {
	// CacheGet++
	i := c.i
	for {
		if c.slots[c.i].rx == rx {
			// CacheHit++
			return c.slots[c.i].pat
		}
		c.i = (c.i + 1) % cacheSize
		if c.i == i { // not found
			// current slot is recently used, so advance to & reuse next slot
			c.i = (c.i + 1) % cacheSize
			pat := Compile(rx)
			c.slots[c.i] = slot{rx: rx, pat: pat}
			return pat
		}
	}
}
