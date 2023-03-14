// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"regexp"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/regex"
)

func TestNumberPat(t *testing.T) {
	test := func (s string, expected bool) {
		t.Helper()
		assert := assert.T(t)
		assert.This(numberPat.Matches(s)).Is(expected)
		if expected == true {
			assert.True(numberPat.Matches("+" + s))
			assert.True(numberPat.Matches("-" + s))
			assert.False(numberPat.Matches("x" + s))
			assert.False(numberPat.Matches(s + "x"))
		}
	}
	test("0", true)
	test("6", true)
	test("007", true)
	test("123", true)
	test("123.", true)
	test(".123", true)
	test("123.465", true)
	test("1e6", true)
	test("1.5e6", true)
	test("1.5e-6", true)
	test("1.5e+6", true)
	test("1.5e-23", true)

	test("", false)
	test(".", false)
	test("+", false)
	test("-", false)
	test("-.", false)
	test("+-.", false)
	test("1.2.3", false)
	test("\n123", false)
	test("123\n", false)
	test("e5", false)
}

func BenchmarkNumberPatRegexp(b *testing.B) {
	numPat := regexp.MustCompile(`^[+-]?(\d+(\.\d*)?|\.\d+)([eE][+-]?\d\d?)?$`)
	big := strings.Repeat("1", 80) + "x"
	for i := 0; i < b.N; i++ {
		numPat.MatchString("0")
		numPat.MatchString("123.45")
		numPat.MatchString("+123.45")
		numPat.MatchString(big)
	}
}

func BenchmarkNumberPatRegex(b *testing.B) {
	numPat := regex.Compile(`\A[+-]?(\d+(\.\d*)?|\.\d+)([eE][+-]?\d\d?)?\Z`)
	big := strings.Repeat("1", 80) + "x"
	for i := 0; i < b.N; i++ {
		numPat.Matches("0")
		numPat.Matches("123.45")
		numPat.Matches("+123.45")
		numPat.Matches(big)
	}
}
