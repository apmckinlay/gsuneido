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

func Hash(s string) uint32 {
	hash := uint32(offset32)
	for i := 0; i < len(s); i++ {
		hash ^= uint32(s[i])
		hash *= prime32
	}
	return hash
}
