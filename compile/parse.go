package compile

import (
	"fmt"
	"strconv"

	"github.com/apmckinlay/gsuneido/compile/ast"
	. "github.com/apmckinlay/gsuneido/lexer"
)

func newParser(src string) *parser {
	lxr := NewLexer(src)
	p := &parser{lxr: lxr, Factory: ast.Folder{ast.Builder{}}}
	p.next()
	return p
}

type parser struct {
	lxr *Lexer

	// Item is the current lexical token etc.
	Item
	
	// Factory is used by expression.go
	// because expressions are shared by both language and queries
	// and generate different types of AST nodes
	ast.Factory

	// newline is true if the current token was preceeded by a newline
	newline bool

	// expectingCompound is used to differentiate control statement body vs. block
	// e.g. if expr {...}
	// set by function.go used by expression.go
	expectingCompound bool

	// className is used for privatization
	className string
}

/*
eval* methods are helpers so you can match/next after evaluating something
match* methods verify that the current is what is expected and then advance
next* methods just advance
*/

func (p *parser) evalMatch(result ast.Node, tok Token) ast.Node {
	p.match(tok)
	return result
}

func (p *parser) evalNext(result ast.Node) ast.Node {
	p.next()
	return result
}

func (p *parser) match(tok Token) {
	p.mustMatch(tok)
	p.next()
}

func (p *parser) matchIdent() {
	if !IsIdent[p.Token] {
		p.error("expecting identifier")
	}
	p.next()
}

func (p *parser) matchIf(tok Token) bool {
	if tok == p.Token {
		p.next()
		return true
	}
	return false
}

func (p *parser) mustMatch(tok Token) {
	if tok != p.Token {
		p.error("expecting ", tok)
	}
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
}

// error panics with "syntax error at " + position
// It claims to return string so it can be called inside panic
// (so compiler knows we don't return)
func (p *parser) error(args ...interface{}) string {
	panic("syntax error at " + strconv.Itoa(int(p.Item.Pos)) + " " +
		fmt.Sprint(args...))
}
