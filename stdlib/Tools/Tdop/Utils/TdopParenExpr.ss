// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
function(children)
	{
	parens = _token().Match(TDOPTOKEN.LPAREN)
	TdopAddChild(children, match: TDOPTOKEN.LPAREN)
	TdopAddChild(children,
		token: parens isnt true ? TdopStmtExpr(expectingCompound:) : _expr())
	TdopAddChild(children, match: TDOPTOKEN.RPAREN, mustMatch: parens)
	}
