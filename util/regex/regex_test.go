// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/ptest"
)

func FuzzRegex(f *testing.F) {
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
		pat := Compile(s)
		pat.Matches("Hello World")
		pat.Matches("")
	})
}

func FuzzRegex2(f *testing.F) {
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
		pat.Matches(s)
	})
}

func FuzzRegexCmp(f *testing.F) {
	sRep := strings.NewReplacer(
		"\r", "r",
		"\n", "n",
	)
	test := func(t *testing.T, r, s string) {
		p1 := Compile(r)
		var result Result
		p1.FirstMatch(s, 0, &result)
		i1, j1 := result[0].Range()

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
				t.Error("pattern:", r, "=>", e)
			}
		}()
		r = fixRegex(r)
		s = sRep.Replace(toAscii(s))
		test(t, r, s)
		test(t, "(?i)"+r, s)
	})
}

func FuzzRegexVsGo(f *testing.F) {
	test := func(t *testing.T, r, s string) {
		p1 := Compile(r)
		var result Result
		p1.FirstMatch(s, 0, &result)
		i1, j1 := result[0].Range()

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

func TestRegexBug(t *testing.T) {
	pat := ".*0|1"
	str := "xxx0"

	p1 := Compile(pat)
	p1.print()
	var result Result
	p1.FirstMatch(str, 0, &result)
	i1, j1 := result[0].Range()
	fmt.Println("SU", i1, j1)

	p2 := regexp.MustCompile(pat)
	i2, j2 := -1, -1
	if m2 := p2.FindStringSubmatchIndex(str); m2 != nil {
		i2, j2 = m2[0], m2[1]
	}
	fmt.Println("GO", i2, j2)

	assert.That(i1 == i2 && j1 == j2)
}

func BenchmarkRegex(b *testing.B) {
	pat := Compile(".+foo")
	var r Result
	s := strings.Repeat("helloworld", 11) + "fo"
	for n := 0; n < b.N; n++ {
		pat.FirstMatch(s, 0, &r)
	}
}

func BenchmarkRegexChars(b *testing.B) {
	pat := Compile("foo")
	var r Result
	s := strings.Repeat("helloworld", 11) + "foo"
	for n := 0; n < b.N; n++ {
		pat.FirstMatch(s, 0, &r)
	}
}

func BenchmarkRegexStart(b *testing.B) {
	pat := Compile(`\Afoo`)
	var r Result
	s := strings.Repeat("helloworld", 11) + "\nfoo"
	for n := 0; n < b.N; n++ {
		pat.FirstMatch(s, 0, &r)
	}
}

func TestRegex(t *testing.T) {
	pat := Compile(".+foo")
	var r Result
	assert := assert.T(t).This
	assert(pat.match("foo", 0, 0, &r)).Is(-1)
	assert(pat.match("", 0, 0, &r)).Is(-1)
	assert(pat.match("hello", 0, 0, &r)).Is(-1)
	assert(pat.match("xfoo", 0, 0, &r)).Is(0)
	assert(pat.match("hifoo", 0, 0, &r)).Is(0)
	assert(pat.match("hifoobar", 0, 0, &r)).Is(0)
}

func TestCapture(t *testing.T) {
	pat := Compile("is")
	s := "now is the time"
	var r Result
	pat.FirstMatch(s, 0, &r)
	assert.T(t).This(r[0].Part(s)).Is("is")
}

func ExamplePattern_ForEachMatch() {
	pat := Compile(`\w+`)
	s := "now is the time"
	pat.ForEachMatch(s, func(r *Result) bool {
		fmt.Println(r[0].Part(s))
		return true
	})
	// Output:
	// now
	// is
	// the
	// time
}

func TestForEachMatch(t *testing.T) {
	s := `one
		two`
	pat := Compile(`^ *`)
	n := 0
	pat.ForEachMatch(s, func(*Result) bool { n++; return true })
	assert.T(t).This(n).Is(2)
}

func TestMatchBug(t *testing.T) {
	pat := Compile("^Date: .*")
	var result Result
	pat.FirstMatch("foo\nDate: Fri, 12 Jul 2019 16:31:35 GMT\r\nbar", 0, &result)
	assert.T(t).This(result[0].pos1).Is(4 + 1)
	assert.T(t).This(result[0].end).Is(39)
}

// ptest support ---------------------------------------------------------------

func TestPtest(t *testing.T) {
	if !ptest.RunFile("regex.test") {
		t.Fail()
	}
}

// pt_match is a ptest for matching
// simple usage is two arguments, string and pattern
// an optional third argument can be "false" for matches that should fail
// or additional arguments can specify expected \0, \1, ...
func ptMatch(args []string, _ []bool) bool {
	pat := Compile(args[1])
	var res Result
	result := pat.FirstMatch(args[0], 0, &res) != -1
	if len(args) > 2 {
		if args[2] == "false" {
			result = !result
		} else {
			for i, e := range args[2:] {
				p := ""
				if res[i].pos1 != 0 {
					p = args[0][res[i].pos1-1 : res[i].end]
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
	var res Result
	result := pat.FirstMatch(s, 0, &res)
	if result == -1 {
		return false
	}
	r := Replace(s, rep, &res)
	pos, end := res[0].Range()
	t := s[:pos] + r + s[end:]
	if t != expected {
		fmt.Println("\t     got:", t, "\n\texpected:", expected)
		return false
	}
	return true
}

var _ = ptest.Add("regex_replace", ptReplace)

// ptest support ---------------------------------------------------------------

func (pat Pattern) print() {
	for i, in := range pat {
		if in.op == branch || in.op == jump {
			in.jump += int16(i)
			in.alt += int16(i)
		}
		fmt.Printf("%d: %s\n", i, in.String())
	}
}
