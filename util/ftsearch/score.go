// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ftsearch

import (
	"container/heap"
	"math"
	"sort"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/slc"
)

type DocScore struct {
	DocId int
	Score float64
}

// scoreTerms returns the top scoring n documents for the given terms.
func scoreTerms(ndocsTotal, avgDocLen int, docLen []int, terms []*term,
	scorefn func(int, int, int, int, uint8) float64) []DocScore {
	// merge iterate the terms, calculating a score for each document
	iters := make([]termIter, len(terms))
	docs := make([]int, len(terms))
	counts := make([]uint8, len(terms))
	for i := range terms {
		iters[i] = terms[i].iterator()
		docs[i], counts[i] = iters[i]()
	}
	var res results
	if len(docs) == 0 {
		return res
	}
	for {
		// find the next (minimum) doc
		d := slc.Min(docs)
		if d == math.MaxInt {
			break
		}
		var score float64
		for i := range terms {
			if docs[i] == d {
				score += scorefn(ndocsTotal, avgDocLen,
					terms[i].ndocsWithTerm, docLen[d], counts[i])
				docs[i], counts[i] = iters[i]()
			}
		}
		res.Add(DocScore{DocId: d, Score: score}, 20)
	}
	sort.Slice(res, func(i, j int) bool { return res[i].Score > res[j].Score })
	return res
}

// bm25 based on
// https://www.elastic.co/blog/practical-bm25-part-2-the-bm25-algorithm-and-its-variables

// defaults from Elasticsearch
const (
	b  float64 = .75
	k1 float64 = 1.2
)

func bm25(ndocsTotal, avgDocLen, ndocsWithTerm, docLen int, countOfTermInDoc uint8) float64 {
	return bm25score(b, k1, ndocsTotal, avgDocLen, ndocsWithTerm, docLen, countOfTermInDoc)
}

func bm25score(b, k1 float64, ndocsTotal, avgDocLen, ndocsWithTerm int, docLen int, countOfTermInDoc uint8) float64 {
	// fmt.Println("b", b, "k1", k1, "countOfTermInDoc", countOfTermInDoc, "docLen", docLen)
	idf := idf(ndocsTotal, ndocsWithTerm)
	// fmt.Println("idf", ndocsTotal, ndocsWithTerm, idf)
	x := float64(countOfTermInDoc) * (k1 + 1)
	y := float64(countOfTermInDoc) + k1*(1-b+b*float64(docLen)/float64(avgDocLen))
	result := idf * x / y
	// fmt.Println(idf, "*", x, "/", y, "=", result)
	return result
}

func idf(ndocsTotal, ndocsWithTerm int) float64 {
	// ndocsTotal is the total number of documents
	// ndocsWithTerm is number of documents that contain the term
	x := float64(ndocsTotal-ndocsWithTerm) + .5
	y := float64(ndocsWithTerm) + .5
	return math.Log(1 + x/y)
}

// results uses a heap to keep the top scores
type results []DocScore

func (r *results) Add(ds DocScore, keep int) {
	if len(*r) < keep {
		heap.Push(r, ds)
	} else if (*r)[0].Score < ds.Score {
		(*r)[0] = ds
		heap.Fix(r, 0)
	}
}

func (r results) Len() int           { return len(r) }
func (r results) Less(i, j int) bool { return r[i].Score < r[j].Score }
func (r results) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

func (r *results) Push(x any) {
	*r = append(*r, x.(DocScore))
}

func (r *results) Pop() any {
	panic(assert.ShouldNotReachHere())
}
