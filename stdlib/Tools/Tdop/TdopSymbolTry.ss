// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopSymbolStmtIdentifier
	{
	// TRYSTMT [TRY, stmts, CATCHSTMT]
	// CATCHSTMT [CATCH, CATCH_COND, stmts]
	// CATCH_COND [LPAREN, IDENTIFIER, COMMA, STRING, RPAREN]
	Std()
		{
		children = Object()
		TdopAddChild(children, token: .Tokenize(.StmtToken))
		TdopAddChild(children, token: TdopStmt(TDOPTOKEN.CATCH))
		TdopAddChild(children, token: .catchStmt())
		return TdopCreateNode(TDOPTOKEN.TRYSTMT, :children)
		}

	catchStmt()
		{
		if not _token().Match(TDOPTOKEN.CATCH)
			return TdopCreateNode(TDOPTOKEN.CATCHSTMT)

		catchChildren = Object()
		TdopAddChild(catchChildren, match: TDOPTOKEN.CATCH, mustMatch:)
		TdopAddChild(catchChildren, token: .catchCond())
		TdopAddChild(catchChildren, token: TdopStmt())
		return TdopCreateNode(TDOPTOKEN.CATCHSTMT, children: catchChildren)
		}

	catchCond()
		{
		if not _token().Match(TDOPTOKEN.LPAREN)
			return TdopCreateNode(TDOPTOKEN.CATCH_COND)

		catchCondChildren = Object()
		TdopAddChild(catchCondChildren, match: TDOPTOKEN.LPAREN, mustMatch:)

		if _token().Match(TDOPTOKEN.IDENTIFIER)
			{
			TdopAddChild(catchCondChildren, match: TDOPTOKEN.IDENTIFIER, mustMatch:)
			if _token().Match(TDOPTOKEN.COMMA)
				{
				TdopAddChild(catchCondChildren, match: TDOPTOKEN.COMMA, mustMatch:)
				TdopAddChild(catchCondChildren, match: TDOPTOKEN.STRING, mustMatch:)
				}
			else
				TdopFillChildren(catchCondChildren,
					tokens: Object(TDOPTOKEN.COMMA, TDOPTOKEN.STRING))
			}
		else
			TdopFillChildren(catchCondChildren,
				tokens: Object(TDOPTOKEN.IDENTIFIER, TDOPTOKEN.COMMA, TDOPTOKEN.STRING))
		TdopAddChild(catchCondChildren, match: TDOPTOKEN.RPAREN, mustMatch:)
		return TdopCreateNode(TDOPTOKEN.CATCH_COND, children: catchCondChildren)
		}
	}
