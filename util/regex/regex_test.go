package regex

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/ptest"
)

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
	Assert(t).That(pat.match("foo", 0, 0, &r), Equals(-1))
	Assert(t).That(pat.match("", 0, 0, &r), Equals(-1))
	Assert(t).That(pat.match("hello", 0, 0, &r), Equals(-1))
	Assert(t).That(pat.match("xfoo", 0, 0, &r), Equals(0))
	Assert(t).That(pat.match("hifoo", 0, 0, &r), Equals(0))
	Assert(t).That(pat.match("hifoobar", 0, 0, &r), Equals(0))
}

func TestCapture(t *testing.T) {
	pat := Compile("is")
	s := "now is the time"
	var r Result
	pat.FirstMatch(s, 0, &r)
	Assert(t).That(r[0].Part(s), Equals("is"))
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
	pat := Compile(args[1])
	var res Result
	result := pat.FirstMatch(args[0], 0, &res) != -1
	if len(args) > 2 {
		if args[2] == "false" {
			result = !result
		} else {
			for i, e := range args[2:] {
				result = result && (e == args[0][res[i].pos1-1:res[i].end])
			}
		}
	}
	return result
}

var _ = ptest.Add("regex_match", pt_match)

// pt_replace is a ptest for regex replace
func pt_replace(args []string, _ []bool) bool {
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

var _ = ptest.Add("regex_replace", pt_replace)
