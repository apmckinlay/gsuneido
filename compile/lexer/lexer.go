// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package lexer implements the lexical scanner for the Suneido language
package lexer

import (
	"strings"

	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// Lexer implements the lexical scanner for Suneido
// It is designed so the sequence of values returned
// forms the complete source.
type Lexer struct {
	keyword func(s string) (tok.Token, string)
	src     string
	ahead   []Item
	si      int
	nlwhite bool
}

// NewLexer returns a new Lexer
func NewLexer(src string) *Lexer {
	return &Lexer{src: src, keyword: keyword}
}

func (lxr *Lexer) Dup() *Lexer {
	return &Lexer{src: lxr.src, keyword: lxr.keyword}
}

func (lxr *Lexer) Source() string {
	return lxr.src
}

// Position will be out of sync if using Ahead
func (lxr *Lexer) Position() int {
	return lxr.si
}

// Item is the return value from Lexer.Next
type Item struct {
	Text  string
	Pos   int32
	Token tok.Token
}

func (it *Item) String() string {
	if it.Text != "" {
		return "'" + it.Text + "'"
	}
	return it.Token.String()
}

// Next returns the next Item
func (lxr *Lexer) Next() Item {
	if len(lxr.ahead) > 0 {
		item := lxr.ahead[0]
		lxr.ahead = lxr.ahead[1:]
		return item
	}
	return lxr.next()
}

// Ahead provides lookahead, 0 is the next item
// Items are buffered so they can be used by Next
func (lxr *Lexer) Ahead(i int) Item {
	for len(lxr.ahead) < i+1 {
		item := lxr.next()
		if item.Token == tok.Eof {
			return item
		}
		lxr.ahead = append(lxr.ahead, item)
	}
	return lxr.ahead[i]
}

// AheadSkip provides lookahead like Ahead
// but skips Whitespace, Newline, and Comment
func (lxr *Lexer) AheadSkip(i int) Item {
	for j := 0; ; j++ {
		switch it := lxr.Ahead(j); it.Token {
		case tok.Whitespace, tok.Newline, tok.Comment:
			continue
		case tok.Eof:
			return it
		default:
			if i <= 0 {
				return it
			}
			i--
		}
	}
}

// Remainder is used by ParseAdmin to get view definitions
func (lxr *Lexer) Remainder() string {
	assert.That(len(lxr.ahead) == 0)
	si := lxr.si
	lxr.si = len(lxr.src)
	return lxr.src[si:]
}

func (lxr *Lexer) next() Item {
	start := lxr.si
	c := lxr.read()
	it := func(token tok.Token) Item {
		// compiler doesn't need Text, but Suneido Scanner does
		return Item{Pos: int32(start), Token: token, Text: lxr.src[start:lxr.si]}
	}
	switch c {
	case eof:
		return it(tok.Eof)
	case '#':
		if p := lxr.peek(); p == '_' || IsLetter(p) {
			lxr.matchIdentTail()
			val := lxr.src[start+1 : lxr.si]
			return Item{Text: val, Pos: int32(start), Token: tok.Symbol}
		}
		return it(tok.Hash)
	case '(':
		return it(tok.LParen)
	case ')':
		return it(tok.RParen)
	case ',':
		return it(tok.Comma)
	case ';':
		return it(tok.Semicolon)
	case '?':
		return it(tok.QMark)
	case '@':
		return it(tok.At)
	case '[':
		return it(tok.LBracket)
	case ']':
		return it(tok.RBracket)
	case '{':
		return it(tok.LCurly)
	case '}':
		return it(tok.RCurly)
	case '~':
		return it(tok.BitNot)
	case ':':
		if lxr.match(':') {
			return it(tok.RangeLen)
		}
		return it(tok.Colon)
	case '=':
		if lxr.match('~') {
			return it(tok.Match)
		}
		return it(tok.Eq)
	case '!':
		if lxr.match('~') {
			return it(tok.MatchNot)
		} else if lxr.match('=') {
			return it(tok.Isnt) //TODO remove when not needed by queries
		}
	case '<':
		if lxr.match('<') {
			if lxr.match('=') {
				return it(tok.LShiftEq)
			}
			return it(tok.LShift)
		}
		if lxr.match('>') {
			return it(tok.Isnt)
		}
		if lxr.match('=') {
			return it(tok.Lte)
		}
		return it(tok.Lt)
	case '>':
		if lxr.match('>') {
			if lxr.match('=') {
				return it(tok.RShiftEq)
			}
			return it(tok.RShift)
		}
		if lxr.match('=') {
			return it(tok.Gte)
		}
		return it(tok.Gt)
	case '|':
		if lxr.match('=') {
			return it(tok.BitOrEq)
		}
		return it(tok.BitOr)
	case '&':
		if lxr.match('=') {
			return it(tok.BitAndEq)
		}
		return it(tok.BitAnd)
	case '^':
		if lxr.match('=') {
			return it(tok.BitXorEq)
		}
		return it(tok.BitXor)
	case '-':
		if lxr.match('-') {
			return it(tok.Dec)
		}
		if lxr.match('=') {
			return it(tok.SubEq)
		}
		return it(tok.Sub)
	case '+':
		if lxr.match('+') {
			return it(tok.Inc)
		}
		if lxr.match('=') {
			return it(tok.AddEq)
		}
		return it(tok.Add)
	case '/':
		if lxr.match('/') {
			return lxr.lineComment(start)
		}
		if lxr.match('*') {
			return lxr.spanComment(start)
		}
		if lxr.match('=') {
			return it(tok.DivEq)
		}
		return it(tok.Div)
	case '*':
		if lxr.match('=') {
			return it(tok.MulEq)
		}
		return it(tok.Mul)
	case '%':
		if lxr.match('=') {
			return it(tok.ModEq)
		}
		return it(tok.Mod)
	case '$':
		if lxr.match('=') {
			return it(tok.CatEq)
		}
		return it(tok.Cat)
	case '`':
		return lxr.rawString(start)
	case '"', '\'':
		return lxr.quotedString(start, c)
	case '.':
		if lxr.match('.') {
			return it(tok.RangeTo)
		}
		if IsDigit(lxr.peek()) {
			return lxr.number(start)
		}
		return it(tok.Dot)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return lxr.number(start)
	default:
		if IsSpace(c) {
			return lxr.whitespace(start, c)
		} else if IsLetter(c) || c == '_' {
			return lxr.identifier(start)
		}
	}
	return it(tok.Error)
}

func it(token tok.Token, pos int, txt string) Item {
	return Item{Text: txt, Pos: int32(pos), Token: token}
}

func (lxr *Lexer) whitespace(start int, c byte) Item {
	result := tok.Whitespace
	for ; IsSpace(c); c = lxr.read() {
		if !lxr.nlwhite && (c == '\n' || c == '\r') {
			result = tok.Newline
		}
	}
	if c != eof {
		lxr.si--
	}
	return it(result, start, lxr.src[start:lxr.si])
}

func (lxr *Lexer) lineComment(start int) Item {
	// does NOT absorb \r\n
loop:
	for {
		switch lxr.read() {
		case eof:
			break loop
		case '\n':
			lxr.si--
			if lxr.src[lxr.si-1] == '\r' {
				lxr.si--
			}
			break loop
		}
	}
	return it(tok.Comment, start, lxr.src[start:lxr.si])
}

func (lxr *Lexer) spanComment(start int) Item {
	return it(tok.Comment, start, lxr.matchUntil(start, "*/"))
}

func (lxr *Lexer) rawString(start int) Item {
	s := lxr.matchUntil(start+1, "`")
	s = strings.TrimSuffix(s, "`")
	return it(tok.String, start, s)
}

func (lxr *Lexer) quotedString(start int, quote byte) Item {
	// if no escapes, return slice of source
	src := lxr.src[lxr.si:]
	for i := 0; ; i++ {
		if i >= len(src) {
			lxr.si += len(src)
			return it(tok.String, start, strings.Clone(src)) // no closing quote
		} else if src[i] == '\\' {
			break
		} else if src[i] == byte(quote) {
			lxr.si += i + 1
			return it(tok.String, start, strings.Clone(src[:i])) // no escapes
		}
	}
	// have escapes so need to build new string
	var buf strings.Builder
	for c := lxr.read(); c != eof && c != quote; c = lxr.read() {
		c = lxr.doesc(c)
		buf.WriteByte(byte(c))
	}
	return Item{Text: buf.String(), Pos: int32(start), Token: tok.String}
}

func (lxr *Lexer) doesc(c byte) byte {
	if c != '\\' {
		return c
	}
	si := lxr.si
	c = lxr.read()
	switch c {
	case 'n':
		return '\n'
	case 't':
		return '\t'
	case 'r':
		return '\r'
	case 'x':
		dig1 := digit(lxr.read(), 16)
		dig2 := digit(lxr.read(), 16)
		if dig1 != -1 && dig2 != -1 {
			return byte(16*dig1 + dig2)
		}
	case '\\', '"', '\'':
		return c
	}
	lxr.si = si
	return '\\'
}

func digit(c byte, radix int) int {
	n := 99
	if IsDigit(c) {
		n = int(c - '0')
	} else if IsHexDigit(c) {
		n = int(10 + ToLower(c) - 'a')
	}
	if n < radix {
		return n
	}
	return -1
}

func isDigitOrUnderscore(c byte) bool {
	return IsDigit(c) || c == '_'
}

func (lxr *Lexer) number(start int) Item {
	// see also string_NumberQ
	if lxr.src[start] == '0' && lxr.matchOneOf("xX") {
		lxr.matchWhile(IsHexDigit)
	} else {
		lxr.matchWhile(isDigitOrUnderscore)
		if lxr.match('.') {
			lxr.matchWhile(isDigitOrUnderscore)
		}
		exp := lxr.si
		if lxr.matchOneOf("eE") {
			lxr.matchOneOf("+-")
			lxr.matchWhile(IsDigit)
			if lxr.si == exp+1 {
				lxr.si = exp
			}
		}
		if lxr.src[lxr.si-1] == '.' && lxr.nonWhiteRemaining() {
			lxr.si-- // don't absorb trailing dot
		}
	}
	numStr := strings.ReplaceAll(lxr.src[start:lxr.si], "_", "")
	return it(tok.Number, start, numStr)
}

func (lxr *Lexer) nonWhiteRemaining() bool {
	for i := lxr.si; i < len(lxr.src); i++ {
		if !IsSpace(lxr.src[i]) {
			return true
		}
	}
	return false
}

func (lxr *Lexer) identifier(start int) Item {
	lxr.matchIdentTail()
	val := lxr.src[start:lxr.si]
	token := tok.Identifier
	if lxr.peek() != ':' || val == "default" || val == "true" || val == "false" {
		token, val = lxr.keyword(val)
	}
	return Item{Text: val, Pos: int32(start), Token: token}
}

// keyword returns the token for a string it is a keyword
// otherwise Identifier and a copy of the string
func keyword(s string) (tok.Token, string) {
	switch len(s) {
	case 2:
		if s == "if" {
			return tok.If, "if"
		}
		if s == "is" {
			return tok.Is, "is"
		}
		if s == "in" {
			return tok.In, "in"
		}
		if s == "or" {
			return tok.Or, "or"
		}
		if s == "do" {
			return tok.Do, "do"
		}
	case 3:
		if s == "and" {
			return tok.And, "and"
		}
		if s == "for" {
			return tok.For, "for"
		}
		if s == "not" {
			return tok.Not, "not"
		}
		if s == "new" {
			return tok.New, "new"
		}
		if s == "try" {
			return tok.Try, "try"
		}
	case 4:
		if s == "true" {
			return tok.True, "true"
		}
		if s == "isnt" {
			return tok.Isnt, "isnt"
		}
		if s == "else" {
			return tok.Else, "else"
		}
		if s == "this" {
			return tok.This, "this"
		}
		if s == "case" {
			return tok.Case, "case"
		}
	case 5:
		if s == "false" {
			return tok.False, "false"
		}
		if s == "super" {
			return tok.Super, "super"
		}
		if s == "class" {
			return tok.Class, "class"
		}
		if s == "throw" {
			return tok.Throw, "throw"
		}
		if s == "catch" {
			return tok.Catch, "catch"
		}
		if s == "while" {
			return tok.While, "while"
		}
		if s == "break" {
			return tok.Break, "break"
		}
	case 6:
		if s == "return" {
			return tok.Return, "return"
		}
		if s == "switch" {
			return tok.Switch, "switch"
		}
	case 7:
		if s == "default" {
			return tok.Default, "default"
		}
		if s == "forever" {
			return tok.Forever, "forever"
		}
	case 8:
		if s == "function" {
			return tok.Function, "function"
		}
		if s == "continue" {
			return tok.Continue, "continue"
		}
	}
	return tok.Identifier, strings.Clone(s)
}

func (lxr *Lexer) matchIdentTail() {
	lxr.matchWhile(isIdentChar)
	if !lxr.match('?') {
		lxr.match('!')
	}
}

func isIdentChar(r byte) bool {
	return r == '_' || IsLetter(r) || IsDigit(r)
}

const eof = 0

func (lxr *Lexer) read() byte {
	if lxr.si >= len(lxr.src) {
		return eof
	}
	c := lxr.src[lxr.si]
	lxr.si++
	return c
}

// peek returns the next character
func (lxr *Lexer) peek() byte {
	if lxr.si >= len(lxr.src) {
		return eof
	}
	return lxr.src[lxr.si]
}

func (lxr *Lexer) match(c byte) bool {
	if c == lxr.peek() {
		lxr.si++
		return true
	}
	return false
}

func (lxr *Lexer) matchOneOf(valid string) bool {
	if -1 != strings.IndexByte(valid, lxr.peek()) {
		lxr.si++
		return true
	}
	return false
}

func (lxr *Lexer) matchWhile(f func(c byte) bool) {
	for ; f(lxr.peek()); lxr.si++ {
	}
}

func (lxr *Lexer) matchUntil(start int, s string) string {
	for lxr.read(); lxr.si < len(lxr.src) && !strings.HasSuffix(lxr.src[:lxr.si], s); lxr.si++ {
	}
	return lxr.src[start:lxr.si]
}
