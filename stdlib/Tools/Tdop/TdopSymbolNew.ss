// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopSymbolReserved
	{
	New(altToken, value, .rbp)
		{
		super(altToken, value)
		}
	// NEWOP [NEW, expr, LPAREN, ATOP|LIST{ARG_ELEM}, RPAREN, BLOCK]
	Nud()
		{
		newToken = .Tokenize(.AltToken)
		e = _expr(.rbp)
		if _token().Match(TDOPTOKEN.LPAREN)
			{
			t = _token()
			_advance(TDOPTOKEN.LPAREN)
			call = t.Led(e)
			TdopAddChild(call.Children, token: newToken, at: 0)
			return TdopCreateNode(TDOPTOKEN.NEWOP, children: call.Children)
			}

		children = Object()
		TdopAddChild(children, token: newToken)
		TdopAddChild(children, token: e)
		TdopAddChild(children, token: TdopCreateNode(TDOPTOKEN.LPAREN))
		TdopAddChild(children, token: TdopCreateNode(TDOPTOKEN.LIST))
		TdopAddChild(children, token: TdopCreateNode(TDOPTOKEN.RPAREN))
		TdopAddChild(children, token: TdopCreateNode(TDOPTOKEN.BLOCK))
		return TdopCreateNode(TDOPTOKEN.NEWOP, :children)
		}
	}