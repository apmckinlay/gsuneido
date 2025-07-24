// Copyright (C) 2017 Axon Development Corporation All rights reserved worldwide.
TdopSymbol
	{
	// RVALUE [LPAREN, expr, RPAREN]
	Nud()
		{
		children = Object()
		TdopAddChild(children, token: this)
		TdopAddChild(children, token: _expr())
		TdopAddChild(children, match: TDOPTOKEN.RPAREN, mustMatch:)
		return TdopCreateNode(TDOPTOKEN.RVALUE, :children)
		}

	// CALL [expr, LPAREN, ATOP|LIST{ARG_ELEM}, RPAREN, BLOCK]
	// ARG_ELEM [KEYARG|ARG, COMMA]
	// KEYARG [STRING|NUMBER, COLON, expr]
	// ARG [expr]
	Led(left)
		{
		children = Object()
		TdopAddChild(children, token: left)
		TdopAddChild(children, token: this)
		TdopAddChild(children, token: .GenerateArgList(TDOPTOKEN.RPAREN))
		TdopAddChild(children, match: TDOPTOKEN.RPAREN, mustMatch:)
		if _token().Token is TDOPTOKEN.LCURLY and
			(_expectingCompound is false or _isNewline() is false)
			{
			t = _token()
			_advance(TDOPTOKEN.LCURLY)
			TdopAddChild(children, token: t.Nud())
			}
		else
			TdopAddChild(children, token: TdopCreateNode(TDOPTOKEN.BLOCK))
		return TdopCreateNode(TDOPTOKEN.CALL, :children)
		}

	GenerateArgList(delim, noAt = false)
		{
		return TdopCreateList()
			{ |list|
			t = _token()
			key = false
			while t isnt _end and not t.Match(delim)
				{
				if t.Match(TDOPTOKEN.AT)
					{
					if list.Size() > 0 or noAt is true
						throw 'Invalid argument list'
					return .handleAtOperator()
					}

				elemChildren = Object()
				if t.Match(TDOPTOKEN.COLON)
					{
					key = true
					arg = .handleKeywordArgShortcut()
					}
				else if .isKeyword()
					{
					key = true
					arg = .handleKeywordArg(delim)
					}
				else if key is true
					throw "un-named arguments must come before named arguments"
				else
					arg = TdopCreateNode(TDOPTOKEN.ARG, children: Object(_expr()))

				TdopAddChild(elemChildren, token: arg)
				TdopAddChild(elemChildren, match: TDOPTOKEN.COMMA)

				list.Add(TdopCreateNode(TDOPTOKEN.ARG_ELEM, children: elemChildren))
				t = _token()
				}
			}
		}

	handleAtOperator()
		{
		t = _token()
		_advance(TDOPTOKEN.AT)
		return t.Nud()
		}

	handleKeywordArgShortcut()
		{
		keyArgChildren = Object()
		TdopAddChild(keyArgChildren, match: TDOPTOKEN.COLON, mustMatch:)
		if not .isJustName()
			throw "Invalid argument list"
		t = _token()
		_advance()
		label = TdopCreateNode(TDOPTOKEN.STRING, value: t.Value)
		TdopAddChild(keyArgChildren, token: label, at: 0)
		TdopAddChild(keyArgChildren, token: t)
		return TdopCreateNode(TDOPTOKEN.KEYARG, children: keyArgChildren)
		}

	isJustName()
		{
		if _token().Token isnt TDOPTOKEN.IDENTIFIER
			return false
		switch (_ahead().Token)
			{
		case TDOPTOKEN.MEMBER, TDOPTOKEN.LPAREN, TDOPTOKEN.LBRACKET, TDOPTOKEN.LCURLY:
			return false
		default:
			return true
			}
		}

	handleKeywordArg(delim)
		{
		keyArgChildren = Object()
		label = _token()
		if label.Match(TDOPTOKEN.IDENTIFIER)
			label.Token = TDOPTOKEN.STRING
		_advance()
		TdopAddChild(keyArgChildren, token: label)
		TdopAddChild(keyArgChildren, match: TDOPTOKEN.COLON, mustMatch:)

		if _token().Match(TDOPTOKEN.COMMA) or .isKeyword() or _token().Match(delim)
			TdopAddChild(keyArgChildren, token: TdopCreateNode(TDOPTOKEN.TRUE))
		else
			TdopAddChild(keyArgChildren, token: _expr())
		return TdopCreateNode(TDOPTOKEN.KEYARG, children: keyArgChildren)
		}

	isKeyword()
		{
		return (TdopAnyName(_token()) or _token().Token is TDOPTOKEN.NUMBER) and
			_ahead().Token is TDOPTOKEN.COLON
		}
	}
