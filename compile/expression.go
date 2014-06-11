package compile

import (
	"math"
	"strconv"
	. "github.com/apmckinlay/gsuneido/lexer"
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
	case NEW:
		return p.newExpr()
	default:
		return p.term(false)
	}
}

func (p *parser) newExpr() T {
	it := p.Item
	p.match(NEW)
	term := p.term(true)
	var args T
	if p.Token == L_PAREN {
		args = p.arguments()
	} else {
		args = p.bld(argList)
	}
	return p.bld(it, term, args)
}

var int_max_str = strconv.Itoa(math.MaxInt32)

func (p *parser) term(newTerm bool) T {
	var preincdec Item
	if p.Token == INC || p.Token == DEC {
		preincdec = p.Item
		p.next()
	}
	term := p.primary()
	for p.Token == DOT || p.Token == L_BRACKET || p.Token == L_PAREN ||
		(p.Token == L_CURLY && !p.expectingCompound) {
		if newTerm && p.Token == L_PAREN {
			return term
		}
		if p.Token == DOT {
			dot := p.Item
			p.nextSkipNL()
			id := p.Item
			p.match(IDENTIFIER)
			term = p.bld(dot, term, p.bld(id))
			if !p.expectingCompound &&
				p.Token == NEWLINE && p.lxr.AheadSkip(0).Token == L_CURLY {
				p.match(NEWLINE)
			}
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
		} else if p.Token == L_PAREN || p.Token == L_CURLY {
			term = p.bld(call, term, p.arguments())
		}
	}
	if preincdec.Token != NIL {
		return p.bld(preincdec, term)
	} else if EQ <= p.Token && p.Token <= BITXOREQ {
		it := p.Item
		p.nextSkipNL()
		expr := p.expr()
		return p.bld(it, term, expr)
	} else if p.Token == INC || p.Token == DEC {
		p.Text = "post"
		return p.evalNext(p.bld(p.Item, term))
	}
	return term
}

var call = Item{Text: "call"}

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
		return p.bld(Item{}, val)
	case L_PAREN:
		p.next()
		return p.evalMatch(p.expr(), R_PAREN)
	case L_CURLY:
		return p.block()
	}
	panic("unexpected '" + p.Text + "'")
}

func (p *parser) arguments() T {
	var args []T
	if p.matchIf(L_PAREN) {
		if p.matchIf(AT) {
			return p.atArgument()
		} else {
			args = p.argumentList(R_PAREN)
		}
	}
	if p.Token == NEWLINE &&
		!p.expectingCompound && p.lxr.AheadSkip(0).Token == L_CURLY {
		p.match(NEWLINE)
	}
	if p.Token == L_CURLY {
		args = append(args, p.bld(blockArg, p.block()))
	}
	return p.bld(argList, args...)
}

var argList = Item{Text: "args"}
var atArg = Item{Text: "atArg"}
var noKeyword = Item{Text: "noKwd"}
var trueItem = Item{Token: TRUE, Text: "true"}
var blockArg = Item{Token: IDENTIFIER, Text: "blockArg"}
var blockItem = Item{Text: "block"}
var blockParams = Item{Text: "blockParams"}
var zeroItem = Item{Token: NUMBER, Text: "0"}

func (p *parser) atArgument() T {
	n := zeroItem
	if p.matchIf(ADD) {
		n = p.Item
		p.match(NUMBER)
	}
	expr := p.expr()
	p.match(R_PAREN)
	return p.bld(atArg, p.bld(n), expr)
}

func (p *parser) argumentList(closing Token) []T {
	var args []T
	keyword := noKeyword
	for p.Token != closing {
		if p.lxr.AheadSkip(0).Token == COLON {
			keyword = p.keyword()
		} else if keyword != noKeyword {
			p.error("un-named arguments must come before named arguments")
		}

		trueDefault := (keyword != noKeyword &&
			(p.Token == COMMA || p.Token == closing ||
				p.lxr.AheadSkip(0).Token == COLON))

		var val T
		if trueDefault {
			val = p.bld(trueItem)
		} else {
			val = p.expr()
		}
		args = append(args, p.bld(keyword, val))
		p.matchIf(COMMA)
	}
	p.match(closing)
	return args
}

func (p *parser) keyword() Item {
	switch p.Token {
	case STRING, IDENTIFIER, NUMBER:
		break
	default:
		p.error("invalid keyword")
	}
	keyword := p.Item
	p.next()
	p.match(COLON)
	return keyword
}

func (p *parser) block() T {
	p.match(L_CURLY)
	params := p.blockParams()
	statements := p.statements()
	p.match(R_CURLY)
	return p.bld(blockItem, params, statements)
}

func (p *parser) blockParams() T {
	var params []T
	if p.matchIf(BITOR) {
		if p.matchIf(AT) {
			params = append(params, p.bld(Item{Text: "@" + p.Text}))
			p.match(IDENTIFIER)
		} else {
			for p.Token == IDENTIFIER {
				params = append(params, p.bld(p.Item))
				p.match(IDENTIFIER)
				p.matchIf(COMMA)
			}
		}
		p.match(BITOR)
	}
	return p.bld(blockParams, params...)
}
