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

// Sketch accumulates a stream of values
type Sketch[T cmp.Ordered] struct {
	k      int // overridden by tests
	count  int
	levels [][]T // level 0 is the input level, weight 1
}

// New creates a new Sketch with the default k value.
func New[T cmp.Ordered]() *Sketch[T] {
	sk := &Sketch[T]{
		k:      200,
		levels: make([][]T, 1, 20),
	}
	sk.levels[0] = make([]T, 0, sk.k+1)
	return sk
}

// Insert adds a new value to the sketch, compacting levels as necessary.
func (sk *Sketch[T]) Insert(value T) {
	sk.count++
	sk.levels[0] = append(sk.levels[0], value)
	for h := 0; h < len(sk.levels) && len(sk.levels[h]) > sk.k; h++ {
		// may need to compact a given level several times
		// if the length is more than twice the capacity (from promotion)
		for len(sk.levels[h]) > sk.k {
			sk.compact(h)
		}
	}
	// invariant - all levels are <= k
	for h := range sk.levels {
		assert.That(len(sk.levels[h]) <= sk.k)
	}
}

func (sk *Sketch[T]) compact(h int) {
	src := sk.levels[h]
	slices.Sort(src)
	sk.ensureLevel(h + 1)
	dst := sk.levels[h+1]

	which := rand.IntN(2)
	for i := 1; i < len(src); i += 2 {
		dst = append(dst, src[i-which])
	}
	sk.levels[h] = src[:0]
	sk.levels[h+1] = dst
}

func (sk *Sketch[T]) ensureLevel(h int) {
	for len(sk.levels) <= h {
		sk.levels = append(sk.levels, make([]T, 0, sk.k))
	}
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
	for h := 0; h < len(sk.levels); h++ {
		weight := weight(h)
		for _, value := range sk.levels[h] {
			items = append(items, weightedItem{value, weight})
		}
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

// weight returns the weight of items at the given level
func weight(h int) int {
	// Each level up represents 2x more items due to compaction
	return 1 << uint(h)
}
