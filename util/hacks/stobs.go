// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package hacks

import (
	"unsafe"
)

// Stobs converts a string to a byte slice without allocating a copy.
// WARNING: The resulting byte slice must not be modified.
// If the byte slice is modified the string will change
// which is illegal since strings are immutable.
// This is an optimization (to avoid allocation) that should not be overused.
func Stobs(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
