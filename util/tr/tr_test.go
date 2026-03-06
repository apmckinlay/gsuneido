// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tr

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func Test_makset(t *testing.T) {
	test := func(s, expected string) {
		t.Helper()
		assert.T(t).This(string(New(s))).Is(expected)
	}
	test("", "")
	test("foo", "foo")
	test("^foo", "^foo")
	test("-foo", "-foo")
	test("foo-", "foo-")
	test("m-p", "mnop")
	test("-0-9-", "-0123456789-")
	test("\xfa-\xff", "\xfa\xfb\xfc\xfd\xfe\xff")
	test("z-a", "")
}

func TestReplace(t *testing.T) {
	test := func(src, from, to, expected string) {
		t.Helper()
		result := Replace(src, New(from), New(to))
		assert.T(t).This(result).Is(expected)
	}
	// empty source
	test("", "", "", "")
	test("", "abc", "ABC", "")
	test("", "^abc", "x", "")

	// delete
	test("abc", "", "", "abc")
	test("abc", "xyz", "", "abc")
	test("zon", "xyz", "", "on")
	test("oyn", "xyz", "", "on")
	test("nox", "xyz", "", "no")
	test("zyx", "xyz", "", "")

	// replace
	test("zon", "xyz", "XYZ", "Zon")
	test("oyn", "xyz", "XYZ", "oYn")
	test("nox", "xyz", "XYZ", "noX")
	test("zyx", "xyz", "XYZ", "ZYX")
	test("zyx", "a-z", "A-Z", "ZYX")

	// allbut delete
	test("a b - c", "^abc", "", "abc")
	test("a b - c", "^a-z", "", "abc")

	// allbut collapse
	test("a  b - c", "^abc", " ", "a b c")
	test("a  b - c", "^a-z", " ", "a b c")

	// literal dash
	test("a-b-c", "-x", "", "abc")
	test("a-b-c", "x-", "", "abc")

	// collapse at end
	test("hello \t\n\n", " \t\n", "\n", "hello\n")

	// complex allbut
	test("abc", "^\t\r\n\x20-\xff", "", "abc")
	test("a\x00b\x1fc\xffd", "^\t\r\n\x20-\xff", "", "abc\xffd")
}

// to run: go test -fuzz=Fuzz_makset -run=Fuzz_makset

func Fuzz_makset(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string) {
		New(s)
	})
}

// to run: go test -fuzz=FuzzReplace -run=FuzzReplace

func FuzzReplace(f *testing.F) {
	f.Fuzz(func(t *testing.T, s1, s2, s3 string) {
		Replace(s1, New(s2), New(s3))
	})
}
