// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopSymbolReserved
	{
	// STRUCTDEF [STRUCT, LCURLY, LIST{STRUCT_MEMBER}, RCURLY]
	Nud()
		{
		children = TdopAddChild(token: .Tokenize(.AltToken))
		TdopAddChild(children, match: TDOPTOKEN.LCURLY, mustMatch:)
		TdopAddChild(children, token: TdopDllEntity.TypeList(.AltToken))
		TdopAddChild(children, match: TDOPTOKEN.RCURLY, mustMatch:)
		return TdopCreateNode(TDOPTOKEN.STRUCTDEF, :children)
		}
	}