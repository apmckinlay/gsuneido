package tuple

import (
	v "github.com/apmckinlay/gsuneido/value"
)

/*
TupleB is an immutable tuple stored in a slice of bytes
using the same format as cSuneido.

It it used in the database for storing data records
and for btree nodes and the keys within them.
It is also used for transferring data records
across the client-server protocol.

[0] is 'c', 's', or 'l' for 8, 16, or 32 bit offsets
[1] is unused
[2,3] is the number of values
followed by the total length (uint8, uint16, or uint32)
followed by the offsets of the fields (uint8, uint16, or uint32)
followed by the contents of the fields
integers are stored little endian (least significant first)
*/
type TupleB []byte

var _ Tuple = TupleM{}

func (t TupleB) Size() int {
	return int(t[2]) + int(t[3])<<8
}

func (t TupleB) Get(i int) v.Value {
	return get(t, i)
}

func (t TupleB) GetRaw(i int) []byte {
	return t[t.offset(i):t.offset(i-1)]
}

func (x TupleB) Compare(y Tuple) int {
	return compare(x, y)
}

func (t TupleB) mode() byte {
	return t[0]
}

func (t TupleB) offset(i int) int {
	const body = 4
	switch t.mode() {
	case 'c':
		return int(t[body+i+1])
	case 's':
		si := body + 2*(i+1)
		return int(t[si]) + (int(t[si+1]) << 8)
	case 'l':
		ii := body + 4*(i+1)
		return int(t[ii]) |
			(int(t[ii+1]) << 8) |
			(int(t[ii+2]) << 16) |
			(int(t[ii+3]) << 24)
	default:
		panic("invalid record type: " + string(t.mode()))
	}
}
