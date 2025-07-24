// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopSymbolStmtIdentifier
	{
	// WHILESTMT [WHILE, LPAREN, expr, RPAREN, stmts]
	Std()
		{
		children = Object()
		TdopAddChild(children, token: .Tokenize(TDOPTOKEN.WHILE))
		TdopParenExpr(children)
		TdopAddChild(children, token: TdopStmt())
		return TdopCreateNode(TDOPTOKEN.WHILESTMT, :children)
		}
	}