// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopSymbolStmtIdentifier
	{
	// FORSTMT [FOR, LPAREN, LIST{EXPR_ELEM}, SEMICOLON, expr, SEMICOLON, LIST{EXPR_ELEM}, RPAREN, stmts]
	// FORINSTMT [FOR, LPAREN, IDENTIFIER, IN, expr|RANGE, RPAREN, stmts]
	// RANGE [expr, RANGETO|RANGELEN, expr]
	Std()
		{
		children = Object()
		TdopAddChild(children, token: .Tokenize(TDOPTOKEN.FOR))

		condInParen = _token().Match(TDOPTOKEN.LPAREN)
		TdopAddChild(children, match: TDOPTOKEN.LPAREN)

		if _token().Match(TDOPTOKEN.IDENTIFIER) and _ahead().Match(TDOPTOKEN.IN)
			{
			stmtType = TDOPTOKEN.FORINSTMT
			.handleForIn(condInParen, children)
			}
		else if not condInParen
			{
			stmtType = TDOPTOKEN.FORINSTMT
			.handleForRange(children)
			}
		else
			{
			stmtType = TDOPTOKEN.FORSTMT
			.handleFor(condInParen, children)
			}
		TdopAddChild(children, match: TDOPTOKEN.RPAREN, mustMatch: condInParen)
		TdopAddChild(children, token: TdopStmt())
		return TdopCreateNode(stmtType, :children)
		}

	GetExpressionList(end?)
		{
		return TdopCreateList()
			{ |list|
			while _token() isnt _end and not end?(_token())
				{
				exprChildren = Object()
				TdopAddChild(exprChildren, token: _expr())
				TdopAddChild(exprChildren, match: TDOPTOKEN.COMMA)
				list.Add(TdopCreateNode(TDOPTOKEN.EXPR_ELEM, children: exprChildren))
				}
			}
		}

	// NOTE: partial duplicate
	GetStatementList(end?)
		{
		return TdopCreateList()
			{ |list|
			while _token() isnt _end and not end?(_token())
				list.Add(_stmt())
			}
		}

	handleForIn(condInParen, children)
		{
		TdopAddChild(children, match: TDOPTOKEN.IDENTIFIER, mustMatch:)
		TdopAddChild(children, match: TDOPTOKEN.IN, mustMatch:)

		expr = _token().Match(TDOPTOKEN.RANGETO)
			? TdopSymbolNumber(0)
			: condInParen
				? _expr()
				: TdopStmtExpr(expectingCompound:)

		if _token().Match(TDOPTOKEN.RANGETO)
			{
			first = expr
			rangeToken = _token()
			_advance()
			second = condInParen ? _expr() : TdopStmtExpr(expectingCompound:)
			expr = TdopCreateNode(TDOPTOKEN.RANGE,
				children: Object(first, rangeToken, second))
			}
		TdopAddChild(children, token: expr)
		}

	handleForRange(children)
		{
		TdopAddChild(children, token: TdopCreateNode(TDOPTOKEN.IDENTIFIER, value: ''))
		TdopAddChild(children, match: TDOPTOKEN.IN)
		rangeChildren = Object(TdopSymbolNumber(0))
		TdopAddChild(rangeChildren, match: TDOPTOKEN.RANGETO, mustMatch:)
		TdopAddChild(rangeChildren, token: TdopStmtExpr(expectingCompound:))
		range = TdopCreateNode(TDOPTOKEN.RANGE, children: rangeChildren)
		TdopAddChild(children, token: range)
		}

	handleFor(condInParen, children)
		{
		if condInParen is false
			throw 'parenthesis required: for(expr; expr; expr)'

		TdopAddChild(children, token: .GetExpressionList(
			{ it.Token is TDOPTOKEN.SEMICOLON }))
		TdopAddChild(children, match: TDOPTOKEN.SEMICOLON, mustMatch:)

		if not _token().Match(TDOPTOKEN.SEMICOLON)
			TdopAddChild(children, token: _expr())
		else
			TdopAddChild(children, token: TdopCreateNode(TDOPTOKEN.TRUE))
		TdopAddChild(children, match: TDOPTOKEN.SEMICOLON, mustMatch:)

		TdopAddChild(children, token: .GetExpressionList({ it.Token is TDOPTOKEN.RPAREN}))
		}
	}
