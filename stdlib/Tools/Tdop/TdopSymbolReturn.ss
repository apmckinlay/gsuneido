// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopSymbolStmtIdentifier
	{
	// RETURNSTMT [RETURN, THROW, LIST{EXPR_ELEM}, SEMICOLON]
	Std()
		{
		children = Object()
		TdopAddChild(children, token: .Tokenize(TDOPTOKEN.RETURN))
		TdopAddChild(children, match: TDOPTOKEN.THROW)
		elems = TdopCreateList()
			{ |list|
			while not .end?()
				{
				elemChildren = Object()
				TdopAddChild(elemChildren, token: TdopStmtExpr())
				TdopAddChild(elemChildren, match: TDOPTOKEN.COMMA)
				list.Add(TdopCreateNode(TDOPTOKEN.EXPR_ELEM, children: elemChildren))
				}
			}
		TdopAddChild(children, token:  elems)
		TdopAddChild(children, match: TDOPTOKEN.SEMICOLON, implicit:)
		return TdopCreateNode(TDOPTOKEN.RETURNSTMT, :children)
		}

	end?()
		{
		return  _token() is _end or _token().Match(TDOPTOKEN.SEMICOLON) or
			_token().Match(TDOPTOKEN.RCURLY) or _isNewline() is true
		}
	}
