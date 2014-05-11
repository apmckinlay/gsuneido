package compile

import (
	"fmt"

	"github.com/apmckinlay/gsuneido/util/verify"
)

func newParser(src string) *parser {
	lxr := NewLexer(src)
	return &parser{lxr: lxr, Item: lxr.Next()}
}

type parser struct {
	lxr *Lexer
	Item
	nest int
	bld  builder // used by expression
}

func (p *parser) match(tok Token) {
	if tok != p.Token && tok != p.Keyword {
		p.error("expecting", tok)
	}
	p.next()
}

func (p *parser) matchIf(tok Token) bool {
	if tok != p.Token && tok != p.Keyword {
		return false
	}
	p.next()
	return true
}

func (p *parser) matchSkipNewlines(tok Token) {
	if tok != p.Token && tok != p.Keyword {
		p.error("expecting", tok)
	}
	for {
		p.next()
		if p.Token != NEWLINE {
			break
		}
	}
}

func (p *parser) evalMatch(result T, tok Token) T {
	p.match(tok)
	return result
}

func (p *parser) evalNext(result T) T {
	p.next()
	return result
}

func (p *parser) nextSkipNewlines() {
	p.next()
	for p.Token == NEWLINE {
		p.next()
	}
}

// next advances to the next non-white token, tracking nesting
func (p *parser) next() {
	for {
		p.Item = p.lxr.Next()
		switch p.Token {
		case COMMENT, WHITESPACE:
			continue
		case L_CURLY, L_PAREN, L_BRACKET:
			p.nest++
		case R_CURLY, R_PAREN, R_BRACKET:
			p.nest--
		}
		break
	}
	verify.That(p.nest >= -1) // final curly on compound will go to -1
	for p.Token == NEWLINE &&
		(p.nest > 0 || binop(p.lxr.Ahead(0))) {
		p.Item = p.lxr.Next()
	}
	if p.Token == STRING && p.Keyword != STRING {
		// make a copy of strings that are slices of the source
		p.Value = " " + p.Value
		p.Value = p.Value[1:]
	}
	//fmt.Println("item:", p.Item)
}

func binop(it Item) bool {
	switch it.KeyTok() {
	// NOTE: not ADD or SUB because they can be unary
	case AND, OR, CAT, MUL, DIV, MOD,
		EQ, ADDEQ, SUBEQ, CATEQ, MULEQ, DIVEQ, MODEQ,
		BITAND, BITOR, BITXOR, BITANDEQ, BITOREQ, BITXOREQ,
		GT, GTE, LT, LTE, LSHIFT, LSHIFTEQ, RSHIFT, RSHIFTEQ,
		IS, ISNT, MATCH, MATCHNOT, Q_MARK:
		return true
	}
	return false
}

// returns string so it can be called inside panic
// so compiler knows we don't return
func (p *parser) error(args ...interface{}) string {
	panic("syntax error" + fmt.Sprint(args...))
}
