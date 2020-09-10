// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package stor

import (
	"github.com/apmckinlay/gsuneido/util/cksum"
)

func (s *Stor) SaveSized(data []byte) Offset {
	n := len(data)
	hn := recHdrLen(n)
	off, buf := s.Alloc(hn + n + cksum.Len)
	recHdr(buf, n)
	copy(buf[hn:], data)
	cksum.Update(buf)
	return off
}

func (s *Stor) DataSized(offset Offset) []byte {
	buf := s.Data(offset)
	return recData(buf)
}

// The following handles saving with a header compatible with runtime.Record
// so we can checksum redirects
// without differentiating data records and index nodes
//	- a single 0 byte means 0 payload length
//	- otherwise there is a two byte header followed by the size
//	- size includes the 2 byte header and 1, 2, or 4 bytes for size

const (
	type8 = iota + 1
	type16
	type32
)

func recData(r []byte) []byte {
	if r[0] == 0 {
		return []byte{}
	}
	const j = 2 // header length
	var hn, rn int
	switch r[0] >> 6 {
	case type8:
		hn, rn = 3, int(r[j])
	case type16:
		hn, rn = 4, (int(r[j])<<8)|int(r[j+1])
	case type32:
		hn, rn = 6, (int(r[j])<<24)|(int(r[j+1])<<16)|
			(int(r[j+2])<<8)|int(r[j+3])
	default:
		panic("invalid record type")
	}
	return r[hn:rn]
}

func recHdrLen(size int) int {
	switch {
	case size == 0:
		return 1
	case size+3 < 0x100:
		return 3
	case size+4 < 0x10000:
		return 4
	default:
		return 6
	}
}

// recHdr initializes buf to have a header compatible with runtime.Record
func recHdr(buf []byte, size int) {
	if size == 0 {
		buf[0] = 0
		return
	}
	buf[1] = 0 // no fields (doesn't occur in records since it would be 0 size)
	switch {
	case size+3 < 0x100:
		buf[0] = 1 << 6
		size += 3
		buf[2] = byte(size)
	case size+4 < 0x10000:
		buf[0] = 2 << 6
		size += 4
		buf[2] = byte(size >> 8)
		buf[3] = byte(size)
	default:
		buf[0] = 3 << 6
		size += 5
		buf[2] = byte(size >> 24)
		buf[3] = byte(size >> 16)
		buf[4] = byte(size >> 8)
		buf[5] = byte(size)
	}
}
