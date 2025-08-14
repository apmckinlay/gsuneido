// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package lrucache provides a small generic LRU cache for <=256 entries.
package lrucache

import (
	"bytes"
	"iter"

	"github.com/apmckinlay/gsuneido/util/shmap"
)

// Cache is a fixed-capacity least-recently-used cache 
// keyed by types that have Hash and Equal methods.
type Cache[K keyT, V any] struct {
	lru     []uint8 // uint8 means max size of 256
	entries []entry[K, V]
	hm      shmap.Map[keyT, uint8, shmap.Meth[keyT]]
	size    int
	hits    int
	misses  int
}

type entry[K keyT, V any] struct {
	key K
	val V
}

type keyT interface {
	Hash() uint64
	Equal(other any) bool
}

// New returns a cache sized to the nearest supported capacity >= `req`.
func New[K keyT, V any](req int) *Cache[K, V] {
	var sizes = []int{6, 13, 27, 55, 111, 223} // 7/8 of ^2
	size := 223                                // max
	for _, n := range sizes {
		if req <= n {
			size = n
			break
		}
	}
	return &Cache[K, V]{size: size,
		lru:     make([]uint8, 0, size),
		entries: make([]entry[K, V], 0, size),
	}
}

// Get returns the cached value for `key`, marking it as most recently used.
func (lc *Cache[K, V]) Get(key K) (val V, ok bool) {
	ei, ok := lc.hm.Get(key)
	if !ok {
		// not in cache
		lc.misses++
		return
	}
	// in cache
	lc.hits++
	li := bytes.IndexByte(lc.lru, uint8(ei))
	if li < lc.size-lc.size/8 {
		// move to the newest (the end)
		copy(lc.lru[li:], lc.lru[li+1:])
		lc.lru[len(lc.lru)-1] = uint8(ei)
	}
	return lc.entries[ei].val, true
}

// Put sets `key` to `val`, evicting the least-recently-used entry if full.
func (lc *Cache[K, V]) Put(key K, val V) {
	ei := len(lc.entries)
	if ei < lc.size {
		lc.entries = append(lc.entries, entry[K, V]{key: key, val: val})
		lc.lru = append(lc.lru, uint8(ei))
	} else { // full
		// replace oldest entry lru[0]
		ei = int(lc.lru[0])
		lc.hm.Del(lc.entries[ei].key)
		lc.entries[ei] = entry[K, V]{key: key, val: val}
		copy(lc.lru, lc.lru[1:])
		lc.lru[lc.size-1] = uint8(ei)
	}
	lc.hm.Put(key, uint8(ei))
}

// GetPut gets `key` or computes it with `getfn` on a miss and caches the result.
func (lc *Cache[K, V]) GetPut(key K, getfn func(key K) V) V {
	val, ok := lc.Get(key)
	if !ok {
		val = getfn(key)
		lc.Put(key, val)
	}
	return val
}

// Entries returns an iterator over current entries.
func (lc *Cache[K, V]) Entries() iter.Seq2[K, V] {
	return func(yield func(key K, val V) bool) {
		for _, e := range lc.entries {
			if !yield(e.key, e.val) {
				return
			}
		}
	}
}

// Stats returns cumulative hit/miss counts since creation or last Reset.
func (lc *Cache[K, V]) Stats() (hits, misses int) {
	return lc.hits, lc.misses
}

// Reset clears the cache and statistics.
func (lc *Cache[K, V]) Reset() {
	lc.hm.Clear()
	lc.lru = lc.lru[:0]
	clear(lc.entries)
	lc.entries = lc.entries[:0]
	lc.hits = 0
	lc.misses = 0
}

// func (lc *lruCache) print() {
// 	fmt.Println("lru")
// 	for li, ei := range lc.lru {
// 		fmt.Println(li, ei)
// 	}
// 	fmt.Println("entries")
// 	for ei, e := range lc.entries {
// 		fmt.Println(ei, e.key, e.val)
// 	}
// 	fmt.Println("hmap")
// 	it := lc.hm.Iter()
// 	for k, x := it(); k != nil; k, x = it() {
// 		fmt.Println(k, x)
// 	}
// }
