// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package slc

import (
	"strings"
	"testing"

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
