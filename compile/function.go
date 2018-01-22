package compile

import . "github.com/apmckinlay/gsuneido/lexer"

// ParseFunction parses a function and returns an AST for it
func ParseFunction(src string) Ast {
	p := newParser(src)
	return p.function()
}

func (p *parser) function() Ast {
	it := p.Item
	p.match(FUNCTION)
	params := p.params()
	body := p.compound()
	return ast(it, params, body)
}

func (p *parser) params() Ast {
	p.match(L_PAREN)
	var params []Ast
	if p.matchIf(AT) {
		//TODO @+1
		params = append(params, ast2("@"+p.Text))
		p.match(IDENTIFIER)
	} else {
		defs := false
		for p.Token != R_PAREN {
			dot := p.matchIf(DOT)
			name := p.Text
			if dot {
				name = "." + name
			}
			p.match(IDENTIFIER)
			if p.matchIf(EQ) {
				defs = true
				def := p.constant()
				params = append(params, astVal(name, def))
			} else {
				if defs {
					p.error("default parameters must come last")
				}
				params = append(params, ast2(name))
			}
			p.matchIf(COMMA)
		}
	}
	p.matchSkipNL(R_PAREN)
	return ast2("params", params...)
}

func (p *parser) compound() Ast {
	p.match(L_CURLY)
	stmts := p.statements()
	p.match(R_CURLY)
	return stmts
}

func (p *parser) statements() Ast {
	list := []Ast{}
	for p.Token != R_CURLY {
		if p.matchIf(NEWLINE) || p.matchIf(SEMICOLON) {
			continue
		}
		stmt := p.statement()
		list = append(list, stmt)
	}
	return ast(code, list...)
}

var code = Item{Token: L_CURLY, Text: "STMTS"}

func (p *parser) statement() Ast {
	if p.Token == NEWLINE {
		p.nextSkipNL()
	}
	if p.Token == L_CURLY {
		return p.compound()
	}
	if p.matchIf(SEMICOLON) {
		return ast(code)
	}
	switch p.Keyword {
	case RETURN:
		return p.returnStmt()
	case IF:
		return p.ifStmt()
	case SWITCH:
		return p.switchStmt()
	case FOREVER:
		return p.foreverStmt()
	case WHILE:
		return p.whileStmt()
	case DO:
		return p.dowhileStmt()
	case FOR:
		return p.forStmt()
	case THROW:
		return p.throwStmt()
	case TRY:
		return p.tryStmt()
	case BREAK, CONTINUE:
		it := p.Item
		p.next()
		p.matchIf(SEMICOLON)
		return ast(it)
	default:
		return p.exprStmt()
	}
}

func (p *parser) ifStmt() Ast {
	it, expr := p.ctrlExpr(IF)
	t := p.statement()
	if p.Keyword == ELSE {
		p.nextSkipNL()
		f := p.statement()
		return ast(it, expr, t, f)
	}
	return ast(it, expr, t)
}

func (p *parser) switchStmt() Ast {
	it := p.Item
	p.nextSkipNL()
	var expr Ast
	if p.Token == L_CURLY {
		expr = ast(Item{Token: TRUE})
	} else {
		expr = p.exprExpecting(true)
		if p.Token == NEWLINE {
			p.nextSkipNL()
		}
	}
	p.nextSkipNL()
	var cases []Ast
	for p.matchIf(CASE) {
		cases = append(cases, p.switchCase())
	}
	result := ast(it, expr, ast2("cases", cases...))
	if p.matchIf(DEFAULT) {
		result.Children = append(result.Children, p.switchBody())
	}
	p.match(R_CURLY)
	return result
}

func (p *parser) switchCase() Ast {
	var values []Ast
	for {
		values = append(values, p.exprAst())
		if !p.matchIf(COMMA) {
			break
		}
	}
	body := p.switchBody()
	return ast(Item{Token: CASE}, ast2("vals", values...), body)
}

func (p *parser) switchBody() Ast {
	p.match(COLON)
	var stmts []Ast
	for p.Token != R_CURLY && p.Keyword != CASE && p.Keyword != DEFAULT {
		stmts = append(stmts, p.statement())
	}
	return ast(code, stmts...)
}

func (p *parser) foreverStmt() Ast {
	it := p.Item
	p.match(FOREVER)
	body := p.statement()
	return ast(it, body)
}

func (p *parser) whileStmt() Ast {
	it, expr := p.ctrlExpr(WHILE)
	body := p.statement()
	return ast(it, expr, body)
}

func (p *parser) dowhileStmt() Ast {
	it := p.Item
	p.match(DO)
	body := p.statement()
	_, expr := p.ctrlExpr(WHILE)
	return ast(it, body, expr)
}

func (p *parser) forStmt() Ast {
	it := p.Item
	forIn := p.isForIn()
	p.match(FOR)
	if forIn {
		return p.forIn(it)
	}
	return p.forClassic(it)
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

func (p *parser) forIn(it Item) Ast {
	it.Text = "for-in"
	parens := p.matchIf(L_PAREN)
	id := p.Text
	p.match(IDENTIFIER)
	p.matchSkipNL(IN)
	if !parens {
		defer func(prev int) { p.nest = prev }(p.nest)
		p.nest = 0
	}
	expr := p.exprExpecting(!parens)
	if parens {
		p.match(R_PAREN)
	} else {
		p.matchIf(NEWLINE)
	}
	body := p.statement()
	return ast(it, ast2(id), expr, body)
}

func (p *parser) forClassic(it Item) Ast {
	p.match(L_PAREN)
	init := p.optExprList(SEMICOLON)
	p.match(SEMICOLON)
	cond := p.exprAst()
	p.match(SEMICOLON)
	incr := p.optExprList(R_PAREN)
	p.match(R_PAREN)
	body := p.statement()
	return ast(it, init, cond, incr, body)
}

func (p *parser) optExprList(after Token) Ast {
	ast := ast2("exprs")
	if p.Token != after {
		for {
			ast.Children = append(ast.Children, p.exprAst())
			if p.Token != COMMA {
				break
			}
			p.next()
		}
	}
	return ast
}

// used by if, while, and do-while
func (p *parser) ctrlExpr(tok Token) (Item, Ast) {
	it := p.Item
	p.matchSkipNL(tok)
	parens := p.matchIf(L_PAREN)
	expr := p.exprExpecting(!parens)
	if parens {
		p.match(R_PAREN)
	} else {
		p.matchIf(NEWLINE)
	}
	return it, expr
}

func (p *parser) exprExpecting(expecting bool) Ast {
	p.expectingCompound = expecting
	expr := p.exprAst()
	p.expectingCompound = false
	return expr
}

func (p *parser) returnStmt() Ast {
	item := p.Item
	p.matchKeepNL(RETURN)
	if p.matchIf(NEWLINE) || p.matchIf(SEMICOLON) || p.Token == R_CURLY {
		return ast(item)
	}
	return ast(item, p.exprStmt())
}

func (p *parser) exprStmt() Ast {
	result := p.exprAst()
	for p.Token == SEMICOLON || p.Token == NEWLINE {
		p.next()
	}
	return result
}

func (p *parser) throwStmt() Ast {
	item := p.Item
	p.matchSkipNL(THROW)
	return ast(item, p.exprStmt())
}

func (p *parser) tryStmt() Ast {
	item := p.Item
	p.matchSkipNL(TRY)
	try := p.statement()
	if p.Keyword != CATCH {
		return ast(item, try)
	}
	catch := p.catch()
	return ast(item, try, catch)
}

func (p *parser) catch() Ast {
	item := p.Item
	p.matchSkipNL(CATCH)
	var children []Ast
	if p.matchIf(L_PAREN) {
		children = append(children, ast(p.Item))
		p.match(IDENTIFIER)
		if p.matchIf(COMMA) {
			children = append(children, ast(p.Item))
			p.match(STRING)
		}
		p.match(R_PAREN)
	}
	children = append(children, p.statement())
	return ast(item, children...)
}

func (p *parser) exprAst() Ast {
	return expression(p, astBuilder).(Ast)
}
