package compile

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
	. "github.com/apmckinlay/gsuneido/lexer"
	tok "github.com/apmckinlay/gsuneido/lexer/tokens"
	. "github.com/apmckinlay/gsuneido/runtime"
)

// function parse a function (starting with the "function" keyword)
func (p *parser) Function() *ast.Function {
	pos := p.Pos
	p.match(tok.Function)
	params := p.params(false)
	body := p.compound()
	return &ast.Function{Pos: pos, Params: params, Body: body}
}

// method parse a class method (without the "function" keyword)
func (p *parser) method() *ast.Function {
	pos := p.Pos
	params := p.params(true)
	body := p.compound()
	return &ast.Function{Pos: pos, Params: params, Body: body}
}

func (p *parser) params(inClass bool) []ast.Param {
	p.match(tok.LParen)
	var params []ast.Param
	if p.matchIf(tok.At) {
		params = append(params,
			ast.Param{Name: "@" + p.Text, Unused: p.unusedAhead(), Pos: p.Pos})
		p.matchIdent()
	} else {
		defs := false
		for p.Token != tok.RParen {
			dot := p.matchIf(tok.Dot)
			name := p.Text
			pos := p.Pos
			if dot {
				if !inClass {
					p.error("dot parameters only allowed in class methods")
				}
				name = "." + name
			}
			unused := p.unusedAhead()
			p.matchIdent()
			p.checkForDupParam(params, name)
			if p.matchIf(tok.Eq) {
				was_string := p.Token == tok.String
				defs = true
				def := p.constant()
				if _, ok := def.(SuStr); ok && !was_string {
					p.error("parameter defaults must be constants")
				}
				params = append(params,
					ast.Param{Name: name, DefVal: def, Unused: unused, Pos: pos})
			} else {
				if defs {
					p.error("default parameters must come last")
				}
				params = append(params,
					ast.Param{Name: name, Unused: unused, Pos: pos})
			}
			p.matchIf(tok.Comma)
		}
	}
	p.match(tok.RParen)
	return params
}

func (p *parser) unusedAhead() bool {
	i := 0
	for ; p.lxr.Ahead(i).Token == tok.Whitespace; i++ {
	}
	return p.lxr.Ahead(i).Text == "/*unused*/"
}

func (p *parser) checkForDupParam(params []ast.Param, name string) {
	for _, a := range params {
		if a.Name == name {
			p.error("duplicate function parameter (" + name + ")")
		}
	}
}

func (p *parser) compound() []ast.Statement {
	p.match(tok.LCurly)
	stmts := p.statements()
	p.match(tok.RCurly)
	return stmts
}

func (p *parser) statements() []ast.Statement {
	list := []ast.Statement{}
	for p.Token != tok.RCurly {
		stmt := p.statement()
		list = append(list, stmt)
	}
	return list
}

var code = Item{Token: tok.LCurly, Text: "STMTS"}

func (p *parser) statement() ast.Statement {
	pos := p.Pos
	stmt := p.statement2()
	stmt.SetPos(int(pos))
	return stmt
}

func (p *parser) statement2() ast.Statement {
	token := p.Token
	switch token {
	case tok.Semicolon:
		p.next()
		return &ast.Compound{Body: []ast.Statement{}}
	case tok.LCurly:
		return &ast.Compound{Body: p.compound()}
	case tok.Return:
		p.next()
		return p.semi(p.returnStmt())
	case tok.If:
		p.next()
		return p.ifStmt()
	case tok.Switch:
		p.next()
		return p.switchStmt()
	case tok.Forever:
		p.next()
		return p.foreverStmt()
	case tok.While:
		p.next()
		return p.whileStmt()
	case tok.Do:
		p.next()
		return p.semi(p.dowhileStmt())
	case tok.For:
		return p.forStmt()
	case tok.Throw:
		p.next()
		return p.semi(&ast.Throw{E: p.expr()})
	case tok.Try:
		p.next()
		return p.tryStmt()
	case tok.Break:
		p.next()
		return p.semi(&ast.Break{})
	case tok.Continue:
		p.next()
		return p.semi(&ast.Continue{})
	default:
		return p.semi(&ast.ExprStmt{E: p.expr()})
	}
}

func (p *parser) semi(stmt ast.Statement) ast.Statement {
	p.matchIf(tok.Semicolon)
	return stmt
}

func (p *parser) ifStmt() *ast.If {
	expr := p.ctrlExpr()
	t := p.statement()
	var f ast.Statement
	if p.matchIf(tok.Else) {
		f = p.statement()
	}
	return &ast.If{Cond: expr, Then: t, Else: f}
}

func (p *parser) switchStmt() *ast.Switch {
	var expr ast.Expr
	if p.Token == tok.LCurly {
		expr = p.Constant(True)
	} else {
		expr = p.exprExpecting(true)
	}
	p.match(tok.LCurly)
	var cases []ast.Case
	for p.matchIf(tok.Case) {
		cases = append(cases, p.switchCase())
	}
	var def []ast.Statement
	if p.matchIf(tok.Default) {
		def = p.switchBody()
	}
	p.match(tok.RCurly)
	return &ast.Switch{E: expr, Cases: cases, Default: def}
}

func (p *parser) switchCase() ast.Case {
	var exprs []ast.Expr
	for {
		exprs = append(exprs, p.expr())
		if !p.matchIf(tok.Comma) {
			break
		}
	}
	body := p.switchBody()
	return ast.Case{Exprs: exprs, Body: body}
}

func (p *parser) switchBody() []ast.Statement {
	p.match(tok.Colon)
	stmts := []ast.Statement{}
	for p.Token != tok.RCurly && p.Token != tok.Case && p.Token != tok.Default {
		stmts = append(stmts, p.statement())
	}
	return stmts
}

func (p *parser) foreverStmt() *ast.Forever {
	body := p.statement()
	return &ast.Forever{Body: body}
}

func (p *parser) whileStmt() *ast.While {
	cond := p.ctrlExpr()
	body := p.statement()
	return &ast.While{Cond: cond, Body: body}
}

func (p *parser) dowhileStmt() *ast.DoWhile {
	body := p.statement()
	p.match(tok.While)
	cond := p.expr()
	return &ast.DoWhile{Body: body, Cond: cond}
}

func (p *parser) forStmt() ast.Statement {
	// easier to check before matching For so everything is ahead
	forIn := p.isForIn()
	p.match(tok.For)
	if forIn {
		return p.forIn()
	}
	return p.forClassic()
}

func (p *parser) isForIn() bool {
	i := 0
	if p.lxr.AheadSkip(i).Token == tok.LParen {
		i++
	}
	if !p.lxr.AheadSkip(i).Token.IsIdent() {
		return false
	}
	return p.lxr.AheadSkip(i+1).Token == tok.In
}

func (p *parser) forIn() *ast.ForIn {
	parens := p.matchIf(tok.LParen)
	id := p.Text
	pos := p.Pos
	p.matchIdent()
	p.match(tok.In)
	expr := p.exprExpecting(!parens)
	if parens {
		p.match(tok.RParen)
	}
	body := p.statement()
	return &ast.ForIn{Var: id, VarPos: pos, E: expr, Body: body}
}

func (p *parser) forClassic() *ast.For {
	p.match(tok.LParen)
	init := p.optExprList(tok.Semicolon)
	p.match(tok.Semicolon)
	var cond ast.Expr
	if p.Token != tok.Semicolon {
		cond = p.expr()
	}
	p.match(tok.Semicolon)
	inc := p.optExprList(tok.RParen)
	p.match(tok.RParen)
	body := p.statement()
	return &ast.For{Init: init, Cond: cond, Inc: inc, Body: body}
}

func (p *parser) optExprList(after tok.Token) []ast.Expr {
	exprs := []ast.Expr{}
	if p.Token != after {
		for {
			exprs = append(exprs, p.expr())
			if p.Token != tok.Comma {
				break
			}
			p.next()
		}
	}
	return exprs
}

// used by if, while, and do-while
func (p *parser) ctrlExpr() ast.Expr {
	parens := p.matchIf(tok.LParen)
	expr := p.exprExpecting(!parens)
	if parens {
		p.match(tok.RParen)
	}
	return expr
}

func (p *parser) exprExpecting(expecting bool) ast.Expr {
	p.expectingCompound = expecting
	expr := p.expr()
	p.expectingCompound = false
	return expr
}

func (p *parser) returnStmt() *ast.Return {
	if p.newline || p.Token == tok.Semicolon || p.Token == tok.RCurly {
		return &ast.Return{}
	}
	expr := p.expr()
	return &ast.Return{E: expr}
}

func (p *parser) throwStmt() *ast.Throw {
	return &ast.Throw{E: p.expr()}
}

func (p *parser) tryStmt() *ast.TryCatch {
	try := p.statement()
	var catchVar string
	var varPos int32
	var catchFilter string
	var catch ast.Statement
	if p.matchIf(tok.Catch) {
		if p.matchIf(tok.LParen) {
			catchVar = p.Text
			varPos = p.Pos
			p.matchIdent()
			if p.matchIf(tok.Comma) {
				catchFilter = p.Text
				p.match(tok.String)
			}
			p.match(tok.RParen)
		}
		catch = p.statement()
	}
	return &ast.TryCatch{Try: try, Catch: catch,
		CatchVar: catchVar, VarPos: varPos, CatchFilter: catchFilter}
}
