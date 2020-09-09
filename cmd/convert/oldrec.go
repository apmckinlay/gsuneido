// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

import . "github.com/apmckinlay/gsuneido/runtime"

/*
OldRec is an immutable record stored in a string
using the pre 2019 format with a four byte header.

[0] is 'c', 's', or 'l' for 8, 16, or 32 bit offsets
[1] is unused
[2,3] is the number of values
followed by the total length (uint8, uint16, or uint32)
followed by the offsets of the fields (uint8, uint16, or uint32)
followed by the contents of the fields
integers are stored little endian (least significant first)
*/
type OldRec string

// Count returns the number of values in the record
func (o OldRec) Count() int {
	return int(o[2]) + int(o[3])<<8
}

// Get returns one of the (usually packed) values
func (o OldRec) Get(i int) string {
	return string(o)[o.offset(i):o.offset(i-1)]
}

// GetVal is a convenience method to get and unpack
func (o OldRec) GetVal(i int) Value {
	return Unpack(o.Get(i))
}

func (o OldRec) mode() byte {
	return o[0]
}

func (o OldRec) offset(i int) int {
	const hdr = 4
	switch o.mode() {
	case 'c':
		return int(o[hdr+i+1])
	case 's':
		si := hdr + 2*(i+1)
		return int(o[si]) + (int(o[si+1]) << 8)
	case 'l':
		ii := hdr + 4*(i+1)
		return int(o[ii]) |
			(int(o[ii+1]) << 8) |
			(int(o[ii+2]) << 16) |
			(int(o[ii+3]) << 24)
	default:
		panic("invalid record type: " + string(o.mode()))
	}
}
