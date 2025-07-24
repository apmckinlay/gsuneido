// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddonIDE
	{
	maxWidth: 0
	Setting: ide_show_line_numbers
	Init()
		{
		.SetMarginTypeN(0, SC.MARGIN_NUMBER)
		.maxWidth = Max(.calcWidth(), .maxWidth)
		.setLineNumbers(.Set)
		}

	marginWidth: 40
	calcWidth()
		{ return ScaleWithDpiFactor(.marginWidth) - .GetMarginWidthN(0) }

	ContextMenu()
		{ #("Show/Hide Line Numbers") }

	On_ShowHide_Line_Numbers()
		{ .setLineNumbers(not .Set) }

	setLineNumbers(.Set)
		{ .SetMarginWidthN(0, .Set ? .maxWidth : 0) }
	}
