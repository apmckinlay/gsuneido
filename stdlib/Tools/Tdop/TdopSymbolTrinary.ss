// Copyright (C) 2017 Axon Development Corporation All rights reserved worldwide.
TdopSymbol
	{
	New(token, lbp)
		{
		super(token, lbp)
		}
	// TRINARYOP [expr, Q_MARK, expr, COLON, expr]
	Led(left)
		{
		children = Object()
		TdopAddChild(children, token: left)
		TdopAddChild(children, token: this)
		_setStmtnest(_getStmtnest() + 1)
		TdopAddChild(children, token: _expr(0))
		TdopAddChild(children, match: TDOPTOKEN.COLON, mustMatch:)
		_setStmtnest(_getStmtnest() - 1)
		TdopAddChild(children, token: _expr(0))
		return TdopCreateNode(TDOPTOKEN.TRINARYOP, :children)
		}
	}
