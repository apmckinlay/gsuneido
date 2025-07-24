// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(position = -1, length = 0)
		{
		children = Object()
		TdopAddChild(children,
			token: TdopCreateNode(TDOPTOKEN.FUNCTION, :position, :length))
		TdopAddChild(children, match: TDOPTOKEN.LPAREN, mustMatch:)
		TdopAddChild(children, token: .handleFuncArgs())
		TdopAddChild(children, match: TDOPTOKEN.RPAREN, mustMatch:)
		TdopAddChild(children, match: TDOPTOKEN.LCURLY, mustMatch:)
		TdopAddChild(children, token: .handleFuncBody())
		TdopAddChild(children, match: TDOPTOKEN.RCURLY, mustMatch:)
		return TdopCreateNode(TDOPTOKEN.FUNCTIONDEF, :children)
		}

	handleFuncArgs()
		{
		if _token().Match(TDOPTOKEN.AT)
			{
			argChildren = Object()
			TdopAddChild(argChildren, match: TDOPTOKEN.AT, mustMatch:)
			TdopAddChild(argChildren, match: TDOPTOKEN.IDENTIFIER, mustMatch:)
			return TdopCreateNode(TDOPTOKEN.PAREM_AT, children: argChildren)
			}

		defaultArg = false
		return TdopCreateList()
			{ |list|
			while not _token().Match(TDOPTOKEN.RPAREN)
				{
				argChildren = Object()
				TdopAddChild(argChildren, match: TDOPTOKEN.DOT)
				TdopAddChild(argChildren, match: TDOPTOKEN.IDENTIFIER, mustMatch:)

				if _token().Match(TDOPTOKEN.EQ)
					{
					defaultArg = true
					TdopAddChild(argChildren, match: TDOPTOKEN.EQ, mustMatch:)
					wasId = _token().Token is TDOPTOKEN.IDENTIFIER
					defaultValue = TdopConstant()
					if '' isnt msg = .validateDefaultValue(defaultValue, wasId)
						throw msg
					TdopAddChild(argChildren, token: defaultValue)

					}
				else if defaultArg is true
					throw 'Default parameters must come last'

				TdopAddChild(argChildren, match: TDOPTOKEN.COMMA)
				list.Add(TdopCreateNode(
					defaultArg ? TDOPTOKEN.PAREM_DEFAULT : TDOPTOKEN.PAREM,
					children: argChildren))
				}
			}
		}

	validateDefaultValue(defaultValue, wasId)
		{
		if defaultValue.Token is TDOPTOKEN.STRING and wasId
			return 'parameter defaults must be constants'
		return ''
		}

	handleFuncBody()
		{
		return _stmts()
		}
	}