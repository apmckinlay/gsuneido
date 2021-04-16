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
	assert(Join("", nil)).Is("")
	assert(Join("", []string{})).Is("")
	assert(Join("", []string{"one", "two", "three"})).Is("onetwothree")
	assert(Join(",", []string{"one", "two", "three"})).Is("one,two,three")
	assert(Join(", ", []string{"one", "two", "three"})).Is("one, two, three")
	assert(Join("()", []string{"one", "two", "three"})).Is("(onetwothree)")
	assert(Join("[::]", []string{"one", "two", "three"})).Is("[one::two::three]")
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

func TestBeforeFirst(t *testing.T) {
	assert := assert.T(t).This
	assert(BeforeFirst("", "")).Is("")
	assert(BeforeFirst("^1234512345$", "^")).Is("")
	assert(BeforeFirst("^1234512345$", "z")).Is("^1234512345$")
	assert(BeforeFirst("^1234512345$", "4")).Is("^123")
	assert(BeforeFirst("^1234512345$", "51")).Is("^1234")
}

func TestAfterFirst(t *testing.T) {
	assert := assert.T(t).This
	assert(AfterFirst("", "")).Is("")
	assert(AfterFirst("^1234512345$", "z")).Is("^1234512345$")
	assert(AfterFirst("^1234512345$", "$")).Is("")
	assert(AfterFirst("^1234512345$", "4")).Is("512345$")
	assert(AfterFirst("^1234512345$", "51")).Is("2345$")
}

func TestBeforeLast(t *testing.T) {
	assert := assert.T(t).This
	assert(BeforeLast("", "")).Is("")
	assert(BeforeLast("^1234512345$", "z")).Is("^1234512345$")
	assert(BeforeLast("^1234512345$", "^")).Is("")
	assert(BeforeLast("^1234512345$", "4")).Is("^12345123")
	assert(BeforeLast("^1234512345$", "51")).Is("^1234")
}

func TestAfterLast(t *testing.T) {
	assert := assert.T(t).This
	assert(AfterLast("", "")).Is("")
	assert(AfterLast("^1234512345$", "z")).Is("^1234512345$")
	assert(AfterLast("^1234512345$", "$")).Is("")
	assert(AfterLast("^1234512345$", "4")).Is("5$")
	assert(AfterLast("^1234512345$", "51")).Is("2345$")
}

func TestEqualCI(t *testing.T) {
	test := func(x, y string) {
		assert.T(t).That(EqualCI(x, y))
		assert.T(t).That(EqualCI(y, x))
	}
	test("", "")
	test("foo bar", "foo bar")
	test("Foo Bar", "foo bar")
	test("Foo Bar", "FOO BAR")
	xtest := func(x, y string) {
		assert.T(t).That(!EqualCI(x, y))
		assert.T(t).That(!EqualCI(y, x))
	}
	xtest("", "xyz")
	xtest("foo", "food")
	xtest("foo", "bar")
}

func BenchmarkEqualCI(b *testing.B) {
	strs := []string{
		"Now Is The Time",
		"For All Good men",
		"To Come To The",
		"Aid Of Their Party",
		"FOO BAR",
		"bar foo",
		"",
		"!@#$$%%^&^&*&*(",
	}
	var a bool
	for i := 0; i < b.N; i++ {
		x := strs[i%len(strs)]
		y := strs[(i+1)%len(strs)]
		a = a || EqualCI(x, y)
	}
	B = a
}

func BenchmarkEqualLower(b *testing.B) {
	strs := []string{
		"Now Is The Time",
		"For All Good men",
		"To Come To The",
		"Aid Of Their Party",
		"FOO BAR",
		"bar foo",
		"",
		"!@#$$%%^&^&*&*(",
	}
	var a bool
	for i := 0; i < b.N; i++ {
		x := strs[i%len(strs)]
		y := strs[(i+1)%len(strs)]
		a = a || ToLower(x) == ToLower(y)
	}
	B = a
}

var B bool
