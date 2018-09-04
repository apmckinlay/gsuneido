/*
Package ptest runs test cases listed in text files.

This is so the same test cases can be shared between cSuneido, jSuneido,
gSuneido, && from within Suneido.

Uses compile.Lexer to parse the test files.
*/
package ptest

import (
	"fmt"
	"io/ioutil"
	"strings"

	c "github.com/apmckinlay/gsuneido/lexer"
)

type parser struct {
	lxr *c.Lexer
	c.Item
	comment string
}

var tdir string

func testdir() string {
	if tdir == "" {
		// first time, read and cache
		file := "ptestdir.txt"
		for i := 0; ; i++ {
			src, err := ioutil.ReadFile(file)
			if err == nil {
				tdir = strings.TrimSpace(string(src))
				break
			}
			if i > 9 {
				panic("can't find ptestdir.txt")
			}
			file = "../" + file
		}
	}
	return tdir
}

func RunFile(filename string) bool {
	src, err := ioutil.ReadFile(testdir() + filename)
	if err != nil {
		panic("can't read " + testdir() + filename)
	}
	lxr := c.NewLexer(string(src))
	p := parser{lxr: lxr}
	p.next(true)
	return p.run()
}

func (p *parser) run() bool {
	ok := true
	for p.Token != c.EOF {
		ok = p.runFixture() && ok
	}
	return ok
}

func (p *parser) runFixture() bool {
	p.match(c.AT, false) // '@'
	name := p.Text
	p.match(c.IDENTIFIER, true)
	fmt.Println(name+":", p.comment)
	test, present := testmap[name]
	if !present {
		fmt.Println("\tMISSING TEST FIXTURE")
		test = nil
	}
	n := 0
	ok := true
	for p.Token != c.EOF && p.Token != c.AT {
		row := []string{}
		for {
			text := p.Text
			if p.Token == c.SUB || p.Token == c.ADD {
				p.next(false)
				text += p.Text
			}
			row = append(row, text)
			p.next(false)
			if p.Text == "," {
				p.next(true)
			}
			if p.Token == c.EOF || p.Token == c.NEWLINE {
				break
			}
		}
		if test != nil {
			if !runCase(test, row) {
				// ok = false
			} else {
				n++
			}
		}
		p.next(true)
	}
	if test != nil {
		fmt.Printf("\t%d passed\n", n)
	}
	return ok
}

func runCase(test testfn, row []string) (ok bool) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("\tFAILED: ", row)
			fmt.Println("\tthrew: ", err)
			ok = false
		}
	}()
	ok = true
	if !test(row) {
		fmt.Println("\tFAILED: ", row)
		ok = false
	}
	return
}

func (p *parser) match(expected c.Token, skip bool) {
	if p.Token != expected && p.Keyword != expected {
		panic("syntax error on " + p.Text)
	}
	p.next(skip)
}

func (p *parser) next(skip bool) {
	p.comment = ""
	nl := false
	for {
		p.Item = p.lxr.Next()
		switch p.Token {
		case c.NEWLINE:
			if !skip {
				return
			}
			nl = true
		case c.WHITESPACE:
			continue
		case c.COMMENT:
			// capture trailing comment on same line
			if !nl {
				p.comment = p.Item.Text
			}
			continue
		default:
			return
		}
	}
}

type testfn func([]string) bool

var testmap = make(map[string]testfn)

// Add is used to add test functions
//
// Other packages normally add tests in their init
// or by e.g. var _ = ptest.Add(...)
func Add(name string, fn testfn) string {
	testmap[name] = fn
	return name
}
