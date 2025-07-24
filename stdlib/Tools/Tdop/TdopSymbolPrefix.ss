// Copyright (C) 2017 Axon Development Corporation All rights reserved worldwide.
TdopSymbol
	{
	New(token, lbp, .rbp)
		{
		super(token, lbp)
		}
	// UNARYOP [op, expr]
	Nud()
		{
		return TdopCreateNode(TDOPTOKEN.UNARYOP,
			children: Object(TdopCreateNode(this), _expr(.rbp)))
		}
	ToWrite()
		{
		if .Token is TDOPTOKEN.NOT
			return '!'
		return super.ToWrite()
		}
	}