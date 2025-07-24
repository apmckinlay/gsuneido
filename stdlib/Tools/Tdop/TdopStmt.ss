// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
function (alt_end = false) // helper
	{
	// SEMICOLON
	// STMTS [LCURLY, LIST{STMT}, RCURLY]
	if _token().Match(TDOPTOKEN.SEMICOLON)
		{
		t = _token()
		_advance()
		return t
		}
	if not _token().Match(TDOPTOKEN.LCURLY)
		return _stmt(alt_end)

	children = Object()
	TdopAddChild(children, match: TDOPTOKEN.LCURLY, mustMatch:)
	TdopAddChild(children, token: _stmts())
	TdopAddChild(children, match: TDOPTOKEN.RCURLY, mustMatch:)
	return TdopCreateNode(TDOPTOKEN.STMTS, :children)
	}
