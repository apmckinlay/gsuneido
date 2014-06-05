package regex

// This is an instantiation of util/cache.go
// Do not modify here.

// LruMapCache is a cache with a crude approximation of LRU.
// Working set should fit in size.
// NOT thread safe.
type LruMapCache struct {
	old   map[string]Pattern
	cur   map[string]Pattern
	size  int
	fn    func(string) Pattern
	nget  int
	nmiss int
}

func NewLruMapCache(size int, fn func(string) Pattern) *LruMapCache {
	return &LruMapCache{cur: make(map[string]Pattern, size),
		old: make(map[string]Pattern), size: size, fn: fn}
}

func (c *LruMapCache) Get(key string) Pattern {
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
		c.cur = make(map[string]Pattern, c.size)
	}
	c.cur[key] = val
	return val
}

func (c *LruMapCache) NGet() int {
	return c.nget
}

func (c *LruMapCache) NMiss() int {
	return c.nmiss
}
