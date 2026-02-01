// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ftsearch

import (
	"math"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/slc"
)

// Builder is an index using array for counts.
type Builder struct {
	Index
	// converted is used by update
	converted []*term
}

func NewBuilder() *Builder {
	return &Builder{Index: *NewIndex()}
}

const boost = 3

func (b *Builder) Add(id int, title, text string) {
	assert.That(id < math.MaxUint16)
	nterms := b.add(id, title, boost)
	nterms += b.add(id, text, 1)
	if id >= len(b.ntermsPerDoc) {
		b.ntermsPerDoc = slc.Allow(b.ntermsPerDoc, id+1)
	}
	b.ntermsPerDoc[id] = nterms
	b.ndocsTotal++
	b.ntermsTotal += nterms
}

func (b *Builder) add(id int, s string, n int) int {
	nterms := 0
	input := NewInput(s)
	for word := input.Next(); word != ""; word = input.Next() {
		trm, ok := b.terms[word]
		if !ok {
			trm = &term{term: word, termCountPerDoc: &array{}}
			b.terms[word] = trm
		} else if _, ok := trm.termCountPerDoc.(list); ok {
			trm.toArray()
			b.converted = append(b.converted, trm)
		}
		trm.add(id, n)
		nterms++
	}
	return nterms
}

func (b *Builder) Delete(id int, title, text string) {
	assert.That(id < math.MaxUint16)
	nterms := b.delete(id, title, boost)
	nterms += b.delete(id, text, 1)
	assert.That(id < len(b.ntermsPerDoc))
	b.ntermsPerDoc[id] = 0
	b.ndocsTotal--
	assert.That(b.ntermsTotal >= nterms)
	b.ntermsTotal -= nterms
}

func (b *Builder) delete(id int, s string, n int) int {
	nterms := 0
	input := NewInput(s)
	for word := input.Next(); word != ""; word = input.Next() {
		trm, ok := b.terms[word]
		assert.That(ok)
		if _, ok := trm.termCountPerDoc.(list); ok {
			trm.toArray()
			b.converted = append(b.converted, trm)
		}
		trm.del(id, n)
		if trm.ndocsWithTerm == 0 {
			delete(b.terms, word)
		}
		nterms++
	}
	return nterms
}

func (b *Builder) backToList() {
	for _, trm := range b.converted {
		trm.toList()
	}
}

func (b *Builder) ToIndex() *Index {
	for _, t := range b.Index.terms {
		t.toList()
	}
	return &b.Index
}

//-------------------------------------------------------------------

// array stores counts per document as a sparse array
// indexed by document id (small dense integers).
// This uses more memory but is faster to update during building.
type array struct {
	counts []uint8
}

func (a *array) add(id int, n int) bool {
	if id >= len(a.counts) {
		a.counts = slc.Allow(a.counts, id+1)
	}
	first := a.counts[id] == 0
	n += int(a.counts[id])
	if n > 255 {
		n = 255 // stick at max
	}
	a.counts[id] = uint8(n)
	return first
}

func (a *array) del(id int, n int) bool {
	assert.That(a.counts[id] >= uint8(n))
	if a.counts[id] < 255 { // stick at max
		a.counts[id] -= uint8(n)
	}
	return a.counts[id] == 0
}

func (a *array) iterator() termIter {
	i := 0
	return func() (int, uint8) {
		for ; i < len(a.counts); i++ {
			if a.counts[i] > 0 {
				i++
				return i - 1, a.counts[i-1]
			}
		}
		return math.MaxInt, 0
	}
}

func (a *array) pack(buf []byte) []byte {
	n := 0
	for _, c := range a.counts {
		if c > 0 {
			n++
		}
	}
	buf = packUint16(buf, n)
	for i, c := range a.counts {
		if c > 0 {
			buf = packUint16(buf, uint16(i))
		}
	}
	for _, c := range a.counts {
		if c > 0 {
			buf = append(buf, c)
		}
	}
	return buf
}
