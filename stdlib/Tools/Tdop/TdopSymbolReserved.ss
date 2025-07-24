// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopSymbolIdentifier
	{
	New(.AltToken, value)
		{
		super(value)
		}
	Nud()
		{
		return .Tokenize(.AltToken)
		}
	Match(token)
		{
		return super.Match(token) or .AltToken is token
		}
	Tokenize(match)
		{
		if .AltToken is match
			return TdopCreateNode(.AltToken, position: .Position, length: .Length)
		return super.Tokenize(match)
		}
	}