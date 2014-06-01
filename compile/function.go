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
	return ast(Item{Token: STATEMENTS}, list...)
}

func (p *parser) statement() Ast {
	if p.Token == L_CURLY {
		return p.compound()
	}
	// TODO other statement types
	switch p.KeyTok() {
	case RETURN:
		return p.returnStmt()
	case IF:
		return p.ifStmt()
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

func (p *parser) whileStmt() Ast {
	it, expr := p.ctrlExpr()
	body := p.statement()
	return ast(it, expr, body)
}

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
	switch p.Token {
	case NEWLINE, SEMICOLON, R_CURLY:
		return ast(item)
	}
	return ast(item, p.exprAst())
}

func (p *parser) exprStmt() Ast {
	return ast(Item{Token: EXPRESSION}, p.exprAst())
}

func (p *parser) exprAst() Ast {
	return expression(p, astBuilder).(Ast)
}
