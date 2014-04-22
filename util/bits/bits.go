package bits

// Nlz returns the number of leading zero bits
// Algorithm from Hacker's Delight by Henry Warren
func Nlz(x uint64) int {
	if x == 0 {
		return 64
	}
	n := 1
	if (x >> 32) == 0 {
		n = n + 32
		x = x << 32
	}
	if (x >> 48) == 0 {
		n = n + 16
		x = x << 16
	}
	if (x >> 56) == 0 {
		n = n + 8
		x = x << 8
	}
	if (x >> 60) == 0 {
		n = n + 4
		x = x << 4
	}
	if (x >> 62) == 0 {
		n = n + 2
		x = x << 2
	}
	n = n - int(x>>63)
	return n
}
