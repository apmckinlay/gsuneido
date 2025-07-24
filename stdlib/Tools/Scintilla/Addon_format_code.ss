// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	ContextMenu()
		{ return #("Format Code") }
	On_Format_Code()
		{
		line = .LineFromPosition()

		text = .get()
		s = FormatCode(text)
		if s is text
			return
		if .SelSize() is 0
			.PasteOverAll(s)
		else
			.Paste(s)

		.GotoLine(line)
		}
	get()
		{
		text = .GetSelText()
		if text is ""
			text = .Get()
		return text
		}
	}