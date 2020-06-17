// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package metadata

type buffer []byte

func (b *buffer) put1(n int) *buffer {
	if n > 1<<8-1 {
		panic("buffer.put1 value too large")
	}
	*b = append(*b, byte(n))
	return b
}

func (b *buffer) put2(n int) *buffer {
	if n > 1<<16-1 {
		panic("buffer.put2 value too large")
	}
	*b = append(*b, byte(n), byte(n>>8))
	return b
}

func (b *buffer) put3(n int) *buffer {
	if n > 1<<24-1 {
		panic("buffer.put3 value too large")
	}
	*b = append(*b, byte(n), byte(n>>8), byte(n>>16))
	return b
}

func (b *buffer) put4(n int) *buffer {
	if n > 1<<32-1 {
		panic("buffer.put4 value too large")
	}
	*b = append(*b, byte(n), byte(n>>8), byte(n>>16), byte(n>>24))
	return b
}

func (b *buffer) put5(n uint64) *buffer {
	if n > 1<<40-1 {
		panic("buffer.put5 value too large")
	}
	*b = append(*b, byte(n), byte(n>>8), byte(n>>16), byte(n>>24), byte(n>>32))
	return b
}

func (b *buffer) get1() int {
	n := int((*b)[0])
	*b = (*b)[1:]
	return n
}

func (b *buffer) get2() int {
	n := int((*b)[0]) + int((*b)[1])<<8
	*b = (*b)[2:]
	return n
}

func (b *buffer) get3() int {
	n := int((*b)[0]) + int((*b)[1])<<8 + int((*b)[2])<<16
	*b = (*b)[3:]
	return n
}

func (b *buffer) get4() int {
	n := int((*b)[0]) + int((*b)[1])<<8 + int((*b)[2])<<16 + int((*b)[3])<<24
	*b = (*b)[4:]
	return n
}

func (b *buffer) get5() uint64 {
	n := uint64((*b)[0]) + uint64((*b)[1])<<8 + uint64((*b)[2])<<16 +
		uint64((*b)[3])<<24 + uint64((*b)[4])<<32
	*b = (*b)[5:]
	return n
}
