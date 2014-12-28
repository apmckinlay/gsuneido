package tr

// This is an instantiation of util/cache.go
// Do not modify here.

// LruMapCache is a cache with a crude approximation of LRU.
// Working set should fit in size.
// NOT thread safe.
type LruMapCache struct {
	old   map[string]trset
	cur   map[string]trset
	size  int
	fn    func(string) trset
	nget  int
	nmiss int
}

func NewLruMapCache(size int, fn func(string) trset) *LruMapCache {
	return &LruMapCache{cur: make(map[string]trset, size),
		old: make(map[string]trset), size: size, fn: fn}
}

func (c *LruMapCache) Get(key string) trset {
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
		c.cur = make(map[string]trset, c.size)
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
