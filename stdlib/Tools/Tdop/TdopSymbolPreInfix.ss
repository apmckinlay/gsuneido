// Copyright (C) 2017 Axon Development Corporation All rights reserved worldwide.
TdopSymbolInfix
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
	}