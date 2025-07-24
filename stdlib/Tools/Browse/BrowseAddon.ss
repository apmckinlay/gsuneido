// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Addon
	{
	CrossTable(query, title)
		{
		if .save()
			ToolDialog(.Parent.Window.Hwnd, Crosstab(query, title))
		}
	Reporter()
		{
		if .save()
			Reporter()
		}
	Summarize()
		{
		if .save()
			SummarizeControl(.Parent, .Parent.GetColumnsSaveName())
		}
	Export(query)
		{
		if .save()
			GlobalExportControl(query)
		}
	save()
		{
		if 0 is .Send("Access_Save")
			if not .Save()
				return false
		return true
		}
	Print(default_title, query)
		{
		if 0 is title = .Send('Browse_GetTitle')
			title = default_title isnt false ? default_title : ''

		cols = .GetVisibleColumns()
		if cols.Empty?()
			cols = .GetColumns()

		// the Delete(0) is to remove the 'listrow_delete' column
		ReportFromQueryAndColumns(
			.getSortedQuery(query), cols.Delete(0), title, .Parent.Window.Hwnd)
		}
	getSortedQuery(query)
		{
		sort = .GetSort()
		if sort is '' or
			not QueryColumns(QueryStripSort(query)).
				Has?(sort.RemovePrefix('reverse '))
			return query
		return QueryStripSort(query) $ ' sort ' $ sort
		}
	Go_To_QueryView(list, query, linkField, value, keyQueryFn)
		{
		sel = .GetSelection()
		if linkField isnt false
			query = QueryAddWhere(query, " where " $ linkField $ " = " $ Display(value))
		if sel.Size() > 0
			{
			rec = list.Get()[sel[0]]
			if rec.listrow_deleted isnt true and rec.Browse_NewRecord isnt true
				query = keyQueryFn(QueryStripSort(query), rec.Browse_OriginalRecord)
			}
		GotoQueryView(query)
		}
	cut_record: false
	CutRecord()
		{
		sel = .GetSelection()
		if (sel.Size() isnt 1)
			{
			Alert("Please select a single record to cut",
				title: 'Context Cut Record', flags: MB.ICONERROR)
			return
			}
		record = .GetRow(sel[0])
		if (record.listrow_deleted is true)
			{
			Alert("Cutting a deleted record is not allowed",
				title: 'Context Cut Record', flags: MB.ICONERROR)
			return
			}
		.cut_record = record.Copy()
		.Send("Browse_CutRecord", .cut_record)
		.DeleteSelection()
		}
	PasteRecord(before = true)
		{
		if (.cut_record is false)
			return
		sel = .GetSelection()
		pasteAt = sel is #() ? false : sel[0]
		if (before is false and pasteAt isnt false)
			++pasteAt
		newrec = .AddRecord(.cut_record, pasteAt)
		.Send("Browse_PasteRecord", newrec)
		.cut_record = false
		}
	SetAttachmentsManager(query, keyField)
		{
		.attachmentsMgr = AttachmentsManager(query, keyField)
		}
	Restore(recQuery, observer_ListRow)
		{
		sel = .GetSelection()
		if (sel is #())
			return
		row = sel[0]
		record = .GetRow(row)

		if (record.Member?("Browse_RecordDirty") and
			record.Browse_NewRecord isnt true)
			{
			if record.listrow_deleted is true
				{
				.AlertInfo('Context Restore', 'The record is marked for ' $
					'deletion. Please use Undelete first.')
				return
				}
			record = Query1(recQuery(record))
			if record is false
				{
				.AlertInfo('Context Restore', "The current record has been deleted.")
				return
				}
			record.Browse_OriginalRecord = record.Copy()
			record._browse = this
			.SetRecordHeaderData(record)
			record.Observer(observer_ListRow)
			.SetRow(row, record)
			.Send("Browse_Restore", record)
			.attachmentsMgr.RestoreOneByKey(record)
			}
		.Refresh_BrowseData()
		}
	RestoreAll(query, columns)
		{
		if YesNo("All of your current changes will be lost. " $
			"Are you sure you want to proceed?", "Restore All?", .Parent.Window.Hwnd,
			MB.ICONWARNING)
			{
			.attachmentsMgr.ProcessQueue(restore?:)
			.SetQuery(query, columns)
			.Send("Browse_RestoreAll")
			}
		}
	CleanupAttachments(restore? = false)
		{
		.attachmentsMgr.ProcessQueue(restore?)
		}
	QueueDeleteAttachmentFile(newFile, oldFile, rec, name, action)
		{
		.attachmentsMgr.QueueDeleteFile(newFile, oldFile, rec, name, action)
		}
	DeleteRecordAttachments(rec)
		{
		.attachmentsMgr.QueueDeleteRecordFiles(rec)
		}
	DeleteNewRecordAttachments(rec)
		{
		.attachmentsMgr.DeleteNewRecordFiles(rec)
		}
	LogColumnsSaveWarning(title)
		{
		warn = 'WARNING: Browse Control is using query as ' $
			'user columns key: ' $ title
		SuneidoLog.Once(warn)
		}
	RecordConflict?(record, cur, query_columns, linkField, quiet? = false)
		{
		if cur is false
			{
			.alertNoOriginalRecord(linkField)
			return true
			}

		return RecordConflict?(record.Browse_OriginalRecord, cur,
			query_columns, .Parent.Window.Hwnd, :quiet?)
		}
	alertNoOriginalRecord(linkField)
		{
		msg = "Browse: can't get record to update.\n" $
			"Another user may have deleted the line.\n"
		if linkField isnt false
			msg $= "Please use Current > Restore and " $
				"re-do your changes if necessary."
		AlertDelayed(msg, title: 'Error', hwnd: .Parent.Window.Hwnd,
			flags: MB.ICONERROR, uniqueId: "BrowseNoOriginalRecord")
		}
	ButtonBar()
		{
		return ButtonBar(#(Restore, 'Delete/Undelete'))
		}
	contextMenu: #("Reason Protected", "", "Customize Columns...",
		"Customize...", "Go To QueryView")
	readonlyMenu: #("Reason Protected", "", "Customize Columns...")
	standaloneMenu: #("Print...", "Reporter...", "Summarize...", "CrossTable...",
		"", "Save All", "Restore All", "", "Export...")
	readonlyStandaloneMenu: #("Print...", "Reporter...", "Summarize...", "CrossTable...",
		"", "Export...")
	BuildMenu(linkField, columnsSaveName, readOnly)
		{
		.menu = readOnly ? .readonlyMenu.Copy() : .contextMenu.Copy()
		if linkField is false
			.menu.Add(@(readOnly ? .readonlyStandaloneMenu : .standaloneMenu))
		if columnsSaveName is false
			.menu = .menu.Difference(#('Customize Columns...', 'Customize...'))
		if Suneido.User isnt 'default'
			.menu.Remove("Go To QueryView")
		}
	SetMenu(.menu)
		{
		}
	HeaderMenu(ctrl, x, y)
		{
		menus = .menu.Copy().Add("", "Reset Columns", at: .menu.Size())
		formatMenu =.addFormatMenu(ctrl, x, y)
		ContextMenu(menus.Append(formatMenu)).ShowCall(ctrl, x, y)
		}
	ListMenu(ctrl, x, y, readOnly)
		{
		customMenu = .buildCustomizableListMenu(readOnly, ctrl)
		editingMenu = Object()
		if readOnly is false
			.addEditingMenu(editingMenu)
		menus = customMenu.Append(editingMenu).Append(.menu.Copy())
		formatMenu =.addFormatMenu(ctrl, x, y)
		ContextMenu(menus.Append(formatMenu)).ShowCall(ctrl, x,y)
		}
	buildCustomizableListMenu(readOnly, ctrl)
		{
		customMenu = Object() //.menu.Copy()
		sendMessage = readOnly is false
			? "Browse_AddToContextMenu"
			: "Browse_AddToReadOnlyContextMenu"
		// you have to redefine browse to use custom context messages
		if 0 isnt ob = .Send(sendMessage)
			for option in ob
				customMenu.Add(option)

		for m in ctrl.Addons.Collect('CurrentMenu')
			customMenu.Add(@m)

		if customMenu.NotEmpty?()
			customMenu.Add("")
		return customMenu
		}
	addEditingMenu(new_menu)
		{
		.addToMenu(new_menu, Object("New", ""), 0)
		selected = .GetSelection()
		if selected isnt #()
			{
			options = Object("Delete/Undelete")
			if .Send("Browse_PreventRestore") is 0
				options.Add("Restore")
			.addToMenu(new_menu, options, 1)
			if (selected.Size() is 1)
				{
				if (.Send("Browse_AllowCutPaste") is true)
					new_menu.Add("Cut Record", "Paste Record Before",
						"Paste Record After", at: 4)
				.addToMenu(new_menu, Object("Edit Field"), 0)
				}
			}
		}
	addToMenu(new_menus, ob, pos)
		{
		for menu in ob.Reverse!()
			new_menus.Add(menu, at: pos)
		}

	addFormatMenu(ctrl, x, y)
		{
		pt = Object(:x, :y)
		ScreenToClient(ctrl.GetList().Hwnd, pt)
		col = Min(.GetColFromX(pt.x), .GetNumCols() - 1)
		fld = .GetCol(col)
		fmt = Datadict(fld).Format[0]
		if not fmt.Suffix?('Format')
			fmt $= 'Format'
		fmt = Global(fmt)
		extra = Object()
		if fmt.Method?('List_ExtraContext')
			if false isnt contextExtra = fmt.List_ExtraContext()
				extra = Object("", contextExtra)
		return extra
		}
	}