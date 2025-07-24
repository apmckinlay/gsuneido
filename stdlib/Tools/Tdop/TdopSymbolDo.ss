// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopSymbolStmtIdentifier
	{
	// DOSTMT [DO, stmts, WHILE, LPAREN, expr, RPAREN]
	Std()
		{
		children = Object()
		TdopAddChild(children, token: .Tokenize(TDOPTOKEN.DO))
		TdopAddChild(children, token: TdopStmt(TDOPTOKEN.WHILE))

		TdopAddChild(children, match: TDOPTOKEN.WHILE)

		TdopParenExpr(children)
		return TdopCreateNode(TDOPTOKEN.DOSTMT, :children)
		}
	}