// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopSymbolIdentifierInfix
	{
	New(altToken, lbp, value)
		{
		super(altToken, lbp, value)
		.rbp = lbp - 1
		}
	// BINARYOP [expr, op, expr]
	Led(left)
		{
		children = Object(left, .Tokenize(.AltToken), _expr(.rbp))
		return TdopCreateNode(TDOPTOKEN.BINARYOP, :children)
		}
	}
