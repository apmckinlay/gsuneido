// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopSymbolStmtIdentifier
	{
	// IFSTMT [IF, LPAREN, expr, RPAREN, stmts, ELSE, stmts]
	Std()
		{
		children = Object()
		TdopAddChild(children, token: .Tokenize(TDOPTOKEN.IF))
		TdopParenExpr(children)
		TdopAddChild(children, token: TdopStmt(TDOPTOKEN.ELSE))

		if _token().Match(TDOPTOKEN.ELSE) is true
			{
			TdopAddChild(children, match: TDOPTOKEN.ELSE, mustMatch:)
			TdopAddChild(children, token: TdopStmt())
			}
		else
			TdopFillChildren(children,
				tokens: Object(TDOPTOKEN.ELSE, TDOPTOKEN.STMTS))
		return TdopCreateNode(TDOPTOKEN.IFSTMT, :children)
		}
	}
