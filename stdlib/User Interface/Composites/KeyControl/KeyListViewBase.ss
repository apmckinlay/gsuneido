// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	ComponentName: 'KeyListViewBaseComponent'
	New(.query, columns, saveInfoName = "", prefix = "", prefixColumn = false,
		keys = false, .field = "", value = "", enableMultiSelect = false,
		.checkBoxColumn = false, .customizeQueryCols = false, .optionalRestrictions = #()
		.startLast = false)
		{
		super(.Layout(query, columns, saveInfoName, prefix, prefixColumn,
			keys, :enableMultiSelect, :checkBoxColumn))
		.initializeList()
		.prevproc = SetWindowProc(.fieldCtrl.Hwnd, GWL.WNDPROC, .Fieldproc)
		.setFieldFromCurrentRec(field, query, prefix, value)
		.NewValue(SelectPrompt(.prefixColumn), .prefixBy)
		.Defer(.set_initial_focus)
		}

	Layout(@args)
		{
		return .BaseLayout(@args)
		}

	BaseLayout(query, columns, saveInfoName, prefix, prefixColumn, keys,
		enableMultiSelect, checkBoxColumn)
		{
		result = .initializeLayoutVariables(query, keys, prefixColumn, columns,
			saveInfoName, prefix)
		keys = result.keys
		colnames = result.colnames

		columnsSaveName = saveInfoName is "" ? .save_query : saveInfoName

		allColumns = .allColumns()
		if columns is false
			columns = allColumns
		if false is UserColumns.GetDefaultColumns(columnsSaveName)
			UserColumns.SaveDefaultColumns(columnsSaveName, columns, allColumns)

		optRestrictionsControl = .optionalRestrictions.NotEmpty?()
			? .optionalRestrictionLayout()
			: ''
		mandatoryFields = keys.Map({ it.Replace('_lower!$') })
		return Object("Vert",
			optRestrictionsControl,
			Object("VirtualList", .baseQuery, allColumns, columnsSaveName,
				headerSelectPrompt:, hideCustomColumns?: not .customizeQueryCols
				menu: #(Copy), :mandatoryFields, :enableMultiSelect,
				:checkBoxColumn, keyField: .field, ymin: 200, xmin: 100, name: 'List',
				excludeSelectFields: .GetExcludeSelectFields(),
				hideColumnsNotSaved?:, preventCustomExpand?:, startLast: .startLast),
			#(Skip, 5),
			Object('Horz',
				#(Skip 2),
				#(Static Locate),
				#(Skip 4),
				#(Field, width: 10, xstretch: 1),
				#(Skip 4),
				#(Static by),
				#(Skip 4),
				(colnames.Size() <= 1
					? Object('Static', SelectPrompt(.prefixColumn),
						name: 'prefixBy')
					: Object('ChooseList', colnames, set: keys[0],
						name: 'prefixBy'))
				#(Skip 6),
				#(Button Go),
				#(Skip 2)
				name: "GoHorz"))
		}

	optionalRestrictionLayout()
		{
		form = Object('Form')
		for field in .optionalRestrictions
			{
			dd = Datadict(field)
			if .ChooseDateCtrl?(dd)
				{
				dateQuery = .BuildChooseDateQueryWhere(field)
				if dd.Prompt.Has?('Terminate')
					.addCheckBoxAndQuery(form, 'terminated', dateQuery)
				else if dd.Prompt.Has?('Expiry')
					.addCheckBoxAndQuery(form, 'expired', dateQuery)
				}
			else if .ChooseListActiveInactive?(dd)
				.addCheckBoxAndQuery(form, 'inactive',
					.BuildChooseListActiveInactiveWhere(field))
			}
		// Initialize .extraWhere with all "Show" checkboxes unchecked
		.updateExtraWhere([])
		return Object('Record', form, name: 'KeyListViewParams')
		}

	BuildChooseDateQueryWhere(field)
		{
		return "(" $ field $ ' is "" or ' $
			field $ ' >= ' $ Display(Date().NoTime()) $ ")"
		}

	BuildChooseListActiveInactiveWhere(field)
		{
		return field $ ' is "active"'
		}

	ChooseDateCtrl?(dd)
		{
		return GetControlClass.FromControl(dd.Control) is ChooseDateControl
		}

	ChooseListActiveInactive?(dd)
		{
		return  GetControlClass.FromControl(dd.Control) is ChooseListControl and
			dd.Control[1].Has?('inactive') and dd.Control[1].Has?('active')
		}

	addCheckBoxAndQuery(form, type, query)
		{
		name = "show_" $ type
		form.Add(Object('CheckBox' text: "Show " $ type.Capitalize(), :name))
		.whereMap[name] = query
		}

	extraWhere: ''
	updateExtraWhere(data)
		{
		whereOb = Object()
		for m in .whereMap.Members().Sort!() // sort so order is consistent for tests
			if data[m] isnt true
				whereOb.Add(.whereMap[m])
		.extraWhere = Opt(' where ', whereOb.Join(' and '))
		}

	getter_whereMap()
		{
		return .whereMap = Object()
		}

	Record_NewValue(field/*unused*/, value/*unused*/, source)
		{
		if source.Name isnt 'KeyListViewParams'
			return
		.updateExtraWhere(source.Get())
		.NewValue(SelectPrompt(.GetPrefixColumn()), .GetPrefixBy())
		}

	SetWhere(where)
		{
		.GetList().SetWhere(where)
		return true
		}

	allColumns()
		{
		// b/c of permission issues per field we can NOT default this to on for all fields
		if not .customizeQueryCols
			return .columns

		// using SelectFields based off Suggestion: 20034
		excludeFields = .GetExcludeSelectFields().Copy()
		if .field isnt "" and GetControlClass.FromField(.field).Base?(IdControl)
			excludeFields.Add(.field)
		allColumns = SelectFields(QueryColumns(.baseQuery),
			joins: false, :excludeFields).Fields.UniqueValues()
		return allColumns
		}

	Startup()
		{
		if .fieldCtrl.Get() isnt ''
			.FieldChange()
		}

	GetExcludeSelectFields()
		{
		return Object()
		}

	savedQuery?: false
	extended: #()
	orig_locateBy: false
	initializeLayoutVariables(query, keys, prefixColumn, columns, saveInfoName, prefix)
		{
		// have to do substr so key on locateby doesn't get too long. Also, some
		// queries may have a number identifier on the end
		// that needs stripped off (unique view names).

		.saveInfoName = saveInfoName
		.save_query = TruncateKey(query)
		.info = .saveInfoName is '' ? false : KeyListViewInfo.Get(saveInfoName)
		// Issue 11413 - fall back to the old way we saved if new way fails,
		// new way will be used to save in 'Destroy'
		if .info is false
			{
			.info = KeyListViewInfo.Get(.save_query)
			.savedQuery? = .info isnt false
			}
		if Object?(.info) and .info.Member?('window_info') and Object?(.info.window_info)
			.orig_locateBy = .info.window_info.GetDefault('locateby_key', false)

		query_columns = QueryColumns(query)
		keys = .initializeKeys(keys, query, query_columns)
		.handlePrefixColumn(prefixColumn, keys, columns, query_columns)
		colnames = keys.Map(SelectPrompt)
		.prefix = prefix
		extend = ''
		if columns isnt false
			{
			.extended = columns.Difference(query_columns)
			extend = .getExtended(.extended)
			}
		.baseQuery = QueryStripSort(query) $ extend $ Opt(" sort ", .prefixColumn)
		return Object(:keys, :colnames)
		}

	getExtended(extended)
		{
		return Opt(' extend ', extended.Join(', '), ' ')
		}

	initializeKeys(keys, query, query_columns)
		{
		if keys is false
			{
			keys = QueryKeys(query)
			.addAbbrevToKeys(keys, query_columns)
			.removeUnwantedKeys(keys)
			}
		else // need this because Access Locate passes e.g. #(Name: biz_name)
			keys = keys.Values()
		if keys.Empty?()
			throw "KeyListView: no usable keys"
		return keys
		}

	addAbbrevToKeys(keys, cols)
		{
		if false isnt (i = keys.FindIf({|c| c.Suffix?('_name') })) and
			cols.Has?(ka = keys[i].Replace('_name$', '_abbrev'))
			keys.Add(ka)
		}

	removeUnwantedKeys(keys)
		{
		keys.RemoveIf() {|k| k =~ ',' or k =~ '_num(_new)?$' }
		keys.RemoveIf() {|k| keys.Has?(k $ '_lower!') }
		keys.Remove("")
		}

	handlePrefixColumn(prefixColumn, keys, columns, query_columns)
		{
		if prefixColumn is false
			prefixColumn = .hasValidSavedKeyChoice?(keys)
				? .info.window_info.locateby_key
				: keys[0]

		.prefixColumn = prefixColumn
		.columns = columns isnt false ? columns.Copy() : query_columns
		// check for invalid key field stored in table
		if query_columns.Has?(.prefixColumn)
			return

		.prefixColumn = keys[0]
		// ensure new key is in columns
		if not .columns.Has?(.prefixColumn)
			.columns.Add(.prefixColumn)
		}

	hasValidSavedKeyChoice?(keys)
		{
		return Object?(.info) and
			.info.Member?('window_info') and
			.info.window_info.Member?('locateby_key') and
			keys.Has?(.info.window_info.locateby_key)
		}

	list: false
	initializeList()
		{
		.list = .FindControl('List')
		.prefixBy = .FindControl('prefixBy')
		.fieldCtrl = .FindControl('Field')
		}

	Fieldproc(hwnd, msg, wparam, lparam)
		{
		_hwnd = .WindowHwnd()
		switch msg
			{
		case WM.GETDLGCODE:
			return .wmGetDlgCode(hwnd, msg, wparam, lparam)
		case WM.KEYDOWN:
			if wparam in (VK.UP, VK.DOWN, VK.NEXT, VK.PRIOR)
				{
				SendMessage(.list.GetGridHwnd(), msg, wparam, lparam)
				SetFocus(.list.GetGridHwnd())
				return 0
				}
		case WM.CHAR:
			return .wmChar(hwnd, msg, wparam, lparam)
		case WM.MOUSEWHEEL:
			SendMessage(.list.GetGridHwnd(), msg, wparam, lparam)
			return 0
		default:
			}
		return CallWindowProc(.prevproc, hwnd, msg, wparam, lparam)
		}

	wmGetDlgCode(hwnd, msg, wparam, lparam)
		{
		if false isnt (m = MSG(lparam)) and
			(m.wParam is VK.RETURN or m.wParam is VK.DOWN) and
			(m.message is WM.CHAR or m.message is WM.KEYDOWN)
			return DLGC.WANTALLKEYS | DLGC.WANTARROWS
		return CallWindowProc(.prevproc, hwnd, msg, wparam, lparam)
		}

	wmChar(hwnd, msg, wparam, lparam)
		{
		if wparam is VK.RETURN
			{
			if GetFocus() is .fieldCtrl.Hwnd
				.FieldEnter()
			return 0
			}
		return CallWindowProc(.prevproc, hwnd, msg, wparam, lparam)
		}

	FieldEnter()
		{
		.FieldChange()
		SetFocus(.list.GetGridHwnd())
		}

	setFieldFromCurrentRec(field, query, prefix, value)
		{
		if (field isnt "" and value isnt "" and
			false isnt rec = Query1(QueryAddWhere(query,
				" where " $ field $ " is " $ Display(value))))
			.fieldCtrl.Set(rec[.prefixColumn])
		else
			.fieldCtrl.Set(prefix)
		.prefixBy.Set(SelectPrompt(.prefixColumn))
		}

	FieldChange()
		{
		if .list is false
			return
		prefix = DatadictEncode(.prefixColumn, .fieldCtrl.Get())
		if .prefixColumn.Suffix?('_lower!')
			prefix = prefix.Lower()
		.list.Seek(.prefixColumn, prefix)
		}

	GetList()
		{
		return .list
		}

	GetField()
		{
		return .fieldCtrl
		}

	GetPrefixColumn()
		{
		return .prefixColumn
		}

	prefixBy: false // NewValue requires this in the middle initialization
	GetPrefixBy()
		{
		return .prefixBy
		}

	GetBaseQuery()
		{
		return .baseQuery
		}

	NewValue(value, source)
		{
		// need this when Access button is used on Browse, sometimes
		// this control gets destroyed but it still tries to do newvalue
		if .Destroyed?() or source isnt .prefixBy
			return

		if value is "" or not .prefixBy.Valid?()
			return

		columns = .allColumns()
		column = field = .field_from_prompt(value)
		if column.Suffix?('_lower!')
			column = column.BeforeLast('_lower!')
		columns.Remove(column, field)
		columns.Add(column, at: 0)
		.prefixColumn = field
		extended = .list.GetColumns().Intersect(.extended)
		.baseQuery = QueryStripSort(.query) $ .getExtended(extended) $ .GetExtraWhere() $
			" sort " $ .prefixColumn
		.list.SetQuery(.baseQuery, columns)
		.list.MoveColumnToFront(column)
		SetFocus(.list.GetGridHwnd())
		}

	GetExtraWhere()
		{
		return .extraWhere
		}

	field_from_prompt(prompt)
		{
		if prompt.Suffix?('*')
			return .field_from_prompt(prompt[..-1]) $ '_lower!'
		for field in .columns
			{
			if field.Prefix?('-')
				field = field[1 ..]
			if SelectPrompt(field) is prompt
				return field
			}
		throw 'Unknown prompt: ' $ prompt
		}

	set_initial_focus()
		{
		if .fieldCtrl.Get() is ""
			SetFocus(.fieldCtrl.Hwnd)
		}

	On_Go()
		{
		.FieldChange()
		SetFocus(.list.GetGridHwnd())
		}

	VirtualList_Tab()
		{
		SetFocus(.fieldCtrl.Hwnd)
		}

	VirtualList_DisableSort?()
		{
		return true
		}

	On_Context_Copy(rec, col)
		{
		if col is false or rec is false
			return
		value = rec[col]
		if not String?(value)
			value = Display(value)
		ClipboardWriteString(value)
		}

	saveLocateByKey()
		{
		prefixByVal = .prefixBy.Get()
		locateby_key = prefixByVal isnt "" and .prefixBy.Valid?()
			? .field_from_prompt(prefixByVal)
			: false
		if locateby_key isnt .orig_locateBy
			.save(locateby_key)
		}

	save(locateby_key)
		{
		name = .saveInfoName is '' ? .save_query : .saveInfoName
		window_info = Object() // KeyListCheckboxView does not save size info
		if false isnt savedInfo = KeyListViewInfo.Get(name)
			window_info = savedInfo.window_info
		window_info.locateby_key = locateby_key
		KeyListViewInfo.Save(name, window_info)
		if .savedQuery? is true and .info isnt false and .saveInfoName isnt ""
			KeyListViewInfo.DeleteRecord(.save_query)
		}

	Destroy()
		{
		.saveLocateByKey()
		super.Destroy()
		ClearCallback(.Fieldproc)
		}
	}
