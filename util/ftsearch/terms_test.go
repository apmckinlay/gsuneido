// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ftsearch

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"golang.org/x/exp/slices"
)

func TestTerm_Iter(t *testing.T) {
	trm := term{termCountPerDoc: &array{[]uint8{5: 5, 11: 11, 13: 13}}}
	iter := trm.iterator()
	i, n := iter()
	assert.T(t).This(i).Is(5)
	assert.T(t).This(n).Is(5)
	i, n = iter()
	assert.T(t).This(i).Is(11)
	assert.T(t).This(n).Is(11)
	i, n = iter()
	assert.T(t).This(i).Is(13)
	assert.T(t).This(n).Is(13)
	i, n = iter()
	assert.T(t).This(i).Is(math.MaxInt)
	assert.T(t).This(n).Is(0)
}

func TestScoreTerms(t *testing.T) {
	trm1 := &term{termCountPerDoc: &array{[]uint8{5: 5, 11: 11, 13: 13}}}
	trm2 := &term{termCountPerDoc: &array{[]uint8{7: 7, 11: 11, 12: 12}}}
	docLens := make([]int, 14)
	docScores := scoreTerms(0, 0, docLens, []*term{trm1, trm2}, testScore)
	assert.T(t).This(fmt.Sprint(docScores)).
		Is("[{11 22} {13 13} {12 12} {7 7} {5 5}]")
}

func testScore(ndocsTotal, avgDocLen, docLen, termCount int, docCount uint8) float64 {
	return float64(docCount)
}

func TestIDF(t *testing.T) {
	// examples from article
	assertApprox(idf(4, 4), "0.10536")
	assertApprox(idf(4, 2), "0.69314")
}

func assertApprox(x float64, y string) {
	assert.Msg(x, y).That(strings.HasPrefix(fmt.Sprint(x), y))
}

func TestResults(t *testing.T) {
	var r results
	const keep = 4
	const n = 50
	scores := make([]float64, 5)
	for i := range scores {
		scores[i] = rand.Float64()
		r.Add(DocScore{Score: scores[i]}, keep)
	}
	assert.This(len(r)).Is(keep)
	slices.Sort(scores)
	sort.Sort(r)
	for i, s := range scores[len(scores)-keep:] {
		assert.This(r[i].Score).Is(s)
	}
}
