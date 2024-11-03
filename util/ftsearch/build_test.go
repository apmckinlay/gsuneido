// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ftsearch

import (
	"fmt"
	"sort"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestBuild(t *testing.T) {
	b := NewBuilder()
	b.Add(1, "Big Trees", "douglas fir")
	b.Add(2, "Small Trees", "apple pear not big")
	b.Add(3, "Pretty Flowers", "in spring time")
	// fmt.Println(b)
	packed := b.Pack()
	assert.T(t).This(len(packed)).Is(199)
	idx := Unpack(packed)

	assert.T(t).This(idx.ndocsTotal).Is(3)
	assert.T(t).This(idx.ntermsTotal).Is(13)

	trm := idx.terms["big"]
	assert.T(t).This(trm.ndocsWithTerm).Is(2)
	assert.T(t).This(len(idx.terms)).Is(11)

	test := func(words ...string) func(any) {
		t.Helper()
		avgtermsPerDoc := idx.ntermsTotal / idx.ndocsTotal
		ts := make([]*term, 0, len(words))
		for _, w := range words {
			if t, ok := idx.terms[w]; ok {
				ts = append(ts, t)
			}
		}
		results := scoreTerms(idx.ndocsTotal, avgtermsPerDoc, idx.ntermsPerDoc,
			ts, testScore)
		return assert.T(t).This(fmt.Sprint(results)).Is
	}
	test("nada")("[]")
	test("fir")("[{1 1}]")
	test("big")("[{1 3} {2 1}]")
	test("big", "pear")("[{1 3} {2 2}]")
	test("flower")("[{3 3}]")

	idx.Update(4, "", "", "New One", "nothing special apple") // add
	assert.T(t).This(idx.ndocsTotal).Is(4)
	test("new", "special")("[{4 4}]")

	idx.Update(2, "Small Trees", "apple pear not big", // update
		"Small Shrubs", "all around the yard")
	assert.T(t).This(idx.ndocsTotal).Is(4)
	test("tree")("[{1 3}]")
	test("yard")("[{2 1}]")

	idx.Update(3, "Pretty Flowers", "in spring time", "", "") // delete
	assert.T(t).This(idx.ndocsTotal).Is(3)
	assert.T(t).This(len(idx.terms)).Is(13)
	test("flower")("[]")
}

// following examples from
// https://www.elastic.co/blog/practical-bm25-part-2-the-bm25-algorithm-and-its-variables

func TestBm25(t *testing.T) {
	bldr := NewBuilder()
	bldr.Add(0, "", "Shane")
	bldr.Add(1, "", "Shane CC")
	bldr.Add(2, "", "Shane PP Connelly")
	bldr.Add(3, "", "Shane Connelly")
	bldr.Add(4, "", "Shane Shane Connelly Connelly")
	bldr.Add(5, "", "Shane Shane Shane Connelly Connelly Connelly")
	avgtermsPerDoc := bldr.ntermsTotal / bldr.ndocsTotal

	b := .5
	k1 := 0.
	bm25 := func(ndocsTotal, avgDocLen, ndocsWithTerm, docLen int, ntermsInDoc uint8) float64 {
		return bm25score(b, k1, ndocsTotal, avgDocLen, ndocsWithTerm, docLen, ntermsInDoc)
	}
	ts := []*term{bldr.terms["shane"]}
	results := scoreTerms(bldr.ndocsTotal, avgtermsPerDoc, bldr.ntermsPerDoc, ts, bm25)
	for _, r := range results {
		assertApprox(r.Score, "0.0741")
	}

	b = 0
	k1 = 10
	results = scoreTerms(bldr.ndocsTotal, avgtermsPerDoc, bldr.ntermsPerDoc, ts, bm25)
	sort.Slice(results, func(i, j int) bool {
		return results[i].DocId < results[j].DocId
	})
	for i := range 4 {
		assertApprox(results[i].Score, "0.0741")
	}
	assertApprox(results[4].Score, "0.1358")
	assertApprox(results[5].Score, "0.1881")

	b = 1
	k1 = 5
	results = scoreTerms(bldr.ndocsTotal, avgtermsPerDoc, bldr.ntermsPerDoc, ts, bm25)
	sort.Slice(results, func(i, j int) bool {
		return results[i].DocId < results[j].DocId
	})
	assertApprox(results[0].Score, "0.1667")
	assertApprox(results[1].Score, "0.1026")
	assertApprox(results[2].Score, "0.0741")
	assertApprox(results[3].Score, "0.1026")
	assertApprox(results[4].Score, "0.1026")
	assertApprox(results[5].Score, "0.1026")
}
