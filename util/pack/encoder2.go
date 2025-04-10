// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package pack

import (
	"bufio"
	"encoding/binary"
	"io"

	"github.com/apmckinlay/gsuneido/util/assert"
)

// Encoder is an interface that abstracts the common functionality
// between Encoder and encoder2 implementations.
type Encoder interface {
	Dup() Encoder
	String() string
	Buffer() []byte
	Put(b []byte) Encoder
	Put1(b byte) Encoder
	Put2(a, b byte) Encoder
	Put4(a, b, c, d byte) Encoder
	PutStr(s string) Encoder
	Move(nbytes, shift int)
	Len() int
	Uint16(n uint16) Encoder
	Uint32(n uint32) Encoder
	VarUint(n uint64) Encoder
	Flush() error
}

var _ Encoder = (*encoder)(nil)
var _ Encoder = (*encoder2)(nil)

type encoder2 struct {
	w *bufio.Writer
}

// NewEncoder2 returns an encoder2 with a buffer of the specified capacity
func NewEncoder2(w io.Writer) Encoder {
	return &encoder2{w: bufio.NewWriter(w)}
}

func (e *encoder2) Put(b []byte) Encoder {
	e.w.Write(b)
	return e
}

func (e *encoder2) Put1(b byte) Encoder {
	e.w.WriteByte(b)
	return e
}

func (e *encoder2) Put2(a, b byte) Encoder {
	e.w.WriteByte(a)
	e.w.WriteByte(b)
	return e
}

func (e *encoder2) Put4(a, b, c, d byte) Encoder {
	e.w.WriteByte(a)
	e.w.WriteByte(b)
	e.w.WriteByte(c)
	e.w.WriteByte(d)
	return e
}

func (e *encoder2) PutStr(s string) Encoder {
	e.w.WriteString(s)
	return e
}

// Big endian (most significant byte first)

func (e *encoder2) Uint16(n uint16) Encoder {
	return e.Put2(byte(n>>8), byte(n))
}

func (e *encoder2) Uint32(n uint32) Encoder {
	return e.Put4(byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
}

func (e *encoder2) VarUint(n uint64) Encoder {
	var buf [binary.MaxVarintLen64]byte
	bytes := binary.PutUvarint(buf[:], n)
	for _, b := range buf[:bytes] {
		e.w.WriteByte(b)
	}
	return e
}

func (e *encoder2) Flush() error {
	return e.w.Flush()
}

// not implemented

func (e *encoder2) Dup() Encoder {
	panic(assert.ShouldNotReachHere())
}

func (e *encoder2) String() string {
	panic(assert.ShouldNotReachHere())
}

func (e *encoder2) Buffer() []byte {
	panic(assert.ShouldNotReachHere())
}

func (e *encoder2) Move(nbytes, shift int) {
	assert.ShouldNotReachHere()
}

func (e *encoder2) Len() int {
	panic(assert.ShouldNotReachHere())
}
