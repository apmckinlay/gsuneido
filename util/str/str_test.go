// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package str

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestCapitalized(t *testing.T) {
	Assert(t).That(Capitalized(""), Equals(false))
	Assert(t).That(Capitalized("a"), Equals(false))
	Assert(t).That(Capitalized("abc"), Equals(false))
	Assert(t).That(Capitalized("?"), Equals(false))
	Assert(t).That(Capitalized("A"), Equals(true))
	Assert(t).That(Capitalized("Abc"), Equals(true))
}

func TestCapitalize(t *testing.T) {
	Assert(t).That(Capitalize(""), Equals(""))
	Assert(t).That(Capitalize("#$%"), Equals("#$%"))
	Assert(t).That(Capitalize("abc"), Equals("Abc"))
	Assert(t).That(Capitalize("a"), Equals("A"))
	Assert(t).That(Capitalize("abC"), Equals("AbC"))
}

func TestUnCapitalize(t *testing.T) {
	Assert(t).That(UnCapitalize(""), Equals(""))
	Assert(t).That(UnCapitalize("#$%"), Equals("#$%"))
	Assert(t).That(UnCapitalize("abc"), Equals("abc"))
	Assert(t).That(UnCapitalize("A"), Equals("a"))
	Assert(t).That(UnCapitalize("AbC"), Equals("abC"))
}

func TestIndexFunc(t *testing.T) {
	f := func(c byte) bool {
		return c == ' '
	}
	Assert(t).That(IndexFunc("", f), Equals(-1))
	Assert(t).That(IndexFunc(" ", f), Equals(0))
	Assert(t).That(IndexFunc("foobar", f), Equals(-1))
	Assert(t).That(IndexFunc("foo bar", f), Equals(3))
}

func TestJoin(t *testing.T) {
	Assert(t).That(Join(""), Equals(""))
	Assert(t).That(Join("", "one", "two", "three"), Equals("onetwothree"))
	Assert(t).That(Join(",", "one", "two", "three"), Equals("one,two,three"))
	Assert(t).That(Join(", ", "one", "two", "three"), Equals("one, two, three"))
	Assert(t).That(Join("()", "one", "two", "three"), Equals("(onetwothree)"))
	Assert(t).That(Join("[::]", "one", "two", "three"), Equals("[one::two::three]"))
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
	Assert(t).That(List{}.Index("five"), Equals(-1))
	list := List{"one", "two", "three", "two", "four"}
	Assert(t).That(list.Index("five"), Equals(-1))
	Assert(t).That(list.Index("one"), Equals(0))
	Assert(t).That(list.Index("two"), Equals(1))
	Assert(t).That(list.Index("four"), Equals(4))
}

func TestList_Without(t *testing.T) {
	Assert(t).That(List{}.Without("five"), Equals([]string{}))
	list := List{"one", "two", "three", "two", "four"}
	Assert(t).That(list.Without("five"), Equals([]string(list)))
	Assert(t).That(list.Without("one"),
		Equals([]string{"two", "three", "two", "four"}))
	Assert(t).That(list.Without("two"),
		Equals([]string{"one", "three", "four"}))
	Assert(t).That(list.Without("four"),
		Equals([]string{"one", "two", "three", "two"}))
}

func TestList_Reverse(t *testing.T) {
	list := []string{}
	List(list).Reverse()
	Assert(t).That(list, Equals([]string{}))
	list = []string{"one", "two", "three"}
	List(list).Reverse()
	Assert(t).That(list, Equals([]string{"three", "two", "one"}))
}
