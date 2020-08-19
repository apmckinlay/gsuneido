// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package bytes

const smallBufferSize = 64

// Grow grows the buffer to guarantee space for n more bytes.
func Grow(buf []byte, n int) []byte {
	if n <= 0 {
		return buf
	}
	l := len(buf)
	// Try to grow by reslice (fast path)
	if n <= cap(buf)-l {
		return buf[:l+n]
	}
	if buf == nil && n <= smallBufferSize {
		return make([]byte, n, smallBufferSize)
	}
	// Not enough space, we need to allocate.
	buf2 := make([]byte, 2*cap(buf) + n)
	copy(buf2, buf)
	return buf2[:l+n]
}
