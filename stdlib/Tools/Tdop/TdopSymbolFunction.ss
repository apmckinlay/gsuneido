// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopSymbolReserved
	{
	// FUNCTIONDEF [FUNCTION, LPAREN, PAREM_AT|LIST{PAREM}, RPAREN, LCURLY, LIST{stmt}, RCURLY]
	// PAREM_AT [AT, IDENTIFIER]
	// PAREM [DOT, IDENTIFIER, COMMA]
	// PAREM_DEFAULT [DOT, IDENTIFIER, EQ, const, COMMA]
	Nud()
		{
		return TdopFunction(position: .Position, length: .Length)
		}
	}