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
	"github.com/apmckinlay/gsuneido/util/slc"
)

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
		assert.T(t).True(Compile(pat).Match(str, &cap))
		// cap.Print(str)
		for i, e := range expected {
			got := ""
			if cap[i*2] >= 0 {
				got = str[cap[i*2]:cap[i*2+1]]
			}
			assert.T(t).This(got).Is(e)
		}
	}
	test("abcd", "(.+)(.+)", "abcd", "abc", "d")
	test("abcd", "(.+?)(.+)", "abcd", "a", "bcd")
	test("abcd", "(.*)(.*)", "abcd", "abcd", "")
	test("abcd", "(.*?)(.*)", "abcd", "", "abcd")

	test("a", "a?", "a")
	test("a", "a??", "")
	test("aaab", "a*", "aaa")
	test("aaab", "a*?", "")
	test("aaab", "a+", "aaa")
	test("aaab", "a+?", "a")
	test("axb", "[\x00-\xff]+", "axb")
	test("foo123", `([a-z]+)([0-9]+)`, "foo123", "foo", "123")
	test("hello there world", `(\w+ )+`, "hello there ", "there ")
	test("hello world", "hello(x?)", "hello", "")
	test("hello world", `(a...b)|(h...o)`, "hello", "", "hello")
	test("-ccc-", `(aaa)|(bbb)|(ccc)`, "ccc", "", "", "ccc")
}

func TestMatch(t *testing.T) {
	var rt rune
	match := func(str string, pat string, expected bool) {
		t.Helper()
		// fmt.Printf("%q =~ %q -> %v\n", str, pat, expected)
		rx := Compile(pat)
		if rt != 0 {
			assert.This(rxType(rx)).Is(rt)
		}
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

	// Additional tests from ptest/regex.test
	rt = 0
	match("", "", true)
	match("abc", "x", false)
	match("ab", "abc", false)
	match("abc", "^...$", true)
	match("ab\n", "...", false)
	match("abde", "abc+de", false)
	match("abcde", "abc+de", true)
	match("abccde", "abc+de", true)
	match("abccd", "abc+de", false)
	match("abde", "abc?de", true)
	match("abcde", "abc?de", true)
	match("abccde", "abc?de", false)
	match("abe", "ab(cd)*ef", false)
	match("abef", "ab(cd)*ef", true)
	match("abcdef", "ab(cd)*ef", true)
	match("abcdcdcdef", "ab(cd)*ef", true)
	match("abcdcdcde", "ab(cd)*ef", false)
	match("abcx", "(ab*c)*x", true)
	match("abbc", "(ab*c)*", true)
	match("abcabc", "(ab*c)*", true)
	match("acabbbc", "(ab*c)*", true)
	match("abbbcac", "(ab*c)*", true)
	match("acabcabbcx", "(ab*c)*x", true)
	match("a", "a|b|c", true)
	match("b", "a|b|c", true)
	match("c", "a|b|c", true)
	match("x", "a|b|c", false)
	match("", "a|b|c", false)
	match("ab", "a?b", true)
	match("ab", "a??b", true)
	match("aaab", "a*b", true)
	match("aaab", "a*?b", true)
	match("aaab", "a+?b", true)
	match("hello\nworld", `\Ahe`, true)
	match("hello\nworld", `\Awo`, false)
	match("hello\nworld", `ld\Z`, true)
	match("hello\nworld", `lo\Z`, false)
	match("(+*)", `^(+*)$`, false)
	match("(+*)", `^(?q)(+*)(?-q)$`, true)
	match("hello", "eLL", false)
	match("hello", "(?i)eLL", true)
	match("hello", "(?i)eL(?-i)L", false)
	match("foobar", `\<foo`, true)
	match("foobar", `\<foo\>`, false)
	match("foo bar", `\<foo\>`, true)
	match("foobar", `bar\>`, true)
	match("foobar", `\<bar\>`, false)
	match("foo bar", `\<bar\>`, true)
	match("foobar", `\<foobar\>`, true)
	match("foo bar", "(?i)bar", true)
	match("foo Bar", "(?i)bar", true)
	match("123x", "(?i)[a-z]", true)
	match("123X", "(?i)[a-z]", true)
	match("hello\nworld", `\Ahello`, true)
	match("hello\nworld", `\Aworld`, false)
	match("hello\nworld", `world\Z`, true)
	match("hello\nworld", `hello\Z`, false)
	match("hello\r\nworld", `\Ahello`, true)
	match("hello\r\nworld", `\Aworld`, false)
	match("hello\r\nworld", `world\Z`, true)
	match("hello\r\nworld", `hello\Z`, false)
	match("one_1 two_2\nthree_3", `\<one_1\>`, true)
	match("one_1 two_2\nthree_3", `\<two_2\>`, true)
	match("one_1 two_2\nthree_3", `\<three_3\>`, true)
	match("one_1 two_2\r\nthree_3", `\<two_2\>`, true)
	match("one_1 two_2\r\nthree_3", `\<three_3\>`, true)
	match("one_1 two_2\nthree_3", `\<one\>`, false)
	match("one_1 two_2\nthree_3", `\<two\>`, false)
	match("hello", "fred", false)
	match("hello", "h.*o", true)
	match("hello", "[a-z]ello", true)
	match("hello", "[^0-9]ello", true)
	match("hello", "ell", true)
	match("hello", "^ell", false)
	match("hello", "ell$", false)
	match("heeeeeeeello", "^he+llo$", true)
	match("heeeeeeeello", "^he*llo*", true)
	match("hllo", "^he*llo$", true)
	match("hllo", "^he?llo$", true)
	match("heello", "^he?llo$", false)
	match("+123.456", `^[+-][0-9]+[.][0123456789]*$`, true)
	match("0123456789", `^\d+$`, true)
	match("0123456789", `\D`, false)
	match("hello_123", `^\w+$`, true)
	match("hello_123", `\W`, false)
	match("hello \t\r\nworld", `^\w+\s+\w+$`, true)
	match("!@#@!# \r\t{},()[];", `^\W+$`, true)
	match("123adjf!@#", `^\S+$`, true)
	match("123adjf!@#", `\s`, false)
	match("()[]", `^\(\)\[\]$`, true)
	match("hello world", `^(hello|howdy) (\w+)$`, true)
	match("ab", "(a|ab)b", true)
	match("abc", "x*c", true)
	match("abc", "x*$", true)
	match("abc", "x?$", true)
	match("abc", "^x?", true)
	match("abc", "^x*", true)
	match("aBcD", "abcd", false)
	match("aBcD", "(?i)abcd", true)
	match("aBCd", "a(?i)bc(?-i)d", true)
	match("aBCD", "a(?i)bc(?-i)D", true)
	match("ABCD", "a(?i)bc(?-i)d", false)
	match("abc", "a.c", true)
	match("a.c", "(?q)a.c", true)
	match("abc", "(?q)a.c", false)
	match("a.cd", "(?q)a.c(?-q).", true)
	match("abcd", "(?q)a.c(?-q).", false)
	match("abc", "(?q)(", false)
	match("ABC", "(?i)[A-Z]", true)
	match("ABC", "(?i)[a-z]", true)
	match("abc", "(?i)[A-Z]", true)
	match("abc", "(?i)[a-z]", true)
	match("a", "[abc]", true)
	match("b", "[abc]", true)
	match("c", "[abc]", true)
	match("b", "[^abc]", false)
	match("x", "[^abc]", true)
	match("c", `\w`, true)
	match(" ", `\W`, true)
	match(" ", `\w`, false)
	match(" ", `\s`, true)
	match("c", `\S`, true)
	match("c", `\s`, false)
	match("c", `[\w]`, true)
	match(" ", `[\W]`, true)
	match(" ", `[\w]`, false)
	match(" ", `[\s]`, true)
	match("c", `[\S]`, true)
	match("c", `[\s]`, false)
	match("b", "[[:alpha:]]", true)
	match("b", "[[:alnum:]]", true)
	match("b", "[[:print:]]", true)
	match("b", "[[:graph:]]", true)
	match("b", "[[:lower:]]", true)
	match("b", "[[:upper:]]", false)
	match("B", "[[:upper:]]", true)
	match("5", "[[:digit:]]", true)
	match("5", "[[:alnum:]]", true)
	match("5", "[[:alpha:]]", false)
	match("5", "[[:lower:]]", false)
	match("5", "[[:upper:]]", false)
	match("aBc", "[aBc]+", true)
	match("aBc", "(?i)[ABC]+", true)
	match("ABC", "(?i)ABC", true)
	match("ABC", "(?i)abc", true)
	match("abc", "(?i)ABC", true)
	match("abc", "(?i)abc", true)
	match("abc", "(?i)ark", false)
	match("b", "(?i)[abc]", true)
	match("b", "(?i)[ABC]", true)
	match("B", "(?i)[abc]", true)
	match("B", "(?i)[ABC]", true)
	match("@", "(?i)[@#]", true)
	match("m", "[a-z]", true)
	match("-", "[-z]", true)
	match("m", "[-z]", false)
	match("-", "[a-]", true)
	match("m", "[a-]", false)
	match("aZ", "[a-Z]", false)
	match("\r\n", "^\n", false)
	match("xyz\r\n\r\nxyz", `^[^x].*$`, false)
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

func TestReplace(t *testing.T) {
	test := func(s, pat, rep, expected string) {
		t.Helper()
		rx := Compile(pat)
		var cap Captures
		assert.T(t).True(rx.Match(s, &cap))
		r := Replacement(s, rep, &cap)
		got := s[:cap[0]] + r + s[cap[1]:]
		assert.T(t).This(got).Is(expected)
	}

	test("now is the time", "now", "never", "never is the time")
	test("now is the time", "now", `\=never`, "never is the time")
	test("now is the time", "is", "&&", "now isis the time")
	test("now is the time", "is", `&\n`, "now is\n the time")
	test("now is the time", "is", `\&`, "now & the time")
	test("now is the time", "the", `\U&`, "now is THE time")
	test("now is the time", "the", `\u&`, "now is The time")
	test("NOW IS THE TIME", "THE", `\L&`, "NOW IS the TIME")
	test("NOW IS THE TIME", "THE", `\l&`, "NOW IS tHE TIME")
	test("now is the time", `(\w+) (\w+)`, "\\2-\\1", "is-now the time")
	test("now is the time", ` (\w+) (\w+) `, "\\2\\1", "nowtheistime")
}

func BenchmarkMatch(b *testing.B) {
	s := strings.Repeat("helloworld", 1000)
	pat := Compile("x|y|z")
	b.ResetTimer()
	for b.Loop() {
		pat.match(s, 0, &Captures{}, false)
	}
}
