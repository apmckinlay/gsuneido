// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package blkbloom

import (
	"math/rand/v2"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestAddContains(t *testing.T) {
	bf := New(1000)
	r := rand.New(rand.NewPCG(1234, 5678))
	for range 1000 {
		bf.Add(r.Uint64())
	}
	r = rand.New(rand.NewPCG(1234, 5678))
	for range 1000 {
		assert.T(t).True(bf.Has(r.Uint64()))
	}
}

func TestNoFalseNegatives(t *testing.T) {
	const n = 100_000
	bf := New(n)
	r := rand.New(rand.NewPCG(7, 21))
	for range n {
		bf.Add(r.Uint64())
	}
	r = rand.New(rand.NewPCG(7, 21))
	for range n {
		assert.T(t).True(bf.Has(r.Uint64()))
	}
}

func TestFalsePositiveRate(t *testing.T) {
	const n = 100_000
	const probes = 1_000_000
	bf := New(n)
	r := rand.New(rand.NewPCG(42, 99))
	for range n {
		bf.Add(r.Uint64())
	}
	fp := 0
	for range probes {
		if bf.Has(r.Uint64()) {
			fp++
		}
	}
	rate := float64(fp) / probes
	t.Logf("false positives: %d/%d = %.3f%%", fp, probes, 100*rate)
	// ~10 bits/item with 4 hashes gives roughly 1-3% false positives
	assert.T(t).That(rate < 0.03)
}

func TestFalsePositiveRate2(t *testing.T) {
	const n = 100_000
	const probes = 1_000_000
	bf := New(n)
	for i := range n {
		bf.Add(uint64(i))
	}
	fp := 0
	for i := range probes {
		if bf.Has(uint64(i + n)) {
			fp++
		}
	}
	rate := float64(fp) / probes
	t.Logf("in order false positives: %d/%d = %.3f%%", fp, probes, 100*rate)
	// ~10 bits/item with 4 hashes gives roughly 1-3% false positives
	assert.T(t).That(rate < 0.03)
}

func TestSmallSize(t *testing.T) {
	bf := New(0)
	assert.T(t).This(len(bf)).Is(uint64(1))
	bf.Add(12345)
	assert.T(t).True(bf.Has(12345))
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
	bf := New(benchN)
	for i := uint64(0); b.Loop(); i++ {
		if i%benchN == 0 {
			bf = New(benchN)
		}
		bf.Add(benchValue(i))
	}
}

func BenchmarkProbe(b *testing.B) {
	bf := New(benchN)
	for i := uint64(0); i < benchN; i++ {
		bf.Add(benchValue(i))
	}
	for i := uint64(0); b.Loop(); i++ {
		bf.Has(benchValue(i % benchN))
	}
}
