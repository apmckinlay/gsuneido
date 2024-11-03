// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package str

// These are similar to the Go string (Last)IndexAny functions.
// But the Go ones are utf8 and only optimize for 7 bit ascii.
// Also, these ones allow negated sets and ranges like Tr.

func Find1of(s, chars string) int {
	if len(chars) == 0 {
		return -1
	}
	b := makeBits(chars)
	for i := range len(s) {
		if b.contains(s[i]) {
			return i
		}
	}
	return -1
}

func FindLast1of(s, chars string) int {
	if len(chars) == 0 {
		return -1
	}
	b := makeBits(chars)
	for i := len(s) - 1; i >= 0; i-- {
		if b.contains(s[i]) {
			return i
		}
	}
	return -1
}

type bits [4]uint64

func (b bits) contains(c byte) bool {
	return b[c/64]&(1<<(c%64)) != 0
}

func makeBits(chars string) bits {
	var b bits
	negated := chars[0] == '^'
	if negated {
		chars = chars[1:]
	}
	n := len(chars)
	for i := 0; i < n; i++ {
		if i+2 < n && chars[i+1] == '-' { // range
			// need to use int to prevent 0xff from wrapping around
			for c := int(chars[i]); c <= int(chars[i+2]); c++ {
				b[c/64] |= 1 << (c % 64)
			}
			i += 2
		} else {
			c := chars[i]
			if c == '-' && i > 0 && i+1 < n { // range
			} else {
				b[c/64] |= 1 << (c % 64)
			}
		}
	}
	if negated {
		for i := range b {
			b[i] = ^b[i]
		}
	}
	return b
}
