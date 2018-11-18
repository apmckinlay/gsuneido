package compile

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
	. "github.com/apmckinlay/gsuneido/lexer"
	. "github.com/apmckinlay/gsuneido/runtime"
)

// expression parses a Suneido expression and builds an AST
func (p *parser) expr() ast.Expr {
	return p.pcExpr(1)
}

// pcExpr implements precedence climbing
// each call processes at least one atom
// a given call processes everything >= minprec
// it recurses to process the right hand side of each operator
func (p *parser) pcExpr(minprec int8) ast.Expr {
	e := p.atom()
	// fmt.Println("pcExpr minprec", minprec, "atom", e)
	for p.Token != EOF {
		kt := p.KeyTok()
		prec := precedence[kt]
		// fmt.Println("loop ", p.Item, "prec", prec)
		if prec < minprec {
			break
		}
		if p.newline {
			break
		}
		p.next()
		switch {
		case kt == DOT:
			id := p.Text
			p.match(IDENTIFIER)
			e = p.Mem(e, p.Constant(SuStr(id)))
			if p.Token == L_CURLY && !p.expectingCompound { // a.F { }
				e = p.Call(e, p.arguments(p.Token))
			}
		case kt == INC || kt == DEC: // postfix
			ckLvalue(e)
			e = p.Unary(kt+1, e) // +1 must be POSTINC/DEC
		case kt == IN:
			e = p.in(e)
		case kt == NOT:
			p.match(IN)
			e = p.Unary(NOT, p.in(e))
		case kt == L_BRACKET:
			var expr ast.Expr
			if p.Token == RANGETO || p.Token == RANGELEN {
				expr = nil
			} else {
				expr = p.expr()
			}
			if p.Token == RANGETO || p.Token == RANGELEN {
				rtype := p.Token
				p.next()
				var expr2 ast.Expr
				if p.Token == R_BRACKET {
					expr2 = nil
				} else {
					expr2 = p.expr()
				}
				if rtype == RANGETO {
					e = &ast.RangeTo{E: e, From: expr, To: expr2}
				} else {
					e = &ast.RangeLen{E: e, From: expr, Len: expr2}
				}
			} else {
				e = p.Mem(e, expr)
			}
			p.match(R_BRACKET)
		case ASSIGN_START < kt && kt < ASSIGN_END:
			ckLvalue(e)
			rhs := p.expr()
			e = p.Binary(e, kt, rhs)
		case kt == Q_MARK:
			t := p.expr()
			p.match(COLON)
			f := p.expr()
			e = p.Trinary(e, t, f)
		case kt == L_PAREN: // function call
			e = p.Call(e, p.arguments(kt))
		case ASSOC_START < kt && kt < ASSOC_END:
			// for associative operators, collect a list of contiguous
			es := []ast.Expr{e}
			listtype := flip(kt)
			for {
				rhs := p.pcExpr(prec + 1) // +1 for left associative
				// invert SUB and DIV to combine as ADD and MUL
				if kt == SUB || kt == DIV {
					rhs = p.Unary(kt, rhs)
				}
				es = append(es, rhs)

				kt = p.KeyTok()
				if !p.same(listtype, kt) {
					break
				}
				p.next()
			}
			e = p.Nary(listtype, es)
		default: // other left associative binary operators
			rhs := p.pcExpr(prec + 1) // +1 for left associative
			e = p.Binary(e, kt, rhs)
		}
	}
	return e
}

func (p *parser) in(e ast.Expr) ast.Expr {
	list := []ast.Expr{}
	p.match(L_PAREN)
	for p.Token != R_PAREN {
		list = append(list, p.expr())
		if p.Token == COMMA {
			p.next()
		}
	}
	p.match(R_PAREN)
	return p.In(e, list)
}

func ckLvalue(e ast.Expr) {
	switch e := e.(type) {
	case *ast.Mem:
		return
	case *ast.Ident:
		if isLocal(e.Name) {
			return
		}
	}
	panic("syntax error: lvalue required")
}

func flip(tok Token) Token {
	switch tok {
	case SUB:
		return ADD
	case DIV:
		return MUL
	default:
		return tok
	}
}

func (p *parser) same(listtype Token, next Token) bool {
	if p.newline {
		return false
	}
	return next == listtype ||
		(next == SUB && listtype == ADD) || (next == DIV && listtype == MUL)
}

func (p *parser) atom() ast.Expr {
	switch it := p.Item.KeyTok(); it {
	case TRUE, FALSE, NUMBER, STRING, HASH:
		return p.Constant(p.constant())
	case L_PAREN:
		p.next()
		e := p.expr()
		p.match(R_PAREN)
		return e
	case L_CURLY:
		return p.block()
	case L_BRACKET:
		return p.record()
	case ADD, SUB, NOT, BITNOT:
		p.next()
		return p.Unary(it, p.pcExpr(precedence[L_PAREN]))
	case INC, DEC:
		p.next()
		e := p.pcExpr(precedence[DOT])
		ckLvalue(e)
		return p.Unary(it, e)
	case DOT: // unary, i.e. implicit "this"
		// does not absorb DOT
		p.newline = false
		return p.Ident("this")
	case FUNCTION:
		return p.function()
	case CLASS:
		return p.Constant(p.class())
	case NEW:
		p.next()
		expr := p.pcExpr(precedence[DOT])
		var args []ast.Arg
		if p.matchIf(L_PAREN) {
			args = p.arguments(L_PAREN)
		} else {
			args = []ast.Arg{}
		}
		return p.Call(expr, args)
	case IDENTIFIER, THIS:
		// MyClass { ... } => class
		if !p.expectingCompound &&
			okBase(p.Text) && p.lxr.AheadSkip(0).Token == L_CURLY {
			return p.Constant(p.class())
		}
		e := p.Ident(p.Text)
		p.next()
		return e
	}
	panic(p.error("syntax error: unexpected " + p.Item.String()))
}

var precedence = [Ntokens]int8{
	Q_MARK:    2,
	OR:        3,
	AND:       4,
	IN:        5,
	NOT:       5, // not in
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

func (p *parser) arguments(opening Token) []ast.Arg {
	var args []ast.Arg
	if opening == L_PAREN {
		if p.matchIf(AT) {
			return p.atArgument()
		}
		args = p.argumentList(R_PAREN)
	}
	if p.Token == L_CURLY && !p.expectingCompound {
		args = append(args, ast.Arg{Name: blockArg, E: p.block()})
	}
	return args
}

var atArg = SuStr("@")
var at1Arg = SuStr("@+1")
var blockArg = SuStr("block")

func (p *parser) atArgument() []ast.Arg {
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
	return []ast.Arg{ast.Arg{Name: which, E: expr}}
}

func (p *parser) argumentList(closing Token) []ast.Arg {
	var args []ast.Arg
	var keyword Value
	for p.Token != closing {
		var val ast.Expr
		if p.matchIf(COLON) {
			keyword = SuStr(p.Text)
			val = p.Ident(p.Text)
			p.match(IDENTIFIER)
		} else {
			if p.isKeyword() {
				keyword = p.keyword()
			} else if keyword != nil {
				p.error("un-named arguments must come before named arguments")
			}
			if keyword != nil &&
				(p.Token == COMMA || p.Token == closing || p.isKeyword()) {
				val = p.Constant(True)
			} else {
				val = p.expr()
			}
		}
		if keyword != nil {
			for _, a := range args {
				if keyword.Equal(a.Name) {
					p.error("duplicate argument name (" + keyword.String() + ")")
				}
			}
		}
		args = append(args, ast.Arg{Name: keyword, E: val})
		p.matchIf(COMMA)
	}
	p.match(closing)
	return args
}

func (p *parser) isKeyword() bool {
	return (p.Token == STRING || p.Token == IDENTIFIER || p.Token == NUMBER) &&
		p.lxr.AheadSkip(0).Token == COLON
}

func (p *parser) keyword() Value {
	it := p.Item
	p.next()
	p.match(COLON)
	switch it.Token {
	case STRING, IDENTIFIER:
		return SuStr(it.Text)
	case NUMBER:
		return NumFromString(it.Text)
	}
	panic(p.error("invalid keyword: " + p.Token.String()))
}

func (p *parser) record() ast.Expr {
	p.match(L_BRACKET)
	args := p.argumentList(R_BRACKET)
	return p.Call(p.Ident("Record"), args)
}

func (p *parser) block() *ast.Block {
	p.match(L_CURLY)
	params := p.blockParams()
	body := p.statements()
	p.match(R_CURLY)
	return &ast.Block{ast.Function{Params: params, Body: body}}
}

func (p *parser) blockParams() []ast.Param {
	var params []ast.Param
	if p.matchIf(BITOR) {
		if p.matchIf(AT) {
			params = append(params, ast.Param{Name: "@" + p.Text})
			p.match(IDENTIFIER)
		} else {
			for p.Token == IDENTIFIER {
				params = append(params, ast.Param{Name: p.Text})
				p.match(IDENTIFIER)
				p.matchIf(COMMA)
			}
		}
		p.match(BITOR)
	}
	return params
}
