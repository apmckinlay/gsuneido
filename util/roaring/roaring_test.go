// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package roaring

import (
	"math/rand/v2"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestBasic(t *testing.T) {
	b := &Bitmap{}

	// Test adding and checking values
	b.Add(100)
	assert.T(t).This(b.Has(100)).Is(true)
	assert.T(t).This(b.Has(101)).Is(false)

	b.Add(101)
	assert.T(t).This(b.Has(100)).Is(true)
	assert.T(t).This(b.Has(101)).Is(true)
	assert.T(t).This(b.Has(102)).Is(false)
}

func TestMultipleContainers(t *testing.T) {
	b := &Bitmap{}

	// Add values in different containers (different high 32 bits)
	b.Add(100)
	b.Add(65536 + 200)
	b.Add(2*65536 + 300)

	assert.T(t).This(b.Has(100)).Is(true)
	assert.T(t).This(b.Has(65536 + 200)).Is(true)
	assert.T(t).This(b.Has(2*65536 + 300)).Is(true)
	assert.T(t).This(b.Has(101)).Is(false)
	assert.T(t).This(b.Has(65536 + 201)).Is(false)
	assert.T(t).This(b.Has(1 << 15)).Is(false)
}

func TestArrayToBitmapConversion(t *testing.T) {
	b := &Bitmap{}

	// Add enough values to trigger conversion to bitmap
	for i := uint64(0); i < 4100; i++ {
		b.Add(i)
	}

	// Verify all values are present
	for i := uint64(0); i < 4100; i++ {
		assert.T(t).This(b.Has(i)).Is(true)
	}
	assert.T(t).This(b.Has(5000)).Is(false)

	// Verify it was converted to bitmap
	assert.T(t).This(b.data[0].bitmap).Is(true)
}

func TestDuplicates(t *testing.T) {
	b := &Bitmap{}

	b.Add(100)
	b.Add(100)
	b.Add(100)

	assert.T(t).This(b.Has(100)).Is(true)
	assert.T(t).This(len(b.data)).Is(1)
	assert.T(t).This(len(b.data[0].data)).Is(1)
}

func TestBitOperations(t *testing.T) {
	blk := make([]uint16, 4096)

	addBit(blk, 0)
	assert.T(t).This(hasBit(blk, 0)).Is(true)
	assert.T(t).This(hasBit(blk, 1)).Is(false)

	addBit(blk, 15)
	assert.T(t).This(hasBit(blk, 15)).Is(true)

	addBit(blk, 16)
	assert.T(t).This(hasBit(blk, 16)).Is(true)

	addBit(blk, 65535)
	assert.T(t).This(hasBit(blk, 65535)).Is(true)
	assert.T(t).This(hasBit(blk, 65534)).Is(false)
}

func BenchmarkAdd(b *testing.B) {
	limit := 16 * 64 * 1024
	rnd := func() uint64 {
		// log distribution
		return uint64(rand.IntN(limit)) >> uint64(rand.IntN(4))
	}
	bm := &Bitmap{}
	for b.Loop() {
		bm = &Bitmap{}
		for range 100_000 {
			bm.Add(rnd())
		}
	}
	var bitmap, array int
	for _, d := range bm.data {
		if d.bitmap {
			bitmap++
		} else {
			array++
		}
	}
	// b.ReportMetric(float64(bitmap), "bitmap")
	// b.ReportMetric(float64(array), "array")
}

func BenchmarkHas(b *testing.B) {
	limit := 16 * 64 * 1024
	rnd := func() uint64 {
		// log distribution
		return uint64(rand.IntN(limit)) >> uint64(rand.IntN(8))
	}
	bm := &Bitmap{}
	for range 200_000 {
		bm.Add(rnd())
	}
	var hit, miss int
	for b.Loop() {
		if bm.Has(rnd()) {
			hit++
		} else {
			miss++
		}
	}
	var bitmap, array int
	for _, d := range bm.data {
		if d.bitmap {
			bitmap++
		} else {
			array++
		}
	}
	b.ReportMetric(float64(bitmap), "bitmap")
	b.ReportMetric(float64(array), "array")
	b.ReportMetric(float64(hit)/float64(b.N), "hit")
	b.ReportMetric(float64(miss)/float64(b.N), "miss")
}
