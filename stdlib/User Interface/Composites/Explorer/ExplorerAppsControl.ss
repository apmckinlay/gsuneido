// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "ExplorerApps"

	New(@args)
		{
		super(@args)
		.Explorer = .FindControl('Explorer')
		.status = .FindControl('Status')
		.Redir('On_New_Folder', .Explorer)
		.Redir('On_New_Item', .Explorer)
		.Redir('On_Save', .Explorer)
		.Redir('On_Delete_Item', .Explorer)
		.Explorer.Redir("Status", this)
		}

	Commands()
		{
		commands = .commands.Copy()
		for c in .Val_or_func(#More_commands)
			commands.Add(c)
		return commands
		}
	commands:
		(
		(New_Folder,				"",			"Create a new folder")
		(New_Item,					"",			"Create a new item")
		(Delete_Item,				"",
			"Delete the selected item or folder", "Delete_item")
		(Print,						"Ctrl+P",	"Print the current item")

		(Undo,						"Ctrl+Z",	"Undo the last action")
		(Redo,						"Ctrl+Y",	"Redo the last action")
		(Cut,						"Ctrl+X",	"Cut the selected text to the clipboard")
		(Copy,						"Ctrl+C",	"Copy the selected text to the clipboard")
		(Paste,						"Ctrl+V",	"Insert the contents of the clipboard")
		(Delete,					"Del",
			"Delete the selected text or next character")
		(Select_All					"Ctrl+A")

		(Find_in_Folders, 			"Shift+Ctrl+F", "Find text in the libraries")
		(Find_Next_in_Folders,		"",
			"Find the next occurrence in the libraries", Find_Next)
		(Find_Previous_in_Folders,	"",
			"Find the previous occurrence in the libraries", Find_Previous)

		(Users_Manual,		"F1")
		(About_Suneido)
		(Close, 					"",			"Close the Window")
		)
	More_commands: ()

	MenuSelect(tip)
		{
		if .status isnt false
			.status.Set(tip)
		}
	Status(status)
		{
		if .status isnt false
			.status.Set(status)
		}
	Inactivate()
		{
		if not .Destroyed?()
			.Explorer.Inactivate()
		}
	}
