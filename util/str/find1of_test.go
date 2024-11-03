// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package str

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestFind1of(t *testing.T) {
	test := func(s, chars string, expected ...int) {
		t.Helper()
		first := expected[0]
		last := first
		if len(expected) > 1 {
			last = expected[1]
		}
		assert.T(t).This(Find1of(s, chars)).Is(first)
		assert.T(t).This(FindLast1of(s, chars)).Is(last)
	}
	test("", "", -1)
	test("", "x", -1)
	test("x", "", -1)
	test("x", "x", 0)
	test("xyz", "x", 0)
	test("xyz", "y", 1)
	test("xyz", "z", 2)
	test("xyz", "xyz", 0, 2)
	test("now is the time", "xyz", -1)

	test("now is the time", "i-k", 4, 12) // beginning of range
	test("now is the time", "a-i", 4, 14) // end of range
	test("now is the time", "^m-z", 3, 14)
	test("0", "\x00-\xff", 0)
	test("0", "^\x00-\xff", -1)
	test("hello", "z-a", -1)
	test("hello", "^z-a", 0, 4)
}

func FuzzFind1of(f *testing.F) {
	f.Fuzz(func(t *testing.T, s, chars string) {
		if len(chars) > 0 {
			makeBits(chars)
		}
		Find1of(s, chars)
		FindLast1of(s, chars)
	})
}

func TestMakeBits(t *testing.T) {
	test := func(chars string, expected string) {
		t.Helper()
		assert.T(t).This(makeBits(chars).String()).Is(expected)
	}
	test("x", "x")
	test("xyz", "xyz")
	test("x^z", "^xz")
	test("xyz^", "^xyz")
	test("-xyz-", "-xyz")
	test("a-e", "abcde")
	test("a-e^", "^abcde")
	test("^x", "!x")
	test("^a-e", "!abcde")
	test("^z-a", "!")
}

func (b bits) String() string {
	var s string
	for i := range b {
		for j := range 64 {
			if b[i]&(1<<uint(j)) != 0 {
				s += string(rune(i*64 + j))
			}
		}
	}
	t := "!"
	for i := range b {
		for j := range 64 {
			if b[i]&(1<<uint(j)) == 0 {
				t += string(rune(i*64 + j))
			}
		}
	}
	if len(s) < len(t) {
		return s
	}
	return t
}
