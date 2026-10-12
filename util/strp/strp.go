// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package strp provides a small string optimization.
// It represents strings as *byte
// with the length in the first byte.
// This allows storing in an interface
// without the overhead of allocating a string header.
package strp

import "unsafe"

// MaxLen is the maximum length of a SuStrP
const MaxLen = 255

// MaxFromGo is the maximum length of a Go string
// that should be converted to a strp.
const MaxFromGo = 23

// empty is the shared byte for zero-length strings.
// (cannot be initialized with From, which would be a cycle)
// @immutable
var empty byte

var single = func() (t [512]byte) {
	for c := range 256 {
		t[2*c] = 1
		t[2*c+1] = byte(c)
	}
	return
}()

func Single(c byte) *byte {
	return (*byte)(unsafe.Pointer(&single[2*int(c)]))
}

func From(s string) *byte {
	n := len(s)
	if n <= 1 {
		if n == 0 {
			return &empty
		}
		return Single(s[0])
	}
	if n > MaxLen {
		panic("string too large for SuStrP")
	}
	b := make([]byte, n+1) // noscan, no pointers
	b[0] = byte(n)
	copy(b[1:], s)
	return &b[0]
}

// Concat returns the concatenation of the strings as a *byte (length-prefixed).
// The total length must be <= MaxLen.
func Concat(ss ...string) *byte {
	n := 0
	for _, s := range ss {
		n += len(s)
	}
	if n <= 1 {
		if n == 0 {
			return &empty
		}
		for _, s := range ss {
			if len(s) == 1 {
				return Single(s[0])
			}
		}
	}
	if n > MaxLen {
		panic("string too large for SuStrP")
	}
	b := make([]byte, n+1) // noscan, no pointers
	b[0] = byte(n)
	pos := 0
	for _, s := range ss {
		pos += copy(b[1+pos:], s)
	}
	return &b[0]
}

func Len(sp *byte) int {
	return int(*sp)
}

func String(sp *byte) string {
	n := int(*sp)
	if n == 0 {
		return ""
	}
	return unsafe.String((*byte)(unsafe.Add(unsafe.Pointer(sp), 1)), n)
}
