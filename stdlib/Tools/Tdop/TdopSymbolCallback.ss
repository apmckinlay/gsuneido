// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopSymbolReserved
	{
	// CALLBACKDEF [CALLBACK, LPAREN, LIST{CALLBACK_PAREM}, RPAREN]
	Nud()
		{
		children = TdopAddChild(token: .Tokenize(.AltToken))
		TdopAddChild(children, match: TDOPTOKEN.LPAREN, mustMatch:)
		TdopAddChild(children, token: TdopDllEntity.TypeList(.AltToken))
		TdopAddChild(children, match: TDOPTOKEN.RPAREN, mustMatch:)
		return TdopCreateNode(TDOPTOKEN.CALLBACKDEF, :children)
		}
	}