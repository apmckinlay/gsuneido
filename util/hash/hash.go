/*
Package hash provides a string hash function
that does not require copying to []byte
and returns the result as an integer.

Currently the hash function is FNV-1a
based on the standard Go fnv package

Only the first 64 bytes are included if longer.
*/
package hash

import "github.com/apmckinlay/gsuneido/util/ints"

const (
	offset32 = 2166136261
	prime32  = 16777619
	maxlen   = 64
)

func HashString(s string) uint32 {
	n := ints.Min(len(s), maxlen)
	hash := uint32(offset32)
	for i := 0; i < n; i++ {
		hash ^= uint32(s[i])
		hash *= prime32
	}
	return hash
}

func HashBytes(bytes []byte) uint32 {
	n := ints.Min(len(bytes), maxlen)
	hash := uint32(offset32)
	for i := 0; i < n; i++ {
		hash ^= uint32(bytes[i])
		hash *= prime32
	}
	return hash
}
