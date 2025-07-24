// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	styleLevel: 80
	Init()
		{
		.indic_word = .IndicatorIdx(level: .styleLevel)
		}

	Styling()
		{
		return [
			[level: .styleLevel,
				indicator: [INDIC.BOX, fore: 0x00ff00]]]
		}

	Highlight(target)
		{
		text = .Get()
		text.ForEachMatch("\<(?q)" $ target $  "(?-q)\>")
			{|m|
			m = m[0]
			.SetIndicator(.indic_word, m[0], m[1])
			text = text[.. m[0]]
			}
		}
	}
