package tuple

import (
	"math"

	. "github.com/apmckinlay/gsuneido/base"
	"github.com/apmckinlay/gsuneido/util/verify"
)

/*
TupleM is an in-memory mutable tuple.

Mutability is limited to appending values.
*/
type TupleM struct {
	offs []int
	buf  []byte
}

var _ Tuple = TupleM{}

func (t TupleM) Size() int {
	return len(t.offs)
}

func (t TupleM) Get(i int) Value {
	return get(t, i)
}

func (t TupleM) GetRaw(i int) []byte {
	off := t.offs[i]
	return t.buf[off : off+t.fieldLength(i)]
}

func (t TupleM) Compare(t2 Tuple) int {
	return compare(t, t2)
}

func (t TupleM) fieldLength(i int) int {
	var end int
	if i+1 < len(t.offs) {
		end = t.offs[i+1]
	} else {
		end = len(t.buf)
	}
	return end - t.offs[i]
}

// AddRaw appends a packed value to the tuple
func (t *TupleM) AddRaw(raw []byte) {
	t.offs = append(t.offs, len(t.buf))
	t.buf = append(t.buf, raw...)
}

// Add appends a value to the tuple
func (t *TupleM) Add(val Packable) {
	t.offs = append(t.offs, len(t.buf))
	t.buf = Ensure(t.buf, val.PackSize())
	t.buf = val.Pack(t.buf)
}

func (t TupleM) ToTupleB() TupleB {
	length := t.packSize()
	buf := make([]byte, 0, length)
	return TupleB(t.pack(buf, length))
}

func (t TupleM) packSize() int {
	nfields := len(t.offs)
	datasize := len(t.buf)
	return tblength(nfields, datasize)
}

func tblength(nfields int, datasize int) int {
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

func (t TupleM) pack(dst []byte, length int) []byte {
	dst = t.packHeader(dst, length)
	nfields := len(t.offs)
	for i := nfields - 1; i >= 0; i-- {
		dst = append(dst, t.GetRaw(i)...)
	}
	return dst
}

func (t TupleM) packHeader(dst []byte, length int) []byte {
	mode := mode(length)
	dst = append(dst, mode, 0)
	nfields := len(t.offs)
	verify.That(nfields <= math.MaxInt16)
	dst = pack16(dst, nfields)
	dst = t.packOffsets(dst, length, mode)
	return dst
}

func pack16(dst []byte, n int) []byte {
	return append(dst, byte(n), byte(n>>8))
}

func pack32(dst []byte, n int) []byte {
	return append(dst, byte(n), byte(n>>8), byte(n>>16), byte(n>>24))
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

func (t TupleM) packOffsets(dst []byte, length int, mode byte) []byte {
	nfields := len(t.offs)
	offset := length
	switch mode {
	case 'c':
		dst = append(dst, byte(offset))
		for i := 0; i < nfields; i++ {
			offset -= t.fieldLength(i)
			dst = append(dst, byte(offset))
		}
	case 's':
		dst = pack16(dst, offset)
		for i := 0; i < nfields; i++ {
			offset -= t.fieldLength(i)
			dst = pack16(dst, offset)
		}
	case 'l':
		dst = pack32(dst, offset)
		for i := 0; i < nfields; i++ {
			offset -= t.fieldLength(i)
			dst = pack32(dst, offset)
		}
	}
	return dst
}
