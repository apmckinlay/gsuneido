// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package varint

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestVarint(t *testing.T) {
	assert := assert.T(t).This
	assert(EncodeUint32(0, []byte{})).Is([]byte{0})
	assert(EncodeInt32(0, []byte{})).Is([]byte{0})
	assert(EncodeUint32(45, []byte{})).Is([]byte{45})
	assert(EncodeInt32(45, []byte{})).Is([]byte{90})
	assert(EncodeInt32(-1, []byte{})).Is([]byte{1})
	assert(EncodeUint32(256, []byte{})).Is([]byte{128, 2})

	for _, n := range []int32{0, 1, -1, 123, -123, 999999, -999999} {
		n2, _ := DecodeInt32(EncodeInt32(n, []byte{}), 0)
		assert(n2).Msg("signed").Is(n)
		if n > 0 {
			n2, _ := DecodeUint32(EncodeUint32(uint32(n), []byte{}), 0)
			assert(n2).Msg("unsigned").Is(uint32(n))
		}
	}
}
