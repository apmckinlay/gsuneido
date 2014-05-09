package compile

/*
expression parses a Suneido expression
and builds an AST using the supplied builder

it takes an existing parser since it is embedded
and not used standalone

expression is used by both function and query parsers
*/
func expression(p *parser, b builder) T {
	p.bld = b
	return p.expr()
}

type builder func(Item, ...T) T

type T interface{}

func (p *parser) expr() T {
	return p.bld(Item{Token: EXPRESSION}, p.addexpr())
}

func (p *parser) addexpr() T {
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
		expr := p.expr()
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
		return p.evalMatch(p.expr(), R_PAREN)
	}
	panic("unexpected " + p.Value)
}
