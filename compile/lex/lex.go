// Package lex implements the lexical scanner for Suneido
package lex

import (
	"bytes"
	"strings"
	"unicode"
	"unicode/utf8"
)

type lexer struct {
	src   string
	si    int
	start int
	width int
	value string
}

type Item struct {
	token   Token
	keyword Token
	value   string
}

func Lexer(src string) *lexer {
	return &lexer{src: src}
}

func (lxr *lexer) Next() Item {
	token := lxr.next()
	var value string
	if lxr.value != "" {
		value = lxr.value
	} else {
		value = lxr.src[lxr.start:lxr.si]
	}
	keyword := NIL
	if token == IDENTIFIER && lxr.peek() != ':' {
		keyword = Keyword(value)
	}
	return Item{token, keyword, value}
}

func (lxr *lexer) next() Token {
	lxr.start = lxr.si
	lxr.value = ""
	c := lxr.read()
	switch c {
	case eof:
		return EOF
	case '#':
		return HASH
	case '(':
		return L_PAREN
	case ')':
		return R_PAREN
	case ',':
		return COMMA
	case ';':
		return SEMICOLON
	case '?':
		return Q_MARK
	case '@':
		return AT
	case '[':
		return L_BRACKET
	case ']':
		return R_BRACKET
	case '{':
		return L_CURLY
	case '}':
		return R_CURLY
	case '~':
		return BITNOT
	case ':':
		if lxr.match(':') {
			return RANGELEN
		} else {
			return COLON
		}
	case '=':
		if lxr.match('=') {
			return IS
		} else {
			if lxr.match('~') {
				return MATCH
			} else {
				return EQ
			}
		}
	case '!':
		if lxr.match('=') {
			return ISNT
		} else {
			if lxr.match('~') {
				return MATCHNOT
			} else {
				return NOT
			}
		}
	case '<':
		if lxr.match('<') {
			if lxr.match('=') {
				return LSHIFTEQ
			} else {
				return LSHIFT
			}
		} else if lxr.match('>') {
			return ISNT
		} else if lxr.match('=') {
			return LTE
		} else {
			return LT
		}
	case '>':
		if lxr.match('>') {
			if lxr.match('=') {
				return RSHIFTEQ
			} else {
				return RSHIFT
			}
		} else if lxr.match('=') {
			return GTE
		} else {
			return GT
		}
	case '|':
		if lxr.match('|') {
			return OR
		} else if lxr.match('=') {
			return BITOREQ
		} else {
			return BITOR
		}
	case '&':
		if lxr.match('&') {
			return AND
		} else if lxr.match('=') {
			return BITANDEQ
		} else {
			return BITAND
		}
	case '^':
		if lxr.match('=') {
			return BITXOREQ
		} else {
			return BITXOR
		}
	case '-':
		if lxr.match('-') {
			return DEC
		} else if lxr.match('=') {
			return SUBEQ
		} else {
			return SUB
		}
	case '+':
		if lxr.match('+') {
			return INC
		} else if lxr.match('=') {
			return ADDEQ
		} else {
			return ADD
		}
	case '/':
		if lxr.match('/') {
			return lxr.lineComment()
		} else if lxr.match('*') {
			return lxr.spanComment()
		} else if lxr.match('=') {
			return DIVEQ
		} else {
			return DIV
		}
	case '*':
		if lxr.match('=') {
			return MULEQ
		} else {
			return MUL
		}
	case '%':
		if lxr.match('=') {
			return MODEQ
		} else {
			return MOD
		}
	case '$':
		if lxr.match('=') {
			return CATEQ
		} else {
			return CAT
		}
	case '`':
		return lxr.rawString()
	case '"', '\'':
		return lxr.quotedString(c)
	case '.':
		if lxr.match('.') {
			return RANGETO
		} else if isDigit(lxr.peek()) {
			return lxr.number()
		} else {
			return DOT
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return lxr.number()
	default:
		if isSpace(c) {
			return lxr.whitespace(c)
		} else if unicode.IsLetter(c) || c == '_' {
			return lxr.identifier()
		}
	}
	return ERROR
}

func (lxr *lexer) whitespace(c rune) Token {
	result := WHITESPACE
	for ; isSpace(c); c = lxr.read() {
		if c == '\n' || c == '\r' {
			result = NEWLINE
		}
	}
	lxr.backup()
	return result
}

func (lxr *lexer) lineComment() Token {
	for c := lxr.read(); c != eof && c != '\n'; c = lxr.read() {
	}
	return COMMENT
}

func (lxr *lexer) spanComment() Token {
	lxr.matchUntil(func() bool { return strings.HasSuffix(lxr.src[:lxr.si], "*/") })
	return COMMENT
}

func (lxr *lexer) rawString() Token {
	for c := lxr.read(); c != eof && c != '`'; c = lxr.read() {
	}
	lxr.value = lxr.src[lxr.start+1 : lxr.si-1]
	return STRING
}

func (lxr *lexer) quotedString(quote rune) Token {
	var buf bytes.Buffer
	lxr.match(quote)
	for c := lxr.read(); c != eof && c != quote; c = lxr.read() {
		buf.WriteRune(lxr.doesc(c))
	}
	lxr.value = buf.String()
	return STRING
}

func (lxr *lexer) doesc(c rune) rune {
	if c != '\\' {
		return c
	}
	save := lxr.si
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
			return rune(16*dig1 + dig2)
		}
	case '\\':
	case '"':
	case '\'':
		return c
	default:
		dig1 := digit(lxr.read(), 8)
		dig2 := digit(lxr.read(), 8)
		dig3 := digit(lxr.read(), 8)
		if dig1 != -1 && dig2 != -1 && dig3 != -1 {
			return rune(64*dig1 + 8*dig2 + dig3)
		}
	}
	lxr.si = save
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

func (lxr *lexer) number() Token {
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
	return NUMBER
}

func (lxr *lexer) identifier() Token {
	lxr.matchWhile(isIdentChar)
	if !lxr.match('?') {
		lxr.match('!')
	}
	return IDENTIFIER
}

const eof = -1

func (lxr *lexer) read() rune {
	if lxr.si >= len(lxr.src) {
		lxr.width = 0
		return eof
	}
	c, w := utf8.DecodeRuneInString(lxr.src[lxr.si:])
	lxr.si += w
	lxr.width = w
	return c
}

func (lxr *lexer) backup() {
	lxr.si -= lxr.width
}

func (lxr *lexer) peek() rune {
	c := lxr.read()
	lxr.backup()
	return c
}

func (lxr *lexer) match(c rune) bool {
	if c == lxr.read() {
		return true
	}
	lxr.backup()
	return false
}

func (lxr *lexer) matchOneOf(valid string) bool {
	if strings.ContainsRune(valid, lxr.read()) {
		return true
	}
	lxr.backup()
	return false
}

func (lxr *lexer) matchRunOf(valid string) {
	for strings.ContainsRune(valid, lxr.read()) {
	}
	lxr.backup()
}

func (lxr *lexer) matchWhile(f func(c rune) bool) {
	for c := lxr.read(); f(c); c = lxr.read() {
	}
	lxr.backup()
}

func (lxr *lexer) matchUntil(f func() bool) {
	for c := lxr.read(); c != eof && !f(); c = lxr.read() {
	}
}

func isIdentChar(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

const hexDigits = "0123456789abcdefABCDEF"

func isSpace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}
