// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "TestEditor"
	Menu:
		(
		("&File",
			"E&xit")
		("&Edit",
			"&Undo", "", "Cu&t", "&Copy", "&Paste")
		)
	Controls:
		(Vert
			(Toolbar, Cut, Copy, Paste, "", Undo)
			Editor
			Statusbar
			)
	Commands:
		(
		(Exit,			"",			"")
		(Undo,			"Ctrl+Z",	"Undo the last action")
		(Redo,			"Ctrl+Y",	"Redo the last action")
		(Cut,			"Ctrl+X",	"Cut the selected text to the clipboard")
		(Copy,			"Ctrl+C",	"Copy the selected text to the clipboard")
		(Paste,			"Ctrl+V",	"Insert the contents of the clipboard")
		)
	}
