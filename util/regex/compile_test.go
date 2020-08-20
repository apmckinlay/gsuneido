// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestCompile(t *testing.T) {
	test := func(input, expected string) {
		t.Helper()
		//fmt.Println("input '" + input + "'")
		p := Compile(input).String()
		assert.T(t).Msg(input).This(strings.TrimSpace(p[5 : len(p)-7])).
			Is(expected)
	}
	test("", "")
	test(".", ".")
	test("a", "'a'")
	test("abc", "'abc'")
	test("^xyz", "^ 'xyz'")
	test("abc$", "'abc' $")
	test("^xyz$", "^ 'xyz' $")
	test("?ab", "'?ab'")
	test("*ab", "'*ab'")
	test("+ab", "'+ab'")
	test("a?", "Branch(1, 2) 'a'")
	test("a??", "Branch(2, 1) 'a'")
	test("abc?de", "'ab' Branch(1, 2) 'c' 'de'")
	test(".*", "^ Branch(1, 3) . Branch(-1, 1)")
	test(".+", "^ . Branch(-1, 1)")
	test("abc+de", "'ab' 'c' Branch(-1, 1) 'de'")
	test("abc*de", "'ab' Branch(1, 3) 'c' Branch(-1, 1) 'de'")
	test("ab\\.?cd", "'ab' Branch(1, 2) '.' 'cd'")
	test("(ab+c)+x", "Left1 'a' 'b' Branch(-1, 1) 'c' Right1 Branch(-6, 1) 'x'")
	test("ab|cd",
		"Branch(1, 3) 'ab' Jump(2) 'cd'")
	test("ab|cd|ef",
		"Branch(1, 3) 'ab' Jump(3) Branch(1, 3) 'cd' Jump(2) 'ef'")
	test("abc\\Z", "'abc' \\Z")
	test("[a]", "'a'")
	test("[\\a]", "'a'")
	test(`[\s]`, "[ \t\r\n]")
	test("[ace]", "[ace]")
	test("[a-cx-z]", "[abcxyz]")
	test("[a-Z]", "[]")
	test("(?i)x.y(?-i)z", "i'x' . i'y' 'z'")

	test("(?q).*(?-q)def", "'.*def'")

	test("\\<Foo\\>", "\\< 'Foo' \\>")

	test("a(bc)d", "'a' Left1 'bc' Right1 'd'")

	test(".\\5.", ". \\5 .")
	test("(?i)\\5", "i\\5")

	test("a[.]b", "'a.b'")

	test("a(?q).(?-q).c(?q).(?-q).", "'a.' . 'c.' .")

	test("\\", "'\\'")

	assert.T(t).This(func() { Compile("(abc") }).Panics("missing ')'")
	assert.T(t).This(func() { Compile("abc)def") }).
		Panics("closing ) without opening (")
}
