// Copyright (C) 2017 Axon Development Corporation All rights reserved worldwide.
TdopSymbol
	{
	// DATE
	// SYMBOL
	// TdopConstantObject
	Nud()
		{
		t = _token()
		if t.Token is TDOPTOKEN.NUMBER
			{
			t.Token = TDOPTOKEN.DATE
			t.Position = .Position
			t.Length += .Length
			_advance()
			return t
			}
		else if TdopAnyName(t)
			{
			t.Token = TDOPTOKEN.SYMBOL
			t.Position = .Position
			t.Length += .Length
			_advance()
			return t
			}
		else if not (t.Token in (TDOPTOKEN.LBRACKET, TDOPTOKEN.LCURLY, TDOPTOKEN.LPAREN))
			throw "invalid literal following " $ Display(t)

		return TdopConstantObject(position: .Position, length: .Length)
		}
	}
