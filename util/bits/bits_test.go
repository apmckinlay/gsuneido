// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package bits

import (
	"math"
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
