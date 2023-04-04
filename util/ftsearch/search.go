// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

/*
Package search implements a simple full text search like Lucene.

Terminology:
- document e.g. a page in the help or wiki
- term is roughly a word
*/
package ftsearch

func (idx *Index) Search(query string) []DocScore {
	ts := make([]*term, 0, 8)
	input := newInput(query)
	for word := input.Next(); word != ""; word = input.Next() {
		if trm, ok := idx.terms[word]; ok {
			ts = append(ts, trm)
		}
	}
	avgtermsPerDoc := idx.ntermsTotal / idx.ndocsTotal
	return scoreTerms(idx.ndocsTotal, avgtermsPerDoc, idx.ntermsPerDoc,
		ts, bm25)
}
