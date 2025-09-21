// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/core"
)

// Function parses a function (starting with the "function" keyword)
func (p *Parser) Function() *ast.Function {
	p.Match(tok.Function)
	return p.function(false)
}

// function parses a function or method (without the "function" keyword)
func (p *Parser) function(inClass bool) (result *ast.Function) {
	funcInfoSave := p.funcInfo
	p.InitFuncInfo()
	params := p.params(inClass)
	pos1 := p.EndPos
	p.Match(tok.LCurly)
	pos2 := p.EndPos
	body := p.statements()
	p.Match(tok.RCurly)
	p.processFinal()
	fn := &ast.Function{Params: params, Body: body, Final: p.final,
		HasBlocks: p.hasBlocks, Pos1: pos1, Pos2: pos2}
	p.funcInfo = funcInfoSave
	return fn
}

func (p *Parser) processFinal() {
	finalConst := false
	for k, n := range p.final {
		if n == 1 && p.assignConst[k] {
			finalConst = true
			break
		}
	}
	if !finalConst {
		p.final = nil
		return
	}
	for k, n := range p.final {
		if n != 1 {
			delete(p.final, k)
		}
	}
}

func (p *Parser) params(inClass bool) []ast.Param {
	p.Match(tok.LParen)
	var params []ast.Param
	addParam := func(name string, pos int32, unused bool, def Value) {
		if name == "unused" || name == "@unused" {
			unused = true
		}
		if def != nil {
			for i := range params {
				x := params[i].DefVal
				if def.Equal(x) && def.Type() == x.Type() {
					def = x // reuse value
				}
			}
		}
		param := mkParam(name, pos, p.EndPos, unused, def)
		params = append(params, param)
	}
	if p.Token == tok.At {
		pos := p.Pos
		p.Match(tok.At)
		name := p.Text
		unused := p.unusedAhead()
		p.MatchIdent()
		addParam("@"+name, pos, unused, nil)
		p.final[name] = disqualified
	} else {
		defs := false
		for p.Token != tok.RParen {
			dot := p.MatchIf(tok.Dot)
			name := p.Text
			p.final[unDyn(name)] = disqualified
			pos := p.Pos
			if dot {
				if !inClass {
					p.Error("dot parameters only allowed in class methods")
				}
				name = "." + name
			}
			unused := p.unusedAhead()
			p.MatchIdent()
			p.checkForDupParam(params, name)
			if p.MatchIf(tok.Eq) {
				wasString := p.Token == tok.String || p.Token == tok.Symbol
				defs = true
				def := p.constant()
				if _, ok := def.(SuStr); ok && !wasString {
					p.Error("parameter defaults must be constants")
				}
				addParam(name, pos, unused, def)
				p.MatchIf(tok.Comma)
			} else {
				if defs {
					p.Error("default parameters must come last")
				}
				addParam(name, pos, unused, nil)
				p.MatchIf(tok.Comma)
			}
		}
	}
	p.Match(tok.RParen)
	return params
}

func mkParam(name string, pos, end int32, unused bool, def Value) ast.Param {
	if name == "unused" || name == "@unused" {
		unused = true
	}
	return ast.Param{Name: ast.Ident{Name: name, Pos: pos},
		DefVal: def, Unused: unused, End: end}
}

// unDyn removes the leading underscore from dynamic parameters
func unDyn(id string) string {
	if len(id) > 0 && id[0] == '_' {
		return id[1:]
	}
	return id
}

// unusedAhead detects whether /*unused*/ is next
func (p *Parser) unusedAhead() bool {
	i := 0
	for ; p.Lxr.Ahead(i).Token == tok.Whitespace; i++ {
	}
	return p.Lxr.Ahead(i).Text == "/*unused*/"
}

func (p *Parser) checkForDupParam(params []ast.Param, name string) {
	for _, a := range params {
		if a.Name.Name == name {
			p.Error("duplicate function parameter (" + name + ")")
		}
	}
}

func (p *Parser) compound() []ast.Statement {
	p.Match(tok.LCurly)
	stmts := p.statements()
	p.Match(tok.RCurly)
	return stmts
}

func (p *Parser) statements() []ast.Statement {
	list := []ast.Statement{}
	for p.Token != tok.RCurly {
		stmt := p.statement()
		list = append(list, stmt)
	}
	return list
}

func (p *Parser) statement() (result ast.Statement) {
	defer func(org int32) {
		SetPos(result, org, p.EndPos)
	}(p.Pos)
	token := p.Token
	switch token {
	case tok.Semicolon:
		for p.MatchIf(tok.Semicolon) {
		}
		return &ast.Compound{Body: []ast.Statement{}}
	case tok.LCurly:
		if p.Lxr.AheadSkip(0).Token == tok.BitOr { // block
			return &ast.ExprStmt{E: p.Expression()}
		}
		return &ast.Compound{Body: p.compound()}
	case tok.Return:
		p.Next()
		return p.returnStmt()
	case tok.If:
		p.Next()
		return p.ifStmt()
	case tok.Switch:
		p.Next()
		return p.switchStmt()
	case tok.Forever:
		p.Next()
		return p.foreverStmt()
	case tok.While:
		p.Next()
		return p.whileStmt()
	case tok.Do:
		p.Next()
		return p.semi(p.dowhileStmt())
	case tok.For:
		return p.forStmt()
	case tok.Throw:
		p.Next()
		return &ast.Throw{E: p.trailingExpr()}
	case tok.Try:
		p.Next()
		return p.tryStmt()
	case tok.Break:
		p.Next()
		return p.semi(&ast.Break{})
	case tok.Continue:
		p.Next()
		return p.semi(&ast.Continue{})
	default:
		exprs := p.exprList()
		if len(exprs) == 1 {
			return &ast.ExprStmt{E: exprs[0]}
		}
		last := exprs[len(exprs)-1]
		bin, ok := last.(*ast.Binary)
		if !ok || bin.Tok != tok.Eq {
			p.Error()
		}
		exprs[len(exprs)-1] = bin.Lhs
		if !p.areLocals(exprs) {
			p.Error()
		}
		if _, ok = bin.Rhs.(*ast.Call); !ok {
			p.Error()
		}
		for _, expr := range exprs {
			p.final[expr.(*ast.Ident).Name] = disqualified
		}
		return &ast.MultiAssign{Lhs: exprs, Rhs: bin.Rhs}
	}
}

// exprList gets a comma-separated list of expressions.
// It is similar to trailingExpr
func (p *Parser) exprList() []ast.Expr {
	exprs := make([]Expr, 0, 1)
	for {
		expr := p.Expression()
		exprs = append(exprs, expr)
		switch p.Token {
		case tok.Comma:
			p.Next()
			continue
		case tok.Semicolon:
			p.Next()
			return exprs
		case tok.RCurly, tok.Return, tok.If, tok.Else, tok.Switch, tok.Forever,
			tok.While, tok.Do, tok.For, tok.Throw, tok.Try, tok.Catch,
			tok.Break, tok.Continue, tok.Case, tok.Default:
			return exprs
		default:
			if !p.newline {
				p.Error()
			}
			return exprs
		}
	}
}

// areLocals checks if all expressions in the list are local variables
func (p *Parser) areLocals(exprs []ast.Expr) bool {
	for _, expr := range exprs {
		if id, ok := expr.(*ast.Ident); !ok ||
			!isLocal(id.Name) || id.Name == "this" || id.Name == "super" {
			return false
		}
	}
	return true
}

// trailingExpr gives a syntax error for two expressions side by side
// without either a semicolon or newline separator
// because this is not very readable and often a mistake
// cSuneido and jSuneido only allowed catch, while, or else to follow
// but here we are more lenient and allow any statement keyword
func (p *Parser) trailingExpr() ast.Expr {
	expr := p.Expression()
	switch p.Token {
	case tok.Semicolon:
		p.Next()
	case tok.RCurly, tok.Return, tok.If, tok.Else, tok.Switch, tok.Forever,
		tok.While, tok.Do, tok.For, tok.Throw, tok.Try, tok.Catch,
		tok.Break, tok.Continue, tok.Case, tok.Default:
		// ok
	default:
		if !p.newline {
			p.Error()
		}
	}
	return expr
}

func (p *Parser) semi(stmt ast.Statement) ast.Statement {
	p.MatchIf(tok.Semicolon)
	return stmt
}

func (p *Parser) ifStmt() *ast.If {
	expr := p.ctrlExpr()
	t := p.statement()
	stmt := &ast.If{Cond: expr, Then: t}
	if p.MatchIf(tok.Else) {
		stmt.ElseEnd = p.EndPos
		stmt.Else = p.statement()
	}
	return stmt
}

func (p *Parser) switchStmt() (result *ast.Switch) {
	var expr ast.Expr
	if p.Token == tok.LCurly {
		expr = p.Constant(True)
	} else {
		expr = p.exprExpecting(true)
	}
	pos1 := p.EndPos
	p.Match(tok.LCurly)
	pos2 := p.EndPos
	var cases []ast.Case
	for p.Token == tok.Case {
		cases = append(cases, p.switchCase())
	}
	var def []ast.Statement
	var posdef int32
	if p.MatchIf(tok.Default) {
		p.Match(tok.Colon)
		posdef = p.EndPos
		def = p.switchBody()
	}
	p.Match(tok.RCurly)
	return &ast.Switch{E: expr, Cases: cases, Default: def,
		Pos1: pos1, Pos2: pos2, PosDef: posdef}
}

func (p *Parser) switchCase() ast.Case {
	pos := p.Pos
	p.Match(tok.Case)
	var exprs []ast.Expr
	for {
		pos := p.Pos
		expr := p.Expression()
		exprs = append(exprs, p.exprPos(expr, pos, p.EndPos))
		if !p.MatchIf(tok.Comma) {
			break
		}
	}
	p.Match(tok.Colon)
	end := p.EndPos
	body := p.switchBody()
	c := ast.Case{Exprs: exprs, Body: body}
	c.SetPos(pos, end)
	return c
}

func (p *Parser) switchBody() []ast.Statement {
	stmts := []ast.Statement{}
	for p.Token != tok.RCurly && p.Token != tok.Case && p.Token != tok.Default {
		stmts = append(stmts, p.statement())
	}
	return stmts
}

func (p *Parser) foreverStmt() *ast.Forever {
	body := p.statement()
	return &ast.Forever{Body: body}
}

func (p *Parser) whileStmt() *ast.While {
	cond := p.ctrlExpr()
	body := p.statement()
	return &ast.While{Cond: cond, Body: body}
}

func (p *Parser) dowhileStmt() *ast.DoWhile {
	body := p.statement()
	p.Match(tok.While)
	pos := p.Pos
	expr := p.Expression()
	return &ast.DoWhile{Body: body,
		Cond: p.exprPos(expr, pos, p.EndPos)}
}

func (p *Parser) forStmt() ast.Statement {
	// easier to check before matching For so everything is ahead
	forIn := p.isForIn()
	p.Match(tok.For)
	if forIn {
		return p.forIn()
	} else if p.Token != tok.LParen {
		return p.forRange()
	}
	return p.forClassic()
}

func (p *Parser) isForIn() bool {
	i := 0
	if p.Lxr.AheadSkip(i).Token == tok.LParen {
		i++
	}
	if !p.Lxr.AheadSkip(i).Token.IsIdent() {
		return false
	}
	i++
	if p.Lxr.AheadSkip(i).Token == tok.Comma {
		i++
		if !p.Lxr.AheadSkip(i).Token.IsIdent() {
			return false
		}
		i++
	}
	return p.Lxr.AheadSkip(i).Token == tok.In
}

func (p *Parser) forIn() *ast.ForIn {
	parens := p.MatchIf(tok.LParen)
	id := p.Text
	p.final[id] = disqualified
	pos := p.Pos
	p.MatchIdent()
	var var2 ast.Ident
	if p.MatchIf(tok.Comma) {
		var2.Name = p.Text
		p.final[var2.Name] = disqualified
		var2.Pos = p.Pos
		p.MatchIdent()
	}
	p.Match(tok.In)
	var expr ast.Expr
	if p.Token == tok.RangeTo {
		expr = p.Constant(Zero)
	} else {
		pos := p.Pos
		expr = p.exprExpecting(!parens)
		expr = p.exprPos(expr, pos, p.EndPos)
	}
	var expr2 Expr
	if !parens && p.Token == tok.RangeTo {
		p.Next()
		pos := p.Pos
		expr2 = p.exprExpecting(!parens)
		expr2 = p.exprPos(expr2, pos, p.EndPos)
	}
	if parens {
		p.Match(tok.RParen)
	}
	body := p.statement()
	return &ast.ForIn{Var: ast.Ident{Name: id, Pos: pos}, Var2: var2,
		E: expr, E2: expr2, Body: body}
}

func (p *Parser) forRange() *ast.ForIn {
	p.Match(tok.RangeTo)
	expr := p.Constant(Zero)
	expr2 := p.exprExpecting(true)
	body := p.statement()
	return &ast.ForIn{E: expr, E2: expr2, Body: body}
}

func (p *Parser) forClassic() *ast.For {
	p.Match(tok.LParen)
	init := p.optExprList(tok.Semicolon)
	p.Match(tok.Semicolon)
	var cond ast.Expr
	if p.Token != tok.Semicolon {
		cond = p.Expression()
	}
	p.Match(tok.Semicolon)
	inc := p.optExprList(tok.RParen)
	p.Match(tok.RParen)
	body := p.statement()
	return &ast.For{Init: init, Cond: cond, Inc: inc, Body: body}
}

func (p *Parser) optExprList(after tok.Token) []ast.Expr {
	exprs := []ast.Expr{}
	if p.Token != after {
		for {
			exprs = append(exprs, p.Expression())
			if p.Token != tok.Comma {
				break
			}
			p.Next()
		}
	}
	return exprs
}

// used by if and while
func (p *Parser) ctrlExpr() ast.Expr {
	parens := p.MatchIf(tok.LParen)
	keepParens := p.Token == tok.LParen
	pos := p.Pos
	expr := p.exprExpecting(!parens)
	expr = p.exprPos(expr, pos, p.EndPos)
	if parens {
		p.Match(tok.RParen)
	}
	if keepParens {
		return p.Unary(tok.LParen, expr)
	}
	return expr
}

func (p *Parser) exprExpecting(expecting bool) ast.Expr {
	p.expectingCompound = expecting
	expr := p.Expression()
	p.expectingCompound = false
	return expr
}

func (p *Parser) returnStmt() *ast.Return {
	if p.newline || p.MatchIf(tok.Semicolon) || p.Token == tok.RCurly {
		return &ast.Return{}
	}
	returnThrow := false
	if p.MatchIf(tok.Throw) {
		return &ast.Return{Exprs: []Expr{p.trailingExpr()}, ReturnThrow: true}
	}
	exprs := p.returnExprs()
	return &ast.Return{Exprs: exprs, ReturnThrow: returnThrow}
}

func (p *Parser) returnExprs() []ast.Expr {
	exprs := make([]Expr, 0, 1)
	for {
		exprs = append(exprs, p.Expression())
		switch p.Token {
		case tok.Comma:
			p.Next()
			continue
		case tok.Semicolon:
			p.Next()
			return exprs
		case tok.RCurly, tok.Return, tok.If, tok.Else, tok.Switch, tok.Forever,
			tok.While, tok.Do, tok.For, tok.Throw, tok.Try, tok.Catch,
			tok.Break, tok.Continue, tok.Case, tok.Default:
			return exprs
		default:
			if p.newline {
				return exprs
			}
			p.Error()
		}
	}
}

func (p *Parser) tryStmt() *ast.TryCatch {
	if p.inTry {
		p.Error("nested try not supported")
	}
	p.funcInfo.inTry = true
	try := p.statement()
	p.funcInfo.inTry = false
	var catchVar string
	var varPos int32
	var catchFilter string
	var catch ast.Statement
	var catchEnd int32
	var unused bool
	catchPos := p.Pos
	if p.MatchIf(tok.Catch) {
		if p.MatchIf(tok.LParen) {
			catchVar = p.Text
			p.final[catchVar] = disqualified
			varPos = p.Pos
			unused = p.unusedAhead()
			p.MatchIdent()
			if p.MatchIf(tok.Comma) {
				catchFilter = p.Text
				p.Match(tok.String)
			}
			p.Match(tok.RParen)
		}
		catchEnd = p.EndPos
		catch = p.statement()
	}
	return &ast.TryCatch{Try: try, Catch: catch,
		CatchPos:       catchPos,
		CatchEnd:       catchEnd,
		CatchVar:       ast.Ident{Name: catchVar, Pos: varPos},
		CatchVarUnused: unused,
		CatchFilter:    catchFilter}
}
