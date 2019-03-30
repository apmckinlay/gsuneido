package tuple

import (
	"math"

	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/pack"
	"github.com/apmckinlay/gsuneido/util/verify"
)

/*
TODO: This is probably not the final approach.
If you're adding Packable's, then it's better to defer packing.
That's the case if we're creating a new record.
But if we're updating an existing record, and we've defered unpacking,
then many of the values will be already still packed.
And it probably makes more sense to build from an iterator
rather than constructing an intermediate data structure.
Also probably want to build directly into mmap buffer.
*/

// TupleBuilder is used to construct tuples
type TupleBuilder struct {
	strs []string
}

// Add appends a string (usually a packed value)
func (tb *TupleBuilder) Add(s string) *TupleBuilder {
	tb.strs = append(tb.strs, s)
	return tb
}

// AddVal is a convenience method that packs and adds
func (tb *TupleBuilder) AddVal(p rt.Packable) *TupleBuilder {
	return tb.Add(rt.Pack(p))
}

// Build

func (tb *TupleBuilder) Build() Tuple {
	length := tb.tupleSize()
	buf := pack.NewEncoder(length)
	tb.build(buf, length)
	return Tuple(buf.String())
}

func (tb *TupleBuilder) tupleSize() int {
	nfields := len(tb.strs)
	datasize := 0
	for _, s := range tb.strs {
		datasize += len(s)
	}
	return tblength(nfields, datasize)
}

func tblength(nfields, datasize int) int {
	length := 4 + (1 + nfields) + datasize
	if length < 0x100 {
		return length
	}
	length = 4 + 2*(1+nfields) + datasize
	if length < 0x10000 {
		return length
	}
	return 4 + 4*(1+nfields) + datasize
}

func (tb *TupleBuilder) build(dst *pack.Encoder, length int) {
	tb.buildHeader(dst, length)
	nfields := len(tb.strs)
	for i := nfields - 1; i >= 0; i-- {
		dst.PutStr(tb.strs[i])
	}
}

func (tb *TupleBuilder) buildHeader(dst *pack.Encoder, length int) {
	mode := mode(length)
	dst.Put2(mode, 0)
	nfields := len(tb.strs)
	verify.That(nfields <= math.MaxInt16)
	dst.Put2(byte(nfields), byte(nfields>>8))
	tb.buildOffsets(dst, length, mode)
}

func mode(length int) byte {
	if length < 0x100 {
		return 'c'
	} else if length < 0x10000 {
		return 's'
	} else {
		return 'l'
	}
}

func (tb *TupleBuilder) buildOffsets(dst *pack.Encoder, length int, mode byte) {
	nfields := len(tb.strs)
	verify.That(0 <= length && length <= math.MaxUint16)
	offset := length
	switch mode {
	case 'c':
		dst.Put1(byte(offset))
		for i := 0; i < nfields; i++ {
			offset -= len(tb.strs[i])
			dst.Put1(byte(offset))
		}
	case 's':
		putUint16(dst, offset)
		for i := 0; i < nfields; i++ {
			offset -= len(tb.strs[i])
			putUint16(dst, offset)
		}
	case 'l':
		putUint32(dst, offset)
		for i := 0; i < nfields; i++ {
			offset -= len(tb.strs[i])
			putUint32(dst, offset)
		}
	}
}

func putUint16(dst *pack.Encoder, n int) {
	// little endian
	dst.Put2(byte(n), byte(n>>8))
}

func putUint32(dst *pack.Encoder, n int) {
	// little endian
	dst.Put4(byte(n), byte(n>>8), byte(n>>16), byte(n>>24))
}
