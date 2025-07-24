// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	mainList: false
	selectColumn: 'params_itemselected'
	New(field, values)
		{
		super(.layout(field, values))
		.values = values
		.field = Object?(field) and field.Member?("name") ? field.name : field
		.control = .FindControl(.field)
		.status = .FindControl(#Status)
		.listControl = .FindControl('valuesList')
		.mainList = .FindControl('mainList')
		.setList(.values)
		.loadList()
		}

	layout(field, initialList)
		{
		.listInfo = GetControlListInfo(field)
		ob = Object('Vert')
		if .listInfo.Member?('query') 		// ListControl w/ check boxes
			ob.Add(.queryLayout(.listInfo, field))
		else if not .listInfo.Empty?()
			{
			if .listInfo[0] is 'ChooseDate'	// ChooseManyDates
				ob.Add(.datesLayout(initialList))
			else							// ChooseTwoList
				ob.Add(.objectLayout(.listInfo, initialList))
			}
		else 								// regular field
			ob.Add(.fieldLayout(field))

		ob.Add('Skip',
			#(HorzEqual
				(Button Load) Skip (Button Save) Skip
					'Fill', OkCancel),
			#Status)
		return Object('Record' ob)
		}

	queryLayout(listInfo, field)
		{
		.set = .setEditorList
		columns = listInfo.columns.Copy()
		return Object('Vert'
			Object(KeyListCheckboxView, query: listInfo.query, :columns,
				saveInfoName: field, field: listInfo.field, keys: listInfo.keys,
				enableMultiSelect:, checkBoxColumn: .selectColumn, name: 'mainList',
				customizeQueryCols: listInfo.GetDefault(#customizeQueryCols, false)
				optionalRestrictions: listInfo.GetDefault(#optionalRestrictions, #()),
				excludeSelect: listInfo.GetDefault(#excludeSelect, #()))
			'Skip'
			#(Scintilla readonly:, wrap:, height: 2, ystretch: 0, name: 'valuesList')
			'Skip'
			#(Button 'Clear List' xmin: 100)
			xstretch: 1)
		}

	datesLayout(initialList)
		{
		.set = .setDates
		.get = .getDates
		.clear = .clearDates
		return Object('Vert'
			Object('MonthCalDates' initialList.Join(','), name: 'valuesList')
			'Skip'
			#(Button 'Clear List' xmin: 100)
			)
		}

	objectLayout(listInfo, initialList)
		{
		.set = .setTwoList
		.get = .getTwoList
		return Object('TwoList' listInfo, initialList, xstretch: 1, name: 'valuesList')
		}

	fieldLayout(field)
		{
		.set = .setListBox
		control = .getFieldControl(field)
		return Object('Vert',
			control,
			'Skip',
			#(Horz
				Fill
				(Button 'Add' xmin: 100)
				Fill
				(Button 'Clear List' xmin: 100)
				Fill)
			'Skip'
			#(ListBox xstretch: 1, name: 'valuesList'))
		}

	getFieldControl(field)
		{
		if not String?(field)
			return field

		.fieldCtrl = GetControlListInfo.GetControlFromField(field)
		return Object('Pair', Object('Static', .getFieldPrompt(field)), .fieldCtrl)
		}

	getFieldPrompt(field)
		{
		if field isnt selectPrompt = .getSelectPrompt(field)
			return selectPrompt

		return .getSelectPrompt(field.Replace(`_(name|abbrev)(\>|_)`, `_num\2`)) $
			.getNameAbbrevSuffix(field)
		}

	// extracted for testing
	getSelectPrompt(field)
		{
		return SelectPrompt(field)
		}

	getNameAbbrevSuffix(field)
		{
		return field =~ `_name(\>|_)`
			? ' Name'
			: field =~ `_abbrev(\>|_)`
				? ' Abbrev'
				: ''
		}

	get: false // Determined via the layout methods
	getList()
		{
		return .listControl is false or .get is false
			? .values
			: (.get)()
		}

	getTwoList()
		{
		return .listControl.GetNewList()
		}

	getDates()
		{
		return .listControl.Get()
		}

	set: false // Determined via the layout methods
	setList(values)
		{
		if .listControl is false or .set is false
			return
		(.set)(values)
		.validate(.getList())
		}

	setListBox(values)
		{
		.listControl.DeleteAll()
		for val in ParamsChooseListControl.DisplayValues(values, .field)
			.listControl.AddItem(val)
		}

	setEditorList(values)
		{
		.listControl.Set(.buildEditorDisplayValue(values))
		}

	buildEditorDisplayValue(values)
		{
		displayLimit = 200
		values = ParamsChooseListControl.DisplayValues(values, .field)
		if values.Size() > displayLimit
			text = values[..displayLimit].Join(',') $ '...'
		else
			text = values.Join(',')
		return text
		}

	setTwoList(values)
		{
		.listControl.AllBack()
		.listControl.Set(values.Join(','))
		}

	setDates(values)
		{
		for date in .listControl.Get().Copy()
			.listControl.DateSelectChange(date)
		for date in values.Copy()
			.listControl.DateSelectChange(date)
		.listControl.Set(values)
		}

	loadList()
		{
		if .mainList is false or .values.Empty?()
			return
		.mainList.CheckRecordByKeys(.values)
		}

	ListBox_ContextMenu(x, y)
		{
		list = .FindControl('valuesList')
		if list.GetCurSel() is -1
			return // click outside items
		ContextMenu(#('Delete')).ShowCall(this, x, y)
		}

	On_Context_Delete()	// from queryLayout and fieldLayout
		{
		list = .FindControl('valuesList')
		item = list.GetCurSel()
		list.DeleteItem(item)
		.values.Delete(item)
		.validate(.getList())
		}

	On_Add()	// from fieldLayout
		{
		// control can be false
		// multi-user issue - one user renames/deletes custom field on
		// reporter report and another user accessing it from reporter report menu
		if .control is false
			return
		// get the value
		val = .control.Get()
		if (val is '' or .values.Has?(val) or not .control.Valid?())
			{
			Beep()
			return
			}
		.values.Add(val)
		.setList(.values)
		.control.Set('')
		}

	On_Load()
		{
		if false is vals = LoadSave_InListValues.Load(.Window.Hwnd, .getFieldSaveName())
			return

		.clearList()
		.values = vals
		.setList(.values)
		.loadList()
		}

	On_Save()
		{
		list = .getList()
		if .validate(list) isnt true
			return
		LoadSave_InListValues.Save(.Window.Hwnd, .getFieldSaveName(), list)
		}

	getFieldSaveName()
		{
		suffixes = #(num id)
		for suffix in suffixes
			if .field.Has?('_' $ suffix $ '_')
				return .field.BeforeFirst('_' $ suffix $ '_') $ '_' $ suffix
		return .field
		}

	On_Clear_List()
		{
		.clearList()
		}

	clear: false // Determined via the layout methods
	clearList()
		{
		if .mainList isnt false
			.mainList.UncheckAll()
		if .clear isnt false
			(.clear)()
		.values = Object()
		.setList(.values)
		}

	clearDates()
		{
		for date in .values.Copy()
			.listControl.DateSelectChange(date)
		}

	On_OK()
		{
		.Send("On_OK")
		}
	OK()
		{
		list = .getList()
		if .validate(list) isnt true
			return false
		return .getList()
		}

	limit: 255
	validate(list)
		{
		if list.Size() > .limit
			{
			.status.SetValid(false)
			.status.Set('Cannot select more than ' $ .limit $ ' items')
			return false
			}
		.status.SetValid(true)
		.status.Set('')
		return true
		}

	VirtualList_LeftClick(rec, col)
		{
		if col isnt .selectColumn
			return

		clicked = rec[.listInfo.field]
		if .values.Has?(clicked)
			.values.Remove(clicked)
		else
			.values.Add(clicked)
		.setList(.values)
		}

	VirtualList_Space()
		{
		checked = .mainList.GetCheckedRecords()
		.values = checked.list.Map({ it[.listInfo.field] }).Instantiate()
		.setList(.values)
		}

	VirtualList_ModelChanged()
		{
		.loadList()
		}
	}