// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package str

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestList_Has(t *testing.T) {
	assert := assert.T(t)
	assert.False(List{}.Has("xxx"))
	list := List{"one", "two", "three"}
	assert.True(list.Has("one"))
	assert.True(list.Has("two"))
	assert.True(list.Has("three"))
	assert.False(list.Has("o"))
	assert.False(list.Has("one1"))
	assert.False(list.Has("four"))
}

func TestList_Index(t *testing.T) {
	assert := assert.T(t).This
	assert(List{}.Index("five")).Is(-1)
	list := List{"one", "two", "three", "two", "four"}
	assert(list.Index("five")).Is(-1)
	assert(list.Index("one")).Is(0)
	assert(list.Index("two")).Is(1)
	assert(list.Index("four")).Is(4)
}

func TestList_Without(t *testing.T) {
	assert := assert.T(t).This
	assert(List{}.Without("five")).Is([]string{})
	list := List{"one", "two", "three", "two", "four"}
	assert(list.Without("five")).Is([]string(list))
	assert(list.Without("one")).Is([]string{"two", "three", "two", "four"})
	assert(list.Without("two")).Is([]string{"one", "three", "four"})
	assert(list.Without("four")).Is([]string{"one", "two", "three", "two"})
}

func TestList_Reverse(t *testing.T) {
	list := []string{}
	List(list).Reverse()
	assert.T(t).This(list).Is([]string{})
	list = []string{"one", "two", "three"}
	List(list).Reverse()
	assert.T(t).This(list).Is([]string{"three", "two", "one"})
}

func TestList_HasPrefix(t *testing.T) {
	test := func(slist, slist2 string, expected bool) {
		t.Helper()
		list := strings.Fields(slist)
		list2 := strings.Fields(slist2)
		assert.T(t).This(List(list).HasPrefix(list2)).Is(expected)
	}
	test("", "", true)
	test("a b c", "", true)
	test("", "a", false)
	test("a b c", "a b c", true)
	test("a b c", "a b c d", false)
	test("a b c", "a x c", false)
}
