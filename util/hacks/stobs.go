// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package hacks

import (
	"reflect"
	"unsafe"
)

// Stobs converts a string to a byte slice without allocating a copy.
// WARNING: The resulting byte slice must not be modified.
// If the byte slice is modified the string will change
// which is illegal since strings are immutable.
// This is an optimization (to avoid allocation) that should not be overused.
func Stobs(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&reflect.SliceHeader{
			Data: (*reflect.StringHeader)(unsafe.Pointer(&s)).Data,
			Len:  len(s),
			Cap:  len(s),
		}))
}
