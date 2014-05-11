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
	return p.bld(Item{Token: EXPRESSION}, p.bitorExpr())
}

func (p *parser) bitorExpr() T {
	x := p.bitxorExpr()
	for p.KeyTok() == BITOR {
		it := p.Item
		p.nextSkipNewlines()
		x = p.bld(it, x, p.bitxorExpr())
	}
	return x
}

func (p *parser) bitxorExpr() T {
	x := p.bitandExpr()
	for p.KeyTok() == BITXOR {
		it := p.Item
		p.nextSkipNewlines()
		x = p.bld(it, x, p.bitandExpr())
	}
	return x
}

func (p *parser) bitandExpr() T {
	x := p.isExpr()
	for p.KeyTok() == BITAND {
		it := p.Item
		p.nextSkipNewlines()
		x = p.bld(it, x, p.isExpr())
	}
	return x
}

func (p *parser) isExpr() T {
	x := p.cmpExpr()
	for p.KeyTok() == IS || p.KeyTok() == ISNT ||
		p.Token == MATCH || p.Token == MATCHNOT {
		it := p.Item
		p.nextSkipNewlines()
		x = p.bld(it, x, p.cmpExpr())
	}
	return x
}

func (p *parser) cmpExpr() T {
	x := p.shiftExpr()
	for p.Token == LT || p.Token == LTE || p.Token == GT || p.Token == GTE {
		it := p.Item
		p.nextSkipNewlines()
		x = p.bld(it, x, p.shiftExpr())
	}
	return x
}

func (p *parser) shiftExpr() T {
	x := p.addExpr()
	for p.Token == LSHIFT || p.Token == RSHIFT {
		it := p.Item
		p.nextSkipNewlines()
		x = p.bld(it, x, p.addExpr())
	}
	return x
}

func (p *parser) addExpr() T {
	x := p.mulExpr()
	for p.Token == ADD || p.Token == SUB || p.Token == CAT {
		it := p.Item
		p.nextSkipNewlines()
		x = p.bld(it, x, p.mulExpr())
	}
	return x
}

func (p *parser) mulExpr() T {
	x := p.unary()
	for p.Token == MUL || p.Token == DIV || p.Token == MOD {
		it := p.Item
		p.nextSkipNewlines()
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
		p.nextSkipNewlines()
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
