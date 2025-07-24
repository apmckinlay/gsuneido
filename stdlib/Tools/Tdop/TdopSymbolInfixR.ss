// Copyright (C) 2017 Axon Development Corporation All rights reserved worldwide.
TdopSymbol
	{
	New(token, lbp, rbp = false)
		{
		super(token, lbp)
		.Rbp = rbp is false ? lbp - 1 : rbp
		}
	// BINARYOP [expr, op, expr]
	Led(left)
		{
		children = Object(left, TdopCreateNode(this), _expr(.Rbp))
		return TdopCreateNode(TDOPTOKEN.BINARYOP, :children)
		}
	ToWrite()
		{
		if .Token is TDOPTOKEN.AND
			return '&&'

		if .Token is TDOPTOKEN.OR
			return '||'

		return super.ToWrite()
		}
	}
