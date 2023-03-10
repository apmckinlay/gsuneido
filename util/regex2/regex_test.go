// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex2

import (
	"fmt"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/ptest"
)

func TestPlay(t *testing.T) {
	s := "now is the time for all good men now that"
	pat := Compile(`now\s`)
	// fmt.Println(pat)

	// fmt.Println(">>> part", pat.Matches(s))
	// fmt.Println(">>> full", pat.Match(s, nil))
	var cap Captures
	fmt.Println(">>> full capture", pat.Match(s, &cap))
	fmt.Println(">>> first capture", pat.FirstMatch(s, &cap), cap[0])
	// fmt.Println(">>> match", pat.match(s, &cap, false))
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

// func TestGoRegexp(t *testing.T) {
// 	pat := regexp.MustCompile(`(?m)^[^x].*$`)
// 	cap := pat.FindStringSubmatch("xyz\r\n\r\nxyz")
// 	fmt.Printf("%#v\n", cap)
// }

func TestCapture(t *testing.T) {
	test := func(str, pat string, expected ...string) {
		t.Helper()
		var cap Captures
		Compile(pat).Match(str, &cap)
		// cap.Print(str)
		for i, e := range expected {
			assert.T(t).This(str[cap[i+2]:cap[i+3+1]]).Is(e)
		}
	}
	test("abcd", "(.+)(.+)", "abc", "d")
	test("abcd", "(.+?)(.+)", "a", "bcd")
	test("abcd", "(.*)(.*)", "abcd", "")
	test("abcd", "(.*?)(.*)", "", "abcd")
}

func TestMatch(t *testing.T) {
	match := func(str string, pat string, expected bool) {
		t.Helper()
		assert.T(t).This(Compile(pat).Match(str, nil)).Is(expected)
	}
	matches := func(str string, pat string, expected bool) {
		t.Helper()
		assert.T(t).This(Compile(pat).Matches(str)).Is(expected)
	}
	// literal
	match("a", "a", true)
	matches("a", "a", true)
	matches("abc", "b", true)
	matches("abc", `\Ab`, false)
	matches("abc", `\Aa`, true)
	match("abc", `\Aabc`, true)
	match("a", "", false)
	match("", "a", false)
	match("a", "b", false)

	// one pass
	match("a", ".", true)
	match("", ".", false)
	match("abc", ".bc", true)
	match("abc", ".bx", false)

	// full
	match("a", "a|b", true)
	match("b", "a|b", true)
	match("", "a|b", false)
	match("c", "a|b", false)
	match("", "a?", true)
	match("a", "a?", true)
}

func TestCompile(t *testing.T) {
	test := func(rx, expected string) {
		t.Helper()
		assert.T(t).This(Compile(rx).String()).Like(expected)
	}
	test("xyz",
		`0: Unanchored
		1: Literal "xyz"`)
	test(`\Axyz`,
		`0: Literal "xyz"`)
	test(`.`,
		`0: Unanchored
		1: OnePass
		2: AnyNotNL
		3: DoneSave1`)
	test(`\A.`,
		`0: OnePass
		1: AnyNotNL
		2: DoneSave1`)
	test("a|b",
		`0: Unanchored
		1: SplitFirst 9
		4: Char a
		6: Jump 11
		9: Char b
		11: DoneSave1`)
	test("a|b|c",
		`0: Unanchored
		1: SplitFirst 9
		4: Char a
		6: Jump 19
		9: SplitFirst 17
		12: Char b
		14: Jump 19
		17: Char c
		19: DoneSave1`)
	test("ab?c",
		`0: Unanchored
		1: LitPrefix "a"
		4: Char a
		6: SplitLast 11
		9: Char b
		11: Char c
		13: DoneSave1`)
	test("ab+c",
		`0: Unanchored
		1: LitPrefix "ab"
		5: Char a
		7: Char b
		9: SplitFirst 7
		12: Char c
		14: DoneSave1`)
	test("ab*c",
		`0: Unanchored
		1: LitPrefix "a"
		4: Char a
		6: SplitLast 14
		9: Char b
		11: Jump 6
		14: Char c
		16: DoneSave1`)
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
	// fmt.Println(args)
	s := args[0]
	pat := Compile("(?m)" + args[1])
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
