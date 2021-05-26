// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package strs

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestEqual(t *testing.T) {
	test := func(x, y []string) {
		assert.T(t).That(Equal(x, y))
		assert.T(t).That(Equal(y, x))
	}
	xtest := func(x, y []string) {
		assert.T(t).That(!Equal(x, y))
		assert.T(t).That(!Equal(y, x))
	}
	empty := []string{}
	test(nil, nil)
	test(nil, empty)
	test(empty, nil)
	test(empty, empty)
	one := []string{"one"}
	test(one, one)
	xtest(one, nil)
	xtest(one, empty)
	four := []string{"one", "two", "three", "four"}
	test(four, four)
	xtest(four, nil)
	xtest(four, one)
}

func TestIndex(t *testing.T) {
	assert := assert.T(t).This
	assert(Index([]string{}, "five")).Is(-1)
	list := []string{"one", "two", "three", "two", "four"}
	assert(Index(list, "five")).Is(-1)
	assert(Index(list, "one")).Is(0)
	assert(Index(list, "two")).Is(1)
	assert(Index(list, "four")).Is(4)
}

func TestContains(t *testing.T) {
	assert := assert.T(t)
	assert.False(Contains([]string{}, "xxx"))
	list := []string{"one", "two", "three"}
	assert.True(Contains(list, "one"))
	assert.True(Contains(list, "two"))
	assert.True(Contains(list, "three"))
	assert.False(Contains(list, "o"))
	assert.False(Contains(list, "one1"))
	assert.False(Contains(list, "four"))
}

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

func TestJoin(t *testing.T) {
	assert := assert.T(t).This
	assert(Join("", nil)).Is("")
	assert(Join("", []string{})).Is("")
	assert(Join("", []string{"one", "two", "three"})).Is("onetwothree")
	assert(Join(",", []string{"one", "two", "three"})).Is("one,two,three")
	assert(Join(", ", []string{"one", "two", "three"})).Is("one, two, three")
	assert(Join("()", []string{"one", "two", "three"})).Is("(onetwothree)")
	assert(Join("[::]", []string{"one", "two", "three"})).Is("[one::two::three]")
}
