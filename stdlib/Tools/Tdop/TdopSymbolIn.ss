// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopSymbolIdentifierInfix
	{
	Layout: #(LEFT: 0, LPAREN: 1, LIST: 2, RPAREN: 3)
	// INOP [expr, IN, LPAREN, LIST{EXPR_ELEM}, PRAREN]
	Led(left)
		{
		children = TdopAddChild(token: left)
		TdopAddChild(children, token: .Tokenize(.AltToken))
		TdopAddChild(children, match: TDOPTOKEN.LPAREN, mustMatch:)
		inList = TdopCreateList()
			{ |list|
			while _token() isnt _end and not _token().Match(TDOPTOKEN.RPAREN)
				{
				elemChildren = Object()
				TdopAddChild(elemChildren, token: _expr())
				TdopAddChild(elemChildren, match: TDOPTOKEN.COMMA)
				list.Add(TdopCreateNode(TDOPTOKEN.EXPR_ELEM, children: elemChildren))
				}
			}
		TdopAddChild(children, token: inList)
		TdopAddChild(children, match: TDOPTOKEN.RPAREN, mustMatch:)
		return TdopCreateNode(TDOPTOKEN.INOP, :children)
		}
	Break(unused)
		{
		return _isNewline() and _getStmtnest() is 0
		}
	}
