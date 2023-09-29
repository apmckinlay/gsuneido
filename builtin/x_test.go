// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"regexp"
	"strings"
	"testing"
	"text/scanner"

	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/types"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/regex"
)

func FuzzNumberQ(f *testing.F) {
	// to run: go test -fuzz=FuzzNumberQ -run=FuzzNumberQ
	f.Fuzz(func(t *testing.T, s string) {
		s2 := s
		if strings.HasPrefix(s, "-") || strings.HasPrefix(s, "+") {
			s2 = s2[1:]
		}
		if len(s2) >= 2 && s2[0] == '0' && s2[1] != 'x' {
			return
		}
		if strings.Contains(s, "p") || strings.Contains(s, "P") {
			return
		}
		assert.Msg(s).This(string_NumberQ(SuStr(s)) == True).Is(gonum(s))
	})
}

func TestGonum(t *testing.T) {
	assert.True(gonum("0x_0"))
}

func gonum(src string) bool {
	var s scanner.Scanner
	s.Init(strings.NewReader(src))
	s.Whitespace = 0 // do not skip whitespace
	s.Mode = scanner.ScanInts | scanner.ScanFloats
	err := false
	s.Error = func(s *scanner.Scanner, msg string) {
		// fmt.Println(msg)
		err = true
	}
	tok := s.Scan()
	if tok == '+' || tok == '-' {
		tok = s.Scan()
	}
	if s.Scan() != scanner.EOF || err {
		return false
	}
	return tok == scanner.Int || tok == scanner.Float
}

func TestNumberQ(t *testing.T) {
	// NOTE: same test is in stdlib:StringNumberTest
	compile := func(s string) (result bool) {
		defer func() {
			if e := recover(); e != nil {
				result = false
			}
		}()
		v := compile.Constant(s)
		return v.Type() == types.Number
	}
	numq := func(s string) bool {
		return string_NumberQ(SuStr(s)) == True
	}
	test := func(s string, expected bool) {
		t.Helper()
		assert.T(t).Msg(s).This(compile(s)).Is(expected)
		assert.T(t).Msg(s).This(numq(s)).Is(expected)
		if expected == true {
			assert.Msg("+" + s).True(numq("+" + s))
			assert.Msg("-" + s).True(numq("-" + s))
			assert.Msg("x" + s).False(numq("x" + s))
			assert.Msg(s + "x").False(numq(s + "x"))
			assert.Msg(" " + s).False(numq(" " + s))
			assert.Msg(s + " ").False(numq(s + " "))
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
	test("123e-", false)
	test("123e+", false)
	test("1\x00", false)
	test(".0.", false)

	test("", false)
	test(".", false)
	test("+", false)
	test("-", false)
	test("-.", false)
	test("+-.", false)
	test("1.2.3", false)
	test("e5", false)

	test("_1", false)
	test("1_", false)
	test("1_2", true)
	test("1__2", false)
	test("1_2_3", true)
	test("1_2_3_", false)
	test("1_2_3__", false)
	test("1_2_3__4", false)
	test("1_2_3__4_", false)
	test("1.0", true)
	test("1._0", false)
	test("1_.0", false)
	test("1.0_", false)
	test("1._0_", false)
	test("1_.0_", false)
	test("1.0__", false)
	test("1.__0", false)
	test("1.__0_", false)
	test("1.__0__", false)
	test("1.0__0", false)
	test("__1.0__0_", false)
	test("__1.0__0__", false)

	test("0x", false)
	test("0x", false)
	test("0x_", false)
	test("_0x", false)
	test("_0x_", false)
	test("0x.", false)
	test("0x123", true)
	test("0x123.456", false)
	test("0x1_2_3", true)
	test("0x1_2_3_", false)
	test("0x1_2_3__", false)
	test("0x1_2_3__4", false)
	test("0xZ12", false)
	test("0x_0", true) // consistent with Go
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
