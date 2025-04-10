// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package pack

import (
	"bufio"
	"encoding/binary"
	"io"
	"strings"

	"github.com/apmckinlay/gsuneido/util/hacks"
)

type Decoder interface {
	// Peek returns the next byte without advancing
	Peek() byte

	// Skip advances the decoder by n bytes
	Skip(n int)

	// Get1 returns the next byte and advances the decoder
	Get1() byte

	// Get returns the next n bytes as a string and advances the decoder
	Get(n int) string

	// Remaining returns true if there are bytes left to read
	Remaining() bool

	// Remainder returns the remaining bytes as a string
	Remainder() string

	// Remainder returns the remaining bytes as a string
	TempRemainder() string

	// TempStr returns a string that must not be held after the next read
	TempStr(n int) string

	// Slice
	Slice(n int) Decoder

	// Uint16 reads a 2-byte big-endian unsigned integer
	Uint16() uint16

	// Uint32 reads a 4-byte big-endian unsigned integer
	Uint32() uint32

	// VarUint reads a variable-length encoded unsigned integer
	VarUint() uint64
}

var _ Decoder = (*decoder)(nil)
var _ Decoder = (*decoder2)(nil)

// decoder2 is used to read values from a binary string (created with Encoder)
// It is somewhat similar to strings.Reader
type decoder2 struct {
	r bufio.Reader
}

func NewDecoder2(r io.Reader) *decoder2 {
	return &decoder2{r: *bufio.NewReader(r)}
}

func (d *decoder2) Peek() byte {
	bs, _ := d.r.Peek(1)
	return bs[0]
}

func (d *decoder2) Skip(n int) {
	d.r.Discard(n)
}

func (d *decoder2) Get1() byte {
	b, _ := d.r.ReadByte()
	return b
}

func (d *decoder2) Get(n int) string {
	buf := make([]byte, n)
	n, _ = io.ReadFull(&d.r, buf)
	return hacks.BStoS(buf[:n])
}

func (d *decoder2) Remaining() bool {
	bs, err := d.r.Peek(1)
	return err == nil && len(bs) == 1
}

func (d *decoder2) Remainder() string {
	var b strings.Builder
	io.Copy(&b, &d.r)
	return b.String()
}

func (d *decoder2) TempRemainder() string {
	size := d.r.Size()
	buf, _ := d.r.Peek(size)
	if len(buf) == size {
		panic("Unpack: TempRemainder too large")
	}
	return hacks.BStoS(buf)
}

func (d *decoder2) TempStr(n int) string {
	bs, _ := d.r.Peek(n)
	d.r.Discard(n)
	return hacks.BStoS(bs)
}

func (d *decoder2) Slice(n int) Decoder {
	return nil
}

// Big endian (most significant byte first)

func (d *decoder2) Uint16() uint16 {
	b1, _ := d.r.ReadByte()
	b2, _ := d.r.ReadByte()
	return uint16(b1)<<8 | uint16(b2)
}

func (d *decoder2) Uint32() uint32 {
	b1, _ := d.r.ReadByte()
	b2, _ := d.r.ReadByte()
	b3, _ := d.r.ReadByte()
	b4, _ := d.r.ReadByte()
	return uint32(b1)<<24 | uint32(b2)<<16 | uint32(b3)<<8 | uint32(b4)
}

func (d *decoder2) VarUint() uint64 {
	n, _ := binary.ReadUvarint(&d.r)
	return n
}
