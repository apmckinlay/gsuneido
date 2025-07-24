// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopSymbolStmtIdentifier
	{
	// THROWSTMT [THROW, expr, SEMICOLON]
	Std()
		{
		children = TdopAddChild(token: .Tokenize(.StmtToken))
		TdopAddChild(children, token: TdopStmtExpr())
		TdopAddChild(children, match: TDOPTOKEN.SEMICOLON)
		return TdopCreateNode(TDOPTOKEN.THROWSTMT, :children)
		}
	}
