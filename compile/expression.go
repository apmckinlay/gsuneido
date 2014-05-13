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
	return p.bitorExpr()
}

// TODO ?:
// TODO and
// TODO or
// TODO in

func (p *parser) bitorExpr() T {
	x := p.bitxorExpr()
	for p.Token == BITOR {
		it := p.Item
		p.nextSkipNL()
		x = p.bld(it, x, p.bitxorExpr())
	}
	return x
}

func (p *parser) bitxorExpr() T {
	x := p.bitandExpr()
	for p.Token == BITXOR {
		it := p.Item
		p.nextSkipNL()
		x = p.bld(it, x, p.bitandExpr())
	}
	return x
}

func (p *parser) bitandExpr() T {
	x := p.isExpr()
	for p.Token == BITAND {
		it := p.Item
		p.nextSkipNL()
		x = p.bld(it, x, p.isExpr())
	}
	return x
}

func (p *parser) isExpr() T {
	x := p.cmpExpr()
	for p.KeyTok() == IS || p.KeyTok() == ISNT ||
		p.Token == MATCH || p.Token == MATCHNOT {
		it := p.Item
		p.nextSkipNL()
		x = p.bld(it, x, p.cmpExpr())
	}
	return x
}

func (p *parser) cmpExpr() T {
	x := p.shiftExpr()
	for p.Token == LT || p.Token == LTE || p.Token == GT || p.Token == GTE {
		it := p.Item
		p.nextSkipNL()
		x = p.bld(it, x, p.shiftExpr())
	}
	return x
}

func (p *parser) shiftExpr() T {
	x := p.addExpr()
	for p.Token == LSHIFT || p.Token == RSHIFT {
		it := p.Item
		p.nextSkipNL()
		x = p.bld(it, x, p.addExpr())
	}
	return x
}

func (p *parser) addExpr() T {
	x := p.mulExpr()
	for p.Token == ADD || p.Token == SUB || p.Token == CAT {
		it := p.Item
		p.nextSkipNL()
		x = p.bld(it, x, p.mulExpr())
	}
	return x
}

func (p *parser) mulExpr() T {
	x := p.unary()
	for p.Token == MUL || p.Token == DIV || p.Token == MOD {
		it := p.Item
		p.nextSkipNL()
		x = p.bld(it, x, p.unary())
	}
	return x
}

func (p *parser) unary() T {
	switch p.KeyTok() {
	case ADD, SUB, NOT, BITNOT:
		item := p.Item
		p.next()
		return p.bld(item, p.unary())
	default:
		return p.term()
	}
}

func (p *parser) term() T {
	var preincdec Item
	if p.Token == INC || p.Token == DEC {
		preincdec = p.Item
		p.next()
	}
	term := p.primary()
	if preincdec.Token != NIL {
		return p.bld(preincdec, term)
	} else if EQ <= p.Token && p.Token <= BITXOREQ { // TODO assignment operators
		it := p.Item
		p.nextSkipNL()
		expr := p.expr()
		return p.bld(it, term, expr)
	} else if p.Token == INC || p.Token == DEC {
		p.Token += POSTINC - INC
		return p.evalNext(p.bld(p.Item, term))
	}
	return term
}

func (p *parser) primary() T {
	switch p.Token {
	case IDENTIFIER:
		switch p.Keyword {
		case TRUE, FALSE:
			return p.evalNext(p.bld(p.Item))
		default:
			return p.evalNext(p.bld(p.Item))
		}
	case NUMBER, STRING:
		return p.evalNext(p.bld(p.Item))
	case HASH:
		val := p.constant()
		return p.bld(Item{Token: VALUE}, val)
	case L_PAREN:
		p.next()
		return p.evalMatch(p.expr(), R_PAREN)
	}
	panic("unexpected " + p.Text)
}
