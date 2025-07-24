// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	CallClass(text, title = '', table = false, name = false)
		{
		Window(Object(this, text, table, name), :title, keep_placement: 'CodeViewer')
		}

	New(.text, .table = false, .name = false)
		{
		.editor = .FindControl('Editor')
		.Redir('On_Find', .editor)
		.Redir('On_Find_Next', .editor)
		.Redir('On_Find_Previous', .editor)
		.Redir('On_Find_Next_Selected', .editor)
		.Redir('On_Find_Prev_Selected', .editor)
		.Redir('On_Flag', .editor)
		.Redir('On_Next_Flag', .editor)
		.Redir('On_Previous_Flag', .editor)
		.Redir('On_Go_To_Definition', .editor)
		}

	Startup()
		{
		.editor.SetSelect(0)
		}

	Controls()
		{
		return Object(#Vert,
			#(Toolbar
				Copy, "",
				Find, Find_Next, Find_Previous, "",
				Flag, Next_Flag, Previous_Flag),
			Object(#CodeView, data: [text: .text, table: .table, name: .name],
				readonly:))
		}

	Commands:
		(
		(Copy,				"Ctrl+C",	"Copy the selected text to the clipboard")

		(Find,				"Ctrl+F",	"Find text in the current item")
		(Find_Next,			"F3",		"Find the next occurrence in the current item")
		(Find_Previous,		"Shift+F3",
			"Find the previous occurrence in the current item")
		(Find_Next_Selected,"Ctrl+F3",	"Find the next occurrence of the selected text")
		(Find_Prev_Selected,"Shift+Ctrl+F3",
			"Find the previous occurrence of the selected text")
		)
	}
