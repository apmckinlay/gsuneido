/*
Package lex implements the lexical scanner for Suneido

The Lexer is designed so the sequence of values returned
forms the complete source.
*/
package lex

import (
	"bytes"
	"fmt"
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
	Value   string
	Pos     int32
	Token   Token
	Keyword Token
	// NOTE: put Token's last to reduce padding
}

func (it *Item) KeyTok() Token {
	if it.Keyword != NIL {
		return it.Keyword
	} else {
		return it.Token
	}
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
	it := func(tok Token) Item {
		return Item{lxr.src[start:lxr.si], int32(start), tok, NIL}
	}
	switch c {
	case eof:
		return it(EOF)
	case '#':
		return it(HASH)
	case '(':
		return it(L_PAREN)
	case ')':
		return it(R_PAREN)
	case ',':
		return it(COMMA)
	case ';':
		return it(SEMICOLON)
	case '?':
		return it(Q_MARK)
	case '@':
		return it(AT)
	case '[':
		return it(L_BRACKET)
	case ']':
		return it(R_BRACKET)
	case '{':
		return it(L_CURLY)
	case '}':
		return it(R_CURLY)
	case '~':
		return it(BITNOT)
	case ':':
		if lxr.match(':') {
			return it(RANGELEN)
		} else {
			return it(COLON)
		}
	case '=':
		if lxr.match('=') {
			return it(IS)
		} else {
			if lxr.match('~') {
				return it(MATCH)
			} else {
				return it(EQ)
			}
		}
	case '!':
		if lxr.match('=') {
			return it(ISNT)
		} else {
			if lxr.match('~') {
				return it(MATCHNOT)
			} else {
				return it(NOT)
			}
		}
	case '<':
		if lxr.match('<') {
			if lxr.match('=') {
				return it(LSHIFTEQ)
			} else {
				return it(LSHIFT)
			}
		} else if lxr.match('>') {
			return it(ISNT)
		} else if lxr.match('=') {
			return it(LTE)
		} else {
			return it(LT)
		}
	case '>':
		if lxr.match('>') {
			if lxr.match('=') {
				return it(RSHIFTEQ)
			} else {
				return it(RSHIFT)
			}
		} else if lxr.match('=') {
			return it(GTE)
		} else {
			return it(GT)
		}
	case '|':
		if lxr.match('|') {
			return it(OR)
		} else if lxr.match('=') {
			return it(BITOREQ)
		} else {
			return it(BITOR)
		}
	case '&':
		if lxr.match('&') {
			return it(AND)
		} else if lxr.match('=') {
			return it(BITANDEQ)
		} else {
			return it(BITAND)
		}
	case '^':
		if lxr.match('=') {
			return it(BITXOREQ)
		} else {
			return it(BITXOR)
		}
	case '-':
		if lxr.match('-') {
			return it(DEC)
		} else if lxr.match('=') {
			return it(SUBEQ)
		} else {
			return it(SUB)
		}
	case '+':
		if lxr.match('+') {
			return it(INC)
		} else if lxr.match('=') {
			return it(ADDEQ)
		} else {
			return it(ADD)
		}
	case '/':
		if lxr.match('/') {
			return lxr.lineComment(start)
		} else if lxr.match('*') {
			return lxr.spanComment(start)
		} else if lxr.match('=') {
			return it(DIVEQ)
		} else {
			return it(DIV)
		}
	case '*':
		if lxr.match('=') {
			return it(MULEQ)
		} else {
			return it(MUL)
		}
	case '%':
		if lxr.match('=') {
			return it(MODEQ)
		} else {
			return it(MOD)
		}
	case '$':
		if lxr.match('=') {
			return it(CATEQ)
		} else {
			return it(CAT)
		}
	case '`':
		return lxr.rawString(start)
	case '"', '\'':
		return lxr.quotedString(start, c)
	case '.':
		if lxr.match('.') {
			return it(RANGETO)
		} else if isDigit(lxr.peek()) {
			return lxr.number(start)
		} else {
			return it(DOT)
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
	return it(ERROR)
}

func it(tok Token, pos int, val string) Item {
	return Item{val, int32(pos), tok, NIL}
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
	return it(result, start, lxr.src[start:lxr.si])
}

func (lxr *Lexer) lineComment(start int) Item {
	return it(COMMENT, start, lxr.matchUntil(start, "\n"))
}

func (lxr *Lexer) spanComment(start int) Item {
	return it(COMMENT, start, lxr.matchUntil(start, "*/"))
}

func (lxr *Lexer) rawString(start int) Item {
	return it(STRING, start, lxr.matchUntil(start, "`"))
}

func (lxr *Lexer) quotedString(start int, quote rune) Item {
	// don't use buffer unless there are escapes
	src := lxr.src[lxr.si:]
	i := strings.IndexByte(src, byte(quote))
	if i == -1 {
		lxr.si += len(src)
		return it(STRING, start, src)
	}
	j := strings.IndexByte(src[:i], '\\')
	if j == -1 { // no escapes
		lxr.si += i + 1
		return it(STRING, start, src[:i])
	}
	var buf bytes.Buffer
	lxr.match(quote)
	for c := lxr.read1(); c != eof && c != quote; c = lxr.read1() {
		buf.WriteRune(lxr.doesc(c))
	}
	// keyword set to STRING means not referencing src
	return Item{buf.String(), int32(start), STRING, STRING}
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
	// Is it hex?
	digits := "0123456789"
	if lxr.match('0') && lxr.matchOneOf("xX") {
		fmt.Println("hex")
		digits = hexDigits
	}
	lxr.matchRunOf(digits)
	if lxr.match('.') {
		fmt.Println("dot")
		lxr.matchRunOf(digits)
	}
	if lxr.matchOneOf("eE") {
		fmt.Println("exp")
		lxr.matchOneOf("+-")
		lxr.matchRunOf("0123456789")
	}
	return it(NUMBER, start, lxr.src[start:lxr.si])
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
	return Item{value, int32(start), IDENTIFIER, keyword}
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
