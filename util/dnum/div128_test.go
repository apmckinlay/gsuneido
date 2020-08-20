// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dnum

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestDiv128(t *testing.T) {
	assert := assert.T(t).This
	assert(div128(1, 4)).Is(uint64(2500000000000000))
	assert(div128(1, 3)).Is(uint64(3333333333333333))
	assert(div128(2, 3)).Is(uint64(6666666666666666))
	assert(div128(1, 11)).Is(uint64(909090909090909))
	assert(div128(11, 13)).Is(uint64(8461538461538461))
	assert(div128(1234567890123456, 9876543210987654)).
		Is(uint64(1249999988609374))
}
