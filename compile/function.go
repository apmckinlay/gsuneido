package compile

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
	. "github.com/apmckinlay/gsuneido/lexer"
	. "github.com/apmckinlay/gsuneido/runtime"
)

// ParseFunction parses a function and returns an AST for it
func ParseFunction(src string) *ast.Function {
	p := newParser(src)
	return p.function()
}

func (p *parser) function() *ast.Function {
	p.match(FUNCTION)
	params := p.params(false)
	body := p.compound()
	return &ast.Function{Params: params, Body: body}
}

func (p *parser) functionWithoutKeyword(inClass bool) *ast.Function {
	it := p.Item
	it.Token = IDENTIFIER
	it.Keyword = FUNCTION
	it.Text = "function"
	params := p.params(inClass)
	body := p.compound()
	return &ast.Function{Params: params, Body: body}
}

func (p *parser) params(inClass bool) []ast.Param {
	p.match(L_PAREN)
	var params []ast.Param
	if p.matchIf(AT) {
		params = append(params, ast.Param{Name: "@" + p.Text})
		p.match(IDENTIFIER)
	} else {
		defs := false
		for p.Token != R_PAREN {
			dot := p.matchIf(DOT)
			name := p.Text
			if dot {
				if !inClass {
					p.error("dot parameters only allowed in class methods")
				}
				name = "." + name
			}
			p.match(IDENTIFIER)
			p.checkForDupParam(params, name)
			if p.matchIf(EQ) {
				defs = true
				def := p.constant()
				params = append(params, ast.Param{Name: name, DefVal: def})
			} else {
				if defs {
					p.error("default parameters must come last")
				}
				params = append(params, ast.Param{Name: name})
			}
			p.matchIf(COMMA)
		}
	}
	p.match(R_PAREN)
	return params
}

func (p *parser) checkForDupParam(params []ast.Param, name string) {
	for _, a := range params {
		if a.Name == name {
			p.error("duplicate function parameter (" + name + ")")
		}
	}
}

func (p *parser) compound() []ast.Statement {
	p.match(L_CURLY)
	stmts := p.statements()
	p.match(R_CURLY)
	return stmts
}

func (p *parser) statements() []ast.Statement {
	list := []ast.Statement{}
	for p.Token != R_CURLY {
		if p.matchIf(SEMICOLON) {
			continue
		}
		stmt := p.statement()
		list = append(list, stmt)
	}
	return list
}

var code = Item{Token: L_CURLY, Text: "STMTS"}

func (p *parser) statement() ast.Statement {
	pos := p.Pos
	stmt := p.statement2()
	stmt.SetPos(int(pos))
	return stmt
}

func (p *parser) statement2() ast.Statement {
	tok := p.KeyTok()
	switch tok {
	case SEMICOLON:
		p.next()
		return &ast.Compound{Body: []ast.Statement{}}
	case L_CURLY:
		return &ast.Compound{Body: p.compound()}
	case RETURN:
		p.next()
		return p.semi(p.returnStmt())
	case IF:
		p.next()
		return p.ifStmt()
	case SWITCH:
		p.next()
		return p.switchStmt()
	case FOREVER:
		p.next()
		return p.foreverStmt()
	case WHILE:
		p.next()
		return p.whileStmt()
	case DO:
		p.next()
		return p.semi(p.dowhileStmt())
	case FOR:
		return p.forStmt()
	case THROW:
		p.next()
		return p.semi(&ast.Throw{E: p.expr()})
	case TRY:
		p.next()
		return p.tryStmt()
	case BREAK:
		p.next()
		return p.semi(&ast.Break{})
	case CONTINUE:
		p.next()
		return p.semi(&ast.Continue{})
	default:
		return p.semi(&ast.Expression{E: p.expr()})
	}
}

func (p *parser) semi(stmt ast.Statement) ast.Statement {
	p.matchIf(SEMICOLON)
	return stmt
}

func (p *parser) ifStmt() *ast.If {
	expr := p.ctrlExpr()
	t := p.statement()
	var f ast.Statement
	if p.matchIf(ELSE) {
		f = p.statement()
	}
	return &ast.If{Cond: expr, Then: t, Else: f}
}

func (p *parser) switchStmt() *ast.Switch {
	var expr ast.Expr
	if p.Token == L_CURLY {
		expr = p.Constant(True)
	} else {
		expr = p.exprExpecting(true)
	}
	p.match(L_CURLY)
	var cases []ast.Case
	for p.matchIf(CASE) {
		cases = append(cases, p.switchCase())
	}
	var def []ast.Statement
	if p.matchIf(DEFAULT) {
		def = p.switchBody()
	}
	p.match(R_CURLY)
	return &ast.Switch{E: expr, Cases: cases, Default: def}
}

func (p *parser) switchCase() ast.Case {
	var exprs []ast.Expr
	for {
		exprs = append(exprs, p.expr())
		if !p.matchIf(COMMA) {
			break
		}
	}
	body := p.switchBody()
	return ast.Case{Exprs: exprs, Body: body}
}

func (p *parser) switchBody() []ast.Statement {
	p.match(COLON)
	var stmts []ast.Statement
	for p.Token != R_CURLY && p.Keyword != CASE && p.Keyword != DEFAULT {
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
	p.match(WHILE)
	cond := p.expr()
	return &ast.DoWhile{Body: body, Cond: cond}
}

func (p *parser) forStmt() ast.Statement {
	// easier to check before matching FOR so everything is ahead
	forIn := p.isForIn()
	p.match(FOR)
	if forIn {
		return p.forIn()
	}
	return p.forClassic()
}

func (p *parser) isForIn() bool {
	i := 0
	if p.lxr.AheadSkip(i).Token == L_PAREN {
		i++
	}
	if p.lxr.AheadSkip(i).Token != IDENTIFIER {
		return false
	}
	return p.lxr.AheadSkip(i+1).Keyword == IN
}

func (p *parser) forIn() *ast.ForIn {
	parens := p.matchIf(L_PAREN)
	id := p.Text
	p.match(IDENTIFIER)
	p.match(IN)
	expr := p.exprExpecting(!parens)
	if parens {
		p.match(R_PAREN)
	}
	body := p.statement()
	return &ast.ForIn{Var: id, E: expr, Body: body}
}

func (p *parser) forClassic() *ast.For {
	p.match(L_PAREN)
	init := p.optExprList(SEMICOLON)
	p.match(SEMICOLON)
	var cond ast.Expr
	if p.Token != SEMICOLON {
		cond = p.expr()
	}
	p.match(SEMICOLON)
	inc := p.optExprList(R_PAREN)
	p.match(R_PAREN)
	body := p.statement()
	return &ast.For{Init: init, Cond: cond, Inc: inc, Body: body}
}

func (p *parser) optExprList(after Token) []ast.Expr {
	exprs := []ast.Expr{}
	if p.Token != after {
		for {
			exprs = append(exprs, p.expr())
			if p.Token != COMMA {
				break
			}
			p.next()
		}
	}
	return exprs
}

// used by if, while, and do-while
func (p *parser) ctrlExpr() ast.Expr {
	parens := p.matchIf(L_PAREN)
	expr := p.exprExpecting(!parens)
	if parens {
		p.match(R_PAREN)
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
	if p.newline || p.Token == SEMICOLON || p.Token == R_CURLY {
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
	var catchFilter string
	var catch ast.Statement
	if p.matchIf(CATCH) {
		if p.matchIf(L_PAREN) {
			catchVar = p.Text
			p.match(IDENTIFIER)
			if p.matchIf(COMMA) {
				catchFilter = p.Text
				p.match(STRING)
			}
			p.match(R_PAREN)
		}
		catch = p.statement()
	}
	return &ast.TryCatch{Try: try,
		CatchVar: catchVar, CatchFilter: catchFilter, Catch: catch}
}
