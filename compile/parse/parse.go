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
type builder func(Item, ...T) T

type T interface{}

func Parse(s string) T {
	lxr := NewLexer(s)
	p := parser{lxr: lxr, bld: astBuilder, Item: lxr.Next()}
	return p.parseConstant()
}

type parser struct {
	lxr *Lexer
	bld builder
	Item
	nest int
}

func (p *parser) parseConstant() T {
	return p.function()
}

func (p *parser) function() T {
	it := p.Item
	p.match(FUNCTION)
	p.match(L_PAREN)
	p.match(R_PAREN)
	body := p.compound()
	return p.bld(it, p.bld(Item{Token: PARAMS}), body)
}

func (p *parser) compound() T {
	p.match(L_CURLY)
	return p.evalMatch(p.statements(), R_CURLY)
}

func (p *parser) statements() T {
	list := []T{}
	for p.Token != R_CURLY {
		stmt := p.statement()
		//fmt.Println("stmt:", stmt)
		if stmt != nil {
			list = append(list, stmt)
		}
	}
	return p.bld(Item{Token: STATEMENTS}, list...)
}

func (p *parser) statement() T {
	if p.Token == L_CURLY {
		return p.compound()
	} else if p.matchIf(SEMICOLON) {
		return nil
	}
	// TODO other statement types
	return p.expression()
}

func (p *parser) expression() T {
	x := p.unary()
	for p.Token == ADD || p.Token == SUB || p.Token == CAT {
		it := p.Item
		p.next()
		x = p.bld(it, x, p.unary())
	}
	return x
}

func (p *parser) unary() T {
	switch p.Token {
	case ADD, SUB, NOT, BITNOT:
		it := p.Item
		p.next()
		return p.bld(it, p.unary())
	default:
		return p.term()
	}
}

func (p *parser) term() T {
	term := p.primary()
	if p.Token == EQ {
		it := p.Item
		p.next()
		expr := p.expression()
		term = p.bld(it, term, expr)
	}
	return term
}

func (p *parser) primary() T {
	switch p.Token {
	case NUMBER, STRING, IDENTIFIER:
		return p.evalNext(p.bld(p.Item))
	case L_PAREN:
		p.next()
		return p.evalMatch(p.expression(), R_PAREN)
	}
	panic("unexpected " + p.Value)
}

// support methods -------------------------------------------------------------

func (p *parser) match(tok Token) {
	if tok == p.Token || tok == p.Keyword {
		p.next()
		return
	}
	panic("unexpected " + p.Value)
}

func (p *parser) matchIf(tok Token) bool {
	if tok == p.Token || tok == p.Keyword {
		p.next()
		return true
	}
	return false
}

func (p *parser) evalMatch(result T, tok Token) T {
	p.match(tok)
	return result
}

func (p *parser) evalNext(result T) T {
	p.next()
	return result
}

// next advances to the next non-white token, tracking nesting
// NOTE: it does NOT skip newlines
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
	//fmt.Println("item:", p.Item)
}
