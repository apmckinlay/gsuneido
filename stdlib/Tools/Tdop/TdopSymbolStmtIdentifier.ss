// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopSymbolIdentifier
	{
	New(.StmtToken, value)
		{
		super(value)
		}
	Match(token)
		{
		return super.Match(token) or .StmtToken is token
		}
	Tokenize(match)
		{
		if .StmtToken is match
			return TdopCreateNode(.StmtToken, position: .Position, length: .Length)
		return super.Tokenize(match)
		}
	}
