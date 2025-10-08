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
	assert := assert.T(t).This
	sum1 := 0
	sum2 := 0
	for i := 0; i <= math.MaxUint16; i++ {
		sum1 += i
		sum2 += int(Shuffle16(uint16(i)))
	}
	assert(sum1).Is((math.MaxUint16 * (math.MaxUint16 + 1)) / 2)
	assert(sum1).Is(sum2)
}

func TestShuffle32(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long test")
	}
	assert := assert.T(t).This
	sum1 := 0
	sum2 := 0
	for i := 0; i <= math.MaxUint32; i++ {
		sum1 += i
		sum2 += int(Shuffle32(uint32(i)))
	}
	assert(sum1).Is((math.MaxUint32 * (math.MaxUint32 + 1)) / 2)
	assert(sum1).Is(sum2)
}
