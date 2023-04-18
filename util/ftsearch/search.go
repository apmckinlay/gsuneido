// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

/*
Package search implements a simple full text search like Lucene.

Terminology:
- document e.g. a page in the help or wiki
- term is roughly a word
*/
package ftsearch

func (ix *Index) Search(query string) []DocScore {
	if ix.ndocsTotal == 0 {
		return nil
	}
	ts := make([]*term, 0, 8)
	input := newInput(query)
	for term := input.Next(); term != ""; term = input.Next() {
		if trm, ok := ix.terms[term]; ok {
			ts = append(ts, trm)
		}
	}
	avgtermsPerDoc := ix.ntermsTotal / ix.ndocsTotal
	return scoreTerms(ix.ndocsTotal, avgtermsPerDoc, ix.ntermsPerDoc,
		ts, bm25)
}
