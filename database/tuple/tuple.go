package tuple

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

/*
Tuple is an immutable tuple stored in a string
using the same format as cSuneido and jSuneido.

NOTE: This is the post 2019 format using a two byte header.

It it used in the database for storing data records in the database
It is also used for transferring data records
across the client-server protocol.

An empty Tuple is a single zero byte.

First two bytes are the type and the count of values, high two bits are the type
followed by the total length (uint8, uint16, or uint32)
followed by the offsets of the fields (uint8, uint16, or uint32)
followed by the contents of the fields
integers are stored big endian (most significant first)
*/
type Tuple string

const (
	type8 = iota + 1
	type16
	type32
)
const sizeMask = 0x3ff

const hdrlen = 2

// Count returns the number of values in the tuple
func (t Tuple) Count() int {
	return (int(t[0])<<8 + int(t[1])) & sizeMask
}

// GetVal is a convenience method to get and unpack
func (t Tuple) GetVal(i int) Value {
	return Unpack(t.GetRaw(i))
}

// Get returns one of the (usually packed) values
func (t Tuple) GetRaw(i int) string {
	var pos, end int
	switch t.mode() {
	case type8:
		j := hdrlen + i
		end = int(t[j])
		pos = int(t[j+1])
	case type16:
		j := hdrlen + 2*i
		end = (int(t[j]) << 8) | int(t[j+1])
		pos = (int(t[j+2]) << 8) | int(t[j+3])
	case type32:
		j := hdrlen + 4*i
		end = (int(t[j]) << 24) | (int(t[j+1]) << 16) |
			(int(t[j+2]) << 8) | int(t[j+3])
		pos = (int(t[j+4]) << 24) | (int(t[j+5]) << 16) |
			(int(t[j+6]) << 8) | int(t[j+7])
	default:
		panic("invalid record type")
	}
	return string(t)[pos:end]
}

func (t Tuple) mode() byte {
	return t[0] >> 6
}
