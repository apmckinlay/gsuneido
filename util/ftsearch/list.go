// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ftsearch

import (
	"math"

	"github.com/apmckinlay/gsuneido/util/assert"
)

// list stores counts per document as a list.
// This is more compact than sparse array.
type list struct {
	docIds []uint16
	counts []uint8
}

func (t *term) toList() list {
	a := t.termCountPerDoc.(*array)
	docIds := make([]uint16, t.ndocsWithTerm)
	counts := make([]uint8, t.ndocsWithTerm)
	j := 0
	for i, c := range a.counts {
		if c != 0 {
			docIds[j] = uint16(i)
			counts[j] = c
			j++
		}
	}
	assert.This(j).Is(t.ndocsWithTerm)
	return list{docIds: docIds, counts: counts}
}

func (cl list) add(int, int) bool {
	panic(assert.ShouldNotReachHere())
}

func (cl list) iterator() termIter {
	i := -1
	return func() (int, uint8) {
		i++
		if i < len(cl.counts) {
			return int(cl.docIds[i]), cl.counts[i]
		}
		return math.MaxInt, 0
	}
}

func (cl list) pack(buf []byte) []byte {
	buf = packUint32(buf, len(cl.counts))
	for _, d := range cl.docIds {
		buf = packUint16(buf, d)
	}
	for _, c := range cl.counts {
		buf = append(buf, c)
	}
	return buf
}

func packUint16(buf []byte, n uint16) []byte {
	return append(buf, byte(n>>8), byte(n))
}

func unpackList[T bytes](buf T) (list, T) {
	n := unpackUint32(buf)
    buf = buf[4:]
    docIds := make([]uint16, n)
	for i := range docIds {
        docIds[i] = unpackUint16(buf)
		buf = buf[2:]
    }
    counts := make([]uint8, n)
	for i := range counts {
        counts[i] = buf[0]
		buf = buf[1:]
    }
    return list{docIds: docIds, counts: counts}, buf
}

func unpackUint16[T bytes](buf T) uint16 {
	return uint16(buf[0])<<8 | uint16(buf[1])
}
