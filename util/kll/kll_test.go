// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package kll

import (
	"math/rand/v2"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/bits"
	"github.com/apmckinlay/gsuneido/util/generic/ints"
)

func TestInsertLevel0(t *testing.T) {
	sk := New[int]()
	sk.Insert(1)
	sk.Insert(2)
	assert.This(sk.Count()).Is(2)
	assert.This(sk.levels[0]).Is([]int{1, 2})
}

func TestInsertCompact(t *testing.T) {
	sk := New[int]()
	sk.k = 8
	for i := range 10 {
		sk.Insert(i)
	}
	assert.This(sk.levels[0]).Is([]int{9})
}

func TestQuery(t *testing.T) {
	sk := New[int]()
	for i := range 4000 {
		sk.Insert(int(bits.Shuffle16(uint16(i))))
	}
	m := sk.Query(0.5)
	assert.That(ints.Abs(m-32767) < 500)

	sk = New[int]()
	for range 4000 {
		sk.Insert(rand.IntN(1000))
	}
	m = sk.Query(0.5)
	assert.That(ints.Abs(m-500) < 25)
}
