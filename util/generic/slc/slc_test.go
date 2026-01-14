// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package slc

import (
	"slices"
	"strings"
	"testing"
	"unsafe"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestHasPrefix(t *testing.T) {
	test := func(slist, slist2 string, expected bool) {
		t.Helper()
		list := strings.Fields(slist)
		list2 := strings.Fields(slist2)
		assert.T(t).This(HasPrefix(list, list2)).Is(expected)
	}
	test("", "", true)
	test("a b c", "", true)
	test("", "a", false)
	test("a b c", "a b c", true)
	test("a b c", "a b c d", false)
	test("a b c", "a x c", false)
}

func TestWithout(t *testing.T) {
	assert := assert.T(t).This
	assert(Without([]string{}, "five")).Is([]string{})
	list := []string{"one", "two", "three", "two", "four"}
	assert(Without(list, "five")).Is([]string(list))
	assert(Without(list, "one")).Is([]string{"two", "three", "two", "four"})
	assert(Without(list, "two")).Is([]string{"one", "three", "four"})
	assert(Without(list, "four")).Is([]string{"one", "two", "three", "two"})
}

func TestWith(t *testing.T) {
	var nilList []int
	var emptyList = []int{}
	var list1 = []int{1, 2, 3}
	var list2 = []int{4, 5}

	assert.T(t).This(With(nilList)).Is(emptyList)
	assert.T(t).This(With(nilList, nilList...)).Is(emptyList)
	assert.T(t).This(With(emptyList)).Is(emptyList)
	assert.T(t).This(With(emptyList, emptyList...)).Is(emptyList)
	assert.T(t).This(With(emptyList, list1...)).Is(list1)
	assert.T(t).This(With(list1)).Is(list1)
	assert.T(t).This(With(list1, 4)).Is([]int{1, 2, 3, 4})
	assert.T(t).This(With(list1, 4, 5)).Is([]int{1, 2, 3, 4, 5})
	assert.T(t).This(With(list1, list2...)).Is([]int{1, 2, 3, 4, 5})
}

func TestMap(t *testing.T) {
	assert := assert.T(t).This
	var nill []int
	fn := func(n int) int { return n * 10 }
	assert(MapFn(nill, fn)).Is(nil)
	assert(MapFn([]int{1}, fn)).Is([]int{10})
	assert(MapFn([]int{1, 2, 3}, fn)).Is([]int{10, 20, 30})
}

var X []int

func BenchmarkClipAppend(b *testing.B) {
	for b.Loop() {
		for n := range 10 {
			X = nil
			for j := range n {
				X = append(slices.Clip(X), j)
			}
		}
	}
}

func BenchmarkCloneAppend(b *testing.B) {
	for b.Loop() {
		for n := range 10 {
			X = nil
			for j := range n {
				X = append(slices.Clone(X), j)
			}
		}
	}
}

func BenchmarkNewAppend(b *testing.B) {
	for b.Loop() {
		for n := range 10 {
			X = nil
			for j := range n {
				y := make([]int, len(X)+1)
				copy(y, X)
				y[len(X)] = j
				X = y
			}
		}
	}
}

func BenchmarkWith(b *testing.B) {
	for b.Loop() {
		for n := range 10 {
			X = nil
			for j := range n {
				X = With(X, j)
			}
		}
	}
}

func TestClone(t *testing.T) {
	assert := assert.T(t).This
	var nilList []int
	var emptyList = []int{}
	var list1 = []int{1, 2, 3}

	assert(Clone(nilList) == nil)
	assert(Clone(emptyList) != nil)
	assert(Clone(list1)).Is([]int{1, 2, 3})

	x := Clone(emptyList)
	assert(unsafe.SliceData(x) != unsafe.SliceData(nilList))

	x = slices.Clone(emptyList)
	assert(unsafe.SliceData(x) == unsafe.SliceData(nilList)) // bad for GC
}

func TestWithoutFn(t *testing.T) {
	list := []int{}
	assert.T(t).That(Same(WithoutFn(list, func(n int) bool { return n == 5 }), list))
	list = []int{1, 2, 3, 2, 4}
	assert.T(t).That(Same(WithoutFn(list, func(n int) bool { return n == 5 }), list))
	assert.T(t).This(WithoutFn(list, func(n int) bool { return n == 1 })).
		Is([]int{2, 3, 2, 4})
	assert.T(t).This(WithoutFn(list, func(n int) bool { return n == 2 })).
		Is([]int{1, 3, 4})
	assert.T(t).This(WithoutFn(list, func(n int) bool { return n == 4 })).
		Is([]int{1, 2, 3, 2})
}

func TestGrow(t *testing.T) {
	assert := assert.T(t)

	// Test with n <= 0
	list := []int{1, 2, 3}
	result := Grow(list, 0)
	assert.This(len(result)).Is(len(list))
	assert.That(Same(result, list))

	result = Grow(list, -5)
	assert.This(len(result)).Is(len(list))
	assert.That(Same(result, list))

	// Test with nil slice
	var nilList []int
	result = Grow(nilList, 3)
	assert.This(len(result)).Is(3)
	assert.This(cap(result)).Is(3)
	// New elements should be zero-initialized
	for i := range result {
		assert.This(result[i]).Is(0)
	}

	// Test growing when capacity is sufficient
	list = make([]int, 3, 10)
	list[0], list[1], list[2] = 1, 2, 3
	result = Grow(list, 5)
	assert.This(len(result)).Is(8)
	assert.This(result[:3]).Is([]int{1, 2, 3})
	assert.That(Same(result[:3], list))
	// IMPORTANT: New elements should be zero-initialized
	for i := 3; i < len(result); i++ {
		assert.This(result[i]).Is(0)
	}

	// Test growing when capacity is insufficient
	list = make([]int, 3, 5)
	list[0], list[1], list[2] = 1, 2, 3
	result = Grow(list, 5)
	assert.This(len(result)).Is(8)
	assert.This(cap(result) >= 8)
	assert.This(result[:3]).Is([]int{1, 2, 3})
	assert.That(!Same(result[:3], list)) // Should be reallocated
	// New elements should be zero-initialized
	for i := 3; i < len(result); i++ {
		assert.This(result[i]).Is(0)
	}

	// Test with empty slice
	emptyList := []int{}
	result = Grow(emptyList, 3)
	assert.This(len(result)).Is(3)
	assert.This(cap(result)).Is(3)
	for i := range result {
		assert.This(result[i]).Is(0)
	}

	// Test grow, use, shrink, grow again
	list = Grow([]int{}, 5)
	for i := range list {
		list[i] = i + 100 // Fill with non-zero data
	}
	list = list[:0]            // Shrink to length 0 (but capacity remains)
	result = Grow(list, 3)     // Grow again - should zero-initialize!
	assert.This(len(result)).Is(3)
	for i := range result {
		assert.This(result[i]).Is(0) // These MUST be zero, not the old 100-104 values
	}
}

func TestHasDup(t *testing.T) {
	test := func(list []int, expected bool) {
		t.Helper()
		assert.T(t).This(HasDup(list)).Is(expected)
	}
	test([]int{}, false)
	test([]int{1}, false)
	test([]int{1, 2, 3}, false)
	test([]int{1, 2, 3, 4, 5}, false)
	test([]int{1, 1}, true)
	test([]int{1, 2, 1}, true)
	test([]int{1, 2, 3, 2}, true)
	test([]int{1, 2, 3, 4, 1}, true)
	test([]int{1, 2, 3, 4, 5, 3}, true)

	testStr := func(list []string, expected bool) {
		t.Helper()
		assert.T(t).This(HasDup(list)).Is(expected)
	}
	testStr([]string{}, false)
	testStr([]string{"a"}, false)
	testStr([]string{"a", "b", "c"}, false)
	testStr([]string{"a", "a"}, true)
	testStr([]string{"a", "b", "a"}, true)
	testStr([]string{"a", "b", "c", "b"}, true)
}

func TestPartition(t *testing.T) {
	// Test empty
	data := []int{}
	pivot := func(i int) bool { return data[i] < 5 }
	swap := func(i, j int) { data[i], data[j] = data[j], data[i] }
	i := Partition(0, pivot, swap)
	assert.T(t).This(i).Is(0)

	// Test all true
	data = []int{1, 2, 3}
	pivot = func(i int) bool { return data[i] < 5 }
	swap = func(i, j int) { data[i], data[j] = data[j], data[i] }
	i = Partition(3, pivot, swap)
	assert.T(t).This(i).Is(3)
	assert.T(t).This(data).Is([]int{1, 2, 3})

	// Test all false
	data = []int{6, 7, 8}
	pivot = func(i int) bool { return data[i] < 5 }
	swap = func(i, j int) { data[i], data[j] = data[j], data[i] }
	i = Partition(3, pivot, swap)
	assert.T(t).This(i).Is(0)
	assert.T(t).This(data).Is([]int{6, 7, 8})

	// Test mixed
	data = []int{3, 6, 1, 8, 4}
	pivot = func(i int) bool { return data[i] < 5 }
	swap = func(i, j int) { data[i], data[j] = data[j], data[i] }
	i = Partition(5, pivot, swap)
	assert.T(t).This(i).Is(3)
	// Check that left are <5, right >=5
	for j := 0; j < i; j++ {
		if data[j] >= 5 {
			t.Errorf("data[%d]=%d should be <5", j, data[j])
		}
	}
	for j := i; j < 5; j++ {
		if data[j] < 5 {
			t.Errorf("data[%d]=%d should be >=5", j, data[j])
		}
	}
}
