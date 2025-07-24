// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Offset: 		0
	widths: 		#()
	selectMgr: 		false

	New(columns = #(), .columnsSaveName = false, .headerSelectPrompt = false,
		.mandatoryFields = #(), .checkBoxColumn = false, excludeSelectFields = #(),
		.hideCustomColumns? = false, query = false, .extraFmts = false, customKey = false,
		.enableUserDefaultSelect = false, .disableSelectFilter = false,
		.StretchColumn = false, .hideColumnsNotSaved? = false, option = false,
		.defaultColumns = false)
		{
		.origColumns = .initColumns(query, columns)
		.setup()

		.plugins = MultiViewPlugins(.columns, option is false ? 'VirtualList' : option)
		.setCustomKey(customKey, query, columns)
		.excludeSelectFields = .initExcludeFields(query, excludeSelectFields)
		}

	initColumns(query, columns)
		{
		if query is false
			return columns
		if columns.Empty?()
			columns = .defaultCols(query)
		if .columnsSaveName isnt false and not .hideCustomColumns?
			columns = ListCustomize.AddCustomColumns(query, columns)
		return columns
		}

	capFields: #()
	defaultCols(query)
		{
		columns = QueryColumns(query)
		keys = QueryKeys(query).Map({ it.BeforeFirst(',') })
		columns.Sort!({|x, y| keys.Has?(x) > keys.Has?(y) }) //put key fields in the front

		// get rule/unsaved fields so we can put them last to optimize QueryView display
		.capFields = QueryRuleColumns(query)
		// put rule fields after
		return columns.Sort!({|x, y| .capFields.Has?(x) < .capFields.Has?(y) })
		}

	ensureDefaultColumns(columns)
		{
		if .defaultColumns is false
			return
		usercolumns_title = .GetColumnsSaveName()
		if usercolumns_title is false or usercolumns_title is .TmpSelectName
			return
		UserColumns.EnsureTable()
		if QueryEmpty?('usercolumns', :usercolumns_title, usercolumns_user:'')
			UserColumns.SaveDefaultColumns(
				usercolumns_title, .defaultColumns, columns, deletecol: false)
		}

	CapFieldPrompt(field)
		{
		return .headerSelectPrompt is 'no_prompts' and .capFields.Has?(field)
			? field.Capitalize()
			: false
		}

	formatting: class { Destroy(){} }
	setup()
		{
		if .checkBoxColumn isnt false and not .origColumns.Has?(.checkBoxColumn)
			.origColumns.Add(.checkBoxColumn)
		.columns = .origColumns.Copy()

		.widths = Object()

		.formatting.Destroy()
		.formatting = ListFormatting(fontFixed: false, booleansAsBox: false, tooltip:,
			noFormat: .headerSelectPrompt is 'no_prompts', extraFmts: .extraFmts)
		.formatting.SetFormats(.columns)
		}

	Plugins_Execute(@args)
		{
		if .plugins is false or not args.Member?('pluginType')
			return

		// pluginType = 'Observers', 'AfterField'
		.plugins[args.pluginType](@args)
		}

	GetFormatting()
		{
		return .formatting
		}

	TmpSelectName: 'virtual_list_temp_name'
	getSelectSaveName()
		{
		return .columnsSaveName isnt false and .columnsSaveName isnt ''
			? .columnsSaveName : .TmpSelectName
		}

	initExcludeFields(query, excludeFields)
		{
		return query is false or .headerSelectPrompt is 'no_prompts'
			? excludeFields
			: excludeFields.Copy().MergeUnion(Customizable.GetNonPermissableFields(query))
		}

	GetColumnsSaveName()
		{
		return .columnsSaveName
		}

	GetHeaderSelectPrompt()
		{
		return .headerSelectPrompt
		}

	GetWidths()
		{
		return .widths
		}

	GetStretchCol(col = false)
		{
		if false isnt stretchCol = .FindCol(.StretchColumn)
			if .GetColWidth(stretchCol) in (0, false)
				stretchCol = false
		if stretchCol is false
			stretchCol = .GetLastVisibleCol()
		if col is false
			return stretchCol
		if col > stretchCol
			return .GetLastVisibleCol()
		if col < stretchCol
			return stretchCol
		return .GetLastVisibleCol()
		}

	GetLastVisibleCol()
		{
		return .widths.FindLastIf({ it > 0 })
		}

	GetColWidth(i)
		{
		return .widths.GetDefault(i, false)
		}

	GetColumnWidth(col)
		{
		return .GetColWidth(.FindCol(col))
		}

	SetColWidth(i, width)
		{
		.widths[i] = width
		}

	SetColumns(.columns)
		{
		.widths = Object()
		.formatting.SetFormats(.columns)
		}

	InsertColumn(column, pos)
		{
		.columns.Add(column, at: pos)
		.widths.Add(100 /*= default col width */, at: pos)
		.headerChanged? = true
		.formatting.SetFormats(.columns)
		}

	RemoveColumn(column)
		{
		pos = .columns.Find(column)
		.columns.Remove(column)
		.widths.Delete(pos)
		.headerChanged? = true
		.formatting.SetFormats(.columns)
		}

	ResetColumns(permissableQuery)
		{
		if .columnsSaveName is false
			.setup()
		else
			{
			origCols = .removeHiddenCols(.origColumns)
			UserColumns.Reset(this, .columnsSaveName, origCols, deletecol: false,
				load_visible?:, extraCols: .extraColumn(), :permissableQuery,
				hideColumnsNotSaved?: .hideColumnsNotSaved?)
			}
		}

	GetTotalWidths()
		{
		return .widths.Sum()
		}

	GetColumns()
		{
		return .columns
		}

	GetOriginalColumns()
		{
		return .origColumns
		}

	Get(i)
		{
		return .columns.GetDefault(i, false)
		}

	FindCol(col)
		{
		return .columns.Find(col)
		}

	GetSize()
		{
		return .columns.Size()
		}

	GetColFormat(i)
		{
		.formatting.GetHeaderAlign(.columns[i])
		}

	GetColumnOffset(col)
		{
		left = 0
		for i in .columns.Members()
			{
			if .columns[i] is col
				break
			left += .widths[i]
			}
		return left - .Offset
		}

	GetColByX(x)
		{
		numCols = .GetSize()
		for (col = 0, left = - .Offset; col < numCols; col++)
			if x < left += .GetColWidth(col)
				break

		return .Get(col)
		}

	SetDC(hdc)
		{
		.formatting.SetDC(hdc)
		}

	SetBackgroundBrush(brush, selected = false)
		{
		return .formatting.SetBackgroundBrush(brush, :selected)
		}

	PaintCell(c, rect, rec)
		{
		.formatting.PaintCell(
			.columns[c], rect.GetX(), rect.GetY(), rect.GetWidth(), rect.GetHeight(), rec)
		}

	ReorderColumn(col, newIdx)
		{
		oldColumns = .columns.Copy()
		.columns.Delete(col).Add(oldColumns[col], at: newIdx)
		org = .widths[col]
		.widths.Delete(col).Add(org, at: newIdx)
		}

	SetSelectVals(select_vals)
		{
		.selectMgr.SetSelectVals(select_vals, .GetSelectFields())
		}

	SetSelectMgr(name, ctrl)
		{
		if name is ''
			name = ctrl.Send('OverrideSelectManager?') is true
				? ''
				: .getSelectSaveName()
		.selectMgr = AccessSelectMgr(ctrl.GetDefaultSelect(), :name)
		noUserDefaultSelects? = not .enableUserDefaultSelect
		.selectMgr.LoadSelects(ctrl, :noUserDefaultSelects?)
		}

	UserDefaultSelectEnabled?()
		{
		return .enableUserDefaultSelect
		}

	GetSelectMgr()
		{
		return .selectMgr
		}

	sf: false
	GetSelectFields(extra = #())
		{
		if .sf isnt false
			return .sf
		cols = .disableSelectFilter ? #() : .GetOriginalColumns()
		cols = cols.Copy().MergeUnion(extra)
		excludes = .disableSelectFilter ? #() : .GetExcludeSelectFields()
		return .sf = .BuildSelectFields(cols, excludes, .headerSelectPrompt)
		}

	BuildSelectFields(cols, excludes, headerSelectPrompt)
		{
		return SelectFields(cols, excludes, :headerSelectPrompt, includeMasterNum:)
		}

	HasSelectedVals?()
		{
		if .selectMgr is false
			return false
		vals = .selectMgr.Select_vals()
		return vals.Any?({ it.check is true and
			it.GetDefault(it.condition_field, []).operation isnt '' })
		}

	GetSelectVals()
		{
		return .selectMgr is false ? [] : .selectMgr.Select_vals()
		}

	GetSelectWhere(selectName, view, availableColumnsFn)
		{
		if false is mgr = .GetSelectMgr()
			.SetSelectMgr(selectName, view)
		mgr = .GetSelectMgr()
		select_vals = mgr.Select_vals()
		return .GetWhereSpecs(select_vals, availableColumnsFn).where
		}

	GetWhereSpecs(select_vals, availableColumnsFn)
		{
		sf = .GetSelectFields()
		whereSpecs = SelectRepeatControl.BuildWhere(sf, select_vals)
		whereStr = sf.Joins(whereSpecs.joinflds) $ whereSpecs.where
		whereSpecs.where = .ExtendWithWhere(whereStr, availableColumnsFn)
		return whereSpecs
		}

	suppressed: false
	SuppressSlowWarning(.suppressed = true)
		{
		}

	// need to handle fields with mandatory datadict if we allow editing
	mandatoryColumn?(column)
		{
		if .mandatoryFields.Has?(column)
			return true
		if Object?(.customFields) and .customFields.Member?(column) and
			.customFields[column].GetDefault('mandatory', false) is true
			return true
		ctrl = Datadict(column).Control
		return ctrl.Member?('mandatory') and ctrl.mandatory is true
		}

	GetMandatoryFields()
		{
		return .origColumns.Filter(.mandatoryColumn?)
		}

	headerChanged?: false
	HeaderChanged?()
		{
		return .headerChanged?
		}

	SetHeaderChanged(status /*unused*/ = true)
		{
		.headerChanged? = true
		}

	customFields: false
	customKey: false
	HasCustomFormula?: false
	setCustomKey(.customKey, query, columns)
		{
		.HasCustomFormula? = CustomizeField.HasCustomFieldFormula?(.customKey)
		if .customKey isnt false and .columnsSaveName is false
			{
			.columnsSaveName = .customKey
			.origColumns = .GetAvailableColumns(query)
			.setup()
			}
		.ensureDefaultColumns(columns)
		.customFields = Customizable.GetCustomizedFields(.customKey)
		.loadSavedCols(query)
		}

	extraColumn()
		{
		return .checkBoxColumn isnt false
			? Object(Object(field: .checkBoxColumn, pos: 0, width: 75))
			: #()
		}

	loadSavedCols(permissableQuery)
		{
		if .columnsSaveName is false
			return

		cols = .removeHiddenCols(.origColumns)
		UserColumns.Load(cols, .columnsSaveName, this, load_visible?:,
			initialized?: false, extraCols: .extraColumn(), :permissableQuery,
			hideColumnsNotSaved?: .hideColumnsNotSaved?)
		}

	removeHiddenCols(cols)
		{
		newCols = cols
		if .customFields isnt false
			{
			hiddenCols = Object()
			for m, v in .customFields
				if v.GetDefault('hidden', false)
					hiddenCols.Add(m)
			newCols = newCols.Difference(hiddenCols)
			}
		return newCols
		}

	AddMissingMandatoryCols()
		{
		for col in .GetMandatoryFields().Difference(.columns)
			{
			.columns.Add(col)
			.formatting.AddColumnsFmt(col)
			.widths.Add(100 /*= default width */)
			.headerChanged? = true
			}
		}

	GetCustomKey()
		{
		return .customKey
		}

	GetExcludeSelectFields()
		{
		return .excludeSelectFields
		}

	GetCustomFields()
		{
		return .customFields
		}

	FirstTime: true
	CustomizeColumns(ctrl, query, editable?)
		{
		availableColumns = .hideCustomColumns?
			? .origColumns
			: .GetAvailableColumns(query)
		extraMandatoryColumns = Object()
		if 0 isnt extraMandatory = ctrl.Send(
			'VirtualList_CustomizeColumnAdditionalMandatoryList')
			extraMandatoryColumns.Append(extraMandatory)

		allowHideMandatory? =
			ctrl.Send('VirtualList_CustomizeColumnAllowHideMandatory?') is true
		mandatoryCols = editable? and not allowHideMandatory?
			? .GetMandatoryFields().Append(extraMandatoryColumns)
			:  extraMandatoryColumns
		if not mandatoryCols.Subset?(availableColumns)
			SuneidoLog('ERROR: mandatory columns not included in Virtual List columns',
				calls:)

		removeOb = Customizable.GetNonPermissableFields(query)
		columns = .removeHiddenCols(availableColumns).Remove(@removeOb)
		.FirstTime = true
		CustomizeColumnsDialog(ctrl.Window.Hwnd, ctrl, columns, .GetColumnsSaveName(),
			mandatoryCols, .GetHeaderSelectPrompt())
		}

	GetAvailableColumns(query)
		{
		return ListCustomize.AddCustomColumns(query, .origColumns)
		}

	Customize(ctrl, query, subTitle = '', hasCustomExpand = false)
		{
		columns = .GetAvailableColumns(query)
		sfOb = Object(cols: columns, excludeFields: Object())
		table = QueryGetTable(query, nothrow:)
		if false is custom_fields = ListCustomize.InitCustomFields(table, .columns)
			return false
		dirty = CustomizeDialog(ctrl.Window.Hwnd, .customKey, query, sfOb,
			hasCustomExpand is false, /* browse? */
			Customizable(table), custom_fields, sub_title: subTitle, virtual_list?:,
			tabs: hasCustomExpand is true ? #(custom_tabs: #(Expand)) : false)
		return .handle_customize_result(dirty, custom_fields, query)
		}

	handle_customize_result(dirty, custom_fields, query)
		{
		Assert(.columnsSaveName isnt: false)
		if Object?(dirty) and dirty.screen and
			custom_fields isnt false and custom_fields.NotEmpty?()
			{
			availableColumns = .GetAvailableColumns(query)
			UserColumns.AddCustomFields(
				.columnsSaveName, this, custom_fields, .origColumns, availableColumns)
			// screen will be refreshed
			}
		return not Object?(dirty) or dirty.fields or dirty.screen
		}

	ExtendWithWhere(where, allAvailableCols)
		{
		fields = Object()
		where.ForEachMatch(' where (\()?(\w+) ')
			{ |x|
			fields.AddUnique(where[ x[2][0] :: x[2][1] ])
			}
		// need to strip where clause, because it could contain rule fields
		joins = where.BeforeFirst('where')
		nonExistFields = fields.Difference(allAvailableCols(joins))
		extend = Opt(' extend ', nonExistFields.Join(','), ' ')
		if .suppressed is true
			where $= ' /*SLOWQUERY SUPPRESS*/ '
		return extend $ where
		}

	MeasureWidth(col, rec)
		{
		return .formatting.MeasureWidth(col, rec)
		}

	Destroy()
		{
		ignoreCols = .checkBoxColumn isnt false ? Object().Add(.checkBoxColumn) : #()
		if .columnsSaveName isnt false and .headerChanged?
			UserColumns.Save(.columnsSaveName, this, .origColumns, :ignoreCols)

		if .selectMgr isnt false and .getSelectSaveName() isnt .TmpSelectName
			.selectMgr.SaveSelects()

		.formatting.Destroy()
		}
	}