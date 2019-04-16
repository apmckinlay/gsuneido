package record

import (
	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/pack"
	"github.com/apmckinlay/gsuneido/util/verify"
)

// Builder is used to construct records
// This is the revised format (2019)
type Builder struct {
	vals []rt.Packable
}

const MaxValues = 0x3fff

// Add appends a Packable
func (b *Builder) Add(p rt.Packable) *Builder {
	b.vals = append(b.vals, p)
	return b
}

// AddRaw appends a string containing an already packed value
func (b *Builder) AddRaw(s string) *Builder {
	b.Add(Packed(s))
	return b
}

// Packed is a Packable wrapper for an already packed value
type Packed string

func (p Packed) Pack(buf *pack.Encoder) {
	buf.PutStr(string(p))
}

func (p Packed) PackSize(int) int {
	return len(p)
}

// Build

func (b *Builder) Build() Record {
	if len(b.vals) > MaxValues {
		panic("too many values for record")
	}
	if len(b.vals) == 0 {
		return Record("\x00")
	}
	sizes := make([]int, len(b.vals))
	for i, v := range b.vals {
		sizes[i] = v.PackSize(0)
	}
	length := b.recSize(sizes)
	buf := pack.NewEncoder(length)
	b.build(buf, length, sizes)
	verify.That(len(buf.String()) == length) //TODO remove
	return Record(buf.String())
}

func (b *Builder) recSize(sizes []int) int {
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

func (b *Builder) build(dst *pack.Encoder, length int, sizes []int) {
	b.buildHeader(dst, length, sizes)
	nfields := len(b.vals)
	for i := nfields - 1; i >= 0; i-- {
		b.vals[i].Pack(dst)
	}
}

func (b *Builder) buildHeader(dst *pack.Encoder, length int, sizes []int) {
	mode := mode(length)
	nfields := len(b.vals)
	dst.Uint16(uint16(mode<<14 | nfields))
	b.buildOffsets(dst, length, sizes)
}

func (b *Builder) buildOffsets(dst *pack.Encoder, length int, sizes []int) {
	nfields := len(b.vals)
	offset := length
	switch mode(length) {
	case type8:
		dst.Put1(byte(offset))
		for i := 0; i < nfields; i++ {
			offset -= sizes[i]
			dst.Put1(byte(offset))
		}
	case type16:
		dst.Uint16(uint16(offset))
		for i := 0; i < nfields; i++ {
			offset -= sizes[i]
			dst.Uint16(uint16(offset))
		}
	case type32:
		dst.Uint32(uint32(offset))
		for i := 0; i < nfields; i++ {
			offset -= sizes[i]
			dst.Uint32(uint32(offset))
		}
	}
}

func mode(length int) int {
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
