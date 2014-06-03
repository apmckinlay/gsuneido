package compile

func ParseFunction(src string) Ast {
	p := newParser(src)
	return p.function()
}

func (p *parser) function() Ast {
	it := p.Item
	p.match(FUNCTION)
	p.match(L_PAREN)
	p.match(R_PAREN)
	body := p.compound()
	return ast(it, ast(Item{Token: PARAMS}), body)
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
		//fmt.Println("stmt:", stmt)
		list = append(list, stmt)
	}
	return ast(code, list...)
}

var code = Item{Token: STATEMENTS}

func (p *parser) statement() Ast {
	if p.Token == NEWLINE {
		p.nextSkipNL()
	}
	if p.Token == L_CURLY {
		return p.compound()
	}
	if p.matchIf(SEMICOLON) {
		return ast(Item{})
	}
	// TODO other statement types
	switch p.Keyword {
	case RETURN:
		return p.returnStmt()
	case IF:
		return p.ifStmt()
	case SWITCH:
		return p.switchStmt()
	case WHILE:
		return p.whileStmt()
	default:
		return p.exprStmt()
	}
}

func (p *parser) ifStmt() Ast {
	it, expr := p.ctrlExpr()
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
		expr = p.exprAst()
		if p.Token == NEWLINE {
			p.nextSkipNL()
		}
	}
	p.nextSkipNL()
	var cases []Ast
	for p.matchIf(CASE) {
		cases = p.switchCase(cases)
	}
	result := ast(it, expr, ast2("cases", cases...))
	if p.matchIf(DEFAULT) {
		result.Children = append(result.Children, p.switchBody())
	}
	p.match(R_CURLY)
	return result
}

func (p *parser) switchCase(cases []Ast) []Ast {
	var values []Ast
	for {
		values = append(values, p.exprAst())
		if !p.matchIf(COMMA) {
			break
		}
	}
	body := p.switchBody()
	c := ast(Item{Token: CASE}, ast2("vals", values...), body)
	return append(cases, c)
}

func (p *parser) switchBody() Ast {
	p.match(COLON)
	var stmts []Ast
	for p.Token != R_CURLY && p.Keyword != CASE && p.Keyword != DEFAULT {
		stmts = append(stmts, p.statement())
	}
	return ast(code, stmts...)
}

func (p *parser) whileStmt() Ast {
	it, expr := p.ctrlExpr()
	body := p.statement()
	return ast(it, expr, body)
}

// used by if and while
func (p *parser) ctrlExpr() (Item, Ast) {
	it := p.Item
	p.nextSkipNL()
	expr := p.exprAst()
	if p.Token == NEWLINE {
		p.nextSkipNL()
	}
	return it, expr
}

func (p *parser) returnStmt() Ast {
	item := p.Item
	p.matchKeepNL(RETURN)
	if p.matchIf(NEWLINE) || p.matchIf(SEMICOLON) || p.Token == R_CURLY {
		return ast(item)
	}
	result := ast(item, p.exprAst())
	for p.Token == SEMICOLON || p.Token == NEWLINE {
		p.next()
	}
	return result
}

func (p *parser) exprStmt() Ast {
	result := ast(Item{Token: EXPRESSION}, p.exprAst())
	for p.Token == SEMICOLON || p.Token == NEWLINE {
		p.next()
	}
	return result
}

func (p *parser) exprAst() Ast {
	return expression(p, astBuilder).(Ast)
}
