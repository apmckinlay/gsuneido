// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package lrucache

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

type testKey int

const phi64 = 0x9e3779b97f4a7c15

func (v testKey) Hash() uint64 { return uint64(v) * phi64 }

func (v testKey) Equal(other any) bool { return v == other.(testKey) }

func TestLruCache(t *testing.T) {
	const size = 20
	get := func(key testKey) int { return int(key) }
	lc := New[testKey, int](size)
	for range size {
		n := rand.Uint32()
		k := testKey(n)
		n2 := lc.GetPut(k, get)
		assert.This(n2).Is(n)
		n2 = lc.GetPut(k, get)
		assert.This(n2).Is(n)
	}
	assert.T(t).Msg("misses").This(lc.misses).Is(size)
	assert.T(t).Msg("hits").This(lc.hits).Is(size)
	for range size {
		n := rand.Uint32()
		k := testKey(n)
		n2 := lc.GetPut(k, get)
		assert.This(n2).Is(n)
		n2 = lc.GetPut(k, get)
		assert.This(n2).Is(n)
	}
	assert.T(t).Msg("misses").This(lc.misses).Is(size * 2)
	assert.T(t).Msg("hits").This(lc.hits).Is(size * 2)

	lc = New[testKey, int](size)
	for range 10000 {
		n := rand.IntN(size + 5)
		k := testKey(n)
		n2 := lc.GetPut(k, get)
		assert.This(n2).Is(n)
	}
	lc.check()
	r := (100 * lc.hits) / 10000
	assert.T(t).That(r > 75) // theoretically 20/25 = 80%
}

// check is used by the test
func (lc *Cache[testKey, int]) check() {
	for _, ei := range lc.lru {
		e := lc.entries[ei]
		xi, ok := lc.hm.Get(e.key)
		assert.That(ok)
		assert.That(xi == ei)
	}
	for ei, e := range lc.entries {
		xi, ok := lc.hm.Get(e.key)
		assert.That(ok)
		assert.This(xi).Is(ei)
	}
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

func BenchmarkLruCache(b *testing.B) {
	const size = 223
	lc := New[testKey, int](size)
	r := rand.New(rand.NewPCG(123, 456))
	get := func(key testKey) int { return int(key) }
	for b.Loop() {
		lc.GetPut(testKey(r.IntN(400)), get)
	}
	fmt.Println("hit %", lc.hits*100/(lc.hits+lc.misses))
}
