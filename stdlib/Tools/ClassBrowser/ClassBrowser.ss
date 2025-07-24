// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "ClassView"
	DisableContextMenuOptions: (Dump)
	New()
		{
		.exp = .Vert.Explorer
		.model = .exp.Model
		.tree = .exp.Tree
		.find = Record()
		}

	redirs: #(On_Find, On_Find_Next, On_Find_Previous, On_Go_To_Definition)
	SetRedirs()
		{ .redirs.Each({ .Redir(it, .exp.View.Editor) }) }

	DeleteRedirs()
		{ .redirs.Each(.DeleteRedir) }

	Commands:
		(
		(Find,			"Ctrl+F",	"Find text in the current item")
		(Find_Next,		"F3",		"Find the next occurrence in the current item")
		(Find_Previous,	"Shift+F3",	"Find the previous occurrence in the current item")

		(Find_in_Folders,			"Shift+Ctrl+F",	"Find text in the libraries")
		(Find_Next_in_Folders,		"",		"Find the next occurrence in the libraries",
			Find_Next)
		(Find_Previous_in_Folders,	"",		"Find the next occurrence in the libraries",
			Find_Previous)

		(Users_Manual,				"F1")
		(Close,						"",			"Close this window")
		)
	Controls:
		(Vert,
			(Toolbar,
				Find_in_Folders Find_Next_in_Folders Find_Previous_in_Folders)
			(EtchedLine before: 0)
			(ExplorerMulti, ClassBrowserModel,
				(CodeView, readonly:),
				treeArgs: [inorder:, readonly:, multi?: false])
			)

	ClassOutline_SkipHierarchy?()
		{ return true }

	CurrentTable()
		{ return .exp.Get().table }

	CurrentName()
		{ return .exp.Get().name }

	Menu:
		(
		("&File",
			"&Close")
		("&Edit",
			"&Find", "Find Next", "Find Previous")
		)
	names: #()
	names_i: 0
	On_Find_in_Folders()
		{
		if false isnt editor = .exp.View.Editor
			if '' isnt sel = editor.GetSelText().Trim()
				.find.name = sel
		if false is ToolDialog(.Window.Hwnd, [.findcontrol, .find])
			return
		.names_i = -1
		.names = .model.FindNames(Find.Regex(.find.name, .find)).Sort!()
		.On_Find_Next_in_Folders()
		}

	On_Find_Next_in_Folders()
		{
		if (.names_i + 1 >= .names.Size())
			{
			Beep()
			return
			}
		.findItem(.model.GetPath(.names[++.names_i]))
		}
	On_Find_Previous_in_Folders()
		{
		if (.names_i - 1 < 0)
			{
			Beep()
			return
			}
		.findItem(.model.GetPath(.names[--.names_i]))
		}
	findItem(path)
		{
		item = false
		list = .tree.GetChildren(TVI.ROOT)
		for (i = path.Size() - 1; i >= 0; --i)
			{
			for (it in list)
				{
				if (path[i] is .tree.GetName(it))
					{
					item = it
					break
					}
				}
			if (i is 0)
				break
			.tree.ExpandItem(item)
			list = .tree.GetChildren(item)
			}
		.tree.SelectItem(item)
		}

	findcontrol: Controller // mostly duplicate of FindInFoldersControl :-(
		{
		Title: "Find in Folders"
		New(data)
			{
			.Data.Set(data)
			}
		Controls:
			("Record"
				("Horz"
					("Vert"
						(Pair
							(Static 'Find in name')
							(FieldHistory, font: '@mono', size: '+1', width: 30,
								name: "name"))
						"Skip"
						(Horz
							Skip
							(Vert
								(CheckBox, "Match case", name: "case")
								(CheckBox, "Match whole words", name: "word")
								(CheckBox, "Regular expression", name: "regex")
								)
							)
						)
					"Skip"
					(Vert
						(Button, "Find First", xstretch: 0)
						(Skip 8)
						(Button, "Cancel", xstretch: 0)
						)
					)
				)
		DefaultButton: "Find First"
		On_Find_First()
			{
			.Window.Result(.Data.Get())
			}
		}
	}
