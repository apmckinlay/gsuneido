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
