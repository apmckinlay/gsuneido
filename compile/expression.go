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
	return p.pcExpr(1)
}

var intMaxStr = strconv.Itoa(math.MaxInt32) // used by ranges

// pcExpr implements precedence climbing
// each call processes at least one atom
// a given call processes everything >= minprec
// it recurses to process the right hand side of each operator
func (p *parser) pcExpr(minprec int8) T {
	e := p.atom()
	for {
		kt := p.KeyTok()
		prec := precedence[kt]
		if prec < minprec {
			break
		}
		op := p.Item
		p.next()
		switch {
		case kt == INC || kt == DEC:
			ckLvalue(e)
			op.Text = "post"
			e = p.bld(op, e)
		case kt == IN:
			list := []T{e}
			p.match(L_PAREN)
			for p.Token != R_PAREN {
				list = append(list, p.expr())
				if p.Token == COMMA {
					p.next()
				}
			}
			p.match(R_PAREN)
			e = p.bld(op, list...)
		case kt == L_BRACKET:
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
					expr2 = p.bld(Item{Token: NUMBER, Text: intMaxStr})
				} else {
					expr2 = p.expr()
				}
				expr = p.bld(rtype, expr, expr2)
			}
			p.match(R_BRACKET)
			e = p.bld(op, e, expr)
		case ASSIGN_START < kt && kt < ASSIGN_END:
			ckLvalue(e)
			rhs := p.expr()
			e = p.bld(op, e, rhs)
		case kt == Q_MARK:
			t := p.expr()
			p.match(COLON)
			f := p.expr()
			e = p.bld(op, e, t, f)
		case kt == L_PAREN: // function call
			e = p.bld(call, e, p.arguments(L_PAREN))
		case ASSOC_START < kt && kt < ASSOC_END:
			// for associative operators, collect a list of contiguous
			es := []T{e}
			listtype := flip(op)
			for {
				rhs := p.pcExpr(prec + 1) // +1 for left associative
				// invert SUB and DIV to combine as ADD and MUL
				if op.Token == SUB || op.Token == DIV {
					rhs = p.bld(op, rhs)
				}
				es = append(es, rhs)
				if !same(listtype.Token, p.Token) {
					break
				}
				op = p.Item
				p.next()
			}
			e = p.bld(listtype, es...)
		default: // other left associative binary operators
			rhs := p.pcExpr(prec + 1) // +1 for left associative
			e = p.bld(op, e, rhs)
		}
	}
	return e
}

func ckLvalue(a T) {
	ast := a.(Ast)
	if (ast.Token == IDENTIFIER && isLocal(ast.Text)) ||
		ast.Token == DOT || ast.Token == L_BRACKET {
		return
	}
	panic("syntax error: lvalue required")
}

func flip(op Item) Item {
	switch op.Token {
	case SUB:
		return Item{Token: ADD, Text: "+"}
	case DIV:
		return Item{Token: MUL, Text: "*"}
	default:
		return op
	}
}

func same(listtype Token, next Token) bool {
	return next == listtype ||
		(next == SUB && listtype == ADD) || (next == DIV && listtype == MUL)
}

func (p *parser) atom() T {
	switch p.Token {
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
	case ADD, SUB, NOT, BITNOT:
		it := p.Item
		p.next()
		return p.bld(it, p.atom())
	case INC, DEC:
		it := p.Item
		p.next()
		e := p.pcExpr(precedence[DOT])
		ckLvalue(e)
		return p.bld(it, e)
	case DOT:
		return p.bld(Item{Token: THIS, Text: "this"})
	case IDENTIFIER:
		switch p.Keyword {
		case TRUE, FALSE, THIS:
			return p.evalNext(p.bld(p.Item))
		case NOT:
			it := p.Item
			p.next()
			return p.bld(it, p.atom())
		case FUNCTION:
			return p.function() //TODO
		case NEW:
			it := p.Item
			p.next()
			expr := p.pcExpr(precedence[DOT])
			var args T
			if p.matchIf(L_PAREN) {
				args = p.arguments(L_PAREN)
			} else {
				args = p.bld(argList)
			}
			return p.bld(it, expr, args)
		default:
			return p.evalNext(p.bld(p.Item))
		}
	}
	panic("syntax error: unexpected '" + p.Text + "'")
}

var precedence = [Ntokens]int8{
	Q_MARK:    2,
	OR:        3,
	AND:       4,
	IN:        5,
	BITOR:     6,
	BITXOR:    7,
	BITAND:    8,
	IS:        9,
	ISNT:      9,
	MATCH:     9,
	MATCHNOT:  9,
	LT:        10,
	LTE:       10,
	GT:        10,
	GTE:       10,
	LSHIFT:    11,
	RSHIFT:    11,
	CAT:       12,
	ADD:       12,
	SUB:       12,
	MUL:       13,
	DIV:       13,
	MOD:       13,
	INC:       14,
	DEC:       14,
	L_PAREN:   15,
	DOT:       16,
	L_BRACKET: 16,
	EQ:        16,
	ADDEQ:     16,
	SUBEQ:     16,
	CATEQ:     16,
	MULEQ:     16,
	DIVEQ:     16,
	MODEQ:     16,
	LSHIFTEQ:  16,
	RSHIFTEQ:  16,
	BITOREQ:   16,
	BITANDEQ:  16,
	BITXOREQ:  16,
}

var call = Item{Text: "call"}

func (p *parser) arguments(opening Token) T {
	var args []T
	if opening == L_PAREN {
		if p.matchIf(AT) {
			return p.atArgument()
		}
		args = p.argumentList(R_PAREN)
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
var at1Arg = Item{Text: "at1Arg"}
var noKeyword = Item{Text: "noKwd"}
var trueItem = Item{Token: TRUE, Text: "true"}
var blockArg = Item{Token: IDENTIFIER, Text: "blockArg"}
var blockItem = Item{Text: "block"}
var blockParams = Item{Text: "blockParams"}
var zeroItem = Item{Token: NUMBER, Text: "0"}

func (p *parser) atArgument() T {
	which := atArg
	if p.matchIf(ADD) {
		if p.Item.Text != "1" {
			panic("only @+1 is supported")
		}
		p.match(NUMBER)
		which = at1Arg
	}
	expr := p.expr()
	p.match(R_PAREN)
	return p.bld(which, expr)
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

//TODO validate lvalues i.e. don't allow ++123

//TODO super call
