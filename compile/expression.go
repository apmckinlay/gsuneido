// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

import (
	"strings"

	"github.com/apmckinlay/gsuneido/compile/ast"
	. "github.com/apmckinlay/gsuneido/lexer"
	tok "github.com/apmckinlay/gsuneido/lexer/tokens"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/ascii"
	. "github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/str"
)

// expression parses a Suneido expression and builds an AST
func (p *parser) expr() ast.Expr {
	return p.pcExpr(1)
}

// ------------------------------------------------------------------
// pcExpr implements precedence climbing
// each call processes at least one atom
// a given call processes everything >= minprec
// it recurses to process the right hand side of each operator
func (p *parser) pcExpr(minprec int8) ast.Expr {
	e := p.atom()
	// fmt.Println("pcExpr minprec", minprec, "atom", e)
	for p.Token != tok.Eof {
		token := p.Token
		prec := precedence[token]
		// fmt.Println("loop ", p.Item, "prec", prec)
		if prec < minprec {
			break
		}
		if p.newline {
			break
		}
		p.next()
		switch {
		case token == tok.Dot:
			id := p.Text
			p.matchIdent()
			if e == nil {
				e = p.Ident("this")
				id = p.privatizeRef(id)
			}
			e = p.Mem(e, p.Constant(SuStr(id)))
			if p.Token == tok.LCurly && !p.expectingCompound { // a.F { }
				e = p.Call(e, p.arguments(p.Token))
			}
		case token == tok.Inc || token == tok.Dec: // postfix
			p.ckLvalue(e)
			e = p.Unary(token+1, e) // +1 must be PostInc/Dec
		case token == tok.In:
			e = p.in(e)
		case token == tok.Not:
			p.match(tok.In)
			e = p.Unary(tok.Not, p.in(e))
		case token == tok.LBracket:
			var expr ast.Expr
			if p.Token == tok.RangeTo || p.Token == tok.RangeLen {
				expr = nil
			} else {
				expr = p.expr()
			}
			if p.Token == tok.RangeTo || p.Token == tok.RangeLen {
				rtype := p.Token
				p.next()
				var expr2 ast.Expr
				if p.Token == tok.RBracket {
					expr2 = nil
				} else {
					expr2 = p.expr()
				}
				if rtype == tok.RangeTo {
					e = &ast.RangeTo{E: e, From: expr, To: expr2}
				} else {
					e = &ast.RangeLen{E: e, From: expr, Len: expr2}
				}
			} else {
				e = p.Mem(e, expr)
			}
			p.match(tok.RBracket)
		case tok.AssignStart < token && token < tok.AssignEnd:
			p.ckLvalue(e)
			if id, ok := e.(*ast.Ident); ok && token == tok.Eq {
				p.assignName = id.Name
			}
			rhs := p.expr()
			p.assignName = ""
			e = p.Binary(e, token, rhs)
		case token == tok.QMark:
			t := p.expr()
			p.match(tok.Colon)
			f := p.expr()
			e = p.Trinary(e, t, f)
		case token == tok.LParen: // function call
			e = p.Call(e, p.arguments(token))
		case tok.AssocStart < token && token < tok.AssocEnd:
			// for associative operators, collect a list of contiguous
			es := []ast.Expr{e}
			listtype := flip(token)
			for {
				rhs := p.pcExpr(prec + 1) // +1 for left associative
				// invert Sub and Div to combine as Add and Mul
				if token == tok.Sub || token == tok.Div {
					rhs = p.Unary(token, rhs)
				}
				es = append(es, rhs)

				token = p.Token
				if !p.same(listtype, token) {
					break
				}
				p.next()
			}
			e = p.Nary(listtype, es)
		default: // other left associative binary operators
			rhs := p.pcExpr(prec + 1) // +1 for left associative
			e = p.Binary(e, token, rhs)
		}
	}
	return e
}

func (p *parser) privatizeRef(name string) string {
	if p.className != "" && ascii.IsLower(name[0]) {
		if strings.HasPrefix(name, "getter_") {
			if len(name) <= 7 || !ascii.IsLower(name[7]) {
				p.error("invalid getter (" + name + ")")
			}
		} else {
			name = p.className + "_" + name
		}
		return name
	} else if strings.HasPrefix(name, "Getter_") && len(name) > 7 &&
		!ascii.IsUpper(name[7]) {
		p.error("invalid getter (" + name + ")")
	}
	return name
}

func (p *parser) in(e ast.Expr) ast.Expr {
	list := []ast.Expr{}
	p.match(tok.LParen)
	for p.Token != tok.RParen {
		list = append(list, p.expr())
		if p.Token == tok.Comma {
			p.next()
		}
	}
	p.match(tok.RParen)
	return p.In(e, list)
}

func (p *parser) ckLvalue(e ast.Expr) {
	switch e := e.(type) {
	case *ast.Mem:
		return
	case *ast.Ident:
		if e.Name == "this" || e.Name == "super" {
			p.error("this and super are read-only")
		}
		if isLocal(e.Name) {
			return
		}
	}
	p.error("syntax error: lvalue required")
}

func flip(token tok.Token) tok.Token {
	switch token {
	case tok.Sub:
		return tok.Add
	case tok.Div:
		return tok.Mul
	default:
		return token
	}
}

func (p *parser) same(listtype tok.Token, next tok.Token) bool {
	if p.newline {
		return false
	}
	return next == listtype ||
		(next == tok.Sub && listtype == tok.Add) ||
		(next == tok.Div && listtype == tok.Mul)
}

// ------------------------------------------------------------------
// atom handles atoms and prefix operators
func (p *parser) atom() ast.Expr {
	switch token := p.Token; token {
	case tok.String:
		// don't call p.constant() because it allows concatenation
		s := p.Text
		p.match(tok.String)
		return p.Constant(SuStr(s))
	case tok.True, tok.False, tok.Number, tok.Hash:
		return p.Constant(p.constant())
	case tok.LParen:
		p.next()
		e := p.expr()
		p.match(tok.RParen)
		// need unary for (ob.m)() [not a method call]
		return p.Unary(tok.LParen, e)
	case tok.LCurly:
		return p.block()
	case tok.LBracket:
		return p.record()
	case tok.Add, tok.Sub, tok.Not, tok.BitNot:
		p.next()
		return p.Unary(token, p.pcExpr(precedence[tok.LParen]))
	case tok.Inc, tok.Dec:
		p.next()
		e := p.pcExpr(precedence[tok.Dot])
		p.ckLvalue(e)
		return p.Unary(token, e)
	case tok.Dot: // unary, i.e. implicit "this"
		// does not absorb Dot
		p.newline = false
		return nil // to indicate it should be privatized
	case tok.Function:
		return p.Constant(p.noName(p.functionValue))
	case tok.Class:
		return p.Constant(p.noName(p.class))
	case tok.New:
		p.next()
		expr := p.pcExpr(precedence[tok.Dot])
		var args []ast.Arg
		if p.matchIf(tok.LParen) {
			args = p.arguments(tok.LParen)
		} else {
			args = []ast.Arg{}
		}
		expr = p.Mem(expr, p.Constant(SuStr("*new*")))
		return p.Call(expr, args)
	default:
		if p.Token.IsIdent() {
			if p.Text[0] == '_' && len(p.Text) > 1 && ascii.IsUpper(p.Text[1]) &&
				p.Text[1:] != p.name && p.lxr.AheadSkip(0).Token != tok.Colon {
				p.error("invalid reference to " + p.Text)
			}
			if !p.expectingCompound &&
				okBase(p.Text) && p.lxr.AheadSkip(0).Token == tok.LCurly {
				// MyClass { ... } => class
				return p.Constant(p.noName(p.class))
			}
			if p.Text == "dll" || p.Text == "callback" || p.Text == "struct" {
				p.error("gSuneido does not implement " + p.Text)
			}
			e := p.Ident(p.Text)
			p.next()
			return e
		}
	}
	panic(p.error("syntax error: unexpected " + p.Item.String()))
}

func (p *parser) noName(f func() Value) Value {
	prevName := p.name
	name := p.assignName
	if p.assignName == "" {
		name = "?"
	}
	p.name = str.Opt(p.name, " ") + name
	result := f()
	p.name = prevName
	return result
}

var precedence = [tok.Ntokens]int8{
	tok.QMark:    2,
	tok.Or:       3,
	tok.And:      4,
	tok.In:       5,
	tok.Not:      5, // not in
	tok.BitOr:    6,
	tok.BitXor:   7,
	tok.BitAnd:   8,
	tok.Is:       9,
	tok.Isnt:     9,
	tok.Match:    9,
	tok.MatchNot: 9,
	tok.Lt:       10,
	tok.Lte:      10,
	tok.Gt:       10,
	tok.Gte:      10,
	tok.LShift:   11,
	tok.RShift:   11,
	tok.Cat:      12,
	tok.Add:      12,
	tok.Sub:      12,
	tok.Mul:      13,
	tok.Div:      13,
	tok.Mod:      13,
	tok.Inc:      14,
	tok.Dec:      14,
	tok.LParen:   15,
	tok.Dot:      16,
	tok.LBracket: 16,
	tok.Eq:       16,
	tok.AddEq:    16,
	tok.SubEq:    16,
	tok.CatEq:    16,
	tok.MulEq:    16,
	tok.DivEq:    16,
	tok.ModEq:    16,
	tok.LShiftEq: 16,
	tok.RShiftEq: 16,
	tok.BitOrEq:  16,
	tok.BitAndEq: 16,
	tok.BitXorEq: 16,
}

var call = Item{Text: "call"}

func (p *parser) arguments(opening tok.Token) []ast.Arg {
	var args []ast.Arg
	if opening == tok.LParen {
		if p.matchIf(tok.At) {
			return p.atArgument()
		}
		args = p.argumentList(tok.RParen)
	}
	if p.Token == tok.LCurly && !p.expectingCompound {
		args = append(args, ast.Arg{Name: blockArg, E: p.block()})
	}
	return args
}

var atArg = SuStr("@")
var at1Arg = SuStr("@+1")
var blockArg = SuStr("block")

func (p *parser) atArgument() []ast.Arg {
	which := atArg
	if p.matchIf(tok.Add) {
		if p.Item.Text != "1" {
			panic("only @+1 is supported")
		}
		p.match(tok.Number)
		which = at1Arg
	}
	expr := p.expr()
	p.match(tok.RParen)
	return []ast.Arg{ast.Arg{Name: which, E: expr}}
}

func (p *parser) argumentList(closing tok.Token) []ast.Arg {
	var args []ast.Arg
	haveNamed := false
	unnamed := func(val ast.Expr) {
		if haveNamed {
			p.error("un-named arguments must come before named arguments")
		}
		args = append(args, ast.Arg{E: val})
	}
	named := func(name Value, val ast.Expr) {
		for _, a := range args {
			if name.Equal(a.Name) {
				p.error("duplicate argument name: " + ToStrOrString(name))
			}
		}
		args = append(args, ast.Arg{Name: name, E: val})
		haveNamed = true
	}
	var pending Value
	handlePending := func(val ast.Expr) {
		if pending != nil {
			named(pending, val)
			pending = nil
		}
	}
	for p.Token != closing {
		var expr ast.Expr
		if p.matchIf(tok.Colon) {
			if !IsLower(p.Text[0]) {
				p.error("expecting local variable name")
			}
			handlePending(p.Constant(True))
			named(SuStr(p.Text), p.Ident(p.Text))
			p.matchIdent()
		} else {
			expr = p.expr() // could be name or value
			if name := p.argname(expr); name != nil && p.matchIf(tok.Colon) {
				handlePending(p.Constant(True))
				pending = name // it's a name but don't know value yet
			} else if pending != nil {
				handlePending(expr)
			} else {
				unnamed(expr)
			}
		}
		if p.matchIf(tok.Comma) {
			handlePending(p.Constant(True))
		}
	}
	p.match(closing)
	handlePending(p.Constant(True))
	return args
}

func (p *parser) argname(expr ast.Expr) Value {
	if id, ok := expr.(*ast.Ident); ok {
		return SuStr(id.Name)
	}
	if c, ok := expr.(*ast.Constant); ok {
		return c.Val
	}
	return nil
}

func (p *parser) record() ast.Expr {
	p.match(tok.LBracket)
	args := p.argumentList(tok.RBracket)
	fn := "Record"
	if hasUnnamed(args) {
		fn = "Object"
	}
	// note: ident will have wrong pos
	return p.Call(p.Ident(fn), args)
}

func hasUnnamed(args []ast.Arg) bool {
	return len(args) > 0 && args[0].Name == nil
}

func (p *parser) block() *ast.Block {
	pos := p.Pos
	p.match(tok.LCurly)
	params := p.blockParams()
	body := p.statements()
	p.match(tok.RCurly)
	return &ast.Block{Function: ast.Function{Pos: pos, Params: params, Body: body}}
}

func (p *parser) blockParams() []ast.Param {
	var params []ast.Param
	if p.matchIf(tok.BitOr) {
		if p.matchIf(tok.At) {
			params = append(params, ast.MkParam("@" + p.Text))
			p.matchIdent()
		} else {
			for p.Token.IsIdent() {
				params = append(params, ast.MkParam(p.Text))
				p.matchIdent()
				p.matchIf(tok.Comma)
			}
		}
		p.match(tok.BitOr)
	}
	return params
}
