package compile

import (
	"math"
	"strconv"
)

/*
expression parses a Suneido expression
and builds an AST using the supplied builder

it takes an existing parser since it is embedded
and not used standalone

expression is used by both function and query parsers
(with different builder's)
*/
func expression(p *parser, b builder) T {
	defer func(prev int) { p.nest = prev }(p.nest)
	p.nest = 0
	p.bld = b
	return p.expr()
}

type builder func(Item, ...T) T

type T interface{}

func (p *parser) expr() T {
	return p.qcExpr()
}

func (p *parser) qcExpr() T {
	e := p.orExpr()
	if p.Token != Q_MARK {
		return e
	}
	it := p.Item
	p.nest++
	p.match(Q_MARK)
	t := p.expr()
	p.match(COLON)
	p.nest--
	f := p.expr()
	return p.bld(it, e, t, f)
}

func (p *parser) orExpr() T {
	x := p.andExpr()
	if p.KeyTok() != OR {
		return x
	}
	it := p.Item
	list := []T{x}
	for p.KeyTok() == OR {
		p.nextSkipNL()
		list = append(list, p.andExpr())
	}
	return p.bld(it, list...)
}

func (p *parser) andExpr() T {
	x := p.inExpr()
	if p.KeyTok() != AND {
		return x
	}
	it := p.Item
	list := []T{x}
	for p.KeyTok() == AND {
		p.nextSkipNL()
		list = append(list, p.inExpr())
	}
	return p.bld(it, list...)
}

func (p *parser) inExpr() T {
	x := p.bitorExpr()
	if p.KeyTok() != IN {
		return x
	}
	it := p.Item
	p.next()
	list := []T{x}
	p.match(L_PAREN)
	for p.Token != R_PAREN {
		list = append(list, p.expr())
		if p.Token == COMMA {
			p.next()
		}
	}
	p.match(R_PAREN)
	return p.bld(it, list...)
}

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
	x := p.catExpr()
	for p.Token == LSHIFT || p.Token == RSHIFT {
		it := p.Item
		p.nextSkipNL()
		x = p.bld(it, x, p.catExpr())
	}
	return x
}

// e.g. a $ b $ c => ($ a b c)
func (p *parser) catExpr() T {
	x := p.addExpr()
	if p.Token != CAT {
		return x
	}
	it := p.Item
	list := []T{x}
	for p.Token == CAT {
		p.nextSkipNL()
		list = append(list, p.addExpr())
	}
	return p.bld(it, list...)
}

// convert to a sequence of additions so it's commutative for folding
// e.g. a + b - c + d => (+ a b (- c) d)
func (p *parser) addExpr() T {
	list := []T{p.mulExpr()}
	for p.Token == ADD || p.Token == SUB {
		it := p.Item
		p.nextSkipNL()
		next := p.mulExpr()
		if it.Token == SUB {
			next = p.bld(it, next)
		}
		list = append(list, next)
	}
	if len(list) == 1 {
		return list[0]
	}
	return p.bld(Item{Token: ADD, Text: "+"}, list...)
}

// convert mul & div to a sequence of mul so it's commutative for folding
// e.g. a * b / c d => (* a b (/ c) d)
func (p *parser) mulExpr() T {
	list := []T{p.unary()}
	for p.Token == MUL || p.Token == DIV || p.Token == MOD {
		it := p.Item
		p.nextSkipNL()
		next := p.unary()
		switch it.Token {
		case MUL:
			list = append(list, next)
		case DIV:
			list = append(list, p.bld(it, next))
		case MOD:
			if len(list) > 1 {
				list = []T{p.bld(Item{Token: MUL, Text: "*"}, list...)}
			}
			list[0] = p.bld(it, list[0], next)
		}
	}
	if len(list) == 1 {
		return list[0]
	}
	return p.bld(Item{Token: MUL, Text: "*"}, list...)
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

var int_max_str = strconv.Itoa(math.MaxInt32)

func (p *parser) term() T {
	var preincdec Item
	if p.Token == INC || p.Token == DEC {
		preincdec = p.Item
		p.next()
	}
	term := p.primary()
	for p.Token == DOT || p.Token == L_BRACKET {
		if p.Token == DOT {
			dot := p.Item
			p.nextSkipNL()
			id := p.Item
			p.match(IDENTIFIER)
			term = p.bld(dot, term, p.bld(id))
		} else if p.Token == L_BRACKET {
			sub := p.Item
			p.next()
			var expr T
			if p.Token == RANGETO || p.Token == RANGELEN {
				expr = p.bld(Item{Token: NUMBER, Text: "0"})
			} else {
				expr = p.expr()
			}
			if p.Token == RANGETO || p.Token == RANGELEN {
				rtype := p.Item
				p.next()
				var expr2 T
				if p.Token == R_BRACKET {
					expr2 = p.bld(Item{Token: NUMBER, Text: int_max_str})
				} else {
					expr2 = p.expr()
				}
				expr = p.bld(rtype, expr, expr2)
			}
			term = p.bld(sub, term, expr)
			p.match(R_BRACKET)
		}
	}
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
