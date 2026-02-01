// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package set

import (
	"math/rand/v2"
	"slices"

	"github.com/apmckinlay/gsuneido/util/bits"
)

// NPerms calculates P(n, k) = n! / (n-k)!
func NPerms(n, k int) int {
	if k < 0 || k > n {
		panic("NPerms invalid arguments")
	}
	res := 1
	for i := range k {
		res *= n - i
		if res < 0 {
			panic("NPerms overflow")
		}
	}
	return res
}

// NthPerm returns the nth permutation of k elements from the given set
func NthPerm(elements []string, k int, index int) []string {
	n := len(elements)
	workingSet := slices.Clone(elements)

	result := make([]string, k)
	for i := range k {
		// How many ways can we fill the remaining (k - 1 - i) slots
		// using the remaining (n - 1 - i) elements?
		options := NPerms(n-1-i, k-1-i)

		selectionIdx := index / options
		result[i] = workingSet[selectionIdx]

		// Remove the selected element by shifting (preserving order)
		workingSet = append(workingSet[:selectionIdx], workingSet[selectionIdx+1:]...)

		index %= options
	}
	return result
}

func RandPerm(rnd *rand.Rand, elements []string, k int) []string {
	return NthPerm(elements, k, rnd.IntN(NPerms(len(elements), k)))
}

type Gen struct {
	elements []string
	k        int
	indexGen *bits.Gen
}

func NewGen(rnd *rand.Rand, elements []string, k int) *Gen {
	n := len(elements)
	nPerms := NPerms(n, k)
	return &Gen{
		elements: slices.Clone(elements),
		k:        k,
		indexGen: bits.NewGen(rnd, uint64(nPerms))}
}

func (g *Gen) Next() []string {
	return NthPerm(g.elements, g.k, int(g.indexGen.Next()))
}

