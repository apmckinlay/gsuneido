package cache

import (
	"fmt"
	"strings"
	"testing"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestCache(t *testing.T) {
	c := NewLruMapCache(10, strings.ToUpper)
	Assert(t).That(c.Get("hello"), Equals("HELLO"))
	Assert(t).That(c.NGet(), Equals(1))
	Assert(t).That(c.NMiss(), Equals(1))
	Assert(t).That(c.Get("hello"), Equals("HELLO"))
	Assert(t).That(c.Get("hello"), Equals("HELLO"))
	Assert(t).That(c.NGet(), Equals(3))
	Assert(t).That(c.NMiss(), Equals(1))
}

func TestCache_eviction(t *testing.T) {
	c := NewLruMapCache(5, strings.ToUpper)
	for i := 0; i < 10; i++ {
		c.Get(fmt.Sprint("key", i))
	}
	for i := 9; i >= 0; i-- {
		c.Get(fmt.Sprint("key", i))
	}
	Assert(t).That(c.NGet(), Equals(20))
	Assert(t).That(c.NMiss(), Equals(14)) // specific to strategy
}

// CacheLM is a cache with a crude approximation of LRU.
// Working set should fit in size.
// NOT thread safe.
type LruMapCache struct {
	old   map[string]string
	cur   map[string]string
	size  int
	fn    func(string) string
	nget  int
	nmiss int
}

func NewLruMapCache(size int, fn func(string) string) *LruMapCache {
	return &LruMapCache{cur: make(map[string]string, size),
		old: make(map[string]string), size: size, fn: fn}
}

func (c *LruMapCache) Get(key string) string {
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
		c.cur = make(map[string]string, c.size)
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
