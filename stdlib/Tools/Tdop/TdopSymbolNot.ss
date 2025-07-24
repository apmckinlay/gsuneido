// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopSymbolIdentifierInfix
	{
	New(.altToken, .lbp, .rbp, value)
		{
		super(altToken, lbp, value)
		}
	// UNARYOP [op, expr]
	Nud()
		{
		children = TdopAddChild(token: .Tokenize(TDOPTOKEN.NOT))
		TdopAddChild(children, token: _expr(.rbp))
		return TdopCreateNode(TDOPTOKEN.UNARYOP, :children)
		}
	// NOTINOP [expr, NOT, IN, LPAREN, LIST{EXPR_ELEM}, PRAREN]
	Led(left)
		{
		t = _token()
		_advance(TDOPTOKEN.IN)
		children = t.Led(left).Children
		TdopAddChild(children, token: .Tokenize(TDOPTOKEN.NOT), at: 1)
		return TdopCreateNode(TDOPTOKEN.NOTINOP, :children)
		}
	Getter_Lbp()
		{
		if _ahead().Match(TDOPTOKEN.IN)
			return .lbp
		return 0
		}
	Break(unused)
		{
		return _isNewline()
		}
	}