/*
Package parse implements parsing for the Suneido language
*/
package parse

import . "github.com/apmckinlay/gsuneido/compile/lex"

/*
builder is used by the parser to build the AST.

It is primarily needed because expression parsing
is shared between code and queries, but with a different AST type.
*/
type builder func(Token, string, ...T) T

type T interface{}

func Parse(s string) T {
	lxr := NewLexer(s)
	p := parser{lxr, astBuilder, lxr.Next()}
	return p.parseConstant()
}

type parser struct {
	lxr *Lexer
	bld builder
	it  Item
}

func (p *parser) parseConstant() T {
	return p.function()
}

func (p *parser) function() T {
	p.match(FUNCTION)
	p.match(L_PAREN)
	p.match(R_PAREN)
	p.match(L_CURLY)
	result := p.expression()
	p.match(R_CURLY)
	return result
}

func (p *parser) expression() T {
	x := p.primary()
	for p.it.Token == ADD || p.it.Token == SUB || p.it.Token == CAT {
		it := p.it
		p.next(false)
		x = p.bld(it.Token, it.Value, x, p.primary())
	}
	return x
}

func (p *parser) primary() T {
	switch p.it.Token {
	case NUMBER:
		return p.matchReturn(p.bld(NUMBER, p.it.Value), NUMBER)
	case STRING:
		return p.matchReturn(p.bld(STRING, p.it.Value), STRING)
	default:
		panic("unexpected " + p.it.Value)
	}
}

// support methods -------------------------------------------------------------

func (p *parser) match(tok Token) {
	if tok == p.it.Token || tok == p.it.Keyword {
		p.next(false)
		return
	}
	panic("unexpected " + p.it.Value)
}

func (p *parser) matchIf(tok Token) bool {
	if tok == p.it.Token || tok == p.it.Keyword {
		p.next(false)
		return true
	}
	return false
}

func (p *parser) matchReturn(result T, tok Token) T {
	p.match(tok)
	return result
}

func (p *parser) next(skipnewlines bool) {
	for {
		p.it = p.lxr.Next()
		if p.it.Token == COMMENT || p.it.Token == WHITESPACE ||
			(skipnewlines && p.it.Token == NEWLINE) {
			continue
		}
		break
	}
}
