// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex

import (
	"fmt"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
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
func pt_match(args []string, _ []bool) bool {
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
