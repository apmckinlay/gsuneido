// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

import (
	"strings"

	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/ascii"
)

// Expression parses a Suneido expression and builds an AST
func (p *Parser) Expression() ast.Expr {
	return p.pcExpr(1)
}

// pcExpr implements precedence climbing.
// Each call processes at least one atom.
// A given call processes everything >= minprec.
// It recurses to process the right hand side of each operator.
func (p *Parser) pcExpr(minprec int8) ast.Expr {
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
		p.Next()
		switch {
		case token == tok.Dot:
			pos := p.Pos
			id := p.MatchIdent()
			if e == nil {
				e = &ast.Ident{Name: "this", Pos: pos, Implicit: true}
				id = p.privatizeRef(id)
			}
			e = &ast.Mem{E: e, M: p.Constant(SuStr(id)), DotPos: pos}
			if p.Token == tok.LCurly && !p.expectingCompound { // a.F { }
				e = p.Call(e, p.arguments(p.Token), 0)
			}
		case token == tok.Inc || token == tok.Dec: // postfix
			p.ckLvalue(e)
			if id := localVar(e); id != "" {
				p.final[id] = disqualified // modified
			}
			e = p.Unary(token+1, e) // +1 must be PostInc/Dec
		case token == tok.In:
			e = p.in(e)
		case token == tok.Not:
			p.Match(tok.In)
			e = p.Unary(tok.Not, p.in(e))
		case token == tok.LBracket:
			var expr ast.Expr
			if p.Token == tok.RangeTo || p.Token == tok.RangeLen {
				expr = nil
			} else {
				expr = p.Expression()
			}
			if p.Token == tok.RangeTo || p.Token == tok.RangeLen {
				rtype := p.Token
				p.Next()
				var expr2 ast.Expr
				if p.Token == tok.RBracket {
					expr2 = nil
				} else {
					expr2 = p.Expression()
				}
				if rtype == tok.RangeTo {
					e = &ast.RangeTo{E: e, From: expr, To: expr2}
				} else {
					e = &ast.RangeLen{E: e, From: expr, Len: expr2}
				}
			} else {
				e = &ast.Mem{E: e, M: expr}
			}
			p.Match(tok.RBracket)
		case tok.AssignStart < token && token < tok.AssignEnd:
			p.ckLvalue(e)
			assignToLocal := ""
			if id, ok := e.(*ast.Ident); ok {
				name := id.Name
				if token == tok.Eq {
					p.assignName = name
					if ascii.IsLower(name[0]) {
						assignToLocal = name
						p.final[name]++
					}
				} else {
					p.final[name] = disqualified
				}
			}
			rhs := p.Expression()
			if assignToLocal != "" {
				if _, ok := rhs.(*ast.Constant); ok {
					p.assignConst[assignToLocal] = true
				}
			}
			p.assignName = ""
			e = p.Binary(e, token, rhs)
		case token == tok.QMark:
			t := p.Expression()
			p.Match(tok.Colon)
			f := p.Expression()
			e = p.Trinary(e, t, f)
		case token == tok.LParen: // function call
			e = p.Call(e, p.arguments(token), p.endPos)
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
				p.Next()
			}
			e = p.Nary(listtype, es)
		default: // other left associative binary operators
			rhs := p.pcExpr(prec + 1) // +1 for left associative
			e = p.Binary(e, token, rhs)
		}
	}
	return e
}

func (p *Parser) privatizeRef(name string) string {
	if p.className != "" && ascii.IsLower(name[0]) {
		if strings.HasPrefix(name, "getter_") {
			if len(name) <= 7 || !ascii.IsLower(name[7]) {
				p.Error("invalid getter (" + name + ")")
			}
		} else {
			name = p.privatize(name, p.className)
		}
		return name
	} else if strings.HasPrefix(name, "Getter_") &&
		len(name) > 7 && !ascii.IsUpper(name[7]) {
		p.Error("invalid getter (" + name + ")")
	}
	return name
}

func (p *Parser) in(e ast.Expr) ast.Expr {
	list := []ast.Expr{}
	p.Match(tok.LParen)
	for p.Token != tok.RParen {
		list = append(list, p.Expression())
		if p.Token == tok.Comma {
			p.Next()
		}
	}
	p.Match(tok.RParen)
	return p.In(e, list)
}

func (p *Parser) ckLvalue(e ast.Expr) {
	switch e := e.(type) {
	case *ast.Mem:
		return
	case *ast.Ident:
		if e.Name == "this" || e.Name == "super" {
			p.Error("this and super are read-only")
		}
		if isLocal(e.Name) {
			return
		}
	}
	p.Error("lvalue required")
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

func (p *Parser) same(listtype tok.Token, next tok.Token) bool {
	if p.newline {
		return false
	}
	return next == listtype ||
		(next == tok.Sub && listtype == tok.Add) ||
		(next == tok.Div && listtype == tok.Mul)
}

// ------------------------------------------------------------------

// atom handles atoms and prefix operators
func (p *Parser) atom() (result ast.Expr) {
	defer func(pos int32) {
		SetPos(result, pos, p.endPos)
	}(p.Pos)
	switch token := p.Token; token {
	case tok.String:
		// don't call p.constant() because it allows concatenation
		s := SuStr(p.Text)
		p.Next()
		return p.Constant(s)
	case tok.Symbol:
		s := SuStr(p.Text)
		p.Next()
		return p.Symbol(s)
	case tok.True, tok.False, tok.Number, tok.Hash:
		return p.Constant(p.constant())
	case tok.LParen:
		p.Next()
		e := p.Expression()
		p.Match(tok.RParen)
		// need unary for (ob.m)() [not a method call]
		return p.Unary(tok.LParen, e)
	case tok.LCurly:
		defer p.noName()()
		return p.block()
	case tok.LBracket:
		return p.record()
	case tok.Add, tok.Sub, tok.Not, tok.BitNot:
		p.Next()
		return p.Unary(token, p.pcExpr(precedence[tok.LParen]))
	case tok.Inc, tok.Dec: // prefix
		p.Next()
		e := p.pcExpr(precedence[tok.Dot])
		p.ckLvalue(e)
		if id := localVar(e); id != "" {
			p.final[id] = disqualified // modified
		}
		return p.Unary(token, e)
	case tok.Dot: // unary, i.e. implicit "this"
		// does not absorb Dot
		p.newline = false
		return nil // to indicate it should be privatized
	case tok.Function:
		defer p.noName()()
		return p.Constant(p.functionValue())
	case tok.Class:
		defer p.noName()()
		return p.Constant(p.class())
	case tok.New:
		p.Next()
		expr := p.pcExpr(precedence[tok.Dot])
		var args []ast.Arg
		if p.MatchIf(tok.LParen) {
			args = p.arguments(tok.LParen)
		} else {
			args = []ast.Arg{}
		}
		expr = &ast.Mem{E: expr, M: p.Constant(SuStr("*new*"))}
		return p.Call(expr, args, 0)
	default:
		if p.Token.IsIdent() {
			if p.Text == "_" {
				p.Error("invalid identifier: '_'")
			}
			if p.Text[0] == '_' && len(p.Text) > 1 && ascii.IsUpper(p.Text[1]) &&
				p.Text[1:] != p.name && p.Lxr.AheadSkip(0).Token != tok.Colon {
				p.Error("invalid reference to " + p.Text)
			}
			if !p.expectingCompound &&
				okBase(p.Text) && p.Lxr.AheadSkip(0).Token == tok.LCurly {
				// MyClass { ... } => class
				defer p.noName()()
				return p.Constant(p.class())
			}
			if p.Text == "it" {
				p.itUsed = true
			}
			if p.Text == "dll" || p.Text == "callback" || p.Text == "struct" {
				p.Error("gSuneido does not implement " + p.Text)
			}
			e := &ast.Ident{Name: p.Text}
			p.Next()
			return e
		}
	}
	panic(p.Error("unexpected " + p.Item.String()))
}

// noName assigns names to anonymous functions and classes.
// Usage: defer p.noName()()
func (p *Parser) noName() func() {
	prevName := p.name
	// assignName is set by pcExpr
	p.name = strings.TrimSpace(p.name + " " + p.assignName)
	return func() { p.name = prevName }
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

func (p *Parser) arguments(opening tok.Token) []ast.Arg {
	var args []ast.Arg
	if opening == tok.LParen {
		if p.MatchIf(tok.At) {
			return p.atArgument()
		}
		args = p.argumentList(tok.RParen)
	}
	if p.Token == tok.LCurly && !p.expectingCompound {
		defer p.noName()()
		args = append(args, ast.Arg{Name: blockArg, E: p.block()})
	}
	return args
}

var atArg = SuStr("@")
var at1Arg = SuStr("@+1")
var blockArg = SuStr("block")

func (p *Parser) atArgument() []ast.Arg {
	which := atArg
	if p.MatchIf(tok.Add) {
		if p.Item.Text != "1" {
			panic("only @+1 is supported")
		}
		p.Match(tok.Number)
		which = at1Arg
	}
	expr := p.Expression()
	p.Match(tok.RParen)
	return []ast.Arg{{Name: which, E: expr}}
}

func (p *Parser) argumentList(closing tok.Token) []ast.Arg {
	var args []ast.Arg
	haveNamed := false
	unnamed := func(val ast.Expr) {
		if haveNamed {
			p.Error("un-named arguments must come before named arguments")
		}
		args = append(args, ast.Arg{E: val})
	}
	named := func(name Value, val ast.Expr, pos, end int32) {
		for _, a := range args {
			if name.Equal(a.Name) {
				p.Error("duplicate argument name: " + ToStrOrString(name))
			}
		}
		arg := ast.Arg{Name: name, E: val}
		arg.SetPos(pos, end)
		args = append(args, arg)
		haveNamed = true
	}
	var pending Value
	var pos, pendingPos int32
	handlePending := func(val ast.Expr, end int32) {
		if pending != nil {
			named(pending, val, pendingPos, end)
			pending = nil
		}
	}
	for p.Token != closing {
		pos = p.Pos
		endPos := p.endPos
		if p.MatchIf(tok.Colon) { // :name shortcut
			if !p.Token.IsIdent() || !ascii.IsLower(p.Text[0]) {
				p.Error("expecting local variable name")
			}
			handlePending(p.Constant(True), endPos)
			name := p.MatchIdent()
			named(SuStr(name), &ast.Ident{Name: name}, pos, p.endPos)
		} else {
			expr := p.Expression() // could be name or value
			if name := p.argname(expr); name != nil && p.MatchIf(tok.Colon) {
				handlePending(p.Constant(True), endPos)
				pending = name // it's a name but don't know value yet
				pendingPos = pos
			} else if pending != nil {
				handlePending(expr, p.endPos)
			} else {
				unnamed(expr)
			}
		}
		if p.MatchIf(tok.Comma) {
			handlePending(p.Constant(True), p.endPos)
		} else if p.newline && pending == nil {
			switch p.Token {
			case tok.LParen, tok.LBracket, tok.LCurly, tok.Dot, tok.Add, tok.Sub:
				p.CheckResult(int(p.Pos), "ERROR: missing comma")
			}
		}
	}
	handlePending(p.Constant(True), p.endPos)
	p.Match(closing)
	return args
}

func (p *Parser) argname(expr ast.Expr) Value {
	if id, ok := expr.(*ast.Ident); ok {
		return SuStr(id.Name)
	}
	if c, ok := expr.(*ast.Constant); ok {
		return c.Val
	}
	return nil
}

func (p *Parser) record() ast.Expr {
	pos := p.Pos
	p.Match(tok.LBracket)
	args := p.argumentList(tok.RBracket)
	fn := "Record"
	if hasUnnamed(args) {
		fn = "Object"
	}
	id := &ast.Ident{Name: fn, Pos: pos, Implicit: true}
	return p.Call(id, args, 0)
}

func hasUnnamed(args []ast.Arg) bool {
	return len(args) > 0 && args[0].Name == nil
}

func (p *Parser) block() *ast.Block {
	p.hasBlocks = true
	itUsedPrev := p.itUsed
	p.itUsed = false
	pos := p.Pos
	p.Match(tok.LCurly)
	params := p.blockParams()
	body := p.statements()
	p.Match(tok.RCurly)
	if p.itUsed && len(params) == 0 {
		params = append(params, mkParam("it", pos, pos, false, nil))
		p.final["it"] = disqualified
	}
	p.itUsed = itUsedPrev
	return &ast.Block{Name: p.name,
		Function: ast.Function{Params: params, Body: body}}
}

func (p *Parser) blockParams() []ast.Param {
	var params []ast.Param
	if p.MatchIf(tok.BitOr) {
		pos := p.Pos
		if p.MatchIf(tok.At) {
			params = append(params, mkParam("@"+p.Text,
				pos, p.Pos+int32(len(p.Text)), p.unusedAhead(), nil))
			p.final[p.Text] = disqualified
			p.MatchIdent()
		} else {
			for p.Token.IsIdent() {
				params = append(params, mkParam(p.Text,
					p.Pos, p.Pos+int32(len(p.Text)), p.unusedAhead(), nil))
				p.final[unDyn(p.Text)] = disqualified
				p.MatchIdent()
				p.MatchIf(tok.Comma)
			}
		}
		p.Match(tok.BitOr)
	}
	return params
}

func localVar(node ast.Node) string {
	if id, ok := node.(*ast.Ident); ok && ascii.IsLower(id.Name[0]) {
		return id.Name
	}
	return ""
}
