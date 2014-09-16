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
			if i > 8 {
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
		ok = ok && p.run1()
	}
	return ok
}

func (p *parser) run1() bool {
	p.match(c.AT, false) // '@'
	name := p.Text
	p.match(c.IDENTIFIER, true)
	fmt.Println(name + ":")
	test, present := testmap[name]
	if !present {
		fmt.Println("\tMISSING")
		test = func(args []string) bool { return true }
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
		if !test(row) {
			ok = false
			fmt.Println("\tFAILED: ", row)
		}
		p.next(true)
		n++
	}
	if ok {
		fmt.Printf("\tok (%d)\n", n)
	}
	return ok
}

func (p *parser) match(expected c.Token, skip bool) {
	if p.Token != expected && p.Keyword != expected {
		panic("syntax error on " + p.Text)
	}
	p.next(skip)
}

func (p *parser) next(skip bool) {
	for {
		p.Item = p.lxr.Next()
		switch p.Token {
		case c.NEWLINE:
			if !skip {
				return
			}
		case c.WHITESPACE, c.COMMENT:
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
