// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// LIST{STRUCT_MEMBER|CALLBACK_PAREM|DLL_PAREM|}
	// STRUCT_MEMBER [DLL_POINTER|DLL_ARRAY|DLL_NORMAL, IDENTIFIER, SEMICOLON]
	// CALLBACK_PAREM [DLL_POINTER|DLL_ARRAY|DLL_NORMAL, IDENTIFIER, COMMA]
	// DLL_PAREM [DLL_IN, DLL_POINTER|DLL_ARRAY|DLL_NORMAL, IDENTIFIER, COMMA]
	// DLL_IN [LBRACKET, IN, RBRACKET]
	// DLL_NORMAL [IDENTIFIER]
	// DLL_POINTER [IDENTIFIER, MUL]
	// DLL_ARRAY [IDENTIFIER, LBRACKET, NUMBER, RBRACKET]
	TypeList(type)
		{
		Assert(Object(TDOPTOKEN.DLL, TDOPTOKEN.STRUCT, TDOPTOKEN.CALLBACK) has: type)
		return TdopCreateList()
			{ |list|
			while .hasAnItem(type)
				{
				children = Object()
				if type is TDOPTOKEN.DLL
					if _token().Match(TDOPTOKEN.LBRACKET)
						{
						inChildren = Object()
						TdopAddChild(inChildren, match: TDOPTOKEN.LBRACKET, mustMatch:)
						TdopAddChild(inChildren, match: TDOPTOKEN.IN, mustMatch:)
						TdopAddChild(inChildren, match: TDOPTOKEN.RBRACKET, mustMatch:)
						TdopAddChild(children, token: TdopCreateNode(TDOPTOKEN.DLL_IN,
							children: inChildren))
						}
					else
						TdopAddChild(children, token: TdopCreateNode(TDOPTOKEN.DLL_IN))


				typeChildren = Object()
				t = TDOPTOKEN.DLL_NORMAL
				TdopAddChild(typeChildren, match: TDOPTOKEN.IDENTIFIER, mustMatch:)
				if _token().Match(TDOPTOKEN.MUL)
					{
					t = TDOPTOKEN.DLL_POINTER
					TdopAddChild(typeChildren, match: TDOPTOKEN.MUL, mustMatch:)
					}
				else if _token().Match(TDOPTOKEN.LBRACKET)
					{
					t = TDOPTOKEN.DLL_ARRAY
					TdopAddChild(typeChildren, match: TDOPTOKEN.LBRACKET, mustMatch:)
					TdopAddChild(typeChildren, match: TDOPTOKEN.NUMBER, mustMatch:)
					TdopAddChild(typeChildren, match: TDOPTOKEN.RBRACKET, mustMatch:)
					}
				TdopAddChild(children, token: TdopCreateNode(t, children: typeChildren))
				TdopAddChild(children, match: TDOPTOKEN.IDENTIFIER, mustMatch:)

				if type is TDOPTOKEN.STRUCT
					TdopAddChild(children, match: TDOPTOKEN.SEMICOLON, implicit:,
						mustMatch: .hasAnItem(type))
				else
					TdopAddChild(children, match: TDOPTOKEN.COMMA,
						mustMatch: .hasAnItem(type))

				list.Add(TdopCreateNode(.nodeType(type), :children))
				}
			}
		}

	hasAnItem(type)
		{
		switch (type)
			{
		case TDOPTOKEN.STRUCT:
			return not _token().Match(TDOPTOKEN.RCURLY)
		case TDOPTOKEN.DLL, TDOPTOKEN.CALLBACK:
			return not _token().Match(TDOPTOKEN.RPAREN)
			}
		}

	nodeType(type)
		{
		switch (type)
			{
		case TDOPTOKEN.STRUCT:
			return TDOPTOKEN.STRUCT_MEMBER
		case TDOPTOKEN.CALLBACK:
			return TDOPTOKEN.CALLBACK_PAREM
		case TDOPTOKEN.DLL:
			return TDOPTOKEN.DLL_PAREM
			}
		}
	}
