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
	factory := ast.Folder{Factory: ast.Builder{}}
	p := &parser{parserBase: parserBase{lxr: lxr, Factory: factory},
		funcInfo: funcInfo{final: map[string]int{}}}
	p.next()
	return p
}

type parserBase struct {
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
}

type parser struct {
	parserBase

	// funcInfo is information gathered specific to a function
	// it must be saved/reset/restored for nested functions
	funcInfo

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

	// expectingCompound is used to differentiate control statement body vs. block
	// e.g. if expr {...}
	// set by function.go used by expression.go
	expectingCompound bool

	// itUsed records whether an "it" variable is used
	// to know whether to add an automatic "it" parameter to blocks
	itUsed bool
}

type funcInfo struct {
	// final is used to identify final variables
	final map[string]int

	// compoundNest is the compound nesting level, used for final
	compoundNest int

	// hasBlocks is whether the function has any blocks
	hasBlocks bool
}

// disqualified is a special value for final
const disqualified = -1


func (p *parserBase) match(token tok.Token) {
	p.mustMatch(token)
	p.next()
}

func (p *parserBase) matchIdent() string {
	text := p.Text
	if !p.Token.IsIdent() {
		p.error("expecting identifier")
	}
	p.next()
	return text
}

func (p *parserBase) matchIf(token tok.Token) bool {
	if token == p.Token {
		p.next()
		return true
	}
	return false
}

func (p *parserBase) mustMatch(token tok.Token) {
	if token != p.Token {
		p.error("expecting ", token)
	}
}

// next advances to the next token, setting p.Item
func (p *parserBase) next() {
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
func (p *parserBase) error(args ...interface{}) string {
	p.errorAt(p.Item.Pos, args...)
	return ""
}

func (p *parserBase) errorAt(pos int32, args ...interface{}) string {
	panic("syntax error @" + strconv.Itoa(int(pos)) + " " + fmt.Sprint(args...))
}

func (p *parser) Ident(name string) ast.Expr {
	return p.Factory.Ident(name, p.Pos)
}
