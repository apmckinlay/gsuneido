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

const benchN = 1_000_000

// benchValue produces 32 bit values on the fly so the benchmarks share the
// same value domain without allocating a large slice.
func benchValue(i uint64) uint64 {
	return mix(i) & (1<<32 - 1)
}

// mix is the splitmix64 finalizer.
func mix(h uint64) uint64 {
	h += 0x9e3779b97f4a7c15
	h = (h ^ (h >> 30)) * 0xbf58476d1ce4e5b9
	h = (h ^ (h >> 27)) * 0x94d049bb133111eb
	return h ^ (h >> 31)
}

func BenchmarkBuild(b *testing.B) {
	bf := New(Calc(benchN, .01))
	for i := uint64(0); b.Loop(); i++ {
		if i%benchN == 0 {
			bf = New(Calc(benchN, .01))
		}
		bf.Add(benchValue(i))
	}
}

func BenchmarkProbe(b *testing.B) {
	bf := New(Calc(benchN, .01))
	for i := uint64(0); i < benchN; i++ {
		bf.Add(benchValue(i))
	}
	for i := uint64(0); b.Loop(); i++ {
		bf.Test(benchValue(i % benchN))
	}
}
