// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package hash

import (
	"hash/maphash"
)

const maxlen = 64 // ???

var seed maphash.Seed = maphash.MakeSeed()

// String uses maphash.String to hash the first 64 bytes of s
func String(s string) uint64 {
	if len(s) > maxlen {
		return 31*maphash.String(seed, s[:maxlen]) + uint64(len(s))
	}
	return maphash.String(seed, s)
}

// FullString uses maphash.FullString to hash all of s
func FullString(s string) uint64 {
	return maphash.String(seed, s)
}

// Bytes uses maphash.Bytes to hash the first 64 bytes of b
func Bytes(b []byte) uint64 {
	if len(b) > maxlen {
		return 31*maphash.Bytes(seed, b[:maxlen]) + uint64(len(b))
	}
	return maphash.Bytes(seed, b)
}

// HashString is the old version, used for compatibility with db checksums.
func HashString(s string) uint32 {
	const offset32 = 2166136261
	const prime32 = 16777619
	n := min(len(s), maxlen)
	hash := uint32(offset32)
	for i := range n {
		hash ^= uint32(s[i])
		hash *= prime32
	}
	return hash
}
