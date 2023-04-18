// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ftsearch

import (
	"fmt"
	"math"
	"strings"

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
	// ndocsWithTerm is the number of documents with this term
	ndocsWithTerm int
	// termCountPerDoc is how many times the term appears in each document
	termCountPerDoc termCountsPerDoc
}

// termCountsPerDoc is the interface for array and list
type termCountsPerDoc interface {
	// add returns whether the count was previously zero
	add(docId, n int) bool
	// del returns whether the count becomes zero
	del(docId, n int) bool
	pack([]byte) []byte
	iterator() termIter // used by search
}

func (t *term) add(docId, n int) {
	if t.termCountPerDoc.add(docId, n) {
		t.ndocsWithTerm++
	}
}

func (t *term) del(docId, n int) {
	if t.termCountPerDoc.del(docId, n) {
		t.ndocsWithTerm--
	}
}

type termIter func() (docId int, count uint8)

func (t *term) iterator() termIter {
	return t.termCountPerDoc.iterator()
}

func (ix *Index) Pack() []byte {
	buf := make([]byte, 0, 1024)
	buf = packUint32(buf, ix.ndocsTotal)
	buf = packUint32(buf, ix.ntermsTotal)
	buf = packUint32(buf, len(ix.ntermsPerDoc))
	for _, nd := range ix.ntermsPerDoc {
		buf = packUint32(buf, nd)
	}
	buf = packUint32(buf, len(ix.terms))
	for k, t := range ix.terms {
		assert.That(k == t.term)
		buf = packString(buf, t.term)
		buf = packUint32(buf, t.ndocsWithTerm)
		buf = t.termCountPerDoc.pack(buf)
	}
	return buf
}

func packUint32(buf []byte, n int) []byte {
	assert.That(n >= 0 && n <= math.MaxUint32)
	return append(buf, byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
}

func packString(buf []byte, s string) []byte {
	buf = append(buf, byte(len(s)))
	return append(buf, s...)
}

type bytes interface {
	~string | ~[]byte
}

func Unpack[T bytes](buf T) *Index {
	var idx Index
	idx.ndocsTotal, buf = unpackUint32(buf)
	idx.ntermsTotal, buf = unpackUint32(buf)
	n, buf := unpackUint32(buf)
	idx.ntermsPerDoc = make([]int, n)
	for i := range idx.ntermsPerDoc {
		idx.ntermsPerDoc[i], buf = unpackUint32(buf)
	}
	n, buf = unpackUint32(buf)
	idx.terms = make(map[string]*term, n)
	for i := 0; i < n; i++ {
		var t term
		t.term, buf = unpackString(buf)
		t.ndocsWithTerm, buf = unpackUint32(buf)
		t.termCountPerDoc, buf = unpackList(buf)
		idx.terms[t.term] = &t
	}
	return &idx
}

func unpackUint32[T bytes](buf T) (int, T) {
	return int(buf[0])<<24 | int(buf[1])<<16 | int(buf[2])<<8 | int(buf[3]),
		buf[4:]
}

func unpackString[T bytes](buf T) (string, T) {
	n := buf[0]
	return string(buf[1 : n+1]), buf[n+1:]
}

func (ix *Index) String() string {
	totalDocsPerTerm := 0
	maxDocsPerTerm := 0
	for _, t := range ix.terms {
		totalDocsPerTerm += t.ndocsWithTerm
		if t.ndocsWithTerm > maxDocsPerTerm {
			maxDocsPerTerm = t.ndocsWithTerm
		}
	}
	avgDocsPerTerm := 0
	if len(ix.terms) > 0 {
		avgDocsPerTerm = totalDocsPerTerm / len(ix.terms)
	}
	avgTermsPerDoc := 0
	if ix.ndocsTotal > 0 {
		avgTermsPerDoc = ix.ntermsTotal / ix.ndocsTotal
	}
	return fmt.Sprint("Ftsearch{ndocs: ", ix.ndocsTotal,
		", ntermsTotal: ", ix.ntermsTotal,
		",\n\tavgTermsPerDoc: ", avgTermsPerDoc,
		", avgDocsPerTerm: ", avgDocsPerTerm,
		", maxDocsPerTerm: ", maxDocsPerTerm, "}")
}

func (ix *Index) WordInfo(s string) string {
	input := newInput(s)
	var sb strings.Builder
	for term := input.Next(); term != ""; term = input.Next() {
		if t, ok := ix.terms[term]; ok {
			n := 0
			it := t.iterator()
			for _, count := it(); count != 0; _, count = it() {
				n += int(count)
			}
			sb.WriteString(fmt.Sprintln(term, "occurs", n,
				"times in", t.ndocsWithTerm, "documents"))
		} else {
			sb.WriteString(fmt.Sprintln(term, "not found"))
		}
	}
	return sb.String()
}

// update ---------------------------------------------------------

func (ix *Index) Update(id int, oldTitle, oldText, newTitle, newText string) {
	b := Builder{Index: *ix}
	if len(oldTitle)+len(oldText) > 0 {
		assert.That(ix.ndocsTotal > 0)
		b.Delete(id, oldTitle, oldText)
	}
	if len(newTitle)+len(newText) > 0 {
		b.Add(id, newTitle, newText)
	}
	b.backToList()
	*ix = b.Index
}
