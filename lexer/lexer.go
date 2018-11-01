package lexer

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

/*
Lexer implements the lexical scanner for Suneido

It is designed so the sequence of values returned
forms the complete source.
*/
type Lexer struct {
	src   string
	si    int
	ahead []Item
}

// NewLexer returns a new instance
func NewLexer(src string) *Lexer {
	return &Lexer{src: src}
}

// Item is the return value from Lexer.Next
//
// For keywords, Token is IDENTIFIER, Keyword is the particular keyword token.
// When Token is STRING, if Keyword is STRING it means Text is not a slice of the source.
type Item struct {
	Text    string
	Pos     int32
	Token   Token
	Keyword Token
	// NOTE: put Token's last to reduce padding
}

// KeyTok returns Keyword if set, else Token
func (it *Item) KeyTok() Token {
	if it.Keyword != NIL {
		return it.Keyword
	}
	return it.Token
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

// Ahead provides lookahead, 0 is the next item.
//
// Items are buffered so they can be used by Next
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

// AheadSkip provides lookahead like Ahead
// but skips WHITESPACE, NEWLINE, and COMMENT
func (lxr *Lexer) AheadSkip(i int) Item {
	for j := 0; ; j++ {
		switch it := lxr.Ahead(j); it.Token {
		case WHITESPACE, NEWLINE, COMMENT:
			continue
		case EOF:
			return it
		default:
			if i <= 0 {
				return it
			}
			i--
		}
	}
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
		if p := lxr.peek(); p == '_' || unicode.IsLetter(p) {
			start++
			lxr.matchWhile(isIdentChar)
			if !lxr.match('?') {
				lxr.match('!')
			}
			val := lxr.src[start:lxr.si]
			return Item{val, int32(start), STRING, NIL}
		}
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
		}
		return it(COLON)
	case '=':
		if lxr.match('=') {
			return it(IS)
		}
		if lxr.match('~') {
			return it(MATCH)
		}
		return it(EQ)
	case '!':
		if lxr.match('=') {
			return it(ISNT)
		}
		if lxr.match('~') {
			return it(MATCHNOT)
		}
		return it(NOT)
	case '<':
		if lxr.match('<') {
			if lxr.match('=') {
				return it(LSHIFTEQ)
			}
			return it(LSHIFT)
		}
		if lxr.match('>') {
			return it(ISNT)
		}
		if lxr.match('=') {
			return it(LTE)
		}
		return it(LT)
	case '>':
		if lxr.match('>') {
			if lxr.match('=') {
				return it(RSHIFTEQ)
			}
			return it(RSHIFT)
		}
		if lxr.match('=') {
			return it(GTE)
		}
		return it(GT)
	case '|':
		if lxr.match('|') {
			return it(OR)
		}
		if lxr.match('=') {
			return it(BITOREQ)
		}
		return it(BITOR)
	case '&':
		if lxr.match('&') {
			return it(AND)
		}
		if lxr.match('=') {
			return it(BITANDEQ)
		}
		return it(BITAND)
	case '^':
		if lxr.match('=') {
			return it(BITXOREQ)
		}
		return it(BITXOR)
	case '-':
		if lxr.match('-') {
			return it(DEC)
		}
		if lxr.match('=') {
			return it(SUBEQ)
		}
		return it(SUB)
	case '+':
		if lxr.match('+') {
			return it(INC)
		}
		if lxr.match('=') {
			return it(ADDEQ)
		}
		return it(ADD)
	case '/':
		if lxr.match('/') {
			return lxr.lineComment(start)
		}
		if lxr.match('*') {
			return lxr.spanComment(start)
		}
		if lxr.match('=') {
			return it(DIVEQ)
		}
		return it(DIV)
	case '*':
		if lxr.match('=') {
			return it(MULEQ)
		}
		return it(MUL)
	case '%':
		if lxr.match('=') {
			return it(MODEQ)
		}
		return it(MOD)
	case '$':
		if lxr.match('=') {
			return it(CATEQ)
		}
		return it(CAT)
	case '`':
		return lxr.rawString(start)
	case '"', '\'':
		return lxr.quotedString(start, c)
	case '.':
		if lxr.match('.') {
			return it(RANGETO)
		}
		if isDigit(lxr.peek()) {
			return lxr.number(start)
		}
		return it(DOT)
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

func it(tok Token, pos int, txt string) Item {
	return Item{txt, int32(pos), tok, NIL}
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
	// does NOT absorb newline
	for {
		si, c := lxr.read()
		if c == eof || c == '\n' {
			lxr.si = si
			break
		}
	}
	return it(COMMENT, start, lxr.src[start:lxr.si])
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
	var buf strings.Builder
	lxr.match(quote)
	for c := lxr.read1(); c != eof && c != quote; c = lxr.read1() {
		buf.WriteByte(byte(lxr.doesc(c)))
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
	}
	return -1
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
	if lxr.src[start] == '0' && lxr.matchOneOf("xX") {
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
	if lxr.src[lxr.si-1] == '.' {
		lxr.si-- // don't absorb trailing dot
	}
	return it(NUMBER, start, lxr.src[start:lxr.si])
}

func (lxr *Lexer) identifier(start int) Item {
	lxr.matchWhile(isIdentChar)
	if !lxr.match('?') {
		lxr.match('!')
	}
	val := lxr.src[start:lxr.si]
	keyword := NIL
	if lxr.peek() != ':' || val == "default" || val == "true" || val == "false" {
		keyword = Keyword(val)
	}
	return Item{val, int32(start), IDENTIFIER, keyword}
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

// peek returns the next character
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
