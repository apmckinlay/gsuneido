// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddonForChanges
	{
	styleLevel: 30
	Init()
		{
		.indic_word = .IndicatorIdx(level: .styleLevel)
		.IndicSetHoverStyle(.indic_word, INDIC.ROUNDBOX)
		}

	Styling()
		{
		return [[level: .styleLevel, indicator: [INDIC.TEXTFORE, fore: CLR.BLUE]]]
		}

	DoubleClick()
		{
		url = .getCurrentURL()
		if false isnt url
			ShellExecute(.WindowHwnd(), 'open', url, 0, 0, SW.SHOW)
		}

	getCurrentURL()
		{
		org = end = .GetCurrentPos()
		while .WordChars.Has?(.GetAt(org - 1))
			--org
		while .WordChars.Has?(.GetAt(end))
			++end
		if org >= end
			return false

		text = .GetRange(org, end)
		if false is match = .MatchUrl(text)
			return false

		.SetSelect(org + match[0], match[1])
		return text[match[0]::match[1]]
		}

	prefixPattern: "https??://|s??ftps??://|mailto://|www\."
	WordChars: "-_.!~*'();?:@&=+$,%#/" $ 	// special
		"0123456789" $						// digits
		"abcdefghijklmnopqrstuvwxyz" $ 		// a-z
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ"		// A-Z
	MatchUrl(text, pos = false)
		{
		if false is match = text.Match(
			"(" $ .prefixPattern $ ")[" $ .WordChars $ "]+", :pos)
			return false

		return match[0]
		}

	ProcessChunk(text, pos)
		{
		.ClearIndicator(.indic_word, pos, text.Size())
		i = 0
		while false isnt m = .MatchUrl(text, i)
			{
			.mark_word(pos + m[0], m[1])
			i = m[0] + m[1]
			}
		}

	mark_word(start, len)
		{
		.SetIndicator(.indic_word, start, len)
		}
	}
