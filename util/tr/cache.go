// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tr

const cacheSize = 8 // should be power of 2

// Cache is a simple cache for Set.
// Zero value is ready to go
type Cache struct {
	slots [cacheSize]slot
	i     int
}

type slot struct {
	s string
	t trset
}

func (c *Cache) Get(s string) trset {
	i := c.i
	for {
		if c.slots[c.i].s == s {
			return c.slots[c.i].t
		}
		c.i = (c.i + 1) % cacheSize
		if c.i == i { // not found
			// current slot is recently used, so advance to & reuse next slot
			c.i = (c.i + 1) % cacheSize
			t := Set(s)
			c.slots[c.i] = slot{s: s, t: t}
			return t
		}
	}
}
