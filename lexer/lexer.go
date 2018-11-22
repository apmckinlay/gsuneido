package lexer

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/util/ascii"
)

// Lexer implements the lexical scanner for Suneido
// It is designed so the sequence of values returned
// forms the complete source.
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
// For keywords, Token is IDENTIFIER, Keyword is the particular keyword token.
// When Token is STRING, if Keyword is STRING it means Text is not a slice of the source.
type Item struct {
	Text    string
	Pos     int32
	Token   Token
	Keyword Token //TODO maybe make this KeyTok ???
	// NOTE: put Token's last to reduce padding
}

// KeyTok returns Keyword if set, else Token
func (it *Item) KeyTok() Token {
	if it.Keyword != NIL {
		return it.Keyword
	}
	return it.Token
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
	start := lxr.si
	c := lxr.read()
	it := func(tok Token) Item {
		return Item{lxr.src[start:lxr.si], int32(start), tok, NIL}
	}
	switch c {
	case eof:
		return it(EOF)
	case '#':
		if p := lxr.peek(); p == '_' || IsLetter(p) {
			lxr.matchIdentTail()
			val := lxr.src[start+1 : lxr.si]
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
		if lxr.match('~') {
			return it(MATCH)
		}
		return it(EQ)
	case '!':
		if lxr.match('~') {
			return it(MATCHNOT)
		}
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
		if lxr.match('=') {
			return it(BITOREQ)
		}
		return it(BITOR)
	case '&':
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
		if IsDigit(lxr.peek()) {
			return lxr.number(start)
		}
		return it(DOT)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return lxr.number(start)
	default:
		if IsSpace(c) {
			return lxr.whitespace(start, c)
		} else if IsLetter(c) || c == '_' {
			return lxr.identifier(start)
		}
	}
	return it(ERROR)
}

func it(tok Token, pos int, txt string) Item {
	return Item{txt, int32(pos), tok, NIL}
}

func (lxr *Lexer) whitespace(start int, c byte) Item {
	result := WHITESPACE
	for ; IsSpace(c); c = lxr.read() {
		if c == '\n' || c == '\r' {
			result = NEWLINE
		}
	}
	if c != eof {
		lxr.si--
	}
	return it(result, start, lxr.src[start:lxr.si])
}

func (lxr *Lexer) lineComment(start int) Item {
	// does NOT absorb newline
loop:
	for {
		switch lxr.read() {
		case eof:
			break loop
		case '\n':
			lxr.si--
			break loop
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

func (lxr *Lexer) quotedString(start int, quote byte) Item {
	// if no escapes, return slice of source
	src := lxr.src[lxr.si:]
	for i := 0; ; i++ {
		if i >= len(src) {
			lxr.si += len(src)
			return it(STRING, start, src) // no closing quote
		} else if src[i] == '\\' {
			break
		} else if src[i] == byte(quote) {
			lxr.si += i + 1
			return it(STRING, start, src[:i]) // no escapes
		}
	}
	// have escapes so need to build new string
	var buf strings.Builder
	for c := lxr.read(); c != eof && c != quote; c = lxr.read() {
		c = lxr.doesc(c)
		buf.WriteByte(byte(c))
	}
	// keyword set to STRING means *not* referencing src
	return Item{buf.String(), int32(start), STRING, STRING}
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

func (lxr *Lexer) number(start int) Item {
	if lxr.src[start] == '0' && lxr.matchOneOf("xX") {
		lxr.matchWhile(IsHexDigit)
	} else {
		lxr.matchWhile(IsDigit)
		if lxr.match('.') {
			lxr.matchWhile(IsDigit)
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
	return it(NUMBER, start, lxr.src[start:lxr.si])
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
	keyword := NIL
	if lxr.peek() != ':' || val == "default" || val == "true" || val == "false" {
		keyword, val = Keyword(val)
	}
	return Item{val, int32(start), IDENTIFIER, keyword}
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

func (lxr *Lexer) matchRunOf(valid string) {
	for ; -1 != strings.IndexByte(valid, lxr.peek()); lxr.si++ {
	}
}

func (lxr *Lexer) matchWhile(f func(c byte) bool) {
	for ; f(lxr.peek()); lxr.si++ {
	}
}

func (lxr *Lexer) matchUntil(start int, s string) string {
	for lxr.si++; lxr.si < len(lxr.src) && !strings.HasSuffix(lxr.src[:lxr.si], s); lxr.si++ {
	}
	return lxr.src[start:lxr.si]
}
