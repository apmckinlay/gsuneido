/*
Package hash provides a string hash function
that does not require copying to []byte
and returns the result as an integer.

Currently the hash function is FNV-1a
based on the standard Go fnv package
*/
package hash

const (
	offset32 = 2166136261
	prime32  = 16777619
)

func HashString(s string) uint32 {
	// TODO don't hash entire string if it's long
	hash := uint32(offset32)
	for i := 0; i < len(s); i++ {
		hash ^= uint32(s[i])
		hash *= prime32
	}
	return hash
}

func HashBytes(bytes []byte) uint32 {
	hash := uint32(offset32)
	for _, b := range bytes {
		hash ^= uint32(b)
		hash *= prime32
	}
	return hash
}
