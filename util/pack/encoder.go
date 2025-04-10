// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package pack

import (
	"encoding/binary"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/hacks"
)

// encoder is used to build a (usually binary) string.
// encoder values should not be copied.
// We assume a fixed size buffer and do not use append
// to work with memory mapped byte slices.
// Write methods return the encoder so they can be chained.
// Similar to strings.Builder.
type encoder struct {
	buf []byte
}

// Newencoder returns an encoder with a buffer of the specified capacity
func NewEncoder(size int) Encoder {
	return &encoder{make([]byte, 0, size)}
}

func (e *encoder) Dup() Encoder {
	return &encoder{e.buf}
}

// String returns the accumulated data as a string.
func (e *encoder) String() string {
	s := hacks.BStoS(e.buf)
	e.buf = nil // ownership transferred to string
	return s
}

// Buffer returns the accumulated data as a byte slice.
// Not used by v2.
func (e *encoder) Buffer() []byte {
	return e.buf
}

// Put appends a byte slice
func (e *encoder) Put(b []byte) Encoder {
	e.buf = e.buf[:len(e.buf)+len(b)]
	copy(e.buf[len(e.buf)-len(b):], b)
	return e
}

// Put1 appends one byte
func (e *encoder) Put1(b byte) Encoder {
	e.buf = e.buf[:len(e.buf)+1]
	e.buf[len(e.buf)-1] = b
	return e
}

// Put2 appends two bytes
func (e *encoder) Put2(a, b byte) Encoder {
	e.buf = e.buf[:len(e.buf)+2]
	e.buf[len(e.buf)-2] = a
	e.buf[len(e.buf)-1] = b
	return e
}

// Put4 appends four bytes
func (e *encoder) Put4(a, b, c, d byte) Encoder {
	e.buf = e.buf[:len(e.buf)+4]
	e.buf[len(e.buf)-4] = a
	e.buf[len(e.buf)-3] = b
	e.buf[len(e.buf)-2] = c
	e.buf[len(e.buf)-1] = d
	return e
}

// PutStr appends a string
func (e *encoder) PutStr(s string) Encoder {
	e.buf = e.buf[:len(e.buf)+len(s)]
	copy(e.buf[len(e.buf)-len(s):], s)
	return e
}

// Move moves the last nbytes over by shift bytes.
// Not used by v2.
func (e *encoder) Move(nbytes, shift int) {
	n := len(e.buf)
	e.buf = e.buf[:n+shift]
	n -= nbytes
	copy(e.buf[n+shift:], e.buf[n:])
}

// Len returns the number of accumulated bytes
func (e *encoder) Len() int {
	return len(e.buf)
}

// Big endian (most significant byte first)

func (e *encoder) Uint16(n uint16) Encoder {
	return e.Put2(byte(n>>8), byte(n))
}

func (e *encoder) Uint32(n uint32) Encoder {
	return e.Put4(byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
}

func (e *encoder) VarUint(n uint64) Encoder {
	prevlen := len(e.buf)
	bytes := binary.PutUvarint(e.buf[prevlen:cap(e.buf)], n)
	e.buf = e.buf[:prevlen+bytes]
	return e
}

func (e *encoder) Flush() error {
	panic(assert.ShouldNotReachHere())
}