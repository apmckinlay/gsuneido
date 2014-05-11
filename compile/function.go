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
	defer func(prev int) { p.nest = prev }(p.nest)
	p.nest = 0
	return expression(p, astBuilder).(Ast)
}
