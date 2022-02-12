// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package pack

import "github.com/apmckinlay/gsuneido/util/hacks"

// Encoder is used to build a (usually binary) string.
// Encoder values should not be copied.
// We assume a fixed size buffer and do not use append
// to work with memory mapped byte slices.
// Write methods return the Encoder so they can be chained.
// Similar to strings.Builder.
type Encoder struct {
	buf []byte
}

// NewEncoder returns an Encoder with a buffer of the specified capacity
func NewEncoder(size int) *Encoder {
	return &Encoder{make([]byte, 0, size)}
}

// NewMmapEncoder returns an Encoder wrapping a memory mapped byte slice.
// It should NOT be used with other byte slices.
func NewMmapEncoder(buf []byte) *Encoder {
	return &Encoder{buf}
}

// String returns the accumulated data as a string
func (e *Encoder) String() string {
	return hacks.BStoS(e.buf)
}

// Buffer returns the accumulated data as a byte slice
func (e *Encoder) Buffer() []byte {
	return e.buf
}

// Put appends a byte slice
func (e *Encoder) Put(b []byte) *Encoder {
	e.buf = e.buf[:len(e.buf)+len(b)]
	copy(e.buf[len(e.buf)-len(b):], b)
	return e
}

// Put1 appends one or more bytes (or a byte slice)
func (e *Encoder) Put1(b byte) *Encoder {
	e.buf = e.buf[:len(e.buf)+1]
	e.buf[len(e.buf)-1] = b
	return e
}

// Put2 appends one or more bytes (or a byte slice)
func (e *Encoder) Put2(a, b byte) *Encoder {
	e.buf = e.buf[:len(e.buf)+2]
	e.buf[len(e.buf)-2] = a
	e.buf[len(e.buf)-1] = b
	return e
}

// Put4 appends one or more bytes (or a byte slice)
func (e *Encoder) Put4(a, b, c, d byte) *Encoder {
	e.buf = e.buf[:len(e.buf)+4]
	e.buf[len(e.buf)-4] = a
	e.buf[len(e.buf)-3] = b
	e.buf[len(e.buf)-2] = c
	e.buf[len(e.buf)-1] = d
	return e
}

// PutStr appends the contents of a string
func (e *Encoder) PutStr(s string) *Encoder {
	e.buf = e.buf[:len(e.buf)+len(s)]
	copy(e.buf[len(e.buf)-len(s):], s)
	return e
}

// Move moves the last nbytes over by shift bytes
func (e *Encoder) Move(nbytes, shift int) {
	n := len(e.buf)
	e.buf = e.buf[:n+shift]
	n -= nbytes
	copy(e.buf[n+shift:], e.buf[n:])
}
