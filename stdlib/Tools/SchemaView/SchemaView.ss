// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// TODO: allow creating tables & columns
Controller
	{
	Title: SchemaView
	CallClass(table)
		{
		.Goto(table)
		}

	New()
		{
		.Explorer.Redir(#Status, this)
		.FindControl(#SchemaLocate).SetFocus()
		}

	Commands:
		(
		(Find,			"Ctrl+F",	"Find text in the current item")
		(Find_Next,		"F3",		"Find the next occurrence in the current item")
		(Find_Previous,	"Shift+F3",	"Find the previous occurrence in the current item")
		(Refresh,		"F5")
		(Close,			"",			"Close this window")
		(Users_Manual,	"F1")
		)

	Controls()
		{
		return Object(#Vert,
			Object("SchemaLocate")
			Object(#ExplorerMulti,
				#SchemaModel, #(SchemaViewView),
				treeArgs: [inorder:, readonly:, multi?: false],
				besideTabs: #RefreshButton),
			#Statusbar)
		}

	Getter_Explorer()
		{
		return .Explorer = .Vert.Explorer
		}

	getter_tree()
		{
		return .tree = .Explorer.Tree
		}

	CurrentTable()
		{
		return .tree.GetName(.tree.GetSelectedItem())
		}

	redirs: #(On_Go_To_Definition, On_Find, On_Find_Next, On_Find_Previous)
	SetRedirs()
		{
		.On_Refresh()
		.redirs.Each({ .Redir(it, .Explorer.View.Editor) })
		}

	DeleteRedirs()
		{
		.redirs.Each(.DeleteRedir)
		}

	Menu: (
		("&File",
			"&Refresh List",
			"&Close")
		("&Edit",
			"&Find...", "Find &Next", "Find &Previous")
		)
	MenuSelect(tip)
		{
		.Vert.Status.Set(tip)
		}

	Status(status)
		{
		.Vert.Status.Set(status)
		}

	On_Refresh_List()
		{
		.Explorer.Reset()
		}

	Activate()
		{
		.On_Refresh()
		}

	On_Refresh()
		{
		selected = .Explorer.Tabs.GetSelected()
		if selected isnt -1 and .Explorer.TabConstructed?(selected)
			.Explorer.RefreshTab(selected, force:)
		}

	Goto(ref)
		{
		sv = GotoPersistentWindow('SchemaView', SchemaView)
		return sv.Explorer.GotoPath(ref)
		}

	Locate(val)
		{
		.Explorer.GotoPath(val)
		}
	}
