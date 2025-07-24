// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopSymbolReserved
	{
	//CALL [MEMBEROP [SUPER, DOT, IDENTIFIER], LPAREN, LIST, RPAREN, BLOCK]
	Nud()
		{
		superToken = .Tokenize(.AltToken)
		if _token().Match(TDOPTOKEN.LPAREN)
			left = TdopCreateNode(TDOPTOKEN.MEMBEROP,
				children: Object(superToken
					TdopCreateNode(TDOPTOKEN.DOT),
					TdopCreateNode(TDOPTOKEN.IDENTIFIER, value: 'New')))
		else
			{
			dot = _token()
			_advance(TDOPTOKEN.DOT)
			left = dot.Led(superToken)
			}
		lparen = _token()
		_advance(TDOPTOKEN.LPAREN)
		return lparen.Led(left)
		}
	}
