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

func (t *term) toList() {
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
	assert.That(j == t.ndocsWithTerm)
	t.termCountPerDoc = list{docIds: docIds, counts: counts}
}

func (t *term) toArray() {
	lst := t.termCountPerDoc.(list)
	assert.That(len(lst.docIds) == t.ndocsWithTerm)
	a := array{}
	for i, id := range lst.docIds {
		a.add(int(id), int(lst.counts[i]))
	}
	t.termCountPerDoc = &a
}

func (cl list) add(int, int) bool {
	panic(assert.ShouldNotReachHere())
}

func (cl list) del(int, int) bool {
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
	assert.That(len(cl.counts) == len(cl.docIds))
	buf = packUint16(buf, len(cl.counts))
	for _, d := range cl.docIds {
		buf = packUint16(buf, d)
	}
	buf = append(buf, cl.counts...)
	return buf
}

func packUint16[T int | uint16](buf []byte, n T) []byte {
	assert.That(n >= 0 && n <= math.MaxUint16)
	return append(buf, byte(n>>8), byte(n))
}

func unpackList[T bytes](buf T) (list, T) {
	n, buf := unpackUint16(buf)
	docIds := make([]uint16, n)
	for i := range docIds {
		docIds[i], buf = unpackUint16(buf)
	}
	counts := make([]uint8, n)
	for i := range counts {
		counts[i] = buf[0]
		buf = buf[1:]
	}
	return list{docIds: docIds, counts: counts}, buf
}

func unpackUint16[T bytes](buf T) (uint16, T) {
	return uint16(buf[0])<<8 | uint16(buf[1]), buf[2:]
}
