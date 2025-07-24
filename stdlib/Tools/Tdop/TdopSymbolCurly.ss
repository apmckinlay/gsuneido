// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopSymbol
	{
	// BLOCK [LCURLY, BITOR, BPAREM_AT|LIST{BPAREM}, BITOR, LIST{STMT}, RCURLY]
	// BPAREM_AT [AT, IDENTIFIER]
	// BPAREM [IDENTIFIER, COMMA]
	Nud()
		{
		children = Object()
		TdopAddChild(children, token: this)
		if _token().Match(TDOPTOKEN.BITOR)
			{
			TdopAddChild(children, match: TDOPTOKEN.BITOR, mustMatch:)
			TdopAddChild(children, token: .handleBlockArgs())
			TdopAddChild(children, match: TDOPTOKEN.BITOR, mustMatch:)
			}
		else
			{
			TdopAddChild(children, token: TdopCreateNode(TDOPTOKEN.BITOR))
			TdopAddChild(children, token: TdopCreateNode(TDOPTOKEN.LIST))
			TdopAddChild(children, token: TdopCreateNode(TDOPTOKEN.BITOR))
			}

		TdopAddChild(children, token: _stmts())
		TdopAddChild(children, match: TDOPTOKEN.RCURLY, mustMatch:)
		return TdopCreateNode(TDOPTOKEN.BLOCK, :children)
		}

	// CALL [expr, LPAREN, LIST, RPAREN, BLOCK]
	Led(left)
		{
		if left.Token is TDOPTOKEN.IDENTIFIER and TdopIsGlobal(left.Value)
			return TdopClass(left, leftCurly: this)

		children = Object()
		TdopAddChild(children, token: left)
		TdopAddChild(children, token: TdopCreateNode(TDOPTOKEN.LPAREN))
		TdopAddChild(children, token: TdopCreateNode(TDOPTOKEN.LIST))
		TdopAddChild(children, token: TdopCreateNode(TDOPTOKEN.RPAREN))
		TdopAddChild(children, token: .Nud())
		return TdopCreateNode(TDOPTOKEN.CALL, :children)
		}

	// STMTS [LCURLY, LIST{STMT}, RCURLY]
	Std()
		{
		children = Object()
		TdopAddChild(children, token: this)
		TdopAddChild(children, token: _stmts())
		TdopAddChild(children, match: TDOPTOKEN.RCURLY, mustMatch:)
		return TdopCreateNode(TDOPTOKEN.STMTS, :children)
		}

	handleBlockArgs()
		{
		if _token().Match(TDOPTOKEN.AT)
			{
			argChildren = TdopAddChild(match: TDOPTOKEN.AT, mustMatch:)
			TdopAddChild(argChildren, match: TDOPTOKEN.IDENTIFIER, mustMatch:)

			return TdopCreateNode(TDOPTOKEN.BPAREM_AT, children: argChildren)
			}

		return TdopCreateList()
			{ |list|
			while not _token().Match(TDOPTOKEN.BITOR) and _token() isnt _end
				{
				argChildren = TdopAddChild(match: TDOPTOKEN.IDENTIFIER, mustMatch:)
				TdopAddChild(argChildren, match: TDOPTOKEN.COMMA)
				list.Add(TdopCreateNode(TDOPTOKEN.BPAREM, children: argChildren))
				}
			}
		}
	Break(left)
		{
		// when _expectingCompound is true, always break except Global {}
		if _expectingCompound is true and
			not (left.Token is TDOPTOKEN.IDENTIFIER and TdopIsGlobal(left.Value) and
				not _isNewline())
			return true
		// when _expectingCompound is false, break only when newline and left token isnt
		// global identifier or member
		if _expectingCompound is false and _isNewline() and _getStmtnest() is 0 and
			not (left.Token is TDOPTOKEN.IDENTIFIER and TdopIsGlobal(left.Value) or
				left.Token is TDOPTOKEN.MEMBEROP)
			return true
		return false
		}
	}
