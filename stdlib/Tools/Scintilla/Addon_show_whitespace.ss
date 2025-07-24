// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddonIDE
	{
	Setting: ide_show_whitespace
	Init()
		{
		// need this to show in the middle of lines
		.SetWhiteSpaceFore(0x888888 /*= gray color*/)
		.setWhiteSpace(.Set)
		}

	setWhiteSpace(.Set)
		{
		.SetViewWS(.Set ? 1 : 0)
		.SetViewEol(.Set ? 1 : 0)
		}

	ContextMenu()
		{ return #("Show/Hide Whitespace") }

	On_ShowHide_Whitespace()
		{ .setWhiteSpace(.GetViewWS() is 0) }
	}
