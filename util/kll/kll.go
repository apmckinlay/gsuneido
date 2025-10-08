// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package kll implements a KLL-style sketch.
// It give approximate quantile estimates from a stream of values.
// See: Optimal Quantile Approximation in Streams
// https://arxiv.org/pdf/1603.05346
//
// This first version makes all levels the same size (k).
package kll

import (
	"cmp"
	"math/rand/v2"
	"slices"

	"github.com/apmckinlay/gsuneido/util/assert"
)

// caps is the level capacities with k = 200, c = 2/3
var caps = []int{200, 133, 89, 59, 40, 26, 18, 12, 8, 5, 3}

//{3, 5, 8, 12, 18, 26, 40, 59, 89, 133, 200}

// Sketch accumulates a stream of values
type Sketch[T cmp.Ordered] struct {
	count  int
	levels [][]T // last level is the input level, weight 1, level[0] is caps[0]
}

// New creates a new Sketch with the default k value.
func New[T cmp.Ordered]() *Sketch[T] {
	sk := &Sketch[T]{
		levels: make([][]T, 1, 20),
	}
	sk.levels[0] = make([]T, 0, caps[0]*3/2)
	return sk
}

// Insert adds a new value to the sketch, compacting levels as necessary.
func (sk *Sketch[T]) Insert(value T) {
	sk.count++
	h := len(sk.levels) - 1
	sk.levels[h] = append(sk.levels[h], value)
	for ; h >= 0 && len(sk.levels[h]) > caps[h]; h-- {
		// may need to compact a given level several times
		// if the length is more than twice the capacity (from promotion)
		for len(sk.levels[h]) > caps[h] {
			sk.compact(h)
		}
	}
}

func (sk *Sketch[T]) compact(h int) {
	src := sk.levels[h]
	slices.Sort(src)
	if h-1 < 0 {
		// increase the weight of all the existing levels
		for h := range sk.levels {
			sk.compactInPlace(h)
		}
		sk.levels = append(sk.levels,
			make([]T, 0, caps[len(sk.levels)]*3/2))
		return
	}
	dst := sk.levels[h-1]

	which := rand.IntN(2)
	for i := 1; i < len(src); i += 2 {
		dst = append(dst, src[i-which])
	}
	sk.levels[h] = src[:0]
	sk.levels[h-1] = dst
}

func (sk *Sketch[T]) compactInPlace(h int) {
	src := sk.levels[h]
	slices.Sort(src)
	write := 0
	which := rand.IntN(2)
	for i := 1; i < len(src); i += 2 {
		src[write] = src[i-which]
		write++
	}
	sk.levels[h] = src[:write]
}

// Count returns the number of values inserted into the sketch.
func (sk *Sketch[T]) Count() int {
	return sk.count
}

// Query returns the approximate value at the given quantile (0.0 <= q <= 1.0).
// For example, Query(0.5) returns the approximate median.
func (sk *Sketch[T]) Query(q float64) T {
	if sk.count == 0 {
		panic("no data")
	}
	if q < 0.0 || q > 1.0 {
		panic("out of range")
	}

	// Collect all items with their weights
	type weightedItem struct {
		value  T
		weight int
	}

	var items []weightedItem
	weight := 1
	for h := len(sk.levels) - 1; h >= 0; h-- {
		for _, value := range sk.levels[h] {
			items = append(items, weightedItem{value, weight})
		}
		weight *= 2
	}
	assert.That(len(items) > 0)

	// Sort items by value
	slices.SortFunc(items, func(a, b weightedItem) int {
		if a.value < b.value {
			return -1
		} else if a.value > b.value {
			return 1
		}
		return 0
	})

	// Find the item at the target rank
	targetRank := int(float64(sk.count) * q)
	currentRank := 0

	for _, item := range items {
		currentRank += item.weight
		if currentRank >= targetRank {
			return item.value
		}
	}

	// If we get here, return the last item
	return items[len(items)-1].value
}
