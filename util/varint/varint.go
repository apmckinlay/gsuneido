// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package varint

import "math/bits"

// Len returns the number of bytes required to varint encode.
// signed and unsigned are the same size
func Len(n uint64) int {
	if n == 0 {
		return 1
	}
	return (bits.Len64(n) + 6) / 7
}
