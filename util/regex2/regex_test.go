// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex2

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/ptest"
)

func TestPlay(t *testing.T) {
	s := "abcd"
	pat := Compile(`(..)(..)`)
	fmt.Println(pat)
	// pat := Compile(`ab|abcd`) "xyz\r\n\r\nxyz", "^[^x].*$", false

	// fmt.Println(">>> part", pat.Matches(s))
	// fmt.Println(">>> full", pat.Match(s, nil))
	var cap Captures
	fmt.Println(">>> full capture", pat.Match(s, &cap))
	// cap.Print(s)
	// fmt.Println(">>> first capture", pat.FirstMatch(s, &cap))
	// fmt.Println(">>> ", pat[npre:].prefixMatch(s, &cap, false))
	cap.Print(s)
}

func (c *Captures) Print(s string) {
	// fmt.Println(cap)
	for i := 0; i < 20; i += 2 {
		if i == 0 || c[i] > 0 || c[i+1] > 0 {
			fmt.Printf("%d: %q\n", i/2, s[c[i]:c[i+1]])
		}
	}
}

func TestGoRegexp(t *testing.T) {
	pat := regexp.MustCompile(`(?m)^[^x].*$`)
	cap := pat.FindStringSubmatch("xyz\r\n\r\nxyz")
	fmt.Printf("%#v\n", cap)
}

func TestCapture(t *testing.T) {
	test := func(pat, str string, expected ...string) {
		t.Helper()
		var cap Captures
		Compile(pat).Match(str, &cap)
		// cap.Print(str)
		for i, e := range expected {
			assert.T(t).This(str[cap[i+2]:cap[i+3+1]]).Is(e)
		}
	}
	test("(.+)(.+)", "abcd", "abc", "d")
	test("(.+?)(.+)", "abcd", "a", "bcd")
	test("(.*)(.*)", "abcd", "abcd", "")
	test("(.*?)(.*)", "abcd", "", "abcd")
}

func TestMatch(t *testing.T) {
	yes := func(pat string, str string) {
		t.Helper()
		assert.T(t).True(Compile(pat).Match(str, nil))
	}
	no := func(pat string, str string) {
		t.Helper()
		assert.T(t).False(Compile(pat).Match(str, nil))
	}
	yes("a", "a")
	no("a", "")
	no("a", "b")

	yes(".", "a")
	no(".", "")

	yes("a|b", "a")
	yes("a|b", "b")
	no("a|b", "")
	no("a|b", "c")

	yes("a?", "")
	yes("a?", "a")
}

func ExampleCompile() {
	test := func(rx string) {
		fmt.Printf("/%v/\n%v\n", rx, Compile(rx)[7:])
	}
	test("abc")
	test("a|b")
	test("ab?c")
	test("ab+c")
	test("ab*c")
	// Output:
	// /abc/
	// 0: Save 0
	// 2: Char a
	// 4: Char b
	// 6: Char c
	// 8: Done
	//
	// /a|b/
	// 0: Save 0
	// 2: SplitFirst 10
	// 5: Char a
	// 7: Jump 12
	// 10: Char b
	// 12: Done
	//
	// /ab?c/
	// 0: Save 0
	// 2: Char a
	// 4: SplitLast 9
	// 7: Char b
	// 9: Char c
	// 11: Done
	//
	// /ab+c/
	// 0: Save 0
	// 2: Char a
	// 4: Char b
	// 6: SplitFirst 4
	// 9: Char c
	// 11: Done
	//
	// /ab*c/
	// 0: Save 0
	// 2: Char a
	// 4: SplitLast 12
	// 7: Char b
	// 9: Jump 4
	// 12: Char c
	// 14: Done
}

func BenchmarkOnePass(b *testing.B) {
	pat := Compile("abc")
	for i := 0; i < b.N; i++ {
		pat.Match("abc", nil)
	}
}

// ptest support ---------------------------------------------------------------

func TestPtest(t *testing.T) {
	if !ptest.RunFile("regex.test") {
		t.Fail()
	}
}

func TestPtest2(t *testing.T) {
	result := ptMatch([]string{"a", "a?", "a"}, nil)
	fmt.Println(result)
}

// pt_match is a ptest for matching
// simple usage is two arguments, string and pattern
// an optional third argument can be "false" for matches that should fail
// or additional arguments can specify expected \0, \1, ...
func ptMatch(args []string, _ []bool) bool {
	s := args[0]
	pat := Compile(args[1])
	var cap Captures
	result := pat.FirstMatch(s, &cap)
	if len(args) > 2 {
		if args[2] == "false" {
			result = !result
		} else {
			for i, e := range args[2:] {
				p := s[cap[i*2]:cap[i*2+1]]
				result = result && (e == p)
			}
		}
	}
	return result
}

var _ = ptest.Add("regex_match", ptMatch)

/*
// pt_replace is a ptest for regex replace
func ptReplace(args []string, _ []bool) bool {
	s := args[0]
	pat := Compile(args[1])
	rep := args[2]
	expected := args[3]
	var cap Captures
	result := pat.FirstMatch(s, &cap)
	if !result {
		return false
	}
	r := Replacement(s, rep, &cap)
	pos, end := cap[0].Range()
	t := s[:pos] + r + s[end:]
	if t != expected {
		fmt.Println("\t     got:", t, "\n\texpected:", expected)
		return false
	}
	return true
}

var _ = ptest.Add("regex_replace", ptReplace)
*/
