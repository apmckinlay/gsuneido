// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex

import (
	"fmt"
	"regexp"
	"regexp/syntax"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/ptest"
)

func TestPlay(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	// s := "now is the time"
	pat := Compile(`\A_?[[:alpha:]]\w*[!?]?\Z`)
	fmt.Println(pat)
	// var cap Captures
	// ok := pat.LastMatch(s, len(s), &cap)
	// fmt.Println(">>>", ok)
	// fmt.Println("cap", cap)
	// cap.Print(s)
}

func (c *Captures) Print(s string) {
	// fmt.Println(cap)
	for i := 0; i < 20; i += 2 {
		if c[i] > 0 || c[i+1] > 0 {
			fmt.Printf("%d: %q\n", i/2, s[c[i]:c[i+1]])
		}
	}
}

func TestAll(t *testing.T) {
	test := func(pat, s string, expected ...int32) {
		var matches []int32
		for cap := range Compile(pat).All(s) {
			matches = append(matches, cap[0])
			//fmt.Println(cap[0], str.BeforeFirst(s[cap[0]:], "\n"))
		}
		assert.T(t).This(matches).Is(expected)
	}
	test("^", "one\ntwo\r\nthree", 0, 4, 9)
	test(`^ *`, "function (a, b = 1)\n{\nSteppingDebugger(0);\n c = a + b\nSteppingDebugger(1);\n return c\n }",
		0, 20, 22, 43, 54, 75, 85)
}

func TestGoRegexp(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	pat := regexp.MustCompile(`0b|0[bB]`)
	cap := pat.FindStringIndex("0B")
	fmt.Printf("%#v\n", cap)
}

func TestGoParse(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	re, _ := syntax.Parse("^+", 0)
	//re = re.Simplify()
	prog, _ := syntax.Compile(re)
	fmt.Println(prog)
}

func TestLeftAnchored(t *testing.T) {
	test := func(rx string, expected bool) {
		t.Helper()
		// fmt.Printf("`%s`\n%v\n", rx, Compile(rx))
		assert.T(t).This(Compile(rx).leftAnchored()).Is(expected)
	}
	test(``, false)      // literal substr
	test(`foo`, false)   // literal substr
	test(`foo\Z`, false) // literal suffix
	test(`\w`, false)    // one pass
	test(`.*`, false)    // full

	test(`\A`, true)      // literal prefix
	test(`\Afoo`, true)   // literal prefix
	test(`\Afoo\Z`, true) // literal equal
	test(`\A\w`, true)    // one pass
	test(`\A.*`, true)    // full

	pat := Compile(`\Afoo`)
	var cap Captures
	assert.T(t).False(pat.LastMatch("other foo", 8, &cap))
	assert.T(t).True(pat.LastMatch("foo other foo", 8, &cap))
	assert.T(t).This(cap[0]).Is(0)
}

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
	var rt rune
	match := func(str string, pat string, expected bool) {
		t.Helper()
		// fmt.Printf("%q =~ %q -> %v\n", str, pat, expected)
		rx := Compile(pat)
		assert.This(rxType(rx)).Is(rt)
		assert.T(t).This(rx.Match(str, nil)).Is(expected)
	}
	full := func(pat string) string {
		if strings.Contains(pat, "|") {
			return `\A(` + pat + `)\Z`
		}
		return `\A` + pat + `\Z`
	}
	rt = 'L'
	match("a", full("a"), true)
	match("a", "a", true)
	match("abc", "b", true)
	match("abc", `\Ab`, false)
	match("abc", `\Aab`, true)
	match("abc", `bc\Z`, true)
	match("abc", `b\Z`, false)
	match("abc", full("abc"), true)
	match("abc", full("b"), false)
	match("a", full(""), false)
	match("", full("a"), false)
	match("a", full("b"), false)

	rt = '1'
	match("a", full("."), true)
	match("", full("."), false)
	match("abc", full(".bc"), true)
	match("abc", full(".bx"), false)
	match("a", full("a|b"), true)
	match("b", full("a|b"), true)
	match("", full("a|b"), false)
	match("c", full("a|b"), false)
	match("", full("a?"), true)
	match("a", full("a?"), true)

	rt = 'M'
	match("a", "a|b", true)
	match("b", "a|b", true)
	match("", "a|b", false)
	match("c", "a|b", false)
	match("", "a?", true)
	match("a", "a?", true)
}

func rxType(pat Pattern) rune {
	switch opType(pat[0]) {
	case opLiteralSubstr, opLiteralPrefix, opLiteralSuffix, opLiteralEqual:
		return 'L'
	case opOnePass:
		return '1'
	default:
		return 'M'
	}
}

func TestCompile(t *testing.T) {
	test := func(rx, expected string) {
		t.Helper()
		assert.T(t).Msg(rx).This(Compile(rx).String()).Like(expected)
	}
	test("xyz",
		`0: LiteralSubstr "xyz"`)
	test(`\Axyz`,
		`0: LiteralPrefix "xyz"`)
	test(`xyz\Z`,
		`0: LiteralSuffix "xyz"`)
	test(`\Axyz\Z`,
		`0: LiteralEqual "xyz"`)
	test(`.`,
		`0: AnyNotNL
		1: DoneSave1`)
	test(`\A.`,
		`0: OnePass
		1: StrStart
		2: AnyNotNL
		3: DoneSave1`)
	test(`[a]`,
		`0: LiteralSubstr "a"`)
	test(`(?i)[a]`,
		`0: ListSet "Aa"
		4: DoneSave1`)
	test(`[^a]`,
		`0: FullSet
		33: DoneSave1`)
	test("a|b",
		`0: SplitNext 8
		3: Char a
		5: Jump 10
		8: Char b
		10: DoneSave1`)
	test("a|b|c",
		`0: SplitNext 8
		3: Char a
		5: Jump 18
		8: SplitNext 16
		11: Char b
		13: Jump 18
		16: Char c
		18: DoneSave1`)
	test("ab?c",
		`0: Prefix "a"
		3: Char a
		5: SplitNext 10
		8: Char b
		10: Char c
		12: DoneSave1`)
	test("ab+c",
		`0: Prefix "ab"
		4: Char a
		6: Char b
		8: SplitJump 6
		11: Char c
		13: DoneSave1`)
	test("ab*c",
		`0: Prefix "a"
		3: Char a
		5: SplitNext 13
		8: Char b
		10: SplitJump 8
		13: Char c
		15: DoneSave1`)
}

func TestRightAnchor(t *testing.T) {
	test := func(pat string, expected bool) {
		t.Helper()
		co := compile(pat)
		assert.T(t).Msg(pat).This(co.rightAnchor).Is(expected)
	}
	test(`abc`, false)
	test(`abc\Z`, true)
	test(`(abc)\Z`, true)
	test(`a|b\Z`, false)
	test(`(a|b)\Z`, true)
	test(`(a\Z)?`, false)
	test(`(a\Z)+`, false)
	test(`(a\Z)*`, false)
	test(`a\Zb`, false)
	test(`a\Z?`, false)
}

func TestOnePass(t *testing.T) {
	// pat := "^0|1$"
	// fmt.Printf("pat %q\n", pat)
	// co := compile(pat)
	// fmt.Println(Pattern(co.prog))
	// if co.onePass1() {
	// 	fmt.Println("onePass1 true")
	// 	if co.onePass2() {
	// 		fmt.Println("onePass2 true")
	// 		co.onePass3()
	// 		fmt.Println(Pattern(co.prog))
	// 	} else {
	// 		fmt.Println("onePass2 false")
	// 	}
	// } else {
	// 	fmt.Println("onePass1 false")
	// }
	// t.SkipNow()

	test := func(rx string, expected bool) {
		co := compile(rx)
		assert.This(co.onePass()).Is(expected)
	}
	test(`abc`, false)
	test(`\Aabc`, true)
	test(`\Aabc\Z`, true)
	test(`\A(a|b)`, true)
	test(`\A(a|bc|de|f)`, true)
	test(`\Aa?`, false)
	test(`\Aa?\Z`, true)
	test(`\Aa+\Z`, true)
	test(`\Aa*\Z`, true)
	test(`\Aab?c`, true)
	test(`\Aab+c`, true)
	test(`\Aab*c`, true)
	test(`\A(0*)*7`, false)
	test(`\A(((0))*)(((\Z))*1)`, false)
}

var M bool

func BenchmarkOnePass(b *testing.B) {
	pat := Compile("^a+b")
	for b.Loop() {
		M = pat.Match("aaab", nil)
	}
}

// fuzzing ----------------------------------------------------------

func FuzzCompile(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string) {
		defer func() {
			if e := recover(); e != nil {
				if err, ok := e.(string); ok &&
					strings.HasPrefix(err, "regex: ") {
					return
				}
				t.Error("pattern:", s, "=>", e)
			}
		}()
		Compile(s)
	})
}

func FuzzRegex(f *testing.F) {
	f.Fuzz(func(t *testing.T, r, s string) {
		defer func() {
			if e := recover(); e != nil {
				if err, ok := e.(string); ok &&
					strings.HasPrefix(err, "regex: ") {
					return
				}
				t.Error("pattern:", s, "=>", e)
			}
		}()
		pat := Compile(r)
		pat.Match(s, nil)
		pat.Match(s, &Captures{})
	})
}

func FuzzRegexCmp(f *testing.F) {
	sRep := strings.NewReplacer(
		"\r", "r",
		"\n", "n",
	)
	test := func(t *testing.T, r, s string) {
		p1 := Compile(r)
		var cap Captures
		slc.Fill(cap[:], -1)
		p1.Match(s, &cap)
		i1, j1 := int(cap[0]), int(cap[1])

		p2, err := regexp.Compile(r)
		if err != nil {
			return
		}
		i2, j2 := -1, -1
		if m2 := p2.FindStringSubmatchIndex(s); m2 != nil {
			i2, j2 = m2[0], m2[1]
		}

		if i1 != i2 || (i1 != -1 && j1 != j2) {
			t.Errorf("r: %q s: %q Suneido: %d,%d Go: %d,%d", r, s, i1, j1, i2, j2)
		}
	}
	f.Fuzz(func(t *testing.T, r, s string) {
		defer func() {
			if e := recover(); e != nil {
				if err, ok := e.(string); ok &&
					strings.HasPrefix(err, "regex: ") {
					return
				}
				if err, ok := e.(error); ok &&
					strings.Contains(err.Error(), "invalid UTF-8") {
					return
				}
				t.Error("pattern:", r, "=>", e)
			}
		}()
		r = fixRegex(r)
		s = sRep.Replace(toAscii(s))
		test(t, r, s)
		// test(t, "(?i)"+r, s)
	})
}

func FuzzRegexVsGo(f *testing.F) {
	test := func(t *testing.T, r, s string) {
		p1 := Compile(r)
		var cap Captures
		slc.Fill(cap[:], -1)
		p1.Match(s, &cap)
		i1, j1 := int(cap[0]), int(cap[1])

		p2, err := regexp.Compile(r)
		if err != nil {
			return
		}
		i2, j2 := -1, -1
		if m2 := p2.FindStringSubmatchIndex(s); m2 != nil {
			i2, j2 = m2[0], m2[1]
		}

		if i1 != i2 || (i1 != -1 && j1 != j2) {
			t.Errorf("r: %q s: %q Suneido: %d,%d Go: %d,%d", r, s, i1, j1, i2, j2)
		}
	}
	f.Fuzz(func(t *testing.T, r string) {
		defer func() {
			if e := recover(); e != nil {
				if err, ok := e.(string); ok &&
					strings.HasPrefix(err, "regex: ") {
					return
				}
				if err, ok := e.(error); ok &&
					strings.Contains(err.Error(), "invalid UTF-8") {
					return
				}
				t.Error("pattern:", r, "=>", e)
			}
		}()
		r = fixRegex(r)
		test(t, r, "")
		test(t, r, "x")
		test(t, r, "Hello World")
		test(t, "(?i)"+r, "x")
		test(t, "(?i)"+r, "Hello World")
	})
}

var rRep = strings.NewReplacer(
	"{", "(",
	"}", ")",
	"(?", "(",
	"$*", ",",
	"^*", ";",
	"$+", ",",
	"^+", ";",
	"$?", ",",
	"^?", ";",
	`\(`, "(",
	`\`, "",
	"[[:ascii:]]", "a",
	"[[:word]]", "w",
	`\Q`, "Q",
	`\E`, "E",
	`\b`, "b",
	`\B`, "B",
	`\<`, "<",
	`\>`, ">",
	`\p`, "p",
	`\P`, "P",
	`\z`, "z",
	`\Z`, "Z",
	`\a`, "a",
	`\x`, "x",
	`\0`, "0",
	`\1`, "1",
	`\2`, "2",
	`\3`, "3",
	`\4`, "4",
	`\5`, "5",
	`\6`, "6",
	`\7`, "7",
	`\8`, "8",
	`\9`, "9",
)

func fixRegex(s string) string {
	s = toAscii(s)
	for {
		s2 := rRep.Replace(s)
		if s2 == s {
			return s
		}
		s = s2
	}
}

func toAscii(s string) string {
	var sb strings.Builder
	for _, b := range []byte(s) {
		sb.WriteByte(b & 0x7f)
	}
	return sb.String()
}

// replace ----------------------------------------------------------

func TestLiteralRep(t *testing.T) {
	test := func(rep string, expected any) {
		x, ok := LiteralRep(rep)
		if !ok {
			assert.Msg(rep).That(expected == false)
		} else {
			assert.Msg(rep).This(x).Is(expected)
		}
	}
	test("", "")
	test("x", "x")
	test("hello world", "hello world")
	test(`\=hello world`, "hello world")
	test("x&y", false)
	test(`\2 \1`, false)
}

func BenchmarkMatch(b *testing.B) {
	s := strings.Repeat("helloworld", 1000)
	pat := Compile("x|y|z")
	b.ResetTimer()
	for b.Loop() {
		pat.match(s, 0, &Captures{}, false)
	}
}

// ptest support ----------------------------------------------------

func TestPtest(t *testing.T) {
	if !ptest.RunFile("regex.test") {
		t.Fail()
	}
}

func TestPtest2(t *testing.T) {
	result := ptReplace([]string{`now is the time`, `is`, `&&`, `now isis the time`}, nil)
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
	result := pat.Match(s, &cap)
	if len(args) > 2 {
		if args[2] == "false" {
			result = !result
		} else {
			for i, e := range args[2:] {
				p := ""
				if cap[i*2] >= 0 {
					p = s[cap[i*2]:cap[i*2+1]]
				}
				result = result && (e == p)
			}
		}
	}
	return result
}

var _ = ptest.Add("regex_match", ptMatch)

// pt_replace is a ptest for regex replace
func ptReplace(args []string, _ []bool) bool {
	s := args[0]
	pat := Compile(args[1])
	rep := args[2]
	expected := args[3]
	var cap Captures
	result := pat.Match(s, &cap)
	if !result {
		return false
	}
	r := Replacement(s, rep, &cap)
	t := s[:cap[0]] + r + s[cap[1]:]
	if t != expected {
		fmt.Println("\t     got:", t, "\n\texpected:", expected)
		return false
	}
	return true
}

var _ = ptest.Add("regex_replace", ptReplace)
