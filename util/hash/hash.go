// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package hash

import (
	"hash/maphash"

	"github.com/apmckinlay/gsuneido/util/generic/ord"
)

const maxlen = 64

var seed maphash.Seed = maphash.MakeSeed()

func String(s string) uint32 {
	if len(s) > maxlen {
		s = s[:maxlen]
	}
	return uint32(maphash.String(seed, s))
}

func Bytes(b []byte) uint32 {
	if len(b) > maxlen {
		b = b[:maxlen]
	}
	return uint32(maphash.Bytes(seed, b))
}

// HashString is the old version, used for compatibility with db checksums.
func HashString(s string) uint32 {
	const offset32 = 2166136261
	const prime32 = 16777619
	n := ord.Min(len(s), maxlen)
	hash := uint32(offset32)
	for i := 0; i < n; i++ {
		hash ^= uint32(s[i])
		hash *= prime32
	}
	return hash
}
