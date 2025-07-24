// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	Init()
		{
		.DefineStyle(SC.STYLE_INDENTGUIDE, .GetSchemeColor('defaultFore'),
			back: .GetSchemeColor('defaultBack'))
		.SetIndentationGuides(SC.IV_LOOKBOTH)
		}
	}
