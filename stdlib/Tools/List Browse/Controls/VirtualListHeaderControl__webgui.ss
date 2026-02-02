// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Name: 			'VirtualListHeader'
	model:			false
	colModel: 		false
	showSortIndicator: false

	New(headerSelectPrompt = false, .checkBoxColumn = false)
		{
		sortDisabled = .Send('VirtualListHeader_DisableSort?') is true
		if not sortDisabled
			.showSortIndicator = true
		.header = SuJsListHeader(:headerSelectPrompt)
		}

	GetLayout()
		{
		return false
		}

	showExpandBar?: false
	SetColModel(.colModel, sort, .showExpandBar?)
		{
		.grid = .Send('GetGrid')
		.grid.Act(#UpdateStretchColumn, .colModel.StretchColumn)
		.setColumnsWidth(sort)
		}

	setColumnsWidth(sort)
		{
		.header.Clear()

		columns = .colModel.GetColumns()
		for (col in columns.Members())
			{
			width = .colModel.GetColWidth(col)
			_capFieldPrompt = .colModel.CapFieldPrompt(columns[col])
			sortParam = false
			if sort.GetDefault('displayCol', sort.col) is columns[col]
				sortParam = sort.dir
			.header.AddItem(columns[col], width, sort: sortParam)
			format = .colModel.GetColFormat(col)
			.header.SetItemFormat(col, format)
			if width is false
				.colModel.SetColWidth(col, .header.GetItemWidth(col))
			}
		.grid.Act(#UpdateHead, .header.Get(), .showExpandBar?,
			showSortIndicator: .showSortIndicator)
		}

	ResetColumns(permissableQuery, sort)
		{
		.colModel.ResetColumns(permissableQuery)
		if .Send('Editable?') is true and
			.Controller.Send('VirtualList_CustomizeColumnAllowHideMandatory?') isnt true
			.colModel.AddMissingMandatoryCols()
		.setColumnsWidth(sort)
		.Send('VirtualListHeader_ResetColumns')
		}

	RefreshSort(sort)
		{
		.setColumnsWidth(sort)
		}

	HeaderResize(col, width)
		{
		minWidth = .mandatoryCols.Has?(.colModel.Get(col)) ? .mandatoryColMinWidth : 0
		width = Max(width, minWidth)
		.header.SetItemWidth(col, width)
		.colModel.SetColWidth(col, width)
		.colModel.SetHeaderChanged()
		.grid.Act("SetColWidth", col, width)
		}

	HeaderDividerDoubleClick(col, width)
		{
		if .colModel.StretchColumn isnt false and col is .colModel.GetLastVisibleCol()
			return
		field = .colModel.Get(col)
		width = .getMinWidth(field, width)
		.colModel.SetColWidth(col, width)
		.header.SetItemWidth(col, width)
		.colModel.SetHeaderChanged()
		.grid.Act("SetColWidth", col, width)
		}

	mandatoryColMinWidth: 50
	getMinWidth(field, width)
		{
		if width is false or width is 0
			return .mandatoryColMinWidth
		width = Min(width, 1000) /*= max width */
		if .mandatoryCols.Has?(field)
			width = Max(.mandatoryColMinWidth, width)
		return width
		}

	getter_mandatoryCols()
		{
		return .mandatoryCols = .Send('GetMandatoryCols') // once only
		}

	BuildSortMenu()
		{
		menu =  Object()
		if .showSortIndicator
			{
			saveSort? = .Send('SaveSort?')
			sort = .Send('GetPrimarySort')
			if sort.col is false or saveSort? is false
				menu.Add(Object(name: 'Set as Default Sort for Current User',
					state: MFS.DISABLED))
			else
				{
				col = sort.GetDefault('displayCol', sort.col)
				if Internal?(col)
					menu.Add(Object(name: 'Set as Default Sort for Current User',
						state: MFS.DISABLED))
				else
					menu.Add(Object(name: 'Set (' $ .header.GetHeaderText(col) $
						') as Default Sort for Current User',
						cmd: 'Set as Default Sort'))
				}
			menu.Add('Reset Sort to System Default')
			}
		return menu
		}

	MoveColumnToFront(col)
		{
		if false is colIdx = .colModel.FindCol(col)
			return
		if colIdx is 0 // already at the front
			return
		swapIdx = .checkBoxColumn isnt false and .colModel.FindCol(.checkBoxColumn) is 0
			? 1
			: 0
		if colIdx is swapIdx
			return

		.header.Reorder(swapIdx, colIdx)
		.colModel.ReorderColumn(colIdx, swapIdx)
		.Send('VirtualListHeader_HeaderReorder')
		.colModel.SetHeaderChanged()
		}

	Default(@unused) {	}
	}
