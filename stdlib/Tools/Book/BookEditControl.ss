// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "BookEdit"
	New(.table = '')
		{
		.initBook()
		.Clipformat = RegisterClipboardFormat('Suneido_BOOKEDIT')

		.subs = [
			PubSub.Subscribe('BookTreeChange', .reset)
			PubSub.Subscribe('BookRecordChange', .refresh)
			]
		}

	CanPaste?()
		{ return IsClipboardFormatAvailable(.Clipformat) }

	initBook()
		{
		if not Suneido.Member?("EditBooks")
			Suneido.EditBooks = Object()
		Suneido.EditBooks[.table] = this
		}

	Controls()
		{
		.contentType = BookContent.Type(.table)
		toolbar = Object("Toolbar"
			"New_Item", "Delete_Item" "",
			"Undo", "Redo", "", "Cut", "Copy", "Paste", "",
			"Refresh", "Run", "",
			"Find", "",
			"Find_in_Folders", "Find_Next_in_Folders", "Find_Previous_in_Folders",
			"")
		if .contentType is #html
			toolbar.Add("H1", "H2", "H3", "H4", "P", "LI", "DT", "DD", "PRE", "")
		toolbar.Add("Bold", "Italic", "Underline", "Code", "Link", "Add_Image_Tag",
			"Goto", ""
			"Find_References_to_Current", "Version_History")
		extraAddons = .contentType is #html
			? #()
			: #(Addon_html_edit: false, Addon_md_edit:, Addon_html: false, Addon_md:)
		return Object('Vert',
			Object('Horz'
				toolbar,
				Object("BookEditLocate", .table) #(Skip small:)),
			#(EtchedLine before: 0),
			Object('ExplorerMulti',
				Object('BookEditModel', .table),
				Object('BookEditSplit', :extraAddons),
				treeArgs: [inorder:]),
			ystretch: 1)
		}

	reset(args)
		{
		reset = args.Filter({ it.table is .table })
		if reset.NotEmpty?()
			.Explorer.ResetControls(force: reset.Any?({ it.force }))
		}

	refresh(args)
		{
		refresh = args.Filter({ it.table is .table })
		if refresh.NotEmpty?()
			.Explorer.Refresh(refresh)
		}

	viewRedirs: #(Status, MenuSelect)
	editorRedirs: #(On_Find, On_Find_Next, On_Find_Previous, On_Replace,
		On_Go_To_Definition, On_Link, On_Add_Paragraph_Tags, On_H1, On_H2,
		On_H3, On_H4, On_P, On_LI, On_DT, On_DD, On_PRE, On_Bold, On_Italic,
		On_Underline, On_Code, On_Add_Image_Tag)
	SetRedirs()
		{
		.viewRedirs.Each({ .Redir(it, .View) })
		.editorRedirs.Each({ .Redir(it, .Editor) })
		}

	DeleteRedirs()
		{
		.viewRedirs.Each(.DeleteRedir)
		.editorRedirs.Each(.DeleteRedir)
		}

	Commands:
		(
		(Close,				"",				"Close this window")
		(New_Item, 			"",				"Add a new help item", New_Folder)
		(Delete_Item, 		"",				"Delete selected help item", delete_item)

		(Set_Order,			"Ctrl+Alt+O",	"Set the item's \"order\" field")
		(Unorder_Children,	"Ctrl+Alt+U",	"Set order of all item's children to ''")
		(Renumber_Children, "Ctrl+Alt+R",
			"Reset the order of all item's children to whole numbers.")


		(Import_Image		""				"Import an image file")
		(Export_Image		""				"Export an image file")
		(Export_Multiple_Html_Files
							""				"Export HTML pages to individual files")
		(Export_Single_Html_File
							""				"Export HTML pages to one large file")
		(Refresh,			"F5",			"Refresh")

		(Undo,				"Ctrl+Z",		"Undo the last action")
		(Redo,				"Ctrl+Y",		"Redo the last action")
		(Cut,				"Ctrl+X",		"Cut the selected text to the clipboard")
		(Copy,				"Ctrl+C",		"Copy the selected text to the clipboard")
		(Paste,				"Ctrl+V",		"Insert the contents of the clipboard")
		(Delete,			"Del",			"Delete the selected text or next character")

		(Find,				"Ctrl+F",		"Find text in the current item")
		(Find_Next,			"F3",
			"Find the next occurrence in the current item")
		(Find_Previous,		"Shift+F3",
			"Find the previous occurrence in the current item")
		(Replace,			"Ctrl+H",		"Find and replace text in the current item")

		(Bold,				"Ctrl+B",		"Add bold tags around the selected text")
		(Italic, 			"Ctrl+I",		"Add italic tags around the selected text")
		(Underline, 		"Ctrl+U",		"Add underline tags around the selected text")
		(Code,				"Ctrl+Alt+C",	"Add code tags around the selected text")
		(Link,				"Alt+L",		"Link")
		(Locate,			"Ctrl+L",		"Locate")
		(Add_Paragraph_Tags, "",
			"Add paragraph tags and remove line breaks", "P")
		(Add_Image_Tag,		"Ctrl+Alt+I", 	"Insert image tag", I)
		(H1					"Ctrl+F1",		"", '1')
		(H2					"Ctrl+F2",		"", '2')
		(H3					"Ctrl+F3",		"", '3')
		(H4					"Ctrl+F4",		"", '4')
		(Goto				"Ctrl+G", 		"Go to the selected record", G)
		(P					"Ctrl+P",		"", P)
		(LI					"Ctrl+Alt+L",	"", L)
		(DT					"Ctrl+Alt+T",	"", T)
		(DD					"Ctrl+Alt+D",	"", D)
		(PRE				"Ctrl+Alt+P",	"", R)

		(Find_in_Folders,	"Shift+Ctrl+F",	"Find text in the libraries")
		(Find_Next_in_Folders,		"",
			"Find the next occurrence in the libraries", Find_Next)
		(Find_Previous_in_Folders,	"",
			"Find the next occurrence in the libraries", Find_Previous)

		(Import_Records,	"",				"Import records from a text file into a book")
		(Export_Record,		"",				"Append current record to a text file")
		(Insert_File,		"",				"Insert a text file")
		(Run,				"F9",			"Run selected text", '!')
		(Open Book)
		(Build How To Index)
		(Build Ftsearch Index)
		(Users_Manual,		"F1")
		(Version_History, 				"", 			"Version History", H)
		(Find_References_to_Current, 	"Ctrl+R",		"Find References to Current", R)
		)

	Menu()
		{
		file = #("&File",
			"&New Item", "&Delete Item", "",
			"Import Ima&ge...", "Export Image...", "",
			"Export &Single Html File...", "Export &Multiple Html Files...", "",
			"&Import Records..." "&Export Record...", "",
			"&Build How To Index", "Build Ftsearch Index", "",
			"&Close"
			)
		edit = #("&Edit",
			"&Undo", "&Redo", "", "Cu&t", "&Copy", "&Paste", "&Delete", "",
			"&Find...", "Find &Next", "Find &Previous", "R&eplace...", "",
			"Find &in Folders", "Find Next in Folders", "Find Previous in Folders", "",
			"&Insert File..."
			)
		format = Object("F&ormat")
		if .contentType is #html
			format.Add("H&1", "H&2", "H&3", "H&4", "&P", "LI", "D&T", "D&D", "P&RE", "")
		format.Add("&Bold", "&Italic", "&Underline", "&Code", "&Link", "Add Image Tag")
		if .contentType is #html
			format.Add("", "&Add Paragraph Tags")
		tools = #("&Tools"
			"Refresh", "",
			"Set &Order...", "Unorder Children", "&Renumber Children", "",
			"Run", "Open Book"
			)
		return Object(file, edit, format, tools)
		}

	On_Context_New()
		{ .Explorer.On_New_Folder() }

	On_New_Item()
		{ .Explorer.On_New_Folder() }

	On_Delete_Item()
		{
		if .Editor is false
			return
		if '' is msg = .View.Deletable?()
			.Explorer.On_Delete_Item()
		else
			.AlertInfo('BookEdit', msg)
		}

	On_Refresh()
		{
		if .Editor is false
			return
		dirty = .Editor.Dirty?()
		.View.Refresh(force?:)
		.Editor.Dirty?(dirty)
		}

	On_Insert_File()
		{
		if '' isnt (filename = OpenFileName(title: 'Insert File')) and
			false isnt (text = GetFile(filename))
			.Editor.Paste(text)
		}

	On_Set_Order()
		{
		// SetFocus is needed in case label is being edited to commit the edit
		.ClearFocus()
		if .Explorer.CurItem is false
			return

		// update the order
		if false is x = Query1(.table, num: .Tree.GetParam(.Explorer.CurItem))
			return

		order = BookEditOrderControl(.Window.Hwnd, oldorder: x.order)
		QueryApply1(.table, num: .Tree.GetParam(.Explorer.CurItem))
			{ .updateOrder(it, order, it.Transaction()) }
		.Tree.SortChildren(.Tree.GetParent(.Explorer.CurItem))
		}

	updateOrder(rec, order, t = false)
		{
		rec.order = order
		.Explorer.Model.Update(rec, t)
		}

	On_Unorder_Children()
		{
		// Set the order of all an item's children to 0
		if false is sel = .getSelected('Unorder Children')
			return

		Transaction(update:)
			{ |t|
			if sel is -1
				path = ""
			else
				{
				if false is x = t.Query1(.table, num: sel)
					return
				path = x.path $ "/" $ x.name
				}

			t.QueryApply(.table $ ' where path is "' $ path $ '" and order isnt ""')
				{ |x|
				.updateOrder(x, '', t)
				}
			}
		}

	getSelected(title)
		{
		if false is sel = .Tree.GetParam(.Explorer.CurItem)
			.AlertError(title, 'Please select an item.')
		return sel
		}

	On_Renumber_Children()
		{
		if false is sel = .getSelected('Renumber Children')
			return
		if sel is -1
			path = ""
		else
			{
			item = Query1(.table, num: sel)
			path = item.path $ '/' $ item.name
			}
		order = 0
		colIncrement = 10
		QueryApplyMulti(.table $ ' where path is ' $ Display(path) $
			' and order isnt "" sort path, order, name', update:)
			{ |x|
			if ((x.order / colIncrement).Int() isnt (order / colIncrement).Int())
				order = (x.order / colIncrement).Int() * colIncrement
			if x.order isnt order
				.updateOrder(x, order, x.Transaction())
			endOfColumnOrder = 9
			smallIncrement = .001
			order += (order % colIncrement is endOfColumnOrder ? smallIncrement : 1)
			}
		}

	On_Run()
		// post:	attempts to Eval the contents of the current selection
		{
		// redirect printing to console
		prevPrint = Suneido.Print
		Suneido.Print = .Print
		// perform run

		if .Editor is false or '' is sel = .Editor.GetSelText().Trim()
			return
		try
			Print(sel.Eval()) // needs Eval
		catch (x)
			.AlertError('Run Error', x)
		Suneido.Print = prevPrint
		}

	On_Open_Book()
		{ BookControl(.table) }

	Print(s)
		// pre:	s is a string
		// post:	if a console for this exists, s is appended to this' console ELSE
		//		a console for this is created and s is appended to it
		{
		if Suneido.GetDefault('Console', false) is false
			Window(#(Console), x: 0, y: 0, w: .35, h: .5)
		Suneido.Console.Append(s)
		}

	On_Build_How_To_Index()
		{
		.Save()
		if #() isnt bad = BookHowToIndex(.table)
			.AlertError('Build How To Index',
				'bad <!-- option: name(s):\n    ' $ bad.Join('\n    '))
		}

	On_Build_Ftsearch_Index()
		{
		msg = ''
		title = 'Build Ftsearch Index'
		Working(title)
			{
			msg = ServerEval(`IndexHelp`, .table)
			}
		if '' isnt msg
			.AlertInfo(title, msg)
		}

	Save()
		{ .Explorer.On_Save() }

	Inactivate()
		{
		if not .Destroyed?()
			.Explorer.Inactivate()
		}

	Goto(address)
		{ .Explorer.GotoPath("Help" $ ((address =~ "^/") ? "" : "/") $ address) }

	On_Import_Image()
		{
		path = .getSelectedPathOrRes()
		if path.AfterFirst('/res') isnt '' and
				false is ImportImageIntoSubfolder(.Window.Hwnd, path)
			path = '/res'

		files = OpenFileName(title: 'Import Image',
			multi:, hwnd: .Window.Hwnd,
			filter: "Image Files (*.png;*.jpg;*.gif;*.bmp;*.emf;*.svg)" $
				"\x00*.png;*.jpg;*.gif;*.bmp;*.emf;*.svg\x00All Files (*.*)\x00*.*")
		if files is #()
			return

		.EnsureResFolderExists(.table)
		refreshTree? = false
		for file in files
			if ImportSvcTableText(file, .table, path, .Window.Hwnd)
				refreshTree? = true // New record has been output

		if refreshTree?
			SvcTable(.table).Publish(#TreeChange, force:)
		}

	getSelectedPathOrRes()
		{
		x = .Explorer.Get()
		path = x is false
			? ''
			: x.GetDefault(#path, '')
		if path isnt "/res" and not path.Prefix?("/res/")
			return '/res'
		else if not x.name.Has?('.') // if no '.' assume folder
			path $= '/' $ x.name
		return path
		}

	EnsureResFolderExists(table)
		{
		if QueryEmpty?(table, path: '', name: 'res')
			QueryOutput(table, [path: '', name: 'res', num: NextTableNum(table)])
		}

	On_Export_Image()
		{
		if .Explorer.RootSelected?()
			{
			.AlertError('Export Image', "Can't export root folder")
			return
			}
		selected = .Explorer.GetSelected()
		path = .Explorer.Getpath(selected)
		name = path[path.FindLast("/") - path.Size() + 1 ..]
		// strip tablename and record name off path
		path = path[.table.Size()..][.. -(name.Size() + 1)]
		filename = SaveFileName(
			hwnd:	.Window.Hwnd,
			flags:	OFN.PATHMUSTEXIST | OFN.HIDEREADONLY | OFN.NOCHANGEDIR,
			title:	"Export Image to",
			file: name
			)
		if filename isnt ''
			PutFile(filename, .Explorer.Get().text)
		}

	On_Export_Multiple_Html_Files()
		{
		if false isnt dir = Ask('Directory', 'Export Html', .Window.Hwnd,
			#(FieldHistory name: 'BookExport', selectFirst:))
			BookExport(.table, dir)
		}

	On_Export_Single_Html_File()
		{
		filename = SaveFileName(
			hwnd:	.Window.Hwnd,
			flags:	OFN.PATHMUSTEXIST | OFN.HIDEREADONLY | OFN.OVERWRITEPROMPT |
				OFN.NOCHANGEDIR,
			title:	"Export to",
			filter: "HTML Files (*.htm)\x00*.htm\x00All Files (*.*)\x00*.*\x00",
			ext: 'htm'
			)
		if filename isnt ''
			BookExportOne(.table, filename)
		}

	On_Export_Record()
		{
		if .Explorer.RootSelected?()
			{
			.AlertError('Export Record', "Can't export root folder")
			return
			}
		path = .Explorer.Getpath(.Explorer.GetSelected())
		name = path.AfterLast('/')
		// strip tablename and record name off path
		path = path.RemovePrefix(.table).BeforeLast('/')
		filename = SaveFileName(hwnd: .Window.Hwnd, title: "Export (append) to",
			flags: OFN.PATHMUSTEXIST | OFN.HIDEREADONLY | OFN.NOCHANGEDIR)
		if filename isnt ''
			LibIO.Export(.table, name, filename, path, interactive:)
		}

	On_Import_Records()
		{
		filename = OpenFileName(hwnd: .Window.Hwnd, title: "Import from")
		if filename isnt ''
			LibIO.Import(filename, .table, interactive:)
		}

	names: #()
	names_i: 0
	find: false
	On_Find_in_Folders()
		{
		if .find is false
			.find = .Editor is false ? [] : .Editor.FindReplaceData()
		if false is FindInFoldersControl(.find)
			return
		.names = Object()
		.names_i = -1
		Transaction(read:)
			{ |t|
			.search(t, Find.Regex(.find.name, .find), Find.Regex(.find.find, .find))
			}
		.On_Find_Next_in_Folders()
		}

	// TODO search with single query, then sort by path (not recursive)
	search(t, name, text, path = "") // recursive
		{
		t.QueryApply(.table $ " where path is " $ Display(path) $
			" sort order, name")
			{ |x|
			path = x.path $ "/" $ x.name
			if (not BookResource?(path) and
				(name is '' or x.name =~ name) and
				(text is '' or x.text =~ text))
				.names.Add(path)
			.search(t, name, text, path) // do children (if any)
			}
		}

	On_Find_Next_in_Folders()
		{
		if .names_i + 1 < .names.Size()
			if .Explorer.GotoPath(.table $ .names[++.names_i])
				.Defer({ .findText(#Next) }, uniqueID: #BookEditFind)
		}

	On_Find_Previous_in_Folders()
		{
		if .names_i - 1 >= 0
			if .Explorer.GotoPath(.table $ .names[--.names_i])
				.Defer({ .findText(#Previous) }, uniqueID: #BookEditFind)
		}

	findText(dir)
		{
		.Editor.FindReplaceData().Merge(.find)
		.Editor[#On_Find_ $ dir]()
		.Editor.EnsureVisible(.Editor.LineFromPosition())
		}

	On_Goto()
		{ .GotoHelp(.getCurrentLink()) }


	On_Locate()
		{ .BookEditLocate.SetFocus() }

	GotoHelp(text)
		{
		text = text.Tr('"').Trim()
		text = text.Has?('/')
			? .table $ '/' $ text.RemovePrefix('suneido:').RemovePrefix('/').
				RemovePrefix(.table $ '/')
			:.table $ '/res/' $ text
		.Explorer.GotoPath(text)
		}

	getCurrentLink()
		{
		if .Editor is false
			return ''
		if '' isnt text = .Editor.GetSelText()
			return text
		return .FindLinkedHelpPage()
		}

	FindLinkedHelpPage(select? = false)
		{
		if .Editor is false
			return ''
		pos = .Editor.GetCurrentPos()
		lastChar = .Editor.GetSelect().cpMax
		if .Editor.GetAt(lastChar) in ("=", '(', ';')
			return ''
		if .Editor.GetAt(pos) in ('"', "'")
			pos = .Editor.GetSelect().cpMin
		if false is org = .findLink(pos, -1)
			return ''
		if false is end = .findLink(pos, +1)
			return ''
		if org >= end
			return ''
		if select?
			.Editor.SetSelect(org, end - org + 1)
		return .Editor.GetRange(org, end + 1)
		}

	findLink(pos, dir)
		{
		while (('"' isnt at = .Editor.GetAt(pos + dir)) and at isnt "'" and at isnt "")
			{
			if at is '\n' or at is '\x00' or at is '<'
				return false
			pos += dir
			}
		return pos
		}

	On_Version_History()
		{
		if .Explorer.RootSelected?()
			.AlertInfo(.Title, 'No version history for book root')
		else
			VersionHistoryControl(.CurrentTable(), .CurrentName())
		}

	On_Find_References_to_Current()
		{ FindReferencesControl(.CurrentName()) }

	Getter_Explorer()
		{ return .Explorer = .Vert.Explorer }

	Getter_View()
		{ return .Explorer.View }

	Getter_Editor()
		{ return .View is false ? false : .View.Editor }

	Getter_Tree()
		{ return .Tree = .Explorer.Tree }

	Getter_BookEditLocate()
		{
		return .Vert.Horz.BookEditLocate
		}

	ModelTable()
		{ return .Explorer.Model.GetTable() }

	CurrentTable()
		{ return .View.CurrentTable() }

	CurrentName()
		{ return .View.CurrentName() }

	BraceMatch_AdditionalBraces()
		{
		return '<>'
		}

	// TODO: handle VertSplit too
	GetState()
		{
		return Object(splitterpos: .Explorer.HorzSplit.GetSplit(),
			tabs: .Explorer.GetTabsPaths(), table: .table,
			activeTabPath: .Explorer.Getpath(.Explorer.CurItem))
		}

	SetState(state)
		{
		stateTable = state.GetDefault(#table, false)
		if .table is '' and stateTable isnt false
			{
			.Explorer.Reset(Object(BookEditModel, .table = stateTable))
			.BookEditLocate.SetTable(stateTable)
			}
		if .table is stateTable
			.Explorer.RestoreState(state)
		}

	Locate(val)
		{
		.Explorer.GotoPath(val)
		}

	Destroy()
		{
		Suneido.EditBooks.Delete(.table)
		.subs.Each(#Unsubscribe)
		super.Destroy()
		}
	}
