// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Addon_VirtualListViewBase
	{
	LoadSavedFilters(view)
		{
		model = view.GetModel()
		if model is false
			return
		if .fromSmallDialog(view)
			{
			.addTopExtraButtons(view)
			FilterButtonControl.UpdateStatus(view, model.ColModel.HasSelectedVals?())
			.UpdateTotalSelected()
			return
			}
		openState = UserSettings.Get(view.Option $ ' - Split Open', 'no_default')
		select = view.GetDefaultSelect()
		if openState is true or (openState is 'no_default' and select.Size() > 0)
			.OpenFilters(view)
		else
			.addTopExtraButtons(view)
		FilterButtonControl.UpdateStatus(view, model.ColModel.HasSelectedVals?())
		// to ensure extra space for creating new record
		if model.GetInitStartLast() is true
			view.Defer(view.On_VirtualListThumb_ArrowDown,
				uniqueID: 'VirtualListLoadFilter')
		.UpdateTotalSelected()
		}

	fromSmallDialog(view)
		{
		if not (view.Window.Base?(ModalWindow) or view.Window.Base?(Dialog))
			return false
		rect = GetWindowRect(view.Window.Hwnd)
		return rect.bottom - rect.top < ScaleWithDpiFactor(480) /*= small dialog */
		}

	addTopExtraButtons(view)
		{
		filtersWrapper = view.FindControl('select')
		extraLayout = view.Select_ExtraLayout()
		filtersWrapper.Append(Object('VirtualListTopLayout', extraLayout))
		}

	ToggleFilter(view)
		{
		topFilters = .getTopFilters(view)
		if .fromSmallDialog(view)
			{
			if topFilters isnt false
				.removeFilters(view, topFilters)
			view.On_VirtualListThumb_ArrowSelect()
			}
		else if topFilters is false
			.OpenFilters(view)
		else
			if .Select_Apply()
				.removeFilters(view, topFilters)
		.UpdateTotalSelected()
		}

	getTopFilters(view)
		{
		return view.FindControl(.selectRepeatName)
		}

	OpenFilters(view)
		{
		grid = view.GetGrid()
		filtersWrapper = view.FindControl('select')  // possibly call once ???
		filtersWrapper.Remove(0)
		rect = GetWindowRect(grid.Hwnd)
		scroll = view.FindControl('VirtualListScroll')
		scroll.Ymin = (rect.bottom - rect.top) / 2
		UserSettings.Put(view.Option $ ' - Split Open', true)
		colModel = view.GetModel().ColModel

		if .Send('UseSubTableFilters?') is true
			filtersWrapper.Append(
				Object('SelectRepeatSubtables', view, colModel, .selectRepeatName))
		else
			filtersWrapper.Append(Object('SelectRepeat',
				view.GetSelectFields(), view.Select_vals, .selectRepeatName,
				option: view.Option, title: view.GetTitle(), fromFilter:,
				selChanged: view.GetDefault('SelectChanged?', false),
				noUserDefaultSelects?: not colModel.UserDefaultSelectEnabled?()))

		split = view.FindControl('VertSplit')
		split.UpdateSplitter()
		if false is split.SetSplitSaveName(view.Option) // no default
			split.MaximizeSecond()
		if .Model.CheckAboveSortLimit?()
			{
			filtersWrapper.FindControl('buttons').
				Insert(0, Object('VirtualListSortWarningButton', .Parent))
			if false isnt warnButton = .FindControl('VirtualListSortWarningButton')
				.setWarnButtonState(warnButton)
			}
		view.AfterTopFilter("open")
		}

	getter_selectRepeatName()
		{
		return .selectRepeatName = .Model.ColModel.GetSelectMgr().Name()
		}

	setWarnButtonState(warnButton)
		{
		if .showSortWarning?()
			warnButton.InsertWarning()
		else
			warnButton.RemoveWarning()
		}

	showSortWarning?()
		{
		x = .Parent.FindControl(.selectRepeatName).Get()
		if x.Member?('Header')
			x = x.Header

		// do NOT need to worry about SubTable filters here
		// the SortQuery only cares about filters on index fields
		// of the base query (i.e. it cannot use indexes from the sub tables,
		// so it will end up ignoring any filters on the sub tables anyway)
		show? = .Model.QueryAboveSortLimit?(x.conditions)
		return .Model.SetOverSortLimit(show?)
		}

	removeFilters(view, topFilters)
		{
		split = view.FindControl('VertSplit')
		view.SelectChanged? = topFilters.SelectChanged?()
		split.SaveSplit()
		UserSettings.Put(view.Option $ ' - Split Open', false)
		filtersWrapper = view.FindControl('select')  // possibly call once ???
		filtersWrapper.RemoveAll()
		.addTopExtraButtons(view)
		split.UpdateSplitter(remove:)
		scroll = view.FindControl('VirtualListScroll')
		scroll.Parent.Ystretch = 2
		view.AfterTopFilter("close")
		}

	On_Count()
		{
		if .Filtersontop is false or
			(.topFilters isnt false and .Select_Apply()) or
			(.Filtersontop and .topFilters is false)
			AccessControl.GetCount(.GetTitle(), .GetQuery())
		}

	Select_Apply()
		{
		if false is .Send('VirtualList_ApplySelect?')
			return false
		if .SaveFirst() is false
			return false
		if false is where = .filtersWhere()
			return false
		_slowQueryLog = Object(logged: false, from: 'Select_Apply')

		// need to handle conditions being an object of objects if topFilters
		// is a SelectRepeatSubtables
		x = .topFilters.Get()
		if x.Member?('Header')
			conditions = x['Header'].conditions
		else
			conditions = x.conditions
		if .checkSlowSelect(conditions, where) is false
			return false
		return .applySelect(x, where)
		}

	checkSlowSelect(conditions, where)
		{
		.Model.ColModel.SuppressSlowWarning(false)
		queryState = .initializeQueryState(:conditions, :where)
		return .Model.CheckSlowQuery(queryState, .afterSlowWarning)
		}

	initializeQueryState(@queryState)
		{
		if 0 is (queryState.presets = .Send('VirtualList_SlowQueryPresets'))
			queryState.presets = #()
		return queryState.Set_default(false)
		}

	afterSlowWarning(windowResult)
		{
		queryState = windowResult
		.Model.ColModel.SuppressSlowWarning(queryState.filter is false)
		if queryState.filter isnt false
			.addIndexedFilter(queryState.filter)
		else
			.applySelect(queryState, queryState.where) // continue even if slow
		}

	addIndexedFilter(filter)
		{
		if .topFilters is false
			.ToggleFilter(.Parent)
		SlowQuery.AddIndexedFilter(filter, .topFilters)
		}

// TODO: 1034 - conditions needs to be cleaned up, object is inconsistant
	applySelect(conditions, where)
		{
		changed = .topFilters.SelectChanged?()
		preWhere = .GetCurrentSelectWhere()
		preSelectVals = .Select_vals.DeepCopy()
		try
			{
			if not conditions.Empty?()
				{
				if conditions.Member?('Header')
					{
					.SetSelectVals(conditions.Header.conditions)
					.SetExtraSelectVals(conditions.Delete('Header'))
					}
				else
					.SetSelectVals(conditions.conditions)
				}

			curSelectVals = .Select_vals.DeepCopy()
			if false isnt .Send('VirtualList_BeforeApplySelect', where)
				.SetWhere(where)
			else
				{
				.UpdateTopFilters(curSelectVals)
				.SetSelectVals(curSelectVals)
				}
			}
		catch (unused, '*regex')
			{
			.SetSelectVals(preSelectVals)
			.SetWhere(preWhere)
			.AlertInfo('Select', 'Invalid matcher')
			return false
			}

		.topFilters.SetSelectApplied(true)
		if changed
			.Model.ClearStickyFieldValues()
		.SelectControl_Changed()
		.UpdateTotalSelected(true)
		return true
		}

	filtersWhere()
		{
		if .topFilters.Base?(SelectRepeatControl)
			{
			if false is where = .topFilters.Where()
				return false
			return .GetSelectFields().Joins(where.joinflds) $ where.where
			}
		return .topFilters.Where(.GetSelectFields())
		}

	getter_topFilters()
		{
		return .FindControl(.selectRepeatName)
		}

	Getter_Select_vals()
		{
		return .Model.ColModel.GetSelectVals()
		}

	Select_OpenDialog()
		{
		conditions = .topFilters.Get().conditions
		.SetSelectVals(conditions)
		.On_VirtualListThumb_ArrowSelect()
		}

	Filtersontop: false
	VirtualListHeader_HeaderClick(col)
		{
		if not .SaveOutstandingChanges()
			return
		if false is .Send("VirtualList_AllowSort", :col)
			return
		if not .Model.CheckSortable(col)
			return
		if .checkSlowSort(col, .Filtersontop) is false
			return
		.applySort(col)
		}

	checkSlowSort(col, filtersOnTop)
		{
		if filtersOnTop is false
			return true

		conditions = .Model.ColModel.GetSelectVals() // use existing conditions

		if not .Model.SortLimitChecked?() // query not checked yet
			.Model.SetOverSortLimit(.Model.QueryAboveSortLimit?(conditions))

		if .Model.OverSortLimit?()
			{
			.FixDisabledSort()
			return false
			}
		.Model.ColModel.SuppressSlowWarning(false)
		queryState = .initializeQueryState(:conditions, sortCol: col)
		return .Model.CheckSlowQuery(queryState, .afterSlowSortWarning)
		}

	afterSlowSortWarning(windowResult)
		{
		queryState = windowResult
		.Model.ColModel.SuppressSlowWarning(queryState.filter is false)
		if queryState.filter isnt false
			.addIndexedFilter(queryState.filter)
		else
			.applySort(queryState.sortCol) // continue even if slow
		}

	applySort(displayCol)
		{
		.ClearSelect()
		// there are scenarios where the column the user clicked on is not actualy
		// the column we want to sort on. need to look that up and pass it down to
		// VirtualListSortModel (as well as the column the user clicked on.)
		if 0 is dataCol = .Send('VirtualList_GetSortCol', col: displayCol)
			dataCol = false

		// displayCol is the column the user clicked on
		// dataCol - if it exists - is the column we actually want to sort on
		posPreSort = .Model.GetPosition()
		.Send('VirtualList_BeforeSort')
		.Model.SetSort(displayCol, dataCol)
		ctrls = .GetViewControls()
		ctrls.header.RefreshSort(.GetPrimarySort())
		if posPreSort isnt pos = .Model.GetPosition()
			if posPreSort isnt 'middle'
				.Model.SetStartLast(posPreSort isnt 'top')
			else
				ctrls.thumb.SetThumbPosition(pos)
		.Grid.Repaint()
		ctrls.expandBar.ShowEditButtons()
		}

	FixDisabledSort()
		{
		msg = 'SORT DISABLED! Use one of the following fields or presets to reduce the ' $
			'number of records read through in order to allow sorting'

		indexes = .Model.SelectableIndexes()
		queryState = .initializeQueryState(:msg, :indexes)
		SlowQuery.SuggestionWindow(queryState, .afterFixSlowDisabled)
		}

	UpdateTotalSelected(recalc = false)
		{
		if .Model.CheckBoxColModel is false
			return
		if false isnt totalCtrl = .FindControl("totalSelected")
			totalCtrl.Set(.Model.CheckBoxColModel.GetSelectedTotal(:recalc))
		}

	afterFixSlowDisabled(windowResult)
		{
		queryState = windowResult
		if queryState.filter isnt false
			.addIndexedFilter(queryState.filter)
		}

	SelectControl_SetSelectApplied(validFunc)
		{
		if false isnt (warnButton = .FindControl('VirtualListSortWarningButton')) and
			false isnt validFunc(quiet:)
			.setWarnButtonState(warnButton)
		}
	}