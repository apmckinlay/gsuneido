package regex

import (
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/ptest"
)

func TestCompile(t *testing.T) {
	test := func(input, expected string) {
		//fmt.Println("input '" + input + "'")
		p := Compile(input).String()
		Assert(t).That(strings.TrimSpace(p[5:len(p)-7]),
			Equals(expected).Comment(input))
	}
	test("", "")
	test(".", "[...]")
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
	test("(?i)x.y(?-i)z", "i'x' [...] i'y' 'z'")

	test("(?q).*(?-q)def", "'.*def'")

	test("\\<Foo\\>", "\\< 'Foo' \\>")

	test("a(bc)d", "'a' Left1 'bc' Right1 'd'")

	test(".\\5.", "[...] \\5 [...]")
	test("(?i)\\5", "i\\5")

	test("a[.]b", "'a.b'")

	test("a(?q).(?-q).c(?q).(?-q).", "'a.' [...] 'c.' [...]")

	test("\\", "'\\'")

	Assert(t).That(func() { Compile("(abc") }, Panics("missing ')'"))
	Assert(t).That(func() { Compile("abc)def") },
		Panics("closing ) without opening ("))
}

func TestBug(t *testing.T) {
	p := Compile("(?i)[\x9a\xbb]")
	Assert(t).That(p.Matches("\x8a"), Equals(false))
}

func TestForEachMatch(t *testing.T) {
	test := func(s, p string, expected ...string) {
		pat := Compile(p)
		ob := []string{}
		pat.ForEachMatch(s, func(r *Result) bool {
			ob = append(ob, r.Group(s, 0))
			return len(ob) < 4
		})
		Assert(t).That(ob, Equals(expected))
	}
	test("now is the time", `\w+`, "now", "is", "the", "time")
	test("hello", `.`, "h", "e", "l", "l")
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
// or additional arguments can specify \0, \1, ...
func pt_match(args []string, _ []bool) bool {
	var res Result
	pat := Compile(args[1])
	result := pat.FirstMatch(args[0], 0, &res)
	if len(args) > 2 {
		if args[2] == "false" {
			result = !result
		} else {
			for i, e := range args[2:] {
				result = result && (e == args[0][res.pos[i]:res.end[i]])
			}
		}
	}
	return result
}

var _ = ptest.Add("regex_match", pt_match)
