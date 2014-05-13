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

/*
eval* methods are helpers so you can match/next after evaluating something
match* methods verify that the current is what is expected and then advance
next* methods just advance
*/

func (p *parser) evalMatch(result T, tok Token) T {
	p.match(tok)
	return result
}

func (p *parser) evalNext(result T) T {
	p.next()
	return result
}

func (p *parser) match(tok Token) {
	p.mustMatch(tok)
	p.next()
}

func (p *parser) matchIf(tok Token) bool {
	if p.isMatch(tok) {
		p.next()
		return true
	}
	return false
}

func (p *parser) matchKeepNL(tok Token) {
	p.mustMatch(tok)
	p.nextKeepNL()
}

func (p *parser) matchSkipNL(tok Token) {
	p.mustMatch(tok)
	p.nextSkipNL()
}

func (p *parser) mustMatch(tok Token) {
	if !p.isMatch(tok) {
		p.error("expecting", tok)
	}
}

func (p *parser) isMatch(tok Token) bool {
	return tok == p.Token || tok == p.Keyword
}

// next keeps or skips newlines based on nesting
// and whether the next line starts with a binary operator
func (p *parser) next() {
	p.nextKeepNL()
	for p.Token == NEWLINE &&
		(p.nest > 0 || binop(p.lxr.Ahead(0))) {
		p.Item = p.lxr.Next()
	}
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

func (p *parser) nextSkipNL() {
	p.nextKeepNL()
	for p.Token == NEWLINE {
		p.nextKeepNL()
	}
}

// next advances to the next token,
// skipping comments and whitespace (but not newlines),
// and tracking nesting
func (p *parser) nextKeepNL() {
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
	if p.Token == STRING && p.Keyword != STRING {
		// make a copy of strings that are slices of the source
		p.Value = " " + p.Value
		p.Value = p.Value[1:]
	}
	//fmt.Println("item:", p.Item)
}

// returns string so it can be called inside panic
// so compiler knows we don't return
func (p *parser) error(args ...interface{}) string {
	panic("syntax error " + fmt.Sprint(args...))
}
