// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	// Title needs to always be the same i.e. not include library name
	// so keep_placement: works
	Title: 'References'
	CallClass(name)
		{
		if name is ""
			{
			Beep()
			return
			}
		IdeTabbedView('FindReferences', name)
		}
	settings_key: 'FindReferencesControlSettings'
	New(.name)
		{
		settings = UserSettings.Get(.settings_key,
			#(exclude_tests: false, exclude_updates: false, exclude_libs: ''),
			user: 'default')
		.exclude_tests_ctrl = .FindControl('exclude_tests_chk')
		.exclude_updates_ctrl = .FindControl('exclude_updates_chk')
		.exclude_libs_ctrl = .FindControl('exclude_libs_list')
		.exclude_tests_ctrl.Set(settings.exclude_tests)
		.exclude_updates_ctrl.Set(settings.GetDefault(#exclude_updates, false))
		.exclude_libs_ctrl.Set(settings.exclude_libs)
		.currentExcludeLibs = settings.exclude_libs
		.load()
		}
	cols: (findref_location, findref_table, findref_found, findref_folder)
	Controls()
		{
		libsList = LibraryTables().SortWith!('Lower')
		used = Libraries().Intersect(libsList)
		libsList = used.MergeUnion(libsList)
		return ['Vert',
			['VirtualList', query: '',
				columns: .cols, loadAll?:,
				columnsSaveName: 'References',
				name: 'list', disableSelectFilter:,
				stretchColumn: 'findref_found']
			#(Skip 4),
			['Horz',
				#(Skip, 6),
				#(Static, "double-click to go to definition"),
				#Fill,
				#(Static, "Exclude:")
				#Skip
				#('CheckBox', "Tests", name: 'exclude_tests_chk'), 'Skip',
				#('CheckBox', "Updates", name: 'exclude_updates_chk'), 'Skip',
				#(Static, 'Libraries: '),
				['ChooseTwoList', libsList, "Exclude Libraries",
					name: 'exclude_libs_list', noSort:],
				'Skip', 'RefreshButton']]
		}
	load()
		{
		.list = .FindControl(#list)
		.list.Redir('Scintilla_DoubleClick', this)
		.list.Redir('Scintilla_KillFocus', this)
		.list.Redir('Scintilla_SetFocus', this)
		excludeTests = .exclude_tests_ctrl.Get()
		excludeUpdates = .exclude_updates_ctrl.Get()
		excludeLibs = .exclude_libs_ctrl.GetSelectedList()
		refs = FindReferences(.name, :excludeTests, :excludeUpdates, :excludeLibs)
		.basename = refs.basename
		.list.SetQuery('views where view_name is ""', .cols) // fake query
		for i in refs.list.Members()
			{
			line = refs.list[i]
			rec = [view_name: Timestamp()].Merge(line)
			rec.findref_location = line.GetDefault('Location', '')
			rec.findref_table = line.GetDefault('Table', '')
			rec.findref_found = line.GetDefault('Found', '')
			rec.findref_folder = line.GetDefault('Folder', '')
			.list.AddRecord(rec, pos: i)
			}
		}

	NewValue(value, source)
		{
		if source is .exclude_libs_ctrl and .exclude_libs_ctrl.Valid?() and
			.currentExcludeLibs isnt value or
			source is .exclude_tests_ctrl or source is .exclude_updates_ctrl
			{
			.currentExcludeLibs = value
			.load()
			}
		}

	On_Refresh()
		{
		.load()
		}

	VirtualList_AllowSort()
		{
		return false
		}
	VirtualList_DoubleClick(rec, col/*unused*/)
		{
		if rec isnt false
			{
			if BookTables().Has?(rec.Table)
				.On_Context_Edit_Documentation()
			else
				.go_To_Definition()
			}
		return true
		}
	go_To_Definition()
		{
		if false isnt loc = .getSelected()
			.goTo(loc)
		}

	goTo(selected, line = false)
		{
		svcTable = SvcTable(selected.Table)
		if false is rec = svcTable.Get(selected.Location)
			return
		if svcTable.Type is 'lib'
			{
			if line is false
				line = rec.lib_current_text.LineFromPosition(selected.Pos) + 1
			GotoLibView(selected.Location, libs: Object(selected.Table), :line)
			}
		else
			OpenBook(selected.Table, selected.Location, bookedit?:)
		}

	VirtualList_BuildContextMenu(rec)
		{
		menu = Object()
		table = rec isnt false ? rec.Table : ''
		if LibraryTables().Has?(table)
			menu.Add('Go To Definition', 'Find References', 'Version History', '')
		else if BookTables().Has?(table)
			{
			if table is 'suneidoc'
				menu.Add('Go To Documentation')
			menu.Add('Edit Documentation', '')
			}
		menu.Add('Expand All', 'Contract All')
		return menu
		}

	On_Context_Go_To_Definition()
		{
		.go_To_Definition()
		}
	On_Context_Find_References()
		{
		if false isnt loc = .getSelected()
			FindReferencesControl(loc.Location)
		}
	On_Context_Version_History()
		{
		if false isnt loc = .getSelected()
			VersionHistoryControl(loc.Table, loc.Location)
		}

	On_Context_Go_To_Documentation()
		{
		if false isnt loc = .getSelected()
			GotoUserManual(loc.Location)
		}

	On_Context_Edit_Documentation()
		{
		if false isnt loc = .getSelected()
			OpenBook(loc.Table, loc.Table $ loc.Location, bookedit?:)
		}

	getSelected()
		{
		sel = .list.GetSelectedRecords()
		return sel.Size() is 1 and sel[0].GetDefault('Table', '') isnt ''
			? sel[0]
			: false
		}

	VirtualList_Expand(rec)
		{
		if rec.Found is ''
			return Object(ctrl: #(Vert), rows: 0)
		if '' is all = FindReferences.AllOccurrences(
			rec.Table, rec.Location, .basename, .name)
			return Object(ctrl: #(Vert), rows: 0)
		height = all.LineCount() + 1
		rec.all = all
		return Object(ctrl: Object('Record',
			Object('WorkSpaceOutput', margin: 0, readonly:, :height, name: 'all')),
			rows: height)
		}

	VirtualList_AddGlobalMenu?()
		{
		return false
		}

	On_Context_Expand_All()
		{
		.list.ExpandByField(#(''), 'field')
		}

	On_Context_Contract_All()
		{
		.list.ExpandByField(#(''), 'field', collapse?:)
		}

	prevFocus: false
	Scintilla_KillFocus(source)
		{
		.prevFocus = source
		}

	Scintilla_SetFocus(source)
		{
		if .prevFocus isnt false and .prevFocus isnt source
			.prevFocus.SetSelect(0)
		}

	Scintilla_DoubleClick(source)
		{
		.list.SelectExpandRecord(source)
		if false isnt selRec = .list.GetSelectedRecord()
			.goTo(selRec, line: Number(source.GetLine().BeforeFirst(':')))
		}

	Destroy()
		{
		UserSettings.Put(.settings_key, Object(exclude_tests: .exclude_tests_ctrl.Get(),
			exclude_libs: .exclude_libs_ctrl.Get()), user: 'default')
		super.Destroy()
		}
	}