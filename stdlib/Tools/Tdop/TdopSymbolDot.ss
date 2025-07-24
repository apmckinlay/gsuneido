// Copyright (C) 2017 Axon Development Corporation All rights reserved worldwide.
TdopSymbol
	{
	// MEMBEROP [SELFREF, DOT, IDENTIFIER]
	Nud()
		{
		member = .getMember()
		children = Object(TdopCreateNode(TDOPTOKEN.SELFREF), this, member)
		return TdopCreateNode(TDOPTOKEN.MEMBEROP, :children)
		}
	// MEMBEROP [expr, DOT, IDENTIFIER]
	Led(left)
		{
		member =  .getMember()
		children = Object(left, this, member)
		return TdopCreateNode(TDOPTOKEN.MEMBEROP, :children)
		}
	getMember()
		{
		member = _token()
		_advance(TDOPTOKEN.IDENTIFIER)
		return member
		}
	}
