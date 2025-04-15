// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"fmt"
	rand "math/rand/v2"
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestLruCache(t *testing.T) {
	const size = 20
	get := func(key Value) Value { return key }
	lc := newLruCache(size)
	for range size {
		n := IntVal(int(rand.Uint32()))
		n2 := lc.GetPut(n, get)
		assert.This(n2).Is(n)
		n2 = lc.GetPut(n, get)
		assert.This(n2).Is(n)
	}
	assert.T(t).Msg("misses").This(lc.misses).Is(size)
	assert.T(t).Msg("hits").This(lc.hits).Is(size)
	for range size {
		n := IntVal(int(rand.Uint32()))
		n2 := lc.GetPut(n, get)
		assert.This(n2).Is(n)
		n2 = lc.GetPut(n, get)
		assert.This(n2).Is(n)
	}
	assert.T(t).Msg("misses").This(lc.misses).Is(size * 2)
	assert.T(t).Msg("hits").This(lc.hits).Is(size * 2)

	lc = newLruCache(size)
	for range 10000 {
		n := IntVal(rand.IntN(size + 5))
		n2 := lc.GetPut(n, get)
		assert.This(n2).Is(n)
	}
	lc.check()
	r := (100 * lc.hits) / 10000
	assert.T(t).That(r > 75) // theoretically 20/25 = 80%
}

func TestLruCache_concurrent(t *testing.T) {
	const size = 20
	lc := newSuLruCache(size, nil, false)
	a := &SuObject{}
	lc.Insert(Zero, a)
	assert.T(t).This(a.IsConcurrent()).Is(False)
	lc.SetConcurrent()
	assert.T(t).This(a.IsConcurrent()).Is(True)
	b := &SuObject{}
	lc.Insert(One, b)
	assert.T(t).This(b.IsConcurrent()).Is(True)
}

// check is used by the test
func (lc *lruCache) check() {
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
	lc := newLruCache(size)
	r := rand.New(rand.NewPCG(123, 456))
	get := func(key Value) Value { return key }
	for b.Loop() {
		lc.GetPut(IntVal(r.IntN(400)), get)
	}
	fmt.Println(lc.hits * 100 / (lc.hits + lc.misses))
}
