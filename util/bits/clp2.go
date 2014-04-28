package bits

// Clp2 returns the next (ceiling) power of 2
//
// Algorithm from Hacker's Delight by Henry Warren
func Clp2(x uint64) uint64 {
	x = x - 1
	x = x | (x >> 1)
	x = x | (x >> 2)
	x = x | (x >> 4)
	x = x | (x >> 8)
	x = x | (x >> 16)
	x = x | (x >> 32)
	return x + 1
}
