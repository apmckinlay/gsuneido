// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
/*
SvcControl handles everything related to the user interface, except for what is
encapsulated in SvcDisplayControl. It calls the appropriate classes/methods for any
actions that happen in the Version Control window.

It works hand-in-hand with SvcModel, as many of the actions call functions in Svc through
that class.

Notable methods:

set_table() is called every time a library is selected, and it calls SvcModel to
reinstantiate the Svc instance and adds any new changes to the lists.

List_selection() always passes in the name and library of the selection to
SvcDisplayControl, which the decides what code to display in the bottom window(s).

On_Get_Master_Changes() calls Svc to update the local library. On_Get_All_Master_Changes()
first calls On_Get_Master_Changes() using svc_all_libraries, and then calls it again for
each book individually. This is faster than calling On_Get_Master_Changes for each
library and book.

move_conflict() simply moves list items between the Conflicts and Master Changes windows.
It does not modify any code.
*/
Controller
	{
	Name: 'Svc'
	Title: 'Version Control'
	CallClass()
		{
		SvcSocketClient().RetryState()
		GotoPersistentWindow('SvcControl', SvcControl)
		}

	New(table = false)
		{
		.local_list = .FindControl('localList')
		.local_list.SetReadOnly(true, grayOut: false)
		.master_list = .FindControl('masterList')
		.master_list.SetReadOnly(true, grayOut: false)
		.table_list = .FindControl('table')
		.user = .FindControl('user')
		.display = .FindControl('svc_display')

		.model = new SvcModel()
		.setSettings()
		.display.SetModel(.model)

		.local_list.SetMultiSelect(true)
		if table isnt false
			.table_list.Set(table)
		.set_table(table)

		.subs = [
			PubSub.Subscribe('LibraryTreeChange', .treeChanged)
			PubSub.Subscribe('LibraryRecordChange', .runChecksFresh)
			PubSub.Subscribe('SvcSettings_ConnectionModified', .setSettings)
			PubSub.Subscribe('SvcSocketClient_StateChanged',
				{ .Defer(.On_Refresh, uniqueID: 'svccontrol_refresh') })
			]
		}

	Controls()
		{
		fixedymin = 230
		return Object('Vert',
			.headerHorz(),
			#(Skip 5)
			Object('VertSplit'
				Object('FixedYmin', fixedymin,
					Object('Horz',
						.localVert(),
						.masterVert(),
						ystretch: 1))
				// use Vert to allow changing control
				Object('FixedYmin', fixedymin, "SvcDisplay")
				)
			)
		}

	skipChecks?: false
	runChecksFresh()
		{
		if .skipChecks?
			return

		SvcCommitChecker.ClearPreCheck()
		.local_list.Repainter = false
		SvcGetMaster(.curtable, .local_list, .settings)
		}

	listSeperator: '*'
	headerHorz()
		{
		hdrCtrls = Object(#Toolbar, #Refresh, #Test_Runner, #Get_All_Master_Changes).
			RemoveIf({ OptContribution('SvcExcludeControls', #()).Has?(it) })
		return Object('Horz'
			Object('Vert'
				#(EtchedLine before: 0)
				Object('Horz'
					'Skip'
					#(Pair
						(Static User)
						(Field name: user width: 15)
						name: 'UserPair')
					'Skip', 'Skip',
					Object('Pair'
						#(Static 'Library/Book')
						Object('ChooseList', .tables(), width: 17, name: 'table',
							listSeparator: .listSeperator)
						name: 'TablePair')
					'Skip', 'Skip'
					)
				xstretch: 0)
			hdrCtrls,
			'SvcSettingsIcons')
		}
	localVert()
		{
		return Object('Vert'
			#(Horz
				Fill
				(Button 'Send Checked Local Changes...')
				Fill)
			#(ListStretch columns: #(svc_checked, svc_lib, svc_type, svc_date,
				svc_local_date, svc_warning, svc_name),
				noShading:,	name: 'localList', defWidth: false,
				columnsSaveName: 'svc_local', stretchColumn: 'svc_name',
				checkBoxColumn: 'svc_checked')
			#(Skip 3)
			#(HorzEqual pad: 0
				Skip
				(Button 'All')
				Skip
				(Button 'None')
				Fill Skip
				(Button 'Export' tip: 'Export checked records')
				Skip
				(Button 'Copy To' tip: 'Copy checked records to another library')
				Fill Skip
				(Button 'Compare'
					tip: 'Compare local copy to version control master (Alt+C)')
				Skip)
			xmin: 400, xstretch: 1,
			name: 'Local')
		}

	masterVert()
		{
		return #(Vert
			(Horz Fill (Button 'Get Master Changes') Fill)
			(ListStretch columns: #(svc_lib, svc_type, svc_who,
				svc_master_date, svc_local_date, svc_name)
				noShading:, defWidth: false, name: 'masterList',
				columnsSaveName: 'svc_master', stretchColumn: 'svc_name')
			xmin: 400, xstretch: 1,
			name: 'Master')
		}

	getListSelected(list)
		{
		return list.GetSelection().GetDefault(0, false)
		}

	tables()
		{
		return SvcSettings.Set?()
			? .formattedLibraryNames().Add('').Add(@.formattedBookNames())
			: []
		}

	formattedLibraryNames()
		{
		libraryTables = .libraryTables()
		libraries = libraryTables.Copy().
			Filter({ it.BeforeFirst(.listSeperator) isnt '' and it isnt .allLibAlias })
		if .formatTableList(libraries, libraryTables, SvcDisabledLibraries())
			libraryTables[0] = .allLibAlias $ .listSeperator
		return libraryTables
		}

	libraryTables()
		{
		newLibs = .addHeading('Local Libraries', SvcDisabledLibraries())
		usedLibs = .addHeading('Used Libraries', .svcLibraries.Copy()).Remove(@newLibs)
		unusedLibs = .addHeading('Unused Libraries',
			.svcLibraryTables.Copy().Remove(@usedLibs).Remove(@newLibs).SortWith!(#Lower))
		libraries = [.allLibAlias].Add('').Add(@usedLibs)
		if newLibs.NotEmpty?()
			libraries.Add('').Add(@newLibs)
		if unusedLibs.NotEmpty?()
			libraries.Add('').Add(@unusedLibs)
		return libraries
		}

	addHeading(heading, list)
		{
		heading = .listSeperator $ ' ' $ heading $ ' ' $ .listSeperator
		if list.NotEmpty?()
			list.Add(heading, at: 0)
		return list
		}

	getter_svcLibraryTables()
		{
		return .svcLibraryTables = .SvcLibraryTables()
		}

	getter_svcLibraries()
		{
		return .svcLibraries = Libraries().Difference(.SvcExcludeLibraries)
		}

	Getter_SvcExcludeLibraries()
		{
		return #(configlib, Test_lib)
		}

	SvcLibraryTables()
		{
		return LibraryTables().Difference(.SvcExcludeLibraries)
		}

	formatTableList(tablesToCheck, tableList, localTables)
		{
		formatted? = false
		for table in tablesToCheck
			if .hasChanges(table) isnt false
				{
				if not localTables.Has?(table)
					formatted? = true
				tableList.Replace(table, table $ .listSeperator)
				}
		return formatted?
		}

	formattedBookNames()
		{
		bookTables = .bookTables()
		books = bookTables.Copy().Filter({ it.BeforeFirst(.listSeperator) isnt '' })
		.formatTableList(books, bookTables, SvcDisabledBooks())
		return bookTables
		}

	bookTables()
		{
		newBooks = .addHeading('Local Books', SvcDisabledBooks())
		books = .addHeading('Books', BookTables()).Remove(@newBooks)
		if newBooks.NotEmpty?()
			books.Add('').Add(@newBooks)
		return books
		}

	setTableList()
		{
		tableName = .table_list.Get()
		.table_list.SetList(.tables())
		if .table_list.GetList().Has?(tableName)
			.table_list.Set(tableName)
		}

	Commands()
		{
		return #(
			#(Refresh, "F5"),
			#(Close, "", "Close the Window"),
			#(Users_Manual, "F1"),
			#(Test_Runner, "Alt+T", "Open a Test Runner window", "T"),
			#(Get_All_Master_Changes, "Alt+G", "", "G"),
			#(Compare, "Alt+C", "Compare local copy to version control master", "C"))
		}

	Menu:
		(
		("&File",
			"&Settings...",
			"&Refresh",
			"Compare...",
			"",
			"&Get All Master Changes",
			"",
			"&Close")
		)

	curtable: ''
	NewValue(value, source)
		{
		if source.Name isnt 'table'
			return

		value = source.Valid?() ? value : ''
		tableName = .tableName(value)
		if tableName is .curtable
			return

		.set_table(value)
		.display.Reset()

		SvcGetMaster(tableName, .local_list, .settings)
		}

	tableName(table)
		{
		return table isnt false
			? table is .allLibAlias
				? .allLibView
				: table
			: ''
		}

	asof: false
	set_table(table)
		{
		.curtable = .tableName(table)
		.local_list.Clear()
		.master_list.Clear()
		.model.Clear()
		.curSelection = false
		table = table is false ? '' : table
		.Window.SetTitle(.Title = 'Version Control' $ Opt(' - ', table.Trim('()')))

		if table is false or table is ''
			{
			.sortLocalList()
			return
			}

		.model.SetTable(.curtable)
		.asof = .model.SvcTime()
		.createLists()
		}

	createLists()
		{
		for rec in .model.LocalChanges
			.local_list.AddRow(.buildRow(rec, .localFields))
		.sortLocalList()
		for rec in .model.MasterChanges
			.master_list.AddRow(.buildRow(rec, .masterFields))

		for rec in .model.Conflicts
			{
			rec.type = '%'
			.local_list.AddRow(.buildRow(rec, .localFields))
			.master_list.AddRow(.buildRow(rec, .masterFields))
			}
		}

	getter_sort()
		{
		return .sort = UserSettings.Get('VersionControl-SortLocal')
		}

	sortLocalList()
		{
		if .sort is false or .sort is ''
			return
		sortField = .sort.RemovePrefix('reverse ')
		if false is sortIndex = .local_list.GetColumns().Find(sortField)
			return
		sortIndex += 1
		if .sort.Prefix?('reverse ')
			sortIndex = -sortIndex
		.local_list.SetSortCol(sortIndex)
		}

	List_AfterSort()
		{
		.sort = .local_list.GetSort(nonMarkExtraCol?:)
		}

	baseFields: #('svc_name': 'name', 'svc_lib': 'lib')
	localFields: #('svc_type': 'type', 'svc_date': 'modified',
		'svc_local_date': 'committed')
	masterFields: #('svc_type': 'type', 'svc_who': 'who', 'svc_master_date': 'modified',
		'svc_local_date': 'committed')
	conflictFields: #('svc_date': 'localModified', 'svc_who': 'who',
		'svc_local_date': 'committed', 'svc_master_date': 'modified')
	buildRow(rec, fields = #())
		{
		ob = Object()
		fields = .baseFields.Copy().Merge(fields)
		for field in fields.Members()
			ob[field] = rec[fields[field]]
		ob.svc_checked = .previousSelected.Has?(Object(name: rec.name, lib: rec.lib))
		return ob
		}

	On_Settings()
		{ SvcSettings(openDialog:) }

	On_Test_Runner()
		{
		TestRunnerGui()
		}
	On_All()
		{ .checkAll(true) }
	On_None()
		{ .checkAll(false) }
	checkAll(check?)
		{
		for row in .local_list.Get()
			row.svc_checked = check?
		.local_list.Repaint()
		}
	toggleCheck(list, row)
		{
		data = list.GetRow(row)
		data.svc_checked = data.svc_checked isnt true
		list.RepaintRow(row)
		}
	List_SingleClick(row, col, source)
		{
		if row is false
			{
			.local_list.ClearSelectFocus()
			.master_list.ClearSelectFocus()
			return 0
			}

		if 'svc_checked' is source.GetCol(col)
			.toggleCheck(source, row)
		return 0
		}
	List_DoubleClick(row, col, source)
		{
		if row is false
			return 1

		list = source.Name is 'localList'
			? .local_list
			: .master_list

		rec = list.GetRow(row)
		lib = rec.svc_lib
		name = rec.svc_name
		line = .display.GetGoToLine()

		msg = .getWarnings()['msgMap'][lib $ '_' $ name]
		if msg isnt #() and list.GetCol(col) is 'svc_warning'
			.AlertInfo('SVC Checking Status', msg)
		else
			GoToDefinition(name, lib, line)

		return 1
		}
	curSelection: false
	List_Selection(selection, source)
		{
		if selection is false
			return 0
		.curSelection = selection.Copy()
		.curSource = source.Copy()
		selection = selection[0]

		.master_list.ClearHighlight()
		.local_list.ClearHighlight()

		if source.Name is 'masterList'
			.masterListSelection(selection)
		else
			.localListSelection(selection)
		return 0
		}

	masterListSelection(selection)
		{
		.local_list.ClearSelectFocus()
		sel = .master_list.GetRow(selection)

		.highlightSelected(.local_list, sel)
		if .master_list.Get().CountIf(
			{ it.svc_name is sel.svc_name and it.svc_lib is sel.svc_lib }) > 1
			.display.Display(sel.svc_name, sel.svc_lib, sel.svc_type, showComment:,
				masterNewer?:, lib_committed: sel.svc_master_date)
		else
			.display.Display(sel.svc_name, sel.svc_lib, sel.svc_type
				showComment:, masterNewer?:)
		.master_list.SetFocus()
		}

	localListSelection(selection)
		{
		.master_list.ClearSelectFocus()
		// Since selecting an item refreshes, modified date is updated as well
		sel = .local_list.GetRow(selection)

		for idx in .local_list.GetSelection()
			{
			item = .local_list.GetRow(idx)
			if item.svc_type isnt '%'
				continue
			.highlightSelected(.master_list, item)
			}

		changeOb = sel.svc_type is '%' ? .model.Conflicts : .model.LocalChanges
		sel.svc_date = .model.
			UpdateLocalModified(sel.svc_name, sel.svc_lib, changeOb, sel.svc_date)
		if sel.svc_type is '%'
			{
			.local_list.SetRow(selection, sel)
			.display.Display(sel.svc_name, sel.svc_lib, '%', showComment:,
				masterNewer?:)
			.local_list.SetFocus()
			}
		else
			{
			if .refreshRequired(sel)
				{
				.Defer(.On_Refresh)
				return
				}
			.local_list.SetRow(selection, sel)
			// replace type because .display.Display is designed for master changes
			.display.Display(sel.svc_name, sel.svc_lib, sel.svc_type.Tr("+-", "-+"))
			}
		}

	highlightSelected(list, item)
		{
		if false isnt x = .getRowFromList(list, item.svc_name, item.svc_lib)
			{
			list.ScrollRowToView(x)
			list.AddHighlight(x)
			}
		}

	refreshRequired(sel)
		{
		// If a record no longer exists, refresh the list to remove it
		rec = .model.GetLocalRec(sel.svc_lib, sel.svc_name, deleted: sel.svc_type is '-')
		if rec is false
			return true
		// If the record has no actual changes, refresh the list to remove it
		if sel.svc_date is '' and sel.svc_type isnt '-' and sel.svc_type isnt '+'
			return true
		return false
		}

	ResetSelection()
		{
		.List_Selection(.curSelection, .curSource)
		}
	Activate()
		{
		.refreshIfCurrentChanged()
		}
	refreshIfCurrentChanged()
		{
		if .curSelection is false or .curSource is false
			{
			.checkTreeChanged()
			return
			}

		sourceName = .curSource.Name
		sel = .curSource.GetRow(.curSelection[0])
		if .checkTreeChanged()
			.reselectCurrent(sourceName, sel)
		else if sel.svc_type is '%'
			{
			// Only refresh conflict records if the local record is updated
			localIdx = .local_list.Get().
				FindIf({ it.svc_name is sel.svc_name and it.svc_lib is sel.svc_lib })
			if localIdx isnt false
				.refreshCurrent(.local_list.GetRow(localIdx), [localIdx], .local_list)
			}
		else if sourceName isnt 'masterList'
			.refreshCurrent(sel, .curSelection, .curSource)
		}

	checkTreeChanged()
		{
		if not .treeChanged?
			return false

		.On_Refresh()
		return true
		}

	reselectCurrent(sourceName, sel)
		{
		if false is list = .FindControl(sourceName)
			return
		if false isnt idx = list.Get().
			FindIf({ it.svc_name is sel.svc_name and it.svc_lib is sel.svc_lib })
			list.SetSelection(idx)
		}
	getRowFromList(list, name, lib)
		{
		return list.Get().FindIf(
			{ it.svc_name is name and it.svc_lib is lib })
		}

	refreshCurrent(sel, selections, source)
		{
		rec = .model.GetLocalRec(sel.svc_lib, sel.svc_name, deleted: sel.svc_type is '-')
		if rec is false
			.On_Refresh()
		else if rec isnt false and rec.lib_modified isnt sel.svc_date
			.List_Selection(selections, source)
		}

	List_ContextMenu(x, y, source)
		{
		if source.Name isnt 'localList'
			return 0

		ContextMenu(#("Go To Definition", "Export Record", "Find References",
			"Version History", "", "Restore", #(Restore))).ShowCall(this, x, y)
		}

	On_Copy_To()
		{
		checked = .getLocalChecked().Map({ Object(lib: it.lib, name: it.name) })
		if checked.Empty?()
			{
			.info("Please checkmark the records to copy")
			return
			}
		if false is SvcCopyRecordsControl(checked)
			return
		if .curtable is .allLibView
			.On_Refresh(skipChecks:)
		.skipChecks? = true
		.resetLibraries()
		.skipChecks? = false
		.runChecksFresh()
		}

	info(text)
		{
		InfoWindowControl(:text, autoClose: 2, titleSize: 0)
		}

	On_Context_Restore()
		{
		if false is changes = .getHighlightedCheckNotEmpty()
			return

		needsReset? = false
		for change in changes
			{
			name = change.name

			if change.type is '-' or change.type is '+'
				needsReset? = true

			if not .model.Restore(name, change.lib, change.type)
				{
				Alert("Can't restore " $ name, title: 'Restore', flags: MB.ICONERROR)
				return
				}

			.previousSelected = #()
			Print('Restored', change.lib $ ':' $ name)
			}
		.On_Refresh(skipChecks:)
		.skipChecks? = true
		if needsReset?
			.resetLibraries()
		.skipChecks? = false
		.runChecksFresh()
		}

	getHighlightedCheckNotEmpty()
		{
		changes = .getHighlighted()
		if changes.Empty?()
			{
			.info("Please highlight a local change")
			return false
			}
		return changes
		}

	On_Context_Version_History()
		{
		if false is changes = .getHighlightedCheckNotEmpty()
			return
		for change in changes
			VersionHistoryControl(change.lib, change.name)
		}

	On_Merge()
		{
		.move_conflict('#')
		}
	On_Use_Master()
		{
		.move_conflict()
		.setTableList()
		}
	move_conflict(prefix = ' ')
		{
		result = .getSelectedRecs()
		del = Object()
		for x in result
			{
			rec = x.rec

			// no selected record
			if rec is false
				return

			recAdded = .model.MoveConflict(rec.svc_name, rec.svc_lib, prefix is '#')
			.master_list.DeleteRows(.getRowFromList(.master_list, rec.svc_name,
				rec.svc_lib))
			recAdded.Each()
				{
				.master_list.AddRow(Object(
					svc_type: it.type
					svc_name: rec.svc_name,
					svc_who: it.type is '#' ? '' : it.who,
					svc_lib: rec.svc_lib,
					svc_master_date: it.modified,
					svc_local_date: rec.svc_local_date))
				}
			del.Add(.getRowFromList(.local_list, rec.svc_name, rec.svc_lib))
			}
		.local_list.DeleteRows(@del)
		.local_list.ClearSelectFocus()
		.master_list.ClearSelectFocus()
		.display.Reset()
		.curSelection = false
		}
	getSelectedRecs()
		{
		recs = Object()
		for list in Object(.local_list, .master_list)
			{
			if ((selection = list.GetSelection()).Empty?())
				continue

			for idx in selection
				recs.Add(Object(rec: list.GetRow(idx), i: idx))
			}
		return recs
		}

	On_Get_Master_Changes(skipChecks = false)
		{
		if .curtable is ''
			{
			.info("Please choose a Library/Book")
			return 0
			}

		if not .model.Conflicts.Empty?()
			{
			.selectAndHighlightConflict(.model.Conflicts)
			.AlertInfo(.Title, "Please resolve conflicts before getting master changes")
			return 0
			}

		if .model.MasterChanges.Empty?()
			return 0

		if .curtable is .allLibView
			return .getAllMasterChanges(skipChecks)

		updates = .masterChanges()
		.On_Refresh(skipChecks:)
		.skipChecks? = true
		.resetLibraries()
		.skipChecks? = false
		if not skipChecks
			.runChecksFresh()
		return updates
		}

	selectAndHighlightConflict(conflicts)
		{
		// Highlight first conflict
		conflict = conflicts[0] // grab first conflict
		if false isnt row = .getRowFromList(.local_list, conflict.name ,conflict.lib)
			{
			.local_list.ScrollRowToView(row)
			.local_list.SetSelection(row)
			}
		}

	treeChanged?: 	false
	treeChanged()
		{
		.treeChanged? = true
		// Reset getters on tree change (removing old / adding new libraries)
		.svcLibraryTables = .SvcLibraryTables()
		.svcLibraries = Libraries().Difference(.SvcExcludeLibraries)
		.runChecksFresh()
		}

	resetLibraries()
		{
		ResetCaches()
		.treeChanged? = false
		}

	setSettings()
		{
		.model.SetSettings(.settings = SvcSettings.Get())
		if .settings isnt false and .settings.svc_user isnt '' and .user.Get() is ''
			.user.Set(.settings.svc_user)
		.On_Refresh()
		}

	getAllMasterChanges(skipChecks)
		{
		libraries = Object()
		.model.MasterChanges.Each({ libraries.AddUnique(it.lib) })

		n = 0
		for lib in libraries
			n += .model.UpdateLibrary(lib)
		.On_Refresh(skipChecks:)
		.skipChecks? = true
		.resetLibraries()
		.skipChecks? = false
		if not skipChecks
			.runChecksFresh()
		.PostGetChanges()
		return n
		}

	PostGetChanges()
		{
		Contributions('Svc_PostGetChanges').Each()
			{
			try
				it(.curtable)
			catch (e)
				SuneidoLog('ERROR: SVC Post Get Changes encountered: ' $ e, params: [it])
			}
		}

	masterChanges()
		{
		update = .model.UpdateLibrary(.curtable)
		.PostGetChanges()
		return update
		}

	GetModel() // sent by Addon_overwrite_lines
		{
		return .model
		}

	allLibView: 'svc_all_changes'
	allLibAlias: 'All VC Libraries'
	On_Get_All_Master_Changes()
		{
		if .table_list.GetList().Empty?()
			{
			.info("Please fill File > Settings prior to getting All Changes")
			return
			}
		table = .tableName(.table_list.Get())
		allLibs? = table is .allLibView
		conflicts = Object()
		tosend = Object()

		.table_list.Set(.allLibAlias)
		.set_table(.allLibView)
		changes = .model.GetChanges(.curtable)
		conflicts.MergeUnion(changes.conflicts)
		tosend.MergeUnion(changes.local_changes.Map({ it.lib }))
		.On_Get_Master_Changes(skipChecks:)

		if not conflicts.Empty?()
			return

		.getBookChanges(BookTables(), conflicts, tosend)
		table = .checkOutstanding(conflicts, allLibs?, tosend, table)
		if allLibs?
			table = .allLibAlias
		.set_table(table)
		.table_list.Set(table)
		.On_Refresh(skipChecks:)
		.runChecksFresh()
		}

	getBookChanges(books, conflicts, tosend)
		{
		for lib in books
			{
			.set_table(lib)
			changes = .model.GetChanges(lib)
			if .model.Conflicts.Empty?()
				.masterChanges()
			else
				conflicts.Add(lib)
			if .model.SvcExists?(lib) and changes.local_changes.NotEmpty?()
				tosend.Add(lib)
			}
		}

	checkOutstanding(conflicts, allLibs?, tosend, table)
		{
		if not conflicts.Empty?()
			{
			if not allLibs?
				table = conflicts[0]
			.AlertError('Conflicts',
				'Please resolve conflicts in:\n\n' $ conflicts.Join('\n'))
			}
		if not tosend.Empty?()
			{
			tables = .tables().Map({ it.BeforeFirst(.listSeperator) })
			tosend.SortWith!({ tables.Find(it) })
			if not allLibs?
				table = tosend[0]
			Print("Changes to send in: " $ tosend.Join(', '))
			Alert("You have changes to send in:\n\n" $ tosend.Join('\n'),
				"Local Changes", flags: MB.ICONINFORMATION)
			}
		return table
		}

	nSent: 0
	On_Send_Checked_Local_Changes()
		{
		if not .checksPassed?()
			return

		if false is desc = SvcGetDescription(.Window.Hwnd, .changes, .warnings)
			return

		.local_list.ClearSelectFocus()

		Print(desc)

		if .curtable is .allLibAlias
			sendResult = .model.SendLocalChangesFromAll(.changes, desc, .user.Get(),
				.asof, .afterEachSend)
		else
			sendResult = .model.SendLocalChanges(.changes, desc, .user.Get(),
				.asof, .afterEachSend)

		if false is sendResult
			{
			Alert("Someone has sent new changes.\n" $
					"Please refresh and get the changes (and test with them)\n" $
					"before sending your changes.",
					title: 'Send Checked Local Changes',
					flags: MB.ICONINFORMATION)
			return
			}

		.On_Refresh(skipChecks:)
		.skipChecks? = true
		.resetLibraries()
		.skipChecks? = false
		.clearStatus()
		.updateListStatus()
		.PostSendChanges()
		}

	PostSendChanges()
		{
		Contributions('Svc_PostSendChanges').Each()
			{
			try
				it(.curtable)
			catch (e)
				SuneidoLog('ERROR: SVC Post Send Changes encountered: ' $ e, params: [it])
			}
		}

	clearStatus()
		{
		changes = .getSentChanges()
		if changes.Empty?()
			return

		memVals = changes.Map({ it.lib $ '_' $ it.name }).Values()

		.clearError(memVals)

		.clearMsg(memVals)

		.clearOtherChecks(changes)
		}

	updateListStatus()
		{
		list = .local_list.Get()
		errOb = Suneido.SvcCommit_Warnings.errMap
		for mem in errOb.Members()
			{
			if false is row = list.FindIf({ mem.Suffix?(it.svc_lib $ '_' $ it.svc_name) })
				continue

			list[row].svc_warning = errOb[mem]
			}
		.local_list.Repaint()
		}

	getSentChanges()
		{
		return .changes
		}

	clearError(memVals)
		{
		.handleClear('errMap')
			{ |m|
			memVals.Any?({ m.Suffix?(it) })
			}
		}

	clearMsg(memVals)
		{
		.handleClear('msgMap')
			{
			memVals.Has?(it)
			}
		}

	handleClear(map, block)
		{
		warnings = .getWarnings()
		mems = warnings[map].MembersIf(block)
		warnings[map].DeleteIf({ mems.Has?(it) })
		}

	clearOtherChecks(sentChanges)
		{
		warnings = .getWarnings()
		for mem in warnings[.getCurTable()].Members()
			{
			warnings[.getCurTable()][mem].RemoveIf(
				{ |rec| sentChanges.Any?({ it.Project(#(lib, name)) is
					rec.Project(#(lib, name)) }) })
			}
		}

	getWarnings()
		{
		return Suneido.SvcCommit_Warnings
		}

	getCurTable()
		{
		return .curtable
		}

	warnings: ''
	checksPassed?()
		{
		if .curtable is ''
			{
			.info('Please choose a Library/Book')
			return false
			}

		.changes = .getLocalChecked()
		if .changes.Empty?()
			{
			.info('Please checkmark the local changes to send')
			return false
			}

		if not .mandatory_checks(.changes)
			return false // abort send

		.warnings = ''
		if not .model.Library?(.curtable)
			return true

		if false is RetryBool(2 /*= maxRetries*/, 1000 /*= min sleep ms*/,
			{ .checkThreadStopped?() })
			{
			.alertThreadRunning()
			return false
			}

		return true
		}

	checkThreadStopped?()
		{
		return 'Thread_Running' isnt
			.warnings = SvcRunChecks.GetPreCheckResults(.curtable, changes: .changes)
		}

	alertThreadRunning()
		{
		.AlertInfo('Send Checked Local Changes', 'Checking is still running!\n' $
				'Please try again!')
		}

	afterEachSend(name, lib)
		{
		.local_list.DeleteRows(.local_list.Get().FindIf({
			it.svc_name is name and it.svc_lib is lib }))
		++.nSent
		}
	// Send mandatory checks ===================================================

	mandatory_checks(changes)
		{
		scc = new SvcCommitChecker()
		for check in [scc.MandatoryChecks, .have_user?]
			if '' isnt msg = check(:changes, model: .model, table: .curtable)
				{
				.AlertWarn("Send Changes", msg)
				return false // abort send
				}
		return true
		}

	have_user?()
		{
		return .user.Get() is ""
			? 'Please enter user(s)'
			: ''
		}

	//==========================================================================
	On_Context_Go_To_Definition()
		{
		if false is selected = .getListSelected(.local_list)
			return
		sel = .local_list.GetRow(selected)
		lib = sel.svc_lib
		line = .display.GetGoToLine()
		GoToDefinition(sel.svc_name, lib, line)
		}

	On_Context_Export_Record() // context menu
		{
		if false is selected = .getListSelected(.local_list)
			return

		sel = .local_list.GetRow(selected)
		if '' is table = sel.svc_lib
			return

		if '' isnt fileName = SaveFileName(hwnd: .Window.Hwnd,
			flags:	OFN.PATHMUSTEXIST | OFN.HIDEREADONLY | OFN.NOCHANGEDIR,
			title:	"Export (append) to")
			.exportOne(table, sel.svc_name, fileName, sel.svc_type is '-')
		}

	On_Context_Find_References() // context menu
		{
		if false is selected = .getListSelected(.local_list)
			return
		FindReferencesControl(.local_list.GetRow(selected).svc_name)
		}

	On_Export() // button
		{
		checked = .getLocalChecked()

		if checked.Empty?()
			{
			.info("Please checkmark the records to export")
			return
			}
		if '' is (fileName = SaveFileName(hwnd:	.Window.Hwnd,
			flags:	OFN.PATHMUSTEXIST | OFN.HIDEREADONLY | OFN.NOCHANGEDIR,
			title:	"Export (append) " $ checked.Size() $ " records to"))
			return
		for record in checked
			.exportOne(record.lib, record.name, fileName, record.type is '-')
		}

	getLocalChecked()
		{
		return .local_list.Get().Filter({ it.svc_checked is true }).Map({
			Object(name: it.svc_name, lib: it.svc_lib, type: it.svc_type )})
		}

	getHighlighted()
		{
		return .local_list.GetSelection().
			Map({ .local_list.GetRow(it) }).
			Map({ Object(name: it.svc_name, lib: it.svc_lib, type: it.svc_type )})
		}

	exportOne(table, name, fileName, delete = false)
		{
		path = false
		if name.Has?('/') // book record
			{
			path = name.BeforeLast('/')
			name = name.AfterLast('/')
			}
		LibIO.Export(table, name, fileName, path, :delete, interactive:)
		}
	previousSelected: #()
	On_Refresh(skipChecks = false)
		{
		.treeChanged? = false
		SvcSocketClient().RetryState()
		.setTableList()
		.previousSelected = .getLocalChecked().Copy().Map({
			Object(name: it.name, lib: it.lib) })
		if .curtable is ''
			{
			.asof = false
			.model.Clear()
			.table_list.Set('')
			}
		else if .table_list.Valid?()
			.set_table(.curtable is .allLibView ? .allLibAlias : .curtable)
		.display.Reset()
		.curSelection = false
		if not skipChecks
			.runChecksFresh()
		}

	On_Compare()
		{
		if .curtable is ''
			{
			.info('Please choose a Library to compare')
			return
			}
		if .curtable is .allLibView
			.info('Compare is not available when "' $ .allLibAlias $ '" is selected')
		else if .model.SvcCompare(.curtable)
			.On_Refresh()
		}

	Ok_to_CloseWindow?()
		{
		if .nSent is 0
			return true
		table = .tableName(.table_list.Get())
		tosend = .changesToSend()
		if table is .allLibAlias or tosend.Empty?() or tosend is [table]
			return true
		return YesNo("You still have changes to send in:\n\n" $
			'     ' $ tosend.Join(', ') $ '\n\n' $
			'Close anyway?')
		}

	changesToSend()
		{
		return .SvcLibraryTables().Filter(.hasChanges).
			MergeUnion(BookTables().Filter(.hasChanges))
		}

	hasChanges(table)
		{
		return not QueryEmpty?(SvcTable(table).ModifiedQuery() $ ' remove text')
		}

	Destroy()
		{
		if .sort isnt false and .sort isnt ''
			UserSettings.Put('VersionControl-SortLocal', .sort)
		.subs.Each(#Unsubscribe)
		if ViewExists?(.allLibView)
			Database('drop ' $ .allLibView)
		super.Destroy()
		}
	}
