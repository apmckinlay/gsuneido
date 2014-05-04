/*
Package lex implements the lexical scanner for Suneido

The Lexer is designed so the sequence of values returned
forms the complete source.
*/
package lex

import (
	"bytes"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	src   string
	si    int
	ahead []Item
}

// Lexer returns a new instance
func NewLexer(src string) *Lexer {
	return &Lexer{src: src}
}

// Item is the return value from Lexer.Next
type Item struct {
	Token   Token
	Value   string
	Keyword Token
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
//
// items are buffered so they can be used by Next
func (lxr *Lexer) Ahead(i int) Item {
	for len(lxr.ahead) < i+1 {
		item := lxr.next()
		if item.Token == EOF {
			return item
		}
		lxr.ahead = append(lxr.ahead, item)
	}
	return lxr.ahead[i]
}

func (lxr *Lexer) next() Item {
	start, c := lxr.read()
	switch c {
	case eof:
		return Item{EOF, "EOF", NIL}
	case '#':
		return Item{HASH, "#", NIL}
	case '(':
		return Item{L_PAREN, "(", NIL}
	case ')':
		return Item{R_PAREN, ")", NIL}
	case ',':
		return Item{COMMA, ",", NIL}
	case ';':
		return Item{SEMICOLON, ";", NIL}
	case '?':
		return Item{Q_MARK, ";", NIL}
	case '@':
		return Item{AT, "@", NIL}
	case '[':
		return Item{L_BRACKET, "[", NIL}
	case ']':
		return Item{R_BRACKET, "]", NIL}
	case '{':
		return Item{L_CURLY, "{", NIL}
	case '}':
		return Item{R_CURLY, "}", NIL}
	case '~':
		return Item{BITNOT, "~", NIL}
	case ':':
		if lxr.match(':') {
			return Item{RANGELEN, "::", NIL}
		} else {
			return Item{COLON, ":", NIL}
		}
	case '=':
		if lxr.match('=') {
			return Item{IS, "==", NIL}
		} else {
			if lxr.match('~') {
				return Item{MATCH, "=~", NIL}
			} else {
				return Item{EQ, "=", NIL}
			}
		}
	case '!':
		if lxr.match('=') {
			return Item{ISNT, "!=", NIL}
		} else {
			if lxr.match('~') {
				return Item{MATCHNOT, "!~", NIL}
			} else {
				return Item{NOT, "!", NIL}
			}
		}
	case '<':
		if lxr.match('<') {
			if lxr.match('=') {
				return Item{LSHIFTEQ, "<<=", NIL}
			} else {
				return Item{LSHIFT, "<<", NIL}
			}
		} else if lxr.match('>') {
			return Item{ISNT, "<>", NIL}
		} else if lxr.match('=') {
			return Item{LTE, "<=", NIL}
		} else {
			return Item{LT, "<", NIL}
		}
	case '>':
		if lxr.match('>') {
			if lxr.match('=') {
				return Item{RSHIFTEQ, ">>=", NIL}
			} else {
				return Item{RSHIFT, ">>", NIL}
			}
		} else if lxr.match('=') {
			return Item{GTE, ">=", NIL}
		} else {
			return Item{GT, ">", NIL}
		}
	case '|':
		if lxr.match('|') {
			return Item{OR, "||", NIL}
		} else if lxr.match('=') {
			return Item{BITOREQ, "|=", NIL}
		} else {
			return Item{BITOR, "|", NIL}
		}
	case '&':
		if lxr.match('&') {
			return Item{AND, "&&", NIL}
		} else if lxr.match('=') {
			return Item{BITANDEQ, "&=", NIL}
		} else {
			return Item{BITAND, "&", NIL}
		}
	case '^':
		if lxr.match('=') {
			return Item{BITXOREQ, "^=", NIL}
		} else {
			return Item{BITXOR, "^", NIL}
		}
	case '-':
		if lxr.match('-') {
			return Item{DEC, "--", NIL}
		} else if lxr.match('=') {
			return Item{SUBEQ, "-=", NIL}
		} else {
			return Item{SUB, "-", NIL}
		}
	case '+':
		if lxr.match('+') {
			return Item{INC, "++", NIL}
		} else if lxr.match('=') {
			return Item{ADDEQ, "+=", NIL}
		} else {
			return Item{ADD, "+", NIL}
		}
	case '/':
		if lxr.match('/') {
			return lxr.lineComment(start)
		} else if lxr.match('*') {
			return lxr.spanComment(start)
		} else if lxr.match('=') {
			return Item{DIVEQ, "/=", NIL}
		} else {
			return Item{DIV, "/", NIL}
		}
	case '*':
		if lxr.match('=') {
			return Item{MULEQ, "*=", NIL}
		} else {
			return Item{MUL, "*", NIL}
		}
	case '%':
		if lxr.match('=') {
			return Item{MODEQ, "%=", NIL}
		} else {
			return Item{MOD, "%", NIL}
		}
	case '$':
		if lxr.match('=') {
			return Item{CATEQ, "$=", NIL}
		} else {
			return Item{CAT, "$", NIL}
		}
	case '`':
		return lxr.rawString(start)
	case '"', '\'':
		return lxr.quotedString(start, c)
	case '.':
		if lxr.match('.') {
			return Item{RANGETO, "..", NIL}
		} else if isDigit(lxr.peek()) {
			return lxr.number(start)
		} else {
			return Item{DOT, ".", NIL}
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return lxr.number(start)
	default:
		if isSpace(c) {
			return lxr.whitespace(start, c)
		} else if unicode.IsLetter(c) || c == '_' {
			return lxr.identifier(start)
		}
	}
	return Item{ERROR, string(c), NIL}
}

func (lxr *Lexer) whitespace(start int, c rune) Item {
	si := start
	result := WHITESPACE
	for ; isSpace(c); si, c = lxr.read() {
		if c == '\n' || c == '\r' {
			result = NEWLINE
		}
	}
	lxr.si = si
	return Item{result, lxr.src[start:lxr.si], NIL}
}

func (lxr *Lexer) lineComment(start int) Item {

	return Item{COMMENT, lxr.matchUntil(start, "\n"), NIL}
}

func (lxr *Lexer) spanComment(start int) Item {
	return Item{COMMENT, lxr.matchUntil(start, "*/"), NIL}
}

func (lxr *Lexer) rawString(start int) Item {
	return Item{STRING, lxr.matchUntil(start, "`"), NIL}
}

func (lxr *Lexer) quotedString(start int, quote rune) Item {
	var buf bytes.Buffer
	lxr.match(quote)
	for c := lxr.read1(); c != eof && c != quote; c = lxr.read1() {
		buf.WriteRune(lxr.doesc(c))
	}
	return Item{STRING, buf.String(), NIL}
}

func (lxr *Lexer) doesc(c rune) rune {
	if c != '\\' {
		return c
	}
	si, c := lxr.read()
	switch c {
	case 'n':
		return '\n'
	case 't':
		return '\t'
	case 'r':
		return '\r'
	case 'x':
		dig1 := digit(lxr.read1(), 16)
		dig2 := digit(lxr.read1(), 16)
		if dig1 != -1 && dig2 != -1 {
			return rune(16*dig1 + dig2)
		}
	case '\\':
	case '"':
	case '\'':
		return c
	default:
		dig1 := digit(lxr.read1(), 8)
		dig2 := digit(lxr.read1(), 8)
		dig3 := digit(lxr.read1(), 8)
		if dig1 != -1 && dig2 != -1 && dig3 != -1 {
			return rune(64*dig1 + 8*dig2 + dig3)
		}
	}
	lxr.si = si
	return '\\'
}

func digit(c rune, radix int) int {
	n := 99
	if isDigit(c) {
		n = int(c - '0')
	} else if isHexDigit(c) {
		n = int(10 + unicode.ToLower(c) - 'a')
	}
	if n < radix {
		return n
	} else {
		return -1
	}
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isHexDigit(r rune) bool {
	return strings.ContainsRune(hexDigits, r)
}

func (lxr *Lexer) number(start int) Item {
	lxr.matchOneOf("+-")
	// Is it hex?
	digits := "0123456789"
	if lxr.match('0') && lxr.matchOneOf("xX") {
		digits = hexDigits
	}
	lxr.matchRunOf(digits)
	if lxr.match('.') {
		lxr.matchRunOf(digits)
	}
	if lxr.matchOneOf("eE") {
		lxr.matchOneOf("+-")
		lxr.matchRunOf("0123456789")
	}
	return Item{NUMBER, lxr.src[start:lxr.si], NIL}
}

func (lxr *Lexer) identifier(start int) Item {
	lxr.matchWhile(isIdentChar)
	if !lxr.match('?') {
		lxr.match('!')
	}
	value := lxr.src[start:lxr.si]
	keyword := NIL
	if lxr.peek() != ':' {
		keyword = Keyword(value)
	}
	return Item{IDENTIFIER, value, keyword}
}

const eof = -1

func (lxr *Lexer) read() (int, rune) {
	si := lxr.si
	return si, lxr.read1()
}

func (lxr *Lexer) read1() rune {
	if lxr.si >= len(lxr.src) {
		return eof
	}
	c, w := utf8.DecodeRuneInString(lxr.src[lxr.si:])
	lxr.si += w
	return c
}

func (lxr *Lexer) peek() rune {
	si, c := lxr.read()
	lxr.si = si
	return c
}

func (lxr *Lexer) match(c rune) bool {
	si, c2 := lxr.read()
	if c == c2 {
		return true
	}
	lxr.si = si
	return false
}

func (lxr *Lexer) matchOneOf(valid string) bool {
	si, c := lxr.read()
	if strings.ContainsRune(valid, c) {
		return true
	}
	lxr.si = si
	return false
}

func (lxr *Lexer) matchRunOf(valid string) {
	for {
		si, c := lxr.read()
		if !strings.ContainsRune(valid, c) {
			lxr.si = si
			break
		}
	}
}

func (lxr *Lexer) matchWhile(f func(c rune) bool) {
	for {
		si, c := lxr.read()
		if !f(c) {
			lxr.si = si
			break
		}
	}
}

func (lxr *Lexer) matchUntil(start int, s string) string {
	for lxr.read1() != eof && !strings.HasSuffix(lxr.src[:lxr.si], s) {
	}
	return lxr.src[start:lxr.si]
}

func isIdentChar(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

const hexDigits = "0123456789abcdefABCDEF"

func isSpace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}
