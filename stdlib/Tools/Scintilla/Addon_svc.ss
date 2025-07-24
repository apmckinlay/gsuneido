// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	ContextMenu()
		{
		return #('Line History')
		}

	On_Line_History()
		{
		.SelectLine(.LineFromPosition())
		line = .GetLine().RightTrim() // trim new line
		if line.Trim().Size() <= 1
			{
			Alert('Please select a line with more than one character.',
				"Version History", .Window.Hwnd, MB.ICONINFORMATION)
			return
			}

		table = .Send("CurrentTable")
		name = .Send("CurrentName")

		if false isnt ctrl = VersionHistoryControl(table, name) // IdeTabbedView
			ctrl.FindControl('History').FindLineHistory(line)
		}
	}
