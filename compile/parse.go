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
	"github.com/apmckinlay/gsuneido/core"
)

func NewParser(src string) *Parser {
	return newParser(NewLexer(src), &cgAspects{})
}

func CheckParser(src string, t *core.Thread) *Parser {
	a := &cgckAspects{}
	a.Check = check.New(t)
	return newParser(NewLexer(src), a)
}

func AstParser(src string) *Parser {
	return newParser(NewLexer(src), &astAspects{})
}

func QueryParser(src string) *Parser {
	return newParser(NewQueryLexer(src), &actionAspects{})
}

func newParser(lxr *Lexer, a Aspects) *Parser {
	p := &Parser{ParserBase: ParserBase{Lxr: lxr, Aspects: a}}
	p.Next()
	return p
}

type actionAspects struct {
	cgAspects
}

type Expr = ast.Expr

func (aa *actionAspects) Binary(lhs Expr, token tok.Token, rhs Expr) Expr {
	if token.IsAssign() {
		panic("assignment operators are not allowed")
	}
	return aa.cgAspects.Binary(lhs, token, rhs)
}

func (aa *actionAspects) Unary(token tok.Token, expr Expr) Expr {
	if token.IsIncDec() {
		panic("increment/decrement operators are not allowed")
	}
	if token == tok.LParen {
		return expr
	}
	return aa.cgAspects.Unary(token, expr)
}

func (p *Parser) InitFuncInfo() {
	p.funcInfo.final = map[string]uint8{}
	p.funcInfo.assignConst = map[string]bool{}
}

type ParserBase struct {
	Aspects
	Lxr *Lexer

	// Item is the current lexical token etc.
	Item

	// EndPos is the end of the previous token
	EndPos int32

	// newline is true if the current token was preceded by a newline
	newline bool

	// EqToIs treats Eq as Is for queries
	EqToIs bool
}

type Parser struct {

	// prevDef is used by libload for _Name
	prevDef core.Value

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

	ParserBase

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

	// inTry is used to give an error for nested try
	inTry bool
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

// Next skips whitespace and comments to advance to the next token.
// It sets p.Item (Text, Pos, Token), p.EndPos, and p.newline.
func (p *ParserBase) Next() {
	p.newline = false
	p.Item = p.Lxr.Next()
	p.EndPos = int32(p.Item.Pos)
	for {
		if p.Token == tok.Newline {
			if p.Lxr.AheadSkip(0).Token != tok.QMark {
				p.newline = true
			}
		} else if p.Token != tok.Comment && p.Token != tok.Whitespace {
			break
		}
		p.Item = p.Lxr.Next()
	}
	if p.EqToIs && p.Token == tok.Eq {
		p.Token = tok.Is // for queries
	}
}

// Error panics with "syntax Error at " + position
// It claims to return string so it can be called inside panic
// (so compiler knows we don't return)
func (p *ParserBase) Error(args ...any) string {
	p.ErrorAt(p.Pos, args...)
	return ""
}

func (p *ParserBase) ErrorAt(pos int32, args ...any) string {
	panic("syntax error @" + strconv.Itoa(int(pos)) + " " + fmt.Sprint(args...))
}

func (*Parser) Constant(val core.Value) Expr {
	return &ast.Constant{Val: val}
}

func SetPos(result iSetPos, org, end int32) {
	if result != nil {
		result.SetPos(org, end)
	}
}

type iSetPos interface {
	SetPos(org, end int32)
}
