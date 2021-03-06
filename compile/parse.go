// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

import (
	"fmt"
	"strconv"

	"github.com/apmckinlay/gsuneido/compile/ast"
	"github.com/apmckinlay/gsuneido/compile/check"
	. "github.com/apmckinlay/gsuneido/compile/lexer"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	"github.com/apmckinlay/gsuneido/runtime"
)

func NewParser(src string) *Parser {
	return newParser(NewLexer(src), &cgAspects{})
}

func CheckParser(src string, t *runtime.Thread) *Parser {
	a := &cgckAspects{}
	a.Check = check.New(t)
	return newParser(NewLexer(src), a)
}

func AstParser(src string) *Parser {
	return newParser(NewLexer(src), &astAspects{})
}

func GogenParser(src string) *Parser {
	return newParser(NewLexer(src), &gogenAspects{})
}

func QueryParser(src string) *Parser {
	return newParser(NewQueryLexer(src), &cgAspects{})
}

func newParser(lxr *Lexer, a Aspects) *Parser {
	p := &Parser{ParserBase: ParserBase{Lxr: lxr, Aspects: a}}
	p.Next()
	return p
}

func (p *Parser) InitFuncInfo() {
	p.funcInfo.final = map[string]uint8{}
	p.funcInfo.assignConst = map[string]bool{}
}

type ParserBase struct {
	Lxr *Lexer

	// Item is the current lexical token etc.
	Item

	// newline is true if the current token was preceeded by a newline
	newline bool

	EqToIs bool

	Aspects
}

type Parser struct {
	ParserBase

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
	final       map[string]uint8
	assignConst map[string]bool

	// hasBlocks is whether the function has any blocks
	hasBlocks bool
}

// disqualified is a special value for final
const disqualified = 2

func (p *ParserBase) Match(token tok.Token) {
	if token == tok.String && p.Token == tok.Symbol {
		token = tok.Symbol
	}
	p.MustMatch(token)
	p.Next()
}

func (p *ParserBase) MatchIdent() string {
	text := p.Text
	if !p.Token.IsIdent() {
		p.Error("expecting identifier")
	}
	p.Next()
	return text
}

func (p *ParserBase) MatchIf(token tok.Token) bool {
	if token == p.Token {
		p.Next()
		return true
	}
	return false
}

func (p *ParserBase) MustMatch(token tok.Token) {
	if token != p.Token {
		p.Error("expecting ", token)
	}
}

// Next advances to the Next token, setting p.Item
func (p *ParserBase) Next() {
	p.newline = false
	for {
		p.Item = p.Lxr.Next()
		if p.Token == tok.Newline {
			if p.Lxr.AheadSkip(0).Token != tok.QMark {
				p.newline = true
			}
		} else if p.Token != tok.Comment && p.Token != tok.Whitespace {
			break
		}
	}
	if p.EqToIs && p.Token == tok.Eq {
		p.Token = tok.Is
	}
}

// Error panics with "syntax Error at " + position
// It claims to return string so it can be called inside panic
// (so compiler knows we don't return)
func (p *ParserBase) Error(args ...interface{}) string {
	p.ErrorAt(p.Item.Pos, args...)
	return ""
}

func (p *ParserBase) ErrorAt(pos int32, args ...interface{}) string {
	panic("syntax error @" + strconv.Itoa(int(pos)) + " " + fmt.Sprint(args...))
}

func (p *Parser) Ident(name string) *ast.Ident {
	return &ast.Ident{Name: name, Pos: p.Pos}
}
