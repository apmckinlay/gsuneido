// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddonIDE
	{
	Setting: ide_show_line_numbers
	Init()
		{
		.setLineNumbers(.Set)
		}

	ContextMenu()
		{ #("Show/Hide Line Numbers") }

	On_ShowHide_Line_Numbers()
		{ .setLineNumbers(not .Set) }

	setLineNumbers(.Set)
		{ .SetOption('lineNumbers', .Set is true) }
	}