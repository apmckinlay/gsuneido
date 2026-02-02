// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
CommandParent
	{
	Name: "ExplorerListView"
	view: false
	New(model, view, query = false, columns = false,
		.validField = false, readonly = false, .title = '',
		.stickyFields = #(), status = true, .linkField = false, protectField = false,
		.noShading = false, .noHeaderButtons = false, primary_accessobserver = false,
		.disableOptions = #(), .buttonBar = false, .extraButtons = #(),
		.columnsSaveName = false, .Commands = #(), .excludeFields = #(),
		.historyFieldsPrefix = false)
		{
		super(.makeControls(view, query, columns, status))
		.model_spec = model
		.model = Construct(model)
		.vertSplit = .FindControl('VertSplit')
		.view = .vertSplit.Scroll.Border.View
		Assert(.view.Base?(RecordControl), "view not a RecordControl")
		.view.AddObserver(.RecordChanged)
		.view.SetProtectField(.protectField = protectField)
		.query = .base_query = query isnt false ? query : .model.GetQuery()
		.sticky_values = Object()
		.list = .vertSplit.List
		.columns = columns.Copy()
		colsSaveName = .getColumnsSaveName()
		.vertSplit.SetSplitSaveName(colsSaveName, ' - Split')
		UserColumns.Load(.columns, colsSaveName, .list)
		.status = (status) ? .Vert.Status : class {
			SetValid(val /*unused*/ = true) {} Set(val /*unused*/) {} }
		.status.SetValid()
		if (readonly is true or ReadOnlyAccess(this))
			.SetReadOnly(true)
		.Send("Data")
		if (linkField isnt false)
			{
			if (primary_accessobserver is true)
				.Send("AccessObserver", .AccessChanged, 0)
			else
				.Send("AccessObserver", .AccessChanged)
			}
		else
			.Window.AddValidationItem(this)
		}

	getColumnsSaveName()
		{
		if .columnsSaveName isnt false
			return .columnsSaveName

		key = .title is '' ? .base_query : .title
		warn = 'WARNING: ExplorerListView Control is using query as ' $
			'user columns key: ' $ key
		SuneidoLog.Once(warn)
		return key
		}

	Startup()
		{
		super.Startup()
		if .linkField is false and Sys.SuneidoJs?()
			.Load_entries() // if done from New, list isnt scrolled to last item
		}

	first_resize: true
	Resize(x, y, w, h)
		{
		super.Resize(x, y, w, h)
		// load entries if standalone, otherwise Set does
		if .first_resize
			{
			.first_resize = false
			if .linkField is false
				.Load_entries() // if done from New, list isnt scrolled to last item
			}
		}
	NewModel(model)
		{
		.model = Construct(model)
		.query = .model.GetQuery()
		.Load_entries()
		// the following line is needed in case the Load_entries
		// does not load a new record (ex. readonly and no records in new model)
		if .list.Get().Empty?()
			.On_New(save?: false, force:)
		}
	Default(@args)
		{
		method = args[0]
		if .list.Method?(method)
			return .list[method](@+1 args)

		if method.Prefix?('On_Context_') and args.Member?('item')
			{
			.Send('ExplorerListView_ContextMenuItemClicked', args.item)
			return 0
			}
		throw 'method not found in ExplorerListViewControl: ' $ method
		}
	GetModel()
		{ return .model }
	GetView()
		{ return .view }
	GetList()
		{ return .list }
	makeControls(view, query, columns, status)
		{
		title = (.title isnt '') ? .title : (query isnt false) ? query : ""
		v = Object('Record', view.Copy())
		v.name = "View"
		controls = Object('Vert')
		if title isnt ''
			controls.Add(Object('CenterTitle', title))
		controls.Add(Object('VertSplit',
			Object('List', :columns, defWidth: false,
				noShading: .noShading, noHeaderButtons: .noHeaderButtons,
				resetColumns:, alwaysHighlightSelected:),
			Object('Scroll', Object('Border', v, border: 5))
			))

		// exclude Global->Save from the button bar
		// since it is not supported by ExplorerListView
		if .buttonBar
			{
			buttons = ButtonBar(#(Delete, Save, Restore, Print),
				globalExclude: #('Save...', 'Restore...'))
			if not .extraButtons.Empty?()
				buttons.MergeUnion(.extraButtons)
			controls.Add(buttons)
			}
		if (status)
			controls.Add('Status')
		return controls
		}
	update_explorerlist_data()
		{
		if (.linkField is false or .model is false)
			return
		data = .model.GetData()
		.Send("SetField", "explorerlist_data", data)
		.Send("InvalidateFields", "explorerlist_data")
		.Send("ExplorerListData_Changed", data)
		}
	Load_entries(restoreSelect = false)
		{
		prev = .selected
		.ignore_selection? = true
		entries = .model.GetEntries()
		.list.Set(entries)

		if (entries.Size() is 0)
			.On_New(false, setFocus: false)
		else
			.load_entry()
		.ignore_selection? = false
		if restoreSelect is true and prev isnt false and prev[0] < .list.GetNumRows()
			.list.SetSelection(prev[0])

		.Send("ExplorerListView_AfterLoadEntries")
		}
	HideColumn(column)
		{
		if (false isnt (col = .columns.Find(column)))
			.list.SetColWidth(col, 0)
		}
	RecordChanged(member)
		{
		.update_explorerlist_data()
		.validate()
		.list.AllowContextOnly(.invalid?)

		data = .view.Get()

		// refresh list
		if (.columns.Has?(member))
			if (false isnt (row = .list.Get().Find(data)))
				.list.RepaintRow(row)

		.Send("ExplorerListView_RecordChanged", member, data)
		}
	ignore_selection?: false
	On_New(save? = true, force = false, setFocus = true)
		{
		if (save?)
			.SaveRecord()
		if .allowInsert?(force)
			return false

		rec = .model.NewRecord()
		rec.Explorer_NewRecord = true

		if (.linkField isnt false)
			rec[.linkField] = .value
		rec.Merge(.sticky_values)

		.Send("ExplorerListView_AddRecord", rec)

		.view.Set(rec)
		.view.Dirty?(false)
		.invalid? = false
		.list.AllowContextOnly(false)
		last_row = .list.GetNumRows() - 1
		if (last_row > 0)
			.list.ScrollRowToView(last_row)
		.status.Set('')
		.status.SetValid(true)
		.ignore_selection? = true
		.list.AddRow(rec)
		.list.SetSelection(.list.GetNumRows() - 1)
		.ignore_selection? = false
		if setFocus is true
			.FocusFirst(.vertSplit.Scroll.Hwnd)
		.Send("ExplorerListView_AfterAddRecord", rec)
		return rec
		}
	allowInsert?(force)
		{
		rec = .linkField isnt false ? .model.NewRecord(add: false) : Record()
		return (.invalid? or (.readonly and not force) or
			false is AllowInsertRecord?(rec, .protectField))
		}
	// TODO: eliminate duplication between On_New and NewRecord
	NewRecord(output = false)
		{
		// if on new non-dirty record use it, otherwise get one from model
		rec = .view.Get()
		if (rec.Explorer_NewRecord isnt true or
			.view.Dirty?() or rec.Explorer_RecordUpdated is true)
			rec = .model.NewRecord()
		rec.Explorer_NewRecord = true

		if (.linkField isnt false)
			rec[.linkField] = .value

		if (output is true)
			{
			.model.Output(rec)
			rec.Delete("Explorer_NewRecord")
			}
		.ignore_selection? = true
		if (rec isnt .view.Get()) // only add if we got a new record from the model
			.list.AddRow(rec)
		.ignore_selection? = false

		return rec
		}
	UpdateRecord(rec)
		{
		rec.Explorer_RecordUpdated = true
		.model.Update(rec)
		.update_explorerlist_data()
		}
	On_Save()
		{
		rec = .view.Get()
		.SaveRecord()
		if false isnt row = .list.Get().Find(rec)
			.list.SetSelection(row)
		}
	/* Interface with ButtonBar */
	On_Current_Restore()
		{
		.On_Context_Restore()
		}
	On_Current_Delete()
		{
		.On_Context_Delete()
		}
	On_Current_Save()
		{
		.On_Context_Save()
		}
	On_Global_Summarize()
		{
		.On_Context_Summarize()
		}
	On_Global_Reporter()
		{
		.On_Context_Reporter()
		}
	On_Global_CrossTable()
		{
		.On_Context_CrossTable()
		}
	On_Global_Export()
		{
		.On_Context_Export()
		}

	On_Current_Print()
		{
		if (#() isnt (.list.GetSelection()))
			{
			row = .view.Get()
			saved_name = .Send('ExplorerListView_CurrentPrintSavedName')
			if saved_name is 0 or not String?(saved_name)
				saved_name = TruncateKey(.query)
			CurrentPrint(row, .Window.Hwnd, .query, saved_name, .excludeFields)
			}
		else
			Alert('You must select a row', title: 'Current Print',
				flags: MB.ICONERROR)
		}
	On_Context_Reporter()
		{
		if .query is false
			return

		Reporter()
		}
	On_Context_Summarize()
		{
		if (.query is false)
			return
		SummarizeControl(this, .getColumnsSaveName())
		}

	On_Context_CrossTable()
		{
		ToolDialog(.Window.Hwnd, Crosstab(.query, .getColumnsSaveName(), .excludeFields))
		}

	On_Context_Export()
		{
		GlobalExportControl(.query, excludeSelectFields: .excludeFields)
		}

	On_Context_Restore()
		{
		.invalid? = false
		.list.AllowContextOnly(false)
		.status.Set("")
		.status.SetValid(true)
		rec = .view.Get()
		if (rec.Explorer_NewRecord is true)
			{
			.On_Delete()
			return
			}
		// can't restore if record has not been saved to database yet
		if not rec.Member?("Explorer_PreviousData")
			return

		key = QueryKeys(.base_query)[0].Split(',')
		where = " where " $ key.
			Map({ it $ " is " $ Display(rec.Explorer_PreviousData[it]) }).
			Join(' and ')
		record = Query1(QueryAddWhere(.query, where))
		if (record is false)
			{
			Alert("The current record has been deleted.", title: 'Context Restore',
				flags: MB.ICONERROR)
			return
			}
		record.PreSet("Explorer_PreviousData", record.Copy())
		if false is row = .list.Get().Find(rec)
			{
			SuneidoLog('ERROR: ExplorerListView: list cannot find rec'
				params: Object(:rec, :record, list: .list.Get()))
			return
			}
		.list.SetRow(row, record)
		.view.Set(record)
		.Send('ExplorerListView_BeforeRestore', record)
		.model.Update(record)
		.Send('ExplorerListView_AfterRestore', record)
		if .linkField isnt false
			.model.SetRecordHeaderData(record)
		.view.Dirty?(false)
		.list.SetSelection(row)
		}
	On_Context_Restore_All()
		{
		if YesNo("All of your current changes will be lost. " $
			"Are you sure you want to proceed?",
			"Restore All?", .Window.Hwnd, MB.ICONWARNING)
			.Set(.Get())
		.list.AllowContextOnly(false)
		}
	On_Context_New()
		{
		.On_New()
		}
	On_Context_Delete()
		{
		.On_Delete()
		}
	On_Context_Save()
		{
		// if linked we should save all dirty records
		if not .Valid?()
			return
		if .linkField is false
			.SaveRecord()
		else
			.Send("Access_Save")
		}
	On_Delete()
		{
		x = .view.Get()
		if .allowDelete?(x)
			return
		.view.Dirty?(false)
		if (x.Explorer_NewRecord isnt true)
			if (false is .model.DeleteItem(x))
				return
		x.Explorer_RecordDeleted? = true
		.ignore_selection? = true
		listdata = .list.Get()
		if (false isnt (i = listdata.Find(x)))
			.list.DeleteRows(i)
		if (.list.GetNumRows() > 0)
			{
			// try setting selection to the next row, otherwise the last row
			row = i isnt false and listdata.Member?(i) ? i : false
			.load_entry(row)
			}
		else
			.On_New(false, setFocus: false)
		.ignore_selection? = false
		.invalid? = false
		.status.Set('')
		.status.SetValid(true)
		.update_explorerlist_data()
		.Send('ExplorerListView_AfterDelete', rec: x)
		}
	allowDelete?(x)
		{
		return (.readonly or
			not ProtectRuleAllowsDelete?(x, .protectField, x.Explorer_NewRecord) or
			false is .Send("ExplorerListView_AllowDelete", x))
		}
	// TODO: eliminate duplicate code between On_Delete and DeleteRow
	DeleteRow(rec)
		{
		if (false is i = .list.Get().Find(rec))
			return false
		if (false is .model.DeleteItem(rec))
			return false
		rec.Explorer_RecordDeleted? = true
		.ignore_selection? = true
		.list.DeleteRows(i)
		if (.view.Get() is rec) // deleting current selection
			{
			.view.Dirty?(false)
			.invalid? = false
			.ignore_selection? = true
			if (.list.GetNumRows() > 0)
				.load_entry()
			else
				.On_New(false)
			.status.Set('')
			.status.SetValid(true)
			}
		.ignore_selection? = false
		.update_explorerlist_data()
		return true
		}

	// Cannot find Select button ???
	On_Select()
		{
		.Select_vals = Object()
		.SaveRecord()
		.validate()
		if (.invalid? is false)
			SelectControl(this, okbutton:)
		}
	GetFields()
		{ return .list.GetColumns() }
	GetQuery()
		{ return .query }
	GetExcludeSelectFields()
		{ return #() }
	sf: false
	GetSelectFields()
		{
		if .sf is false
			.sf = SelectFields(.GetFields(), .GetExcludeSelectFields())
		return .sf
		}
	SetWhere(where, join? = false) // called by Select
		{
		.select_join? = join?
		if (where > "")
			{
			.status.Set("Select On")
			.select = true
			}
		else
			{
			.status.Set("")
			.select = false
			}
		.query = join?
			? QueryAddWhere("(" $ .base_query, where)
			: QueryAddWhere(.base_query, where)
		if (.linkField isnt false)
			.query = QueryAddWhere(.query, " where " $ .linkField $
				" is " $ Display(.value))
		.model.ChangeQuery(.query)
		.Load_entries()
		return true
		}
	readonly: false
	SetReadOnly(readOnly)
		{
		.readonly = readOnly
		.view.SetReadOnly(readOnly)
		}

	SaveRecord()
		{
		x = .view.Get()
		if (not .view.Dirty?() and x.Explorer_RecordUpdated isnt true)
			{
			.deleteUntouchedNewRec(x)
			return true
			}
		if (not .Valid?(evalRule?:))
			return false
		saveResult = true
		keyExceptionResult = KeyException.TryCatch()
			{
			.Send('ExplorerListView_BeforeSave', x)
			.updateHistoryFields(x)
			saveResult = x.Explorer_NewRecord is true
				? .Output()
				: .Update()
			if saveResult
				.Send('ExplorerListView_AfterSave', x)
			}
		return keyExceptionResult and (saveResult isnt false)
		}

	updateHistoryFields(x)
		{
		if .historyFieldsPrefix is false
			return

		x[.historyFieldsPrefix $ '_date_modified'] = Timestamp()
		x[.historyFieldsPrefix $ '_user_modified'] = Suneido.User
		}

	deleteUntouchedNewRec(x)
		{
		if (x.Explorer_NewRecord is true)
			{
			.ignore_selection? = true
			list_data = .list.Get()
			if (false isnt (i = list_data.Find(x)) and
				list_data.Size() > 1)
				.list.DeleteRows(i)
			.ignore_selection? = false
			}
		}

	Output()
		{
		x = .view.Get()
		Assert(x.Explorer_NewRecord is true, "ExplorerListView: Output New record")
		if (.model.Output(x) is false)
			return false
		// set sticky values
		for (f in .stickyFields)
			.sticky_values[f] = x[f]
		x.Delete("Explorer_NewRecord")
		.view.Dirty?(false)
		return true
		}
	Update()
		{
		if (not .view.Dirty?() or #() is (sel = .list.GetSelection()))
			return false

		// update list
		.list.RepaintRow(sel[0])

		x = .view.Get()
		result = .model.Update(x)

		// the following restore handles getting the new version of record if another
		// user changed it.Most of the time this is unnecessary, but it doesn't require
		// extra code to communicate overwrites to the model, and it doesn't seem to
		// cause any performance issues
		if .linkField is false and result is false
			.On_Context_Restore()

		.view.Dirty?(false)
		return result isnt false
		}

	// from List, return false to prevent list behaviour
	List_RightClick()
		{
		return .disableList()
		}
	List_SingleClick(row /*unused*/, col /*unused*/)
		{
		return .disableList()
		}

	disableList()
		{
		.validate(leaving:)
		.list.AllowContextOnly(.invalid?)
		return .invalid? ? false : 0
		}
	List_DoubleClick(row, col)
		{
		.Send("ExplorerListView_DoubleClick", row, col)
		return .On_New()
		}
	List_ContextMenu(x, y)
		{
		if .readonly
			return
		menu = Object("New", "Delete", "", "Save", "Restore")
		if (.linkField isnt false)
			menu.Add("Restore All")
		else
			menu.Add("", "Reporter...", "Summarize...", "CrossTable...", "Export...")

		if 0 isnt ob = .Send('ExplorerListView_AddToContextMenu')
			{
			menu.Add("")
			for option in ob
				menu.Add(option)
			}
		ContextMenu(menu.Difference(.disableOptions)).ShowCall(this, x, y)
		}
	List_KeyDown(wParam)
		{
		if wParam is VK.INSERT
			.On_New()
		}
	List_DeleteKeyDown()
		{
		.On_Delete()
		return false
		}
	List_DeleteRecord(record /*unused*/)
		{
		return true
		}
	List_WantNewRow()
		{
		return false
		}
	List_WantEditField(col /*unused*/, row /*unused*/, data /*unused*/)
		{
		return false
		}

	List_BeforeClearSelect()
		{
		if (.ignore_selection? or .view is false)
			return
		.SaveRecord()
		}

	List_AfterClearSelect()
		{
		.load_entry()
		}

	selected: false
	List_Selection(selection)
		{
		if not .selectable?(selection)
			return
		x = .list.GetRow(.selected[0])
		.SaveRecord()
		// SaveRecord may delete a new, non-dirty record which could change the
		// selection.  The following line is to retain the selection
		if (false isnt (i = .list.Get().Find(x)) and i isnt .selected[0])
			.list.SetSelection(i)
		if (.selected isnt false)
			{
			.Send('ExplorerListView_Selection', x)
			.view.Set(x)
			if (.protectField isnt false)
				x[.protectField]
			.Send('ExplorerListView_EntryLoaded')
			}
		}

	selectable?(selection)
		{
		if selection is false or selection is .selected
			return false
		.selected = selection
		if .ignore_selection? or .view is false
			return false
		return true
		}

	List_SelectedRowPositionChanged(selection)
		{
		if selection is false or selection is .selected
			return
		.selected = selection
		}

	invalid?: false
	Valid?(evalRule? = false)
		{
		.validate(leaving:, evalRule?: evalRule?)
		return .invalid? is false
		}
	validate(leaving = false, evalRule? = false)
		{
		if (not .view.Dirty?() or .view.GetReadOnly() is true)
			return

		rec = .view.Get()

		// don't want to check after every field on new records
		if ((rec.Explorer_NewRecord is true or .list.GetNumRows() is 0) and
			.invalid? is false and not leaving)
			return

		if not .recordControlsValid?(leaving)
			return

		.checkValidField(evalRule?, rec)
		}
	checkValidField(evalRule?, rec)
		{
		.valid_msg = ''
		if (.view.Dirty?() and .validField isnt false)
			.valid_msg = evalRule? is true
				? rec.Eval(Global('Rule_' $ .validField))
				: rec[.validField]
		.invalid? = (.valid_msg isnt '')
		.status.Set(.valid_msg)
		.status.SetValid(.invalid? is false)
		}
	recordControlsValid?(leaving)
		{
		if (leaving isnt false and (invalid_fields = .view.Valid()) isnt true)
			{
			.invalid? = true
			.status.SetValid(false)
			.status.Set(invalid_fields)
			Beep()
			return false
			}
		return true
		}
	value: false
	Set(value)
		{
		.value = value
		if (.linkField isnt false)
			{
			// save and restore select
			// so it doesn't changed when you click Edit
			sel = .list.GetSelection()
			curRec = sel.NotEmpty?() and sel[0] < .list.GetNumRows()
				? .list.GetRow(sel[0])
				: false
			header_data = .Send("GetData")
			.query = QueryAddWhere(.base_query, " where " $ .linkField $
				" is " $ Display(value))
			keyFields = .model isnt false ? .model.GetKey() : .model_spec[2]
			headerfields = .model isnt false ? .model.GetHeaderFields() : #()
			.NewModel(Object(.model_spec[0], .query, keyFields,
				headerFields: headerfields, headerData: header_data))
			if Object?(curRec) and false isnt (i = .list.Get().FindIf(
				{ |item|
				item.Project(keyFields) is curRec.Project(keyFields)
				}))
				.list.SetSelection(i)

			.Send('DoWithoutDirty')
				{
				.update_explorerlist_data()
				}
			}
		}
	Get()
		{
		return .value
		}

	load_entry(row = false)
		{
		if row is false
			row = .list.GetNumRows() - 1
		if (row >= 0)
			{
			.list.SetSelection(row)
			x = .list.GetRow(row)
			.Send('ExplorerListView_BeforeEntryLoaded', x)
			.view.Set(x)
			if (.protectField isnt false)
				x[.protectField]
			.list.AllowContextOnly(false)
			.Send('ExplorerListView_EntryLoaded')
			}
		else
			.On_New()
		}
	AccessChanged(@args)
		// pre:		event is an AccessControl event string
		// post:	performs ExplorerListView processing for event
		{
		switch (args[0])
			{
		case 'restore':		.restore()
		case 'delete':		return .delete()
		case 'save' :		return .save(args[1])
		default:
			}
		return true
		}
	restore()
		// post:	current query is reissued without saving changes
		{
		.model.Restore()
		.Load_entries()
		}
	delete()
		// post:	checks if OK to delete,
		//			model's data and deletions lists are empty
		// NOTE: 	There should be a foreign key from lines to header
		//			with cascading deletes enabled.
		{
		for rec in .model.GetData()
			if false is ProtectRuleAllowsDelete?(rec, .protectField,
				rec.Explorer_NewRecord, true)
				return false
		.model.Clear()
		.ignore_selection? = true
		.list.Clear()
		.ignore_selection? = false
		return true
		}
	save(t)
		{
		if (not .Valid?())
			return false
		saveResult = .SaveRecord() and .model.Save(t)
		if saveResult is true
			.view.Dirty?(false)
		return saveResult
		}

	ConfirmDestroy()
		{
		.ClearFocus() // ensure current field is updated in RecordControl

		.validate(leaving:)
		if .invalid? is true
			return false
		// if linked, the Access will trigger the save
		if .linkField is false
			return .SaveRecord()
		return true
		}

	Destroy()
		{
		.Send("NoData")
		UserColumns.Save(.getColumnsSaveName(), .list, .columns)
		if .linkField is false
			.Window.RemoveValidationItem(this)
		else
			.Send("RemoveAccessObserver", .AccessChanged)
		super.Destroy()
		}
	}