// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopSymbolStmtIdentifier
	{
	// BREAKCONTINUESTMT [BREAK|CONTINUE, SEMICOLON]
	Std()
		{
		children = TdopAddChild(token: .Tokenize(.StmtToken))
		TdopAddChild(children, match: TDOPTOKEN.SEMICOLON)
		return TdopCreateNode(TDOPTOKEN.BREAKCONTINUESTMT, :children)
		}
	}