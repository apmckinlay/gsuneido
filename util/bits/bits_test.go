// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package bits

import (
	"math"
	"math/rand/v2"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestNextPow2(t *testing.T) {
	assert := assert.T(t).This
	assert(NextPow2(0)).Is(0)
	assert(NextPow2(1)).Is(1)
	assert(NextPow2(2)).Is(2)
	assert(NextPow2(3)).Is(4)
	assert(NextPow2(123)).Is(128)
	assert(NextPow2(65536)).Is(65536)
}

func TestShuffle16(t *testing.T) {
	// Cycle through all 2^16 values starting from 0
	// This verifies that Shuffle16 visits each value exactly once
	start := uint16(0)
	curr := start
	length := 0
	for {
		curr = Shuffle16(curr)
		length++
		if curr == start {
			break
		}
		if length > math.MaxUint16 {
			t.Errorf("Cycle too long: length=%d", length)
			return
		}
	}
	if length != 1<<16 {
		t.Errorf("Expected cycle length %d, got %d", 1<<16, length)
	}
}

func TestShuffle32(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long test")
	}
	// Cycle through all 2^32 values starting from 0
	// This verifies that Shuffle32 visits each value exactly once
	start := uint32(0)
	curr := start
	length := 0
	for {
		curr = Shuffle32(curr)
		length++
		if curr == start {
			break
		}
		if length > math.MaxUint32 {
			t.Errorf("Cycle too long: length=%d", length)
			return
		}
	}
	if length != 1<<32 {
		t.Errorf("Expected cycle length %d, got %d", 1<<32, length)
	}
}

func TestMix(t *testing.T) {
	for range 100 {
		b := rand.IntN(64)
		x := rand.Uint64N(1 << b)
		y := Mix(x, b)
		assert.T(t).That(y < 1<<b)
	}
	for _, b := range []int{3, 7, 17, 23} {
		n := uint64(1 << b)
		sum1 := uint64(0)
		sum2 := uint64(0)
		for i := range n {
			sum1 += i * i
			x := Mix(i, b)
			sum2 += x * x
		}
		assert.Msg(b, n).This(sum2).Is(sum1)
	}
}

func TestCycle(t *testing.T) {
	for range 100 {
		rangeLimit := rand.Uint64()
		x := rand.Uint64N(rangeLimit)
		y := Cycle(x, rangeLimit)
		assert.T(t).That(y < rangeLimit)
	}
	for range 100 {
		rangeLimit := rand.Uint64N(100_000)
		sum1 := uint64(0)
		sum2 := uint64(0)
		for i := range rangeLimit {
			sum1 += i * i
			x := Cycle(i, rangeLimit)
			assert.T(t).That(x < rangeLimit)
			sum2 += x * x
		}
		assert.Msg(rangeLimit).This(sum2).Is(sum1)
	}
}

func TestGen(t *testing.T) {
	rnd := rand.New(rand.NewPCG(1, 1))
	g := NewGen(rnd, 1000)
	for range 100 {
		x := g.Next()
		assert.T(t).That(x < 1000)
	}
}
