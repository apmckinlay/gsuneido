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
	assert := assert.T(t).This
	list := []string{"one", "two", "three", "two", "four"}
	assert(Replace(list, nil, nil)).Is(list)
	from := []string{"two", "five", "one"}
	to := []string{"2", "5", "1"}
	assert(Replace(list, from, to)).Is([]string{"1", "2", "three", "2", "four"})
}