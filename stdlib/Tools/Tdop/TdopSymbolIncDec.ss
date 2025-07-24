// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopSymbol
	{
	New(token, lbp, .rbp)
		{
		super(token, lbp)
		}
	// PREINCDEC [INC|DEC, expr]
	Nud()
		{
		e = _expr(.rbp)
		return TdopCreateNode(TDOPTOKEN.PREINCDEC, Object(this, e))
		}
	// POSTINCDEC [expr, INC|DEC]
	Led(left)
		{
		return TdopCreateNode(TDOPTOKEN.POSTINCDEC, Object(left, this))
		}
	}