// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package str

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestCapitalized(t *testing.T) {
	Assert(t).That(Capitalized(""), Is(false))
	Assert(t).That(Capitalized("a"), Is(false))
	Assert(t).That(Capitalized("abc"), Is(false))
	Assert(t).That(Capitalized("?"), Is(false))
	Assert(t).That(Capitalized("A"), Is(true))
	Assert(t).That(Capitalized("Abc"), Is(true))
}

func TestCapitalize(t *testing.T) {
	Assert(t).That(Capitalize(""), Is(""))
	Assert(t).That(Capitalize("#$%"), Is("#$%"))
	Assert(t).That(Capitalize("abc"), Is("Abc"))
	Assert(t).That(Capitalize("a"), Is("A"))
	Assert(t).That(Capitalize("abC"), Is("AbC"))
}

func TestUnCapitalize(t *testing.T) {
	Assert(t).That(UnCapitalize(""), Is(""))
	Assert(t).That(UnCapitalize("#$%"), Is("#$%"))
	Assert(t).That(UnCapitalize("abc"), Is("abc"))
	Assert(t).That(UnCapitalize("A"), Is("a"))
	Assert(t).That(UnCapitalize("AbC"), Is("abC"))
}

func TestIndexFunc(t *testing.T) {
	f := func(c byte) bool {
		return c == ' '
	}
	Assert(t).That(IndexFunc("", f), Is(-1))
	Assert(t).That(IndexFunc(" ", f), Is(0))
	Assert(t).That(IndexFunc("foobar", f), Is(-1))
	Assert(t).That(IndexFunc("foo bar", f), Is(3))
}

func TestJoin(t *testing.T) {
	Assert(t).That(Join(""), Is(""))
	Assert(t).That(Join("", "one", "two", "three"), Is("onetwothree"))
	Assert(t).That(Join(",", "one", "two", "three"), Is("one,two,three"))
	Assert(t).That(Join(", ", "one", "two", "three"), Is("one, two, three"))
	Assert(t).That(Join("()", "one", "two", "three"), Is("(onetwothree)"))
	Assert(t).That(Join("[::]", "one", "two", "three"), Is("[one::two::three]"))
}

func TestList_Has(t *testing.T) {
	Assert(t).False(List{}.Has("xxx"))
	list := List{"one", "two", "three"}
	Assert(t).True(list.Has("one"))
	Assert(t).True(list.Has("two"))
	Assert(t).True(list.Has("three"))
	Assert(t).False(list.Has("o"))
	Assert(t).False(list.Has("one1"))
	Assert(t).False(list.Has("four"))
}

func TestList_Index(t *testing.T) {
	Assert(t).That(List{}.Index("five"), Is(-1))
	list := List{"one", "two", "three", "two", "four"}
	Assert(t).That(list.Index("five"), Is(-1))
	Assert(t).That(list.Index("one"), Is(0))
	Assert(t).That(list.Index("two"), Is(1))
	Assert(t).That(list.Index("four"), Is(4))
}

func TestList_Without(t *testing.T) {
	Assert(t).That(List{}.Without("five"), Is([]string{}))
	list := List{"one", "two", "three", "two", "four"}
	Assert(t).That(list.Without("five"), Is([]string(list)))
	Assert(t).That(list.Without("one"),
		Is([]string{"two", "three", "two", "four"}))
	Assert(t).That(list.Without("two"),
		Is([]string{"one", "three", "four"}))
	Assert(t).That(list.Without("four"),
		Is([]string{"one", "two", "three", "two"}))
}

func TestList_Reverse(t *testing.T) {
	list := []string{}
	List(list).Reverse()
	Assert(t).That(list, Is([]string{}))
	list = []string{"one", "two", "three"}
	List(list).Reverse()
	Assert(t).That(list, Is([]string{"three", "two", "one"}))
}
