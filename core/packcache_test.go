package core

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestPackCache(t *testing.T) {
	cache := &packCache{}

	test := func(key Value, expected int) {
		t.Helper()
		val, found := cache.Get(key)
		assert.T(t).This(found).Is(expected != -1)
		if found {
			assert.T(t).This(val).Is(expected)
		}
	}
	// Test Put and Get
	cache.Put(SuStr("key1"), 100)
	cache.Put(SuStr("key2"), 200)
	cache.Put(SuStr("key3"), 300)

	test(SuStr("key1"), 100)
	test(SuStr("key2"), 200)
	test(SuStr("key3"), 300)

	// Test cache miss
	test(SuStr("nonexistent"), -1)

	// Test cache eviction
	for i := range packCacheSize {
		cache.Put(SuInt(i), i)
	}
	// This should evict the oldest entry (key1)
	cache.Put(SuStr("newkey"), 1000)

	test(SuStr("key1"), -1)
	test(SuStr("newkey"), 1000)
}
