// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	Margin: false
	Init()
		{
		charSize = .SendMessageTextIn(SCI.TEXTWIDTH, SC.STYLE_DEFAULT, 'M')
		spaceSize = .SendMessageTextIn(SCI.TEXTWIDTH, SC.STYLE_DEFAULT, ' ')
		.SetEdgeMode(charSize is spaceSize
			? SC.EDGE_LINE
			: SC.EDGE_BACKGROUND)
		.SetEdgeColour(.GetSchemeColor('longLineMargin'))
		.SetEdgeColumn(.Margin isnt false ? .Margin : CheckCode.MaxLineLength)
		}
	}