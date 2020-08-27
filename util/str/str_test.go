// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package str

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestCapitalized(t *testing.T) {
	assert := assert.T(t).This
	assert(Capitalized("")).Is(false)
	assert(Capitalized("a")).Is(false)
	assert(Capitalized("abc")).Is(false)
	assert(Capitalized("?")).Is(false)
	assert(Capitalized("A")).Is(true)
	assert(Capitalized("Abc")).Is(true)
}

func TestCapitalize(t *testing.T) {
	assert := assert.T(t).This
	assert(Capitalize("")).Is("")
	assert(Capitalize("#$%")).Is("#$%")
	assert(Capitalize("abc")).Is("Abc")
	assert(Capitalize("a")).Is("A")
	assert(Capitalize("abC")).Is("AbC")
}

func TestUnCapitalize(t *testing.T) {
	assert := assert.T(t).This
	assert(UnCapitalize("")).Is("")
	assert(UnCapitalize("#$%")).Is("#$%")
	assert(UnCapitalize("abc")).Is("abc")
	assert(UnCapitalize("A")).Is("a")
	assert(UnCapitalize("AbC")).Is("abC")
}

func TestIndexFunc(t *testing.T) {
	f := func(c byte) bool {
		return c == ' '
	}
	assert := assert.T(t).This
	assert(IndexFunc("", f)).Is(-1)
	assert(IndexFunc(" ", f)).Is(0)
	assert(IndexFunc("foobar", f)).Is(-1)
	assert(IndexFunc("foo bar", f)).Is(3)
}

func TestJoin(t *testing.T) {
	assert := assert.T(t).This
	assert(Join("")).Is("")
	assert(Join("", "one", "two", "three")).Is("onetwothree")
	assert(Join(",", "one", "two", "three")).Is("one,two,three")
	assert(Join(", ", "one", "two", "three")).Is("one, two, three")
	assert(Join("()", "one", "two", "three")).Is("(onetwothree)")
	assert(Join("[::]", "one", "two", "three")).Is("[one::two::three]")
}

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

func TestCmpLower(t *testing.T) {
	test := func(s1, s2 string, result int) {
		assert.T(t).Msg(s1, "<=>", s2).This(CmpLower(s1, s2)).Is(result)
		assert.T(t).Msg(s2, "<=>", s1).This(CmpLower(s2, s1)).Is(-result)
	}
	test("", "", 0)
	test("", "2", -1)
	test("123", "123", 0)
	test("hello world", "hello world", 0)
	test("Hello World", "hello world", 0)
	test("Hello", "world", -1)
	test("hello", "World", -1)
}

func TestRemovePrefix(t *testing.T) {
	test := func(s, pre, expected string) {
		t.Helper()
		assert.T(t).Msg(s, ",", pre).This(RemovePrefix(s, pre)).Is(expected)
	}
	test("", "", "")
	test("", "abc", "")
	test("abc", "", "abc")
	test("abc", "xyz", "abc")
	test("foobar", "foo", "bar")
	test("foobar", "bar", "foobar")
}

func TestRemoveSuffix(t *testing.T) {
	test := func(s, pre, expected string) {
		t.Helper()
		assert.T(t).Msg(s, ",", pre).This(RemoveSuffix(s, pre)).Is(expected)
	}
	test("", "", "")
	test("", "abc", "")
	test("abc", "", "abc")
	test("abc", "xyz", "abc")
	test("foobar", "foo", "foobar")
	test("foobar", "bar", "foo")
}
