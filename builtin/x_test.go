// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"regexp"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/regex"
	. "github.com/apmckinlay/gsuneido/runtime"
)

func TestNumberQ(t *testing.T) {
	test := func(s string, expected bool) {
		t.Helper()
		assert.T(t).This(string_NumberQ(SuStr(s))).Is(SuBool(expected))
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

	test("_1", false)
	test("1_", true)
	test("1_2", true)
	test("1__2", true)
	test("1_2_3", true)
	test("1_2_3_", true)
	test("1_2_3__", true)
	test("1_2_3__4", true)
	test("1_2_3__4_", true)
	test("1.0", true)
	test("1._0", true)
	test("1_.0", true)
	test("1.0_", true)
	test("1._0_", true)
	test("1_.0_", true)
	test("1.0__", true)
	test("1.__0", true)
	test("1.__0_", true)
	test("1.__0__", true)
	test("1.0__0", true)
	test("__1.0__0_", false)
	test("__1.0__0__", false)

	test("0x123", true)
	test("-0x123", true)
	test("_0x123", false)
	test("0x123.456", false)
	test("0x1_2_3", true)
	test("0x1_2_3_", true)
	test("0x1_2_3__", true)
	test("0x1_2_3__4", true)
	test("0xZ12", false)
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
