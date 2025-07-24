// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopSymbolStmtIdentifier
	{
	// SWITCHSTMT [SWITCH, LPAREN, expr, RPAREN, LCURLY, LIST{CASE_ELEM}, RCURLY]
	// CASE_ELEM [CASE|DEFAULT, LIST{EXPR_ELEM}, COLON, LIST{stmt}]
	// EXPR_ELEM [expr, COMMA]
	Std()
		{
		children = Object()
		TdopAddChild(children, token: .Tokenize(TDOPTOKEN.SWITCH))
		if _token().Match(TDOPTOKEN.LCURLY)
			TdopFillChildren(children,
				tokens: Object(TDOPTOKEN.LPAREN, TDOPTOKEN.TRUE, TDOPTOKEN.RPAREN))
		else
			TdopParenExpr(children)
		TdopAddChild(children, match: TDOPTOKEN.LCURLY, mustMatch:)
		TdopAddChild(children, token: .getCaseList())
		TdopAddChild(children, match: TDOPTOKEN.RCURLY, mustMatch:)
		return TdopCreateNode(TDOPTOKEN.SWITCHSTMT, :children)
		}

	getCaseList()
		{
		return TdopCreateList()
			{ |list|
			hasDefault = false
			while _token() isnt _end and not _token().Match(TDOPTOKEN.RCURLY)
				{
				if not .isCaseOrDefault(_token())
					throw 'Unexpected: ' $ Display(_token())

				caseChildren = Object()
				if _token().Match(TDOPTOKEN.CASE)
					{
					if hasDefault is true
						throw 'Invalid switch: un-reachable case after default'
					TdopAddChild(caseChildren, match: TDOPTOKEN.CASE, mustMatch:)
					TdopAddChild(caseChildren,
						token: TdopSymbolFor.GetExpressionList(
							{ it.Match(TDOPTOKEN.COLON) }))
					}
				else
					{
					TdopAddChild(caseChildren, match: TDOPTOKEN.DEFAULT, mustMatch:)
					TdopAddChild(caseChildren, token: TdopCreateNode(TDOPTOKEN.LIST))
					hasDefault = true
					}
				TdopAddChild(caseChildren, match: TDOPTOKEN.COLON, mustMatch:)
				TdopAddChild(caseChildren, token: TdopSymbolFor.GetStatementList(
					{ .isCaseOrDefault(it) or it.Token is TDOPTOKEN.RCURLY }))
				list.Add(TdopCreateNode(TDOPTOKEN.CASE_ELEM, children: caseChildren))
				}
			}
		}

	isCaseOrDefault(token)
		{
		return token.Match(TDOPTOKEN.CASE) or token.Match(TDOPTOKEN.DEFAULT)
		}
	}
