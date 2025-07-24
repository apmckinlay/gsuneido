// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopSymbolReserved
	{
	New(altToken, .lbp, value)
		{
		super(altToken, value)
		.Delete(#Lbp)
		}
	Nud()
		{
		return this
		}
	// BINARYOP [expr, op, expr]
	Led(left)
		{
		children = TdopAddChild(token: left)
		TdopAddChild(children, token: .Tokenize(.AltToken))
		TdopAddChild(children, token: _expr(.lbp))
		return TdopCreateNode(TDOPTOKEN.BINARYOP, :children)
		}
	Getter_Lbp()
		{
		if _ahead().Token is TDOPTOKEN.COLON
			return 0
		return .lbp
		}
	Break(unused)
		{
		return false
		}
	}