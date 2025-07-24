// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopSymbolReserved
	{
	// DLLDEF [DLL, IDENTIFIER, IDENTIFIER, COLON, STRING, LPAREN, LIST{DLL_PAREM}, RPAREN]
	Nud()
		{
		children = TdopAddChild(token: .Tokenize(.AltToken))

		TdopAddChild(children, match: TDOPTOKEN.IDENTIFIER, mustMatch:)
		TdopAddChild(children, match: TDOPTOKEN.IDENTIFIER, mustMatch:)
		TdopAddChild(children, match: TDOPTOKEN.COLON, mustMatch:)

		userFunctionName = _token()
		_advance(TDOPTOKEN.IDENTIFIER)
		userFunctionName.Token = TDOPTOKEN.STRING

		if _token().Match(TDOPTOKEN.AT)
			{
			_advance(TDOPTOKEN.AT)
			t = _token()
			_advance(TDOPTOKEN.NUMBER)
			userFunctionName.Value $= '@' $ t.Value
			}
		TdopAddChild(children, token: userFunctionName)
		TdopAddChild(children, match: TDOPTOKEN.LPAREN, mustMatch:)
		TdopAddChild(children, token: TdopDllEntity.TypeList(.AltToken))
		TdopAddChild(children, match: TDOPTOKEN.RPAREN, mustMatch:)

		return TdopCreateNode(TDOPTOKEN.DLLDEF, :children)
		}
	}