// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package varint

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestLen(t *testing.T) {
	tests := []struct {
		n    uint64
		want int
	}{
		{0, 1},           // 0 requires 1 byte
		{1, 1},           // Small values up to 127
		{127, 1},         // Maximum for 1 byte
		{128, 2},         // Minimum for 2 bytes
		{16383, 2},       // Maximum for 2 bytes
		{16384, 3},       // Minimum for 3 bytes
		{2097151, 3},     // Maximum for 3 bytes
		{2097152, 4},     // Minimum for 4 bytes
		{268435455, 4},   // Maximum for 4 bytes
		{268435456, 5},   // Minimum for 5 bytes
		{34359738367, 5}, // Maximum for 5 bytes
		{34359738368, 6}, // Minimum for 6 bytes
		{4398046511103, 6}, // Maximum for 6 bytes
		{4398046511104, 7}, // Minimum for 7 bytes
		{562949953421311, 7}, // Maximum for 7 bytes
		{562949953421312, 8}, // Minimum for 8 bytes
		{72057594037927935, 8}, // Maximum for 8 bytes
		{72057594037927936, 9}, // Minimum for 9 bytes
		{9223372036854775807, 9}, // Maximum for 9 bytes (int64 max)
		{18446744073709551615, 10}, // uint64 max
	}

	for _, test := range tests {
		got := Len(test.n)
		assert.T(t).This(got).Is(test.want)
	}
}
