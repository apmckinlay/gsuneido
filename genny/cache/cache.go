package cache

import "github.com/cheekybits/genny/generic"

type K generic.Type
type V generic.Type

// VCache is a cache with a crude approximation of LRU.
// Working set should fit in size.
// NOT thread safe.
type VCache struct {
	old   map[K]V
	cur   map[K]V
	size  int
	fn    func(K) V
	nget  int
	nmiss int
}

func NewVCache(size int, fn func(K) V) *VCache {
	return &VCache{cur: make(map[K]V, size),
		old: make(map[K]V), size: size, fn: fn}
}

func (c *VCache) Get(key K) V {
	c.nget++
	val, ok := c.cur[key]
	if ok {
		return val
	}
	val, ok = c.old[key]
	if !ok {
		c.nmiss++
		val = c.fn(key)
	}
	if len(c.cur) >= c.size {
		c.old = c.cur
		c.cur = make(map[K]V, c.size)
	}
	c.cur[key] = val
	return val
}

func (c *VCache) NGet() int {
	return c.nget
}

func (c *VCache) NMiss() int {
	return c.nmiss
}
