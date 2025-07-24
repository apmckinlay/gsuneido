// Copyright (C) 2017 Axon Development Corporation All rights reserved worldwide.
TdopSymbol
	{
	bracket: ''
	New(text)
		{
		super(TDOPTOKEN.STRING)
		if text.Size() > 0 and text[0] in ("#", "'", "`", '"')
			{
			.bracket = text[0]
			valueSize = text.Size() - (.bracket is '#' ? 1 : 2)
			.Value = text[1::valueSize]
			}
		}
	Nud()
		{
		return this
		}
	ToWrite()
		{
		switch (.bracket)
			{
		case '':
			return .Value
		case '#':
			return .bracket $ .Value
		case "'", '"', '`':
			return .bracket $ .Value $ .bracket
			}
		}
	}