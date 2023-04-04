// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ftsearch

import (
	"fmt"
	"math"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func NewIndex() *Index {
	return &Index{terms: make(map[string]*term)}
}

type Index struct {
	ndocsTotal   int
	ntermsTotal  int
	ntermsPerDoc []int
	terms        map[string]*term
}

type term struct {
	term string
	// termCount is the total number of times
	// this term appears in all the documents
	termCount int
	// ndocsWithTerm is the number of documents with this term
	ndocsWithTerm int
	// termCountPerDoc is how many times the term appears in each document
	termCountPerDoc termCountsPerDoc
}

type termCountsPerDoc interface {
	add(docId, n int) bool // used by build
	pack([]byte) []byte
	iterator() termIter // used by search
}

func (t *term) add(docId, n int) {
	if t.termCountPerDoc.add(docId, n) {
		t.ndocsWithTerm++
	}
}

type termIter func() (docId int, count uint8)

func (t *term) iterator() termIter {
	return t.termCountPerDoc.iterator()
}

func (b *Index) Pack() []byte {
	buf := make([]byte, 0, 1024)
	buf = packUint32(buf, b.ndocsTotal)
	buf = packUint32(buf, len(b.ntermsPerDoc))
	for _, nd := range b.ntermsPerDoc {
		buf = packUint32(buf, nd)
	}
	buf = packUint32(buf, len(b.terms))
	for _, t := range b.terms {
		buf = append(buf, byte(len(t.term)))
		buf = append(buf, t.term...)
		buf = packUint32(buf, t.termCount)
		buf = packUint32(buf, t.ndocsWithTerm)
		buf = t.termCountPerDoc.pack(buf)
	}
	return buf
}

func packUint32(buf []byte, n int) []byte {
	assert.That(n >= 0 && n <= math.MaxUint32)
	return append(buf, byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
}

type bytes interface {
	~string | ~[]byte
}

func Unpack[T bytes](buf T) *Index {
	var idx Index
	idx.ndocsTotal = unpackUint32(buf)
	buf = buf[4:]
	n := unpackUint32(buf)
	buf = buf[4:]
	idx.ntermsPerDoc = make([]int, n)
	for i := range idx.ntermsPerDoc {
		idx.ntermsPerDoc[i] = unpackUint32(buf)
		buf = buf[4:]
	}
	idx.ntermsTotal = unpackUint32(buf)
	buf = buf[4:]
	idx.terms = make(map[string]*term, idx.ntermsTotal)
	for i := 0; i < idx.ntermsTotal; i++ {
		var t term
		n := buf[0]
		t.term = string(buf[1 : n+1])
		buf = buf[n+1:]
		t.termCount = unpackUint32(buf)
		buf = buf[4:]
		t.ndocsWithTerm = unpackUint32(buf)
		buf = buf[4:]
		t.termCountPerDoc, buf = unpackList(buf)
		idx.terms[t.term] = &t
	}
	return &idx
}

func unpackUint32[T bytes](buf T) int {
	return int(buf[0])<<24 | int(buf[1])<<16 | int(buf[2])<<8 | int(buf[3])
}

func (b *Index) String() string {
	totalDocsPerTerm := 0
	maxDocsPerTerm := 0
	for _, t := range b.terms {
		totalDocsPerTerm += t.ndocsWithTerm
		if t.ndocsWithTerm > maxDocsPerTerm {
			maxDocsPerTerm = t.ndocsWithTerm
		}
	}
	avgDocsPerTerm := totalDocsPerTerm / len(b.terms)
	return fmt.Sprint("ndocs: ", b.ndocsTotal,
		", ntermsTotal: ", b.ntermsTotal,
		", unique terms: ", len(b.terms),
		", avgTermsPerDoc: ", b.ntermsTotal/b.ndocsTotal,
		",\n\tavgDocsPerTerm: ", avgDocsPerTerm,
		", maxDocsPerTerm: ", maxDocsPerTerm,
	)
}
