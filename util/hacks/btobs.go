// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package hacks

import (
	"unsafe"
)

// Btobs converts a byte to a byte slice without allocating.
// WARNING: The resulting byte slice must not be modified or escape.
func Btobs(b byte) []byte {
	return unsafe.Slice(&b, 1)
}
