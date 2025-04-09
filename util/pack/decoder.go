// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package pack

import (
	"encoding/binary"

	"github.com/apmckinlay/gsuneido/util/assert"
)

// Decoder is used to read values from a binary string (created with Encoder)
// It is somewhat similar to strings.Reader
type Decoder struct {
	s string
}

func MakeDecoder(s string) Decoder {
	return Decoder{s: s}
}

func (d Decoder) Peek() byte {
	return d.s[0]
}

func (d *Decoder) Skip(n int) {
	d.s = d.s[n:]
}

func (d *Decoder) Get1() byte {
	c := d.s[0]
	d.s = d.s[1:]
	return c
}

func (d *Decoder) Get(n int) string {
	s := d.s[:n]
	d.s = d.s[n:]
	return s
}

func (d *Decoder) Remaining() int {
	return len(d.s)
}

func (d Decoder) Remainder() string {
	return d.s
}

func (d *Decoder) Slice(n int) Decoder {
	d2 := Decoder{s: d.s[:n]}
	d.s = d.s[n:]
	return d2
}

// Big endian (most significant byte first)

func (d *Decoder) Uint16() uint16 {
	n := uint16(d.s[0])<<8 | uint16(d.s[1])
	d.s = d.s[2:]
	return n
}

func (d *Decoder) Uint32() uint32 {
	n := uint32(d.s[0])<<24 | uint32(d.s[1])<<16 | uint32(d.s[2])<<8 | uint32(d.s[3])
	d.s = d.s[4:]
	return n
}

func (d *Decoder) VarUint() uint64 {
	n, bytes := binary.Uvarint([]byte(d.s))
	assert.That(bytes > 0)
	d.s = d.s[bytes:]
	return n
}
