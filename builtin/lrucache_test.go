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
	for i := 0; i < size; i++ {
		n := IntVal(int(rand.Uint32()))
		n2 := lc.GetPut(n, get)
		assert.This(n2).Is(n)
		n2 = lc.GetPut(n, get)
		assert.This(n2).Is(n)
	}
	assert.T(t).Msg("misses").This(lc.misses).Is(size)
	assert.T(t).Msg("hits").This(lc.hits).Is(size)
	for i := 0; i < size; i++ {
		n := IntVal(int(rand.Uint32()))
		n2 := lc.GetPut(n, get)
		assert.This(n2).Is(n)
		n2 = lc.GetPut(n, get)
		assert.This(n2).Is(n)
	}
	assert.T(t).Msg("misses").This(lc.misses).Is(size * 2)
	assert.T(t).Msg("hits").This(lc.hits).Is(size * 2)

	lc = newLruCache(size)
	for i := 0; i < 10000; i++ {
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
