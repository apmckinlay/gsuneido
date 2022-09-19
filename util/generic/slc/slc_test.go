// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package slc contains additions to x/exp/slices
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

func TestReplace(t *testing.T) {
	list := []string{"one", "two", "three", "two", "four"}
	list2 := Replace(list, nil, nil)
	assert.T(t).That(Same(list, list2))
	from := []string{"two", "five", "one"}
	to := []string{"2", "5", "1"}
	list1 := []string{"a", "b", "c"}
	list2 = Replace(list1, from, to)
	assert.T(t).That(Same(list1, list2))
	assert.T(t).This(Replace(list, from, to)).
		Is([]string{"1", "2", "three", "2", "four"})
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
