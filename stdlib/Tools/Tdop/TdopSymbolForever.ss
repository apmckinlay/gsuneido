// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopSymbolStmtIdentifier
	{
	// FOREVERSTMT [FOREVER, stmts]
	Std()
		{
		children = Object()
		TdopAddChild(children, token: .Tokenize(TDOPTOKEN.FOREVER))
		TdopAddChild(children, token: TdopStmt())
		return TdopCreateNode(TDOPTOKEN.FOREVERSTMT, :children)
		}
	}