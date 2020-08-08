// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tabs

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestDetab(t *testing.T) {
	test := func(s, expected string) {
		Assert(t).That(Detab(s), Is(expected))
	}
	test("", "")
	test("foo bar", "foo bar")
	test("  foo", "  foo")
	test("\tfoo", "    foo")
	test("  \tfoo", "    foo")
	test("\t\tfoo", "        foo")
	test(" \t \tfoo", "        foo")
	test("x\ty", "x   y")
	test("\tfoo\n\tbar", "    foo\n    bar")
	test("\tfoo\r\n\tbar", "    foo\r\n    bar")
}

func TestEntab(t *testing.T) {
	test := func(s, expected string) {
		Assert(t).That(Entab(s), Is(expected))
	}
	test("", "")
	test("foo", "foo")
	test("foo bar", "foo bar")
	test("    foo", "\tfoo")
	test("  \tfoo", "\tfoo")
	test(" \t foo", "\t foo")
	test(" \t foo  \t  ", "\t foo")
	test("foo\tbar", "foo\tbar") // only leading converted
	test("    foo\r\n    bar\r\n", "\tfoo\r\n\tbar\r\n")
}
