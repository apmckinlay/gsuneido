// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
// This addon handles Toolbar controls / commands which can function without an .Editor
LibViewAddon
	{
	Commands(cmds)
		{
		cmds.Add(
			#(Version_Control_Settings),

			#(Print_Item, 'Ctrl+P',	'Print the current item', 'Print'),
			#(Print_Library, '', `Print a library's contents`),

			#(Import_Records, '', 'Import records from a text file into a library'),
			#(Export_Record, '', 'Append current record to a text file'),
			#(Restore_Import, '', 'Restore records previously imported'),

			#(New_Library, '', 'Create and use a new library'),
			#(Use_Library, '', 'Use an existing library'),
			#(Unuse_Library, '', 'Un-use a currently loaded library'),

			#(Locate, 'Ctrl+L'),
			)
		}

	Ctrl()
		{
		return Object(order: 0,
			ctrl: Object('Horz',
				Object('MenuButton', 'New', .newmenu())
				Object('Vert'
					Object('Horz',
						.toolbar(),
						'SvcSettingsIcons',
						'LibLocate')
						#(EtchedLine before: 0))))
		}

	newmenu()
		{
		.menuRedirs = Object()
		newmenu = Object('Item', 'Folder', 'Library', '')
		Plugins().ForeachContribution('LibView', 'New')
			{ |c|
			newmenu.Add(c[2])
			.menuRedirs[c[2]] = c[3] /*= 2: 'New' menu option, 3: new text */
			}
		return newmenu
		}

	toolbar()
		{
		cmds = .BuildCommands(sequenced?:).Filter({ it.Member?(#seq) }).Sort!(By(#seq))
		toolbar = .buildToolBar(cmds)
		.PluginTools(
			{ |cmd, icon|
			if icon isnt ''
				toolbar.Add(cmd)
			})
		return toolbar
		}

	// The following lists consists of base, Explorer and Toolbar buttons,
	// add addon remaining buttons after these
	baseToolbar: #('Toolbar', '', 'Delete_Item', '',
		'Run', 'Print_Item', '',
		'Cut', 'Copy', 'Paste', '', 'Undo', 'Redo')
	buildToolBar(cmds)
		{
		toolbar = .baseToolbar.Copy()
		prevSeq = ''
		for cmd in cmds
			{
			if prevSeq isnt curSeq = cmd.seq.RoundDown(0)
				{
				toolbar.Add('')
				prevSeq = curSeq
				}
			toolbar.Add(cmd.Copy().Delete(#seq))
			}
		toolbar.Add('')
		return toolbar
		}

	Init()
		{
		.hwnd = .Window.Hwnd
		}

	On_Export_Record()
		{ .LibExportFile() }

	On_Import_Records(defaultLib = false)
		{
		if false isnt info = LibViewImportRecordControl(.hwnd, defaultLib)
			LibIO.Import(info.fileName, info.lib, interactive:)
		}
	On_Restore_Import()
		{
		LibViewImportRestoreControl(.hwnd)
		}

	On_New(what)
		{
		if .menuRedirs.Member?(what)
			.Explorer.NewItem(false, text: .menuRedirs[what])
		}

	On_New_Library()
		{
		if false is lib = Ask('Name', 'New Library', .hwnd)
			return
		if TableExists?(lib)
			{
			.AlertError('New Library', 'Table ' $ lib $ ' already exists')
			return
			}
		LibTreeModel.Create(lib)
		.use(lib)
		.ResetCtrls()
		}

	On_Print_Item()
		{
		.Save()
		x = .Explorer.Get()
		if x is #()
			return
		lib = .CurrentTable()
		Params.On_Print(Object('Library', x.name, x.text)
			title: lib.Capitalize(),
			name: 'print_library',
			previewWindow: .hwnd)
		}

	On_Print_Library()
		{
		.Save()
		ToolDialog(.hwnd, LibraryReport(.CurrentTable()))
		}

	On_Unuse_Library(option)
		{ .changeLibs({ .unuse(it) }, option, 'Unuse Library') }

	changeLibs(block, option, msg)
		{
		ResetCaches()
		if false is block(option)
			return .AlertError(msg, 'Unable to ' $ msg.Lower() $ ' ' $ option)
		.ResetCtrls()
		LibraryTags.Reset()
		if Sys.Client?()
			Unload()
		}

	unuse(option)
		{ return Sys.Client?() ? ServerEval('Unuse', option) : Unuse(option) }

	On_Use_Library(option)
		{ .changeLibs({ .use(it) }, option, 'Use Library') }

	use(option)
		{ return Sys.Client?() ? ServerEval('Use', option) : Use(option) }

	On_Locate()
		{
		if false isnt locate = .FindControl('LibLocate')
			locate.SetFocus()
		}

	ToolbarMenu(menu)
		{
		tools = Object('&Tools'
			'&Run', '',
			'&Go To Definition'
			'Create Test for &Method'
			'&Show Parameters'
			'&Inspect', '')
		.PluginTools() { |cmd| tools.Add(cmd.Tr('_', ' ')) }

		menu.Add(
			#("&File",
				"&New Library...", ('&Use Library'), ('Unuse Library'), "",
				"New &Folder", "New &Item", "&Delete Item",
				"&View/Restore Item as of...", "",
				"&Import Records...", "&Export Record...", "Restore Import...", "",
				"&Print Item", "Print Library...", "",
				"&Close"),
			#("&Edit",
				"&Undo", "&Redo", "",
				"Cu&t", "&Copy", "&Paste", "&Delete", "Select &All", "",
				"&Find...", "Find &Next", "Find &Previous", "R&eplace...", "",
				"Comment &Lines", "Comment &Selection", "",
				"Flag", "Next Flag", "Previous Flag", "&Go To Line..."),
			tools)
		}
	}
