// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package interleave

import (
	"slices"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestForeachBasic(t *testing.T) {
	var results [][]int
	a := []int{1, 2}
	b := []int{3}

	for seq := range All(a, b) {
		results = append(results, slices.Clone(seq))
	}

	expected := [][]int{
		{1, 2, 3},
		{1, 3, 2},
		{3, 1, 2},
	}

	assert.T(t).This(len(results)).Is(3)
	for _, exp := range expected {
		assert.T(t).True(containsSlice(results, exp))
	}
}

func TestForeachSymmetric(t *testing.T) {
	var results [][]int
	a := []int{1}
	b := []int{2, 3}

	for seq := range All(a, b) {
		results = append(results, slices.Clone(seq))
	}

	expected := [][]int{
		{1, 2, 3},
		{2, 1, 3},
		{2, 3, 1},
	}

	assert.T(t).This(len(results)).Is(3)
	for _, exp := range expected {
		assert.T(t).True(containsSlice(results, exp))
	}
}

func TestForeachEmptySlices(t *testing.T) {
	var results [][]int
	a := []int{}
	b := []int{}

	for seq := range All(a, b) {
		results = append(results, slices.Clone(seq))
	}

	assert.T(t).This(len(results)).Is(1)
	assert.T(t).This(len(results[0])).Is(0)
}

func TestForeachEmptyA(t *testing.T) {
	var results [][]int
	a := []int{}
	b := []int{1, 2}

	for seq := range All(a, b) {
		results = append(results, slices.Clone(seq))
	}

	assert.T(t).This(len(results)).Is(1)
	assert.T(t).This(results[0]).Is([]int{1, 2})
}

func TestForeachEmptyB(t *testing.T) {
	var results [][]int
	a := []int{1, 2}
	b := []int{}

	for seq := range All(a, b) {
		results = append(results, slices.Clone(seq))
	}

	assert.T(t).This(len(results)).Is(1)
	assert.T(t).This(results[0]).Is([]int{1, 2})
}

func TestForeachSingleElements(t *testing.T) {
	var results [][]int
	a := []int{1}
	b := []int{2}

	for seq := range All(a, b) {
		results = append(results, slices.Clone(seq))
	}

	expected := [][]int{
		{1, 2},
		{2, 1},
	}

	assert.T(t).This(len(results)).Is(2)
	for _, exp := range expected {
		assert.T(t).True(containsSlice(results, exp))
	}
}

func TestForeachStrings(t *testing.T) {
	var results [][]string
	a := []string{"a", "b"}
	b := []string{"x"}

	for seq := range All(a, b) {
		results = append(results, slices.Clone(seq))
	}

	expected := [][]string{
		{"a", "b", "x"},
		{"a", "x", "b"},
		{"x", "a", "b"},
	}

	assert.T(t).This(len(results)).Is(3)
	for _, exp := range expected {
		assert.T(t).True(containsSlice(results, exp))
	}
}

func TestForeachLarger(t *testing.T) {
	var results [][]int
	a := []int{1, 2}
	b := []int{3, 4}

	for seq := range All(a, b) {
		results = append(results, slices.Clone(seq))
	}

	// Should have C(4,2) = 6 interleavings
	assert.T(t).This(len(results)).Is(6)

	// Verify each result contains all elements
	for _, result := range results {
		assert.T(t).This(len(result)).Is(4)
		assert.T(t).True(slices.Contains(result, 1))
		assert.T(t).True(slices.Contains(result, 2))
		assert.T(t).True(slices.Contains(result, 3))
		assert.T(t).True(slices.Contains(result, 4))
	}

	// Verify no duplicates
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			assert.T(t).False(slices.Equal(results[i], results[j]))
		}
	}
}

func TestForeachPreservesElements(t *testing.T) {
	var results [][]int
	a := []int{10, 20, 30}
	b := []int{1, 2}

	for seq := range All(a, b) {
		results = append(results, slices.Clone(seq))
	}

	// Should have C(5,2) = 10 interleavings
	assert.T(t).This(len(results)).Is(10)

	// Every result should contain exactly the same elements
	for _, result := range results {
		assert.T(t).This(len(result)).Is(5)

		// Count occurrences of each element
		counts := make(map[int]int)
		for _, val := range result {
			counts[val]++
		}

		assert.T(t).This(counts[10]).Is(1)
		assert.T(t).This(counts[20]).Is(1)
		assert.T(t).This(counts[30]).Is(1)
		assert.T(t).This(counts[1]).Is(1)
		assert.T(t).This(counts[2]).Is(1)
	}
}

func TestForeachCallbackCount(t *testing.T) {
	count := 0
	a := []int{1, 2, 3}
	b := []int{4, 5}

	for range All(a, b) {
		count++
	}

	// Should have C(5,2) = 10 interleavings
	assert.T(t).This(count).Is(10)
}

// Helper function to check if a slice of int slices contains a specific slice
func containsSlice[T comparable](haystack [][]T, target []T) bool {
	for _, s := range haystack {
		if slices.Equal(s, target) {
			return true
		}
	}
	return false
}
