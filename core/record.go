// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"fmt"
	"strings"
)

/*
Record is an immutable record stored in a string
using the same format as cSuneido and jSuneido.

NOTE: This is the post 2019 format using a two byte header.

It is used for storing data records in the database and in dumps
and for transferring data records across the client-server protocol.

An empty Record is a single zero byte.

First two bytes are the type and the count of values, high two bits are the type
followed by the total length (uint8, uint16, or uint32)
followed by the offsets of the fields (uint8, uint16, or uint32)
followed by the contents of the fields
integers are stored big endian (most significant first)
*/
type Record string

const (
	type8 = iota + 1
	type16
	type32
)
const sizeMask = 0x3fff

const hdrlen = 2

// Count returns the number of values in the record
func (r Record) Count() int {
	if r == "" || r[0] == 0 {
		return 0
	}
	return (int(r[0])<<8 + int(r[1])) & sizeMask
}

func (r Record) Len() int {
	if r[0] == 0 {
		return 1
	}
	switch r.mode() {
	case type8:
		j := hdrlen
		return int(r[j])
	case type16:
		j := hdrlen
		return (int(r[j]) << 8) | int(r[j+1])
	case type32:
		j := hdrlen
		return (int(r[j]) << 24) | (int(r[j+1]) << 16) |
			(int(r[j+2]) << 8) | int(r[j+3])
	default:
		panic("invalid record type")
	}
}

func RecLen(r []byte) int {
	if r[0] == 0 {
		return 1
	}
	switch r[0] >> 6 {
	case type8:
		j := hdrlen
		return int(r[j])
	case type16:
		j := hdrlen
		return (int(r[j]) << 8) | int(r[j+1])
	case type32:
		j := hdrlen
		return (int(r[j]) << 24) | (int(r[j+1]) << 16) |
			(int(r[j+2]) << 8) | int(r[j+3])
	default:
		panic("invalid record type")
	}
}

// GetVal is a convenience method to get and unpack
func (r Record) GetVal(i int) Value {
	return Unpack(r.GetRaw(i))
}

// GetStr is a more direct method to get a packed string
func (r Record) GetStr(i int) string {
	s := r.GetRaw(i)
	if s == "" {
		return ""
	}
	if s[0] != PackString {
		panic("Record GetStr not string")
	}
	return s[1:]
}

// GetRaw returns one of the (usually packed) values
func (r Record) GetRaw(i int) string {
	if i < 0 || r.Count() <= i {
		return ""
	}
	var pos, end int
	switch r.mode() {
	case type8:
		j := hdrlen + i
		end = int(r[j])
		pos = int(r[j+1])
	case type16:
		j := hdrlen + 2*i
		end = (int(r[j]) << 8) | int(r[j+1])
		pos = (int(r[j+2]) << 8) | int(r[j+3])
	case type32:
		j := hdrlen + 4*i
		end = (int(r[j]) << 24) | (int(r[j+1]) << 16) |
			(int(r[j+2]) << 8) | int(r[j+3])
		pos = (int(r[j+4]) << 24) | (int(r[j+5]) << 16) |
			(int(r[j+6]) << 8) | int(r[j+7])
	default:
		panic("invalid record type")
	}
	return string(r)[pos:end]
}

func (r Record) mode() byte {
	return r[0] >> 6
}

func (r Record) String() string {
	if r == "" {
		return "{}"
	}
	var sb strings.Builder
	sep := "{"
	for i := range r.Count() {
		sb.WriteString(sep)
		sep = ", "
		sb.WriteString(r.GetVal(i).String())
	}
	sb.WriteString("}")
	return sb.String()
}

// Truncate shortens records to n fields.
// If the record has n or less fields it is returned unchanged.
// Otherwise it builds a new record (and trims it)
// It is used by UpdateTran output and update.
func (r Record) Truncate(n int) Record {
	rn := r.Count()
	if rn <= n {
		return r
	}
	var rb RecordBuilder
	for i := range n {
		rb.AddRaw(r.GetRaw(i))
	}
	return rb.Trim().Build()
}

// ------------------------------------------------------------------

// RecordBuilder is used to construct records. Zero value is ready to use.
type RecordBuilder struct {
	vals []Packable
}

const MaxValues = 0x3fff

// Add appends a Packable
func (b *RecordBuilder) Add(p Packable) *RecordBuilder {
	b.vals = append(b.vals, p)
	return b
}

// AddRaw appends a string containing an already packed value
func (b *RecordBuilder) AddRaw(s string) *RecordBuilder {
	if s == "" {
		b.Add(SuStr(""))
	} else {
		b.Add(Packed(s))
	}
	return b
}

// Packed is a Packable wrapper for an already packed value
type Packed string

var _ Packable = (*Packed)(nil)

func (p Packed) PackSize(*packing) int {
	return len(p)
}

func (p Packed) Pack(pk *packing) {
	pk.PutStr(string(p))
}

// Trim removes trailing empty fields
func (b *RecordBuilder) Trim() *RecordBuilder {
	n := len(b.vals)
	for n > 0 && b.vals[n-1] == SuStr("") {
		n--
	}
	b.vals = b.vals[:n]
	return b
}

// Build

const maxRecordLen = 1_000_000

func (b *RecordBuilder) Build() Record {
	pk := &packing{}
	if len(b.vals) > MaxValues {
		panic("too many values for record")
	}
	if len(b.vals) == 0 {
		return Record("\x00")
	}
	sizes := make([]int, len(b.vals))
	for i, v := range b.vals {
		sizes[i] = v.PackSize(pk)
	}
	length := b.recSize(sizes)
	if length > maxRecordLen {
		panic(fmt.Sprintf("record too large (%d > %d)", length, maxRecordLen))
	}
	*pk = *newPacking(length)
	b.build(pk, length, sizes)
	//assert.That(len(buf.String()) == length)
	return Record(pk.String())
}

func (b *RecordBuilder) recSize(sizes []int) int {
	nfields := len(b.vals)
	datasize := 0
	for _, size := range sizes {
		datasize += size
	}
	return tblength(nfields, datasize)
}

func tblength(nfields, datasize int) int {
	if nfields == 0 {
		return 1
	}
	length := hdrlen + (1 + nfields) + datasize
	if length < 0x100 {
		return length
	}
	length = hdrlen + 2*(1+nfields) + datasize
	if length < 0x10000 {
		return length
	}
	return hdrlen + 4*(1+nfields) + datasize
}

func (b *RecordBuilder) build(pk *packing, length int, sizes []int) {
	b.buildHeader(pk, length, sizes)
	nfields := len(b.vals)
	for i := nfields - 1; i >= 0; i-- {
		b.vals[i].Pack(pk)
	}
}

func (b *RecordBuilder) buildHeader(pk *packing, length int, sizes []int) {
	mode := mode(length)
	nfields := len(b.vals)
	pk.Uint16(uint16(mode<<14 | nfields))
	b.buildOffsets(pk, length, sizes)
}

func (b *RecordBuilder) buildOffsets(pk *packing, length int, sizes []int) {
	nfields := len(b.vals)
	offset := length
	switch mode(length) {
	case type8:
		pk.Put1(byte(offset))
		for i := range nfields {
			offset -= sizes[i]
			pk.Put1(byte(offset))
		}
	case type16:
		pk.Uint16(uint16(offset))
		for i := range nfields {
			offset -= sizes[i]
			pk.Uint16(uint16(offset))
		}
	case type32:
		pk.Uint32(uint32(offset))
		for i := range nfields {
			offset -= sizes[i]
			pk.Uint32(uint32(offset))
		}
	}
}

func mode(length int) int { // length must include header
	if length == 0 {
		return 0
	} else if length < 0x100 {
		return type8
	} else if length < 0x10000 {
		return type16
	} else {
		return type32
	}
}
