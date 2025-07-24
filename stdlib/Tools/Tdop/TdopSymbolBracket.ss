// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopSymbol
	{
	// CALL [IDENTIFIER(Record), LPAREN, LIST{ARG_ELEM}, RPAREN, BLOCK]
	Nud()
		{
		children = Object()
		TdopAddChild(children, token: TdopSymbolIdentifier('Record'))
		TdopAddChild(children, token: this)
		TdopAddChild(children,
			token: TdopSymbolParens.GenerateArgList(TDOPTOKEN.RBRACKET, noAt:))
		TdopAddChild(children, match: TDOPTOKEN.RBRACKET, mustMatch:)
		TdopAddChild(children, token: TdopCreateNode(TDOPTOKEN.BLOCK))
		return TdopCreateNode(TDOPTOKEN.CALL, :children)
		}
	MaxRange: 2147483647
	// SUBSCRIPT [expr, LBRACKET, RANGE|expr, RBRACKET]
	// RANGE [expr, RANGETO|RANGELEN, expr]
	Led(left)
		{
		children = Object()
		TdopAddChild(children, token: left)
		TdopAddChild(children, token: this)

		rangeToken = false
		if _token().Match(TDOPTOKEN.RANGETO) or _token().Match(TDOPTOKEN.RANGELEN)
			first = TdopSymbolNumber(0)
		else
			first = _expr(0)

		if _token().Match(TDOPTOKEN.RANGETO) or _token().Match(TDOPTOKEN.RANGELEN)
			{
			rangeToken = _token()
			_advance()

			if _token().Match(TDOPTOKEN.RBRACKET)
				second = TdopSymbolNumber(.MaxRange)
			else
				second = _expr(0)
			index = TdopCreateNode(TDOPTOKEN.RANGE,
				children: Object(first, rangeToken, second))
			}
		else
			index = first

		TdopAddChild(children, token: index)
		TdopAddChild(children, match: TDOPTOKEN.RBRACKET, mustMatch:)
		return TdopCreateNode(TDOPTOKEN.SUBSCRIPT, :children)
		}
	}