// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"math/rand"
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
		n := IntVal(rand.Intn(size + 5))
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
		x, _ := lc.hm.Get(e.key)
		assert.That(x != nil)
		xi, _ := x.ToInt()
		assert.That(xi == int(ei))
	}
	for ei, e := range lc.entries {
		x, _ := lc.hm.Get(e.key)
		assert.That(x != nil)
		xi, _ := x.ToInt()
		assert.That(xi == int(ei))
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
