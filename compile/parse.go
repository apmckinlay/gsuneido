package compile

import (
	"fmt"
	"strconv"

	. "github.com/apmckinlay/gsuneido/lexer"
)

func newParser(src string) *parser {
	lxr := NewLexer(src)
	p := &parser{lxr: lxr}
	p.next()
	return p
}

type parser struct {
	lxr *Lexer
	// Item is the current lexical token etc.
	Item
	// nest is used by parse.go to track nesting
	// in order to skip newlines within e.g. parenthesis
	nest int
	// bld is used by expression.go
	// it is needed because expressions are shared by both language and queries
	bld builder
	// newline is true if the current token was preceeded by a newline
	newline bool
	// expectingCompound is used to differentiate control statement body vs. block
	// e.g. if expr {...}
	// set by function.go used by expression.go
	expectingCompound bool
}

/*
eval* methods are helpers so you can match/next after evaluating something
match* methods verify that the current is what is expected and then advance
next* methods just advance
*/

func (p *parser) evalMatch(result T, tok Token) T {
	p.match(tok)
	return result
}

func (p *parser) evalNext(result T) T {
	p.next()
	return result
}

func (p *parser) match(tok Token) {
	p.mustMatch(tok)
	p.next()
}

func (p *parser) matchIf(tok Token) bool {
	if p.isMatch(tok) {
		p.next()
		return true
	}
	return false
}

func (p *parser) mustMatch(tok Token) {
	if !p.isMatch(tok) {
		p.error("expecting ", tok)
	}
}

func (p *parser) isMatch(tok Token) bool {
	return tok == p.Token || tok == p.Keyword
}

// next advances to the next token, setting p.Item
func (p *parser) next() {
	p.newline = false
	for {
		p.Item = p.lxr.Next()
		if p.Token == NEWLINE {
			if p.lxr.AheadSkip(0).Token != Q_MARK {
				p.newline = true
			}
		} else if p.Token != COMMENT && p.Token != WHITESPACE {
			break
		}
	}
	if p.Token == STRING && p.Keyword != STRING {
		// make a copy of strings that are slices of the source
		p.Text = " " + p.Text
		p.Text = p.Text[1:]
	}
}

// error panics with "syntax error at " + position
// It claims to return string so it can be called inside panic
// (so compiler knows we don't return)
func (p *parser) error(args ...interface{}) string {
	panic("syntax error at " + strconv.Itoa(int(p.Item.Pos)) + " " +
		fmt.Sprint(args...))
}
