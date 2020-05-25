package cache

import (
	"fmt"
	"strings"
	"testing"
)

func fn(k K) V {
	return strings.ToUpper(k.(string))
}

func TestCache(t *testing.T) {
	c := NewVCache(10, fn)
	assertEquals(t, c.Get("hello"), "HELLO")
	assertEquals(t, c.NGet(), 1)
	assertEquals(t, c.NMiss(), 1)
	assertEquals(t, c.Get("hello"), "HELLO")
	assertEquals(t, c.Get("hello"), "HELLO")
	assertEquals(t, c.NGet(), 3)
	assertEquals(t, c.NMiss(), 1)
}

func TestCache_eviction(t *testing.T) {
	c := NewVCache(5, fn)
	for i := 0; i < 10; i++ {
		c.Get(fmt.Sprint("key", i))
	}
	for i := 9; i >= 0; i-- {
		c.Get(fmt.Sprint("key", i))
	}
	assertEquals(t, c.NGet(), 20)
	assertEquals(t, c.NMiss(), 14) // specific to strategy
}

func assertEquals(t *testing.T, x interface{}, y interface{}) {
	if x != y {
		t.Error("expected", x, "==", y)
	}
}
