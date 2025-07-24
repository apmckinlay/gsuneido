// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	ContextMenu()
		{
		return #('Move Lines Up\tShift+Ctrl+Up',
			'Move Lines Down\tShift+Ctrl+Down')
		}
	On_Move_Lines_Up()
		{
		.MoveSelectedLinesUp()
		}
	On_Move_Lines_Down()
		{
		.MoveSelectedLinesDown()
		}
	}