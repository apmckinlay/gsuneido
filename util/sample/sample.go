// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package sample

import (
	"iter"
	"math/rand/v2"

	"github.com/apmckinlay/gsuneido/util/assert"
)

// From returns a sequence of unique pseudo-random integers in [0, m).
// The sequence yields each value at most once and stops when the caller stops
// iterating; callers should not consume more than m values.
func From(m int) iter.Seq[int] {
	return func(yield func(int) bool) {
		n := 0
		selected := make(map[int]struct{})
		for {
			assert.That(n < m)
			r := rand.IntN(m)
			if _, exists := selected[r]; !exists {
				selected[r] = struct{}{}
				n++
				if !yield(r) {
					return
				}
			}
		}
	}
}

// Take returns a sequence of exactly n unique samples drawn from [0, m).
// It yields pairs (i, r) where i is the 0-based sample index (0..n-1) and r
// is the sampled value. Requires 0 < n < m.
func Take(n, m int) iter.Seq2[int, int] {
	assert.That(n > 0 && n < m)
	return func(yield func(int, int) bool) {
		selected := make(map[int]struct{})
		for i := 0; len(selected) < n; {
			r := rand.IntN(m)
			if _, exists := selected[r]; !exists {
				selected[r] = struct{}{}
				if !yield(i, r) {
					return
				}
				i++
			}
		}
	}
}
