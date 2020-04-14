// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

import (
	"fmt"
	"strconv"

	"github.com/apmckinlay/gsuneido/compile/ast"
	"github.com/apmckinlay/gsuneido/compile/check"
	. "github.com/apmckinlay/gsuneido/lexer"
	tok "github.com/apmckinlay/gsuneido/lexer/tokens"
)

func NewParser(src string) *parser {
	lxr := NewLexer(src)
	p := &parser{lxr: lxr, Factory: ast.Builder{}}
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
	// and also to handle folding.
	ast.Factory

	// newline is true if the current token was preceeded by a newline
	newline bool

	// expectingCompound is used to differentiate control statement body vs. block
	// e.g. if expr {...}
	// set by function.go used by expression.go
	expectingCompound bool

	// lib is passed to named constants for Display
	lib string

	// name is used for Named constants
	name string

	// className is set for privatization
	className string

	// assignName is used to pass the variable name through an assignment
	// e.g. foo = function () { ... }; Name(foo) => "foo"
	assignName string

	// checker is used to add additional checking along with codegen
	checker *check.Check
}

/*
eval* methods are helpers so you can match/next after evaluating something
match* methods verify that the current is what is expected and then advance
next* methods just advance
*/

func (p *parser) evalMatch(result ast.Node, token tok.Token) ast.Node {
	p.match(token)
	return result
}

func (p *parser) evalNext(result ast.Node) ast.Node {
	p.next()
	return result
}

func (p *parser) match(token tok.Token) {
	p.mustMatch(token)
	p.next()
}

func (p *parser) matchIdent() {
	if !p.Token.IsIdent() {
		p.error("expecting identifier")
	}
	p.next()
}

func (p *parser) matchIf(token tok.Token) bool {
	if token == p.Token {
		p.next()
		return true
	}
	return false
}

func (p *parser) mustMatch(token tok.Token) {
	if token != p.Token {
		p.error("expecting ", token)
	}
}

// next advances to the next token, setting p.Item
func (p *parser) next() {
	p.newline = false
	for {
		p.Item = p.lxr.Next()
		if p.Token == tok.Newline {
			if p.lxr.AheadSkip(0).Token != tok.QMark {
				p.newline = true
			}
		} else if p.Token != tok.Comment && p.Token != tok.Whitespace {
			break
		}
	}
}

// error panics with "syntax error at " + position
// It claims to return string so it can be called inside panic
// (so compiler knows we don't return)
func (p *parser) error(args ...interface{}) string {
	p.errorAt(p.Item.Pos, args...)
	return ""
}

func (p *parser) errorAt(pos int32, args ...interface{}) string {
	panic("syntax error @" + strconv.Itoa(int(pos)) + " " + fmt.Sprint(args...))
}

func (p *parser) Ident(name string) ast.Expr {
	return p.Factory.Ident(name, p.Pos)
}
