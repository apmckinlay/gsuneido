// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
	. "github.com/apmckinlay/gsuneido/compile/lexer"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/runtime"
)

// function parse a function (starting with the "function" keyword)
func (p *Parser) Function() *ast.Function {
	p.Match(tok.Function)
	return p.function(false)
}

// function parses a function or method (without the "function" keyword)
func (p *Parser) function(inClass bool) *ast.Function {
	funcInfoSave := p.funcInfo
	p.InitFuncInfo()
	pos := p.Pos
	params := p.params(inClass)
	body := p.compound()
	p.processFinal()
	fn := &ast.Function{Pos: pos, Params: params, Body: body, Final: p.final,
		HasBlocks: p.hasBlocks}
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
	if p.MatchIf(tok.At) {
		params = append(params,
			mkParam("@"+p.Text, p.Pos, p.unusedAhead(), nil))
		p.final[p.Text] = disqualified
		p.MatchIdent()
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
				was_string := p.Token == tok.String || p.Token == tok.Symbol
				defs = true
				def := p.constant()
				if _, ok := def.(SuStr); ok && !was_string {
					p.Error("parameter defaults must be constants")
				}
				params = append(params, mkParam(name, pos, unused, def))
			} else {
				if defs {
					p.Error("default parameters must come last")
				}
				params = append(params, mkParam(name, pos, unused, nil))
			}
			p.MatchIf(tok.Comma)
		}
	}
	p.Match(tok.RParen)
	return params
}

func mkParam(name string, pos int32, unused bool, def Value) ast.Param {
	if name == "unused" || name == "@unused" {
		unused = true
	}
	return ast.Param{Name: ast.Ident{Name: name, Pos: pos},
		DefVal: def, Unused: unused}
}

// unDyn removes the leading underscore from dynamic parameters
func unDyn(id string) string {
	if id[0] == '_' {
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

var code = Item{Token: tok.LCurly, Text: "STMTS"}

func (p *Parser) statement() ast.Statement {
	pos := p.Pos
	stmt := p.statement2()
	stmt.SetPos(int(pos))
	return stmt
}

func (p *Parser) statement2() ast.Statement {
	token := p.Token
	switch token {
	case tok.Semicolon:
		for p.MatchIf(tok.Semicolon) {
		}
		return &ast.Compound{Body: []ast.Statement{}}
	case tok.LCurly:
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
		return &ast.ExprStmt{E: p.trailingExpr()}
	}
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
		stmt.Else = p.statement()
	}
	return stmt
}

func (p *Parser) switchStmt() *ast.Switch {
	var expr ast.Expr
	if p.Token == tok.LCurly {
		expr = p.Constant(True)
	} else {
		expr = p.exprExpecting(true)
	}
	p.Match(tok.LCurly)
	var cases []ast.Case
	for p.MatchIf(tok.Case) {
		cases = append(cases, p.switchCase())
	}
	var def []ast.Statement
	if p.MatchIf(tok.Default) {
		def = p.switchBody()
	}
	p.Match(tok.RCurly)
	return &ast.Switch{E: expr, Cases: cases, Default: def}
}

func (p *Parser) switchCase() ast.Case {
	var exprs []ast.Expr
	for {
		exprs = append(exprs, p.Expression())
		if !p.MatchIf(tok.Comma) {
			break
		}
	}
	body := p.switchBody()
	return ast.Case{Exprs: exprs, Body: body}
}

func (p *Parser) switchBody() []ast.Statement {
	p.Match(tok.Colon)
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
	cond := p.Expression()
	return &ast.DoWhile{Body: body, Cond: cond}
}

func (p *Parser) forStmt() ast.Statement {
	// easier to check before matching For so everything is ahead
	forIn := p.isForIn()
	p.Match(tok.For)
	if forIn {
		return p.forIn()
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
	return p.Lxr.AheadSkip(i+1).Token == tok.In
}

func (p *Parser) forIn() *ast.ForIn {
	parens := p.MatchIf(tok.LParen)
	id := p.Text
	p.final[id] = disqualified
	pos := p.Pos
	p.MatchIdent()
	p.Match(tok.In)
	expr := p.exprExpecting(!parens)
	if parens {
		p.Match(tok.RParen)
	}
	body := p.statement()
	return &ast.ForIn{Var: ast.Ident{Name: id, Pos: pos}, E: expr, Body: body}
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

// used by if, while, and do-while
func (p *Parser) ctrlExpr() ast.Expr {
	parens := p.MatchIf(tok.LParen)
	keepParens := p.Token == tok.LParen
	expr := p.exprExpecting(!parens)
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
	return &ast.Return{E: p.trailingExpr()}
}

func (p *Parser) throwStmt() *ast.Throw {
	return &ast.Throw{E: p.Expression()}
}

func (p *Parser) tryStmt() *ast.TryCatch {
	try := p.statement()
	var catchVar string
	var varPos int32
	var catchFilter string
	var catch ast.Statement
	var unused bool
	catchPos := int(p.Pos)
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
		catch = p.statement()
	}
	return &ast.TryCatch{Try: try, Catch: catch,
		CatchPos:       catchPos,
		CatchVar:       ast.Ident{Name: catchVar, Pos: varPos},
		CatchVarUnused: unused,
		CatchFilter:    catchFilter}
}
