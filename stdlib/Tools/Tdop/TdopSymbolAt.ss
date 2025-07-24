// Copyright (C) 2017 Axon Development Corporation All rights reserved worldwide.
TdopSymbol
	{
	// ATOP [AT, ADD, NUMBER, expr]
	Nud()
		{
		children = Object()
		TdopAddChild(children, token: this)
		add = _token().Match(TDOPTOKEN.ADD)
		TdopAddChild(children, match: TDOPTOKEN.ADD)
		if add is true
			TdopAddChild(children, match: TDOPTOKEN.NUMBER, mustMatch:)
		else
			TdopAddChild(children, token: TdopCreateNode(TDOPTOKEN.NUMBER))
		TdopAddChild(children, token: _expr(0))
		return TdopCreateNode(TDOPTOKEN.ATOP, :children)
		}
	}