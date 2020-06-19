// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package stor

import "encoding/binary"

func (s *Stor) AllocSized(n int) (Offset, []byte) {
	if n <= 0 {
		panic("stor.AllocSized bad size")
	}
	var uv [8]byte
	size := binary.PutUvarint(uv[:], uint64(n))
	off, buf := s.Alloc(n + size)
	copy(buf, uv[:size])
	return off, buf[size:]
}

func (s *Stor) DataSized(offset Offset) []byte {
	buf := s.Data(offset)
	n,size := binary.Uvarint(buf)
	if size <= 0 {
		panic("stor.DataSized bad size, corruption?")
	}
	return buf[size:uint64(size) + n]
}
