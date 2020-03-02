// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package hacks

import "unsafe"

// BStoS converts a byte slice to a string.
// Copied from strings.Builder String method.
// WARNING: this should only be used when the byte slice is final.
// If the byte slice is modified the string may change
// which is illegal since strings are immutable.
// This is an optimization that should not be overused.
func BStoS(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
	// return string(bs)
}
