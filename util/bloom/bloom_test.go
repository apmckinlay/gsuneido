// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package bloom implements a simple Bloom filter
package bloom

import (
	"math/rand/v2"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestCalc(t *testing.T) {
	m, k := Calc(100000, 0.01)
	assert.Msg(m).That(958000 < m && m < 959000)
	assert.This(k).Is(7)
}

const phi64 = 0x9e3779b97f4a7c15

func TestBloom(t *testing.T) {
	bf := New(Calc(1000, .000001))
	r := rand.New(rand.NewPCG(1234, 5678))
	// assume we won't get duplicates in 1000 random 64 bit integers
	for range 1000 {
		n := r.Uint64() * phi64
		assert.That(!bf.Test(n))
		bf.Add(n)
	}
	r = rand.New(rand.NewPCG(1234, 5678))
	for range 1000 {
		n := r.Uint64() * phi64
		assert.That(bf.Test(n))
	}
}
