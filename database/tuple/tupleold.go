package tuple

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

/*
TupleOld is an immutable tuple stored in a string
using the same format as cSuneido and jSuneido.

NOTE: This is the pre 2019 format using a four byte header.

It it used in the database for storing data records in the database
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
type TupleOld string

// Count returns the number of values in the tuple
func (t TupleOld) Count() int {
	return int(t[2]) + int(t[3])<<8
}

// Get returns one of the (usually packed) values
func (t TupleOld) Get(i int) string {
	return string(t)[t.offset(i):t.offset(i-1)]
}

// GetVal is a convenience method to get and unpack
func (t TupleOld) GetVal(i int) Value {
	return Unpack(t.Get(i))
}

func (t TupleOld) mode() byte {
	return t[0]
}

func (t TupleOld) offset(i int) int {
	const hdr = 4
	switch t.mode() {
	case 'c':
		return int(t[hdr+i+1])
	case 's':
		si := hdr + 2*(i+1)
		return int(t[si]) + (int(t[si+1]) << 8)
	case 'l':
		ii := hdr + 4*(i+1)
		return int(t[ii]) |
			(int(t[ii+1]) << 8) |
			(int(t[ii+2]) << 16) |
			(int(t[ii+3]) << 24)
	default:
		panic("invalid record type: " + string(t.mode()))
	}
}
