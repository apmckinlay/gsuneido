// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

import (
	"fmt"

	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/runtime"
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
	pos1 := p.endPos
	p.Match(tok.LCurly)
	pos2 := p.endPos
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
		param := mkParam(name, pos, p.endPos, unused, def)
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
		SetPos(result, org, p.endPos)
	}(p.Pos)
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
		stmt.ElseEnd = p.endPos
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
	pos1 := p.endPos
	p.Match(tok.LCurly)
	pos2 := p.endPos
	var cases []ast.Case
	for p.Token == tok.Case {
		cases = append(cases, p.switchCase())
	}
	var def []ast.Statement
	var posdef int32
	if p.MatchIf(tok.Default) {
		p.Match(tok.Colon)
		posdef = p.endPos
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
		exprs = append(exprs, p.Expression())
		if !p.MatchIf(tok.Comma) {
			break
		}
	}
	p.Match(tok.Colon)
	end := p.endPos
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
	cond := p.Expression()
	return &ast.DoWhile{Body: body, Cond: cond}
}

func (p *Parser) forStmt() ast.Statement {
	// easier to check before matching For so everything is ahead
	forIn := p.isForIn()
	forSlice := p.isForSlice()
	p.Match(tok.For)
	if forSlice {
		return p.forSlice()
	} else if forIn {
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

func (p *Parser) isForSlice() bool {
	i := 0
	if p.Lxr.AheadSkip(i).Token == tok.LParen {
		i++
	}
	if !p.Lxr.AheadSkip(i).Token.IsIdent() {
		return false
	}
	if p.Lxr.AheadSkip(i+1).Token == tok.In {	
		return  p.Lxr.AheadSkip(i+2).Token == tok.Number
	}
	return false
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

func (p *Parser) forSlice() *ast.ForSlice {
	// match this for loop structure:
	// for i in 0..=10 { ... }
	// for i in 0..<10 { ... }
fmt.Println(p.String())
	p.MatchIf(tok.LParen) 	// consume "(" if present
	p.MatchIdent() 			// consume whitespace
	p.Match(tok.In) 		// consume "in"
fmt.Println(p.String())
	p.Match(tok.Number) 	// consume "0" (lower bound)
fmt.Println(p.String())
	p.Match(tok.RangeTo) 	// consume ".."
fmt.Println(p.String())
	p.MatchIf(tok.Eq) 		// consume "=" if present
fmt.Println(p.String())
	p.MatchIf(tok.Lt) 		// consume "<" if present
fmt.Println(p.String())
	p.Match(tok.Number) 	// consume "10" (upper bound)
fmt.Println(p.String())
	body := p.statement() 	// consume within "{ ... }"

	astrepr := &ast.ForSlice{ Body: body}
fmt.Println(astrepr.String())
	return astrepr
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
	astrepr := &ast.For{Init: init, Cond: cond, Inc: inc, Body: body}
fmt.Println(astrepr.String())
	return astrepr
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
	returnThrow := false
	if p.MatchIf(tok.Throw) {
		returnThrow = true
	}
	return &ast.Return{E: p.trailingExpr(), ReturnThrow: returnThrow}
}

func (p *Parser) tryStmt() *ast.TryCatch {
	try := p.statement()
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
		catchEnd = p.endPos
		catch = p.statement()
	}
	return &ast.TryCatch{Try: try, Catch: catch,
		CatchPos:       catchPos,
		CatchEnd:       catchEnd,
		CatchVar:       ast.Ident{Name: catchVar, Pos: varPos},
		CatchVarUnused: unused,
		CatchFilter:    catchFilter}
}
