package tuple

import (
	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/pack"
)

// TupleBuilder is used to construct tuples
// This is the revised format (2019)
type TupleBuilder struct {
	vals []rt.Packable
}

const MaxValues = 0x3ff

// Add appends a Packable
func (tb *TupleBuilder) Add(p rt.Packable) *TupleBuilder {
	tb.vals = append(tb.vals, p)
	return tb
}

// AddRaw appends a string containing an already packed value
func (tb *TupleBuilder) AddRaw(s string) *TupleBuilder {
	tb.Add(Packed(s))
	return tb
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

func (tb *TupleBuilder) Build() Tuple {
	if len(tb.vals) > MaxValues {
		panic("too many values for tuple")
	}
	sizes := make([]int, len(tb.vals))
	for i, v := range tb.vals {
		sizes[i] = v.PackSize(0)
	}
	length := tb.tupleSize(sizes)
	buf := pack.NewEncoder(length)
	tb.build(buf, length, sizes)
	return Tuple(buf.String())
}

func (tb *TupleBuilder) tupleSize(sizes []int) int {
	nfields := len(tb.vals)
	datasize := 0
	for _, size := range sizes {
		datasize += size
	}
	return tblength(nfields, datasize)
}

func tblength(nfields, datasize int) int {
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

func (tb *TupleBuilder) build(dst *pack.Encoder, length int, sizes []int) {
	tb.buildHeader(dst, length, sizes)
	nfields := len(tb.vals)
	for i := nfields - 1; i >= 0; i-- {
		tb.vals[i].Pack(dst)
	}
}

func (tb *TupleBuilder) buildHeader(dst *pack.Encoder, length int, sizes []int) {
	mode := mode(length)
	nfields := len(tb.vals)
	dst.Uint16(uint16(mode<<14 | nfields))
	tb.buildOffsets(dst, length, sizes)
}

func (tb *TupleBuilder) buildOffsets(dst *pack.Encoder, length int, sizes []int) {
	nfields := len(tb.vals)
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
	if length < 0x100 {
		return type8
	} else if length < 0x10000 {
		return type16
	} else {
		return type32
	}
}
