// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package pack

import (
	"encoding/binary"

	"github.com/apmckinlay/gsuneido/util/assert"
)

// Big endian (most significant byte first)

func (e *Encoder) Uint16(n uint16) *Encoder {
	return e.Put2(byte(n>>8), byte(n))
}

func (d *Decoder) Uint16() uint16 {
	n := uint16(d.s[0])<<8 | uint16(d.s[1])
	d.s = d.s[2:]
	return n
}

func (e *Encoder) Int32(n int) *Encoder {
	// complement leading bit to ensure correct unsigned compare
	return e.Put4(byte(n>>24)^0x80, byte(n>>16), byte(n>>8), byte(n))
}

func (d *Decoder) Int32() int {
	n := int32(d.s[0]^0x80)<<24 | int32(d.s[1])<<16 | int32(d.s[2])<<8 | int32(d.s[3])
	d.s = d.s[4:]
	return int(n)
}

func (e *Encoder) Uint32(n uint32) *Encoder {
	return e.Put4(byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
}

func (d *Decoder) Uint32() uint32 {
	n := uint32(d.s[0])<<24 | uint32(d.s[1])<<16 | uint32(d.s[2])<<8 | uint32(d.s[3])
	d.s = d.s[4:]
	return n
}

func (e *Encoder) VarUint(n uint64) *Encoder {
	prevlen := len(e.buf)
	bytes := binary.PutUvarint(e.buf[prevlen:cap(e.buf)], n)
	e.buf = e.buf[:prevlen+bytes]
	return e
}

func (d *Decoder) VarUint() uint64 {
	n, bytes := binary.Uvarint([]byte(d.s))
	assert.That(bytes > 0)
	d.s = d.s[bytes:]
	return n
}
