// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	ContextMenu()
		{ return #("Go To Line\tCtrl+G") }
	On_Go_To_Line()
		{
		if (false is line = Ask("Line number", "Go To Line", .Window.Hwnd))
			return
		.GotoLine(Number(line) - 1)
		}
	}