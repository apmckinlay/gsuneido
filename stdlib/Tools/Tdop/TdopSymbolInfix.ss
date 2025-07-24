// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
TdopSymbol
	{
	// BINARYOP [expr, op, expr]
	Led(left)
		{
		return TdopCreateNode(TDOPTOKEN.BINARYOP,
			children: Object(left, this, _expr(.Lbp)))
		}
	ToWrite()
		{
		if .Token is TDOPTOKEN.ISNT
			return '!='

		if .Token is TDOPTOKEN.IS
			return '=='

		return super.ToWrite()
		}
	}