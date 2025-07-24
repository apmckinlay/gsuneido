// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
WndProc
	{
	Xstretch: 		1
	Ystretch:		0
	Ymin: 			0
	Name: 			'VirtualListHeader'
	colModel: 		false

	New(headerSelectPrompt = false, .checkBoxColumn = false)
		{
		.CreateWindow("SuBtnfaceArrow", "", WS.VISIBLE)
		.SubClass()

		sortDisabled = .Controller.Send('VirtualList_DisableSort?') is true
		if not sortDisabled
			.showSortIndicator = true

		style = HDS.DRAGDROP
		if not sortDisabled
			style |= HDS.BUTTONS
		.header = .Construct(Object('Header', :style, :headerSelectPrompt))
		.Ymin = .header.Ymin
		}

	Startup()
		{
		.header.Startup()
		}

	SetColModel(colModel, sort)
		{
		.colModel = colModel
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
			.header.AddItem(columns[col], width)
			.setColFormat(col, sort, columns)
			if width is false
				.colModel.SetColWidth(col, .header.GetItemWidth(col))
			}
		.synchHeader()
		}

	synchHeader()
		{
		offset = .resize()

		if (offset > 0)
			{
			rect = Object(left: 0 , top: 0, right: offset, bottom: .Ymin)
			FillRect(hdc = GetDC(.Hwnd), rect, GetSysColorBrush(COLOR.BTNFACE))
			ReleaseDC(.Hwnd, hdc)
			}
		}

	RefreshSort(sort)
		{
		if not .showSortIndicator
			return
		columns = .colModel.GetColumns()
		for col in columns.Members()
			.setColFormat(col, sort, columns)
		}

	showSortIndicator: false
	setColFormat(col, sort, columns)
		{
		// sometimes the column we are showing as the sort is not the column we are
		// actually sorting on. sort.col is what we are sorted on, but if we
		// have sort.displayCol, then that is what we want to show as the sorted col
		// instead
		sortCol = sort.GetDefault('displayCol', sort.col)
		fmt = .colModel.GetColFormat(col)
		if .showSortIndicator and sortCol is columns[col] and sortCol isnt false
			fmt |= sort.dir is 1 ? HDF.SORTUP : HDF.SORTDOWN
		.header.SetItemFormat(col, fmt)
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

	ScrollClientRect(movePix)
		{
		rect = .GetClientRect().ToWindowsRect()
		ScrollWindowEx(.Hwnd, movePix, 0, rect, rect, NULL, NULL, SW.INVALIDATE)
		.synchHeader()
		}

	Resize(x, y, w, h)
		{
		super.Resize(x, y, w, h)
		.resize()
		}

	resize()
		{
		if .colModel is false
			return 0

		.adjust()
		offset = -.colModel.Offset
		width = Max(.GetClientRect().GetWidth(), .header.Xmin)
		.header.Resize(offset, 0, width, .Ymin)
		return offset
		}

	Header_AllowTrack(col)
		{
		return .colModel.StretchColumn is false or col isnt .colModel.GetLastVisibleCol()
		}

	HeaderTrack(col, width)
		{
		if 0 is movePix = width - .colModel.GetColWidth(col)
			return

		resizeLeftSide? = false
		if movePix < 0 and .colModel.Offset > 0
			{
			hEnd = .header.GetItemRect(.colModel.GetSize() - 1).right - .colModel.Offset
			wndEnd = .GetRect().GetX2()
			if wndEnd > hEnd
				{
				if movePix < -.colModel.Offset
					movePix = -.colModel.Offset
				.colModel.Offset += movePix
				resizeLeftSide? = true
				.synchHeader()
				}
			}

		.colModel.SetHeaderChanged()
		.Send('VirtualListHeader_HeaderTrack', col, width, movePix, :resizeLeftSide?)
		.adjust(col)
		}

	adjust(col = false)
		{
		if .colModel.StretchColumn is false
			return

		GetClientRect(.Hwnd, rect = Object())
		if 0 is w = rect.right - rect.left
			return

		if col is false and .colModel.FirstTime // first time or customize column
			{
			.fitWindowSize(w, .colModel.GetWidths(), .header)
			.colModel.FirstTime = false
			}
		stretchCol = .colModel.GetStretchCol(col)
		otherwid = 0
		for (i = 0; i < .colModel.GetSize(); ++i)
			if i isnt stretchCol
				otherwid += .colModel.GetColWidth(i)
		newWidth = Max(w - otherwid - 1, 30) /*= min width when stretching*/
		.header.SetItemWidth(stretchCol, newWidth)
		.colModel.SetColWidth(stretchCol, newWidth)
		if col isnt false
			.Send('RepaintGrid')
		}

	fitWindowSize(w, widths, header)
		{
		totalWidth = widths.Sum()
		if totalWidth > w
			{
			ratio = w / totalWidth
			widths.Map!({ (it *  ratio).Int() })
			for (i = 0; i < widths.Size(); ++i)
				header.SetItemWidth(i, widths[i])
			}
		}

	HeaderResize(col, width)
		{
		.colModel.SetColWidth(col, width)
		}

	HeaderReorder(oldIdx, newIdx)
		{
		if oldIdx isnt newIdx
			.colModel.ReorderColumn(oldIdx, newIdx)
		.Send('VirtualListHeader_HeaderReorder')
		.colModel.SetHeaderChanged()
		}

	HeaderClick(iItem, iButton /*unused*/)
		{
		col = .colModel.Get(iItem)
		.Send('VirtualListHeader_HeaderClick', col)
		}

	HeaderDividerDoubleClick(iItem)
		{
		if .colModel.StretchColumn isnt false and iItem is .colModel.GetLastVisibleCol()
			return
		field = .colModel.Get(iItem)
		width = .getMinWidth(iItem, field)
		.colModel.SetColWidth(iItem, width)
		.header.SetItemWidth(iItem, width)
		.colModel.SetHeaderChanged()
		.adjust(iItem)
		.Send('RepaintGrid', checkTotal:)
		}

	getMinWidth(iItem, field)
		{
		width = .Send('VirtualListHeader_MeasureWidth', colIdx: iItem, :field)
		if width is false or width is 0
			return .mandatoryColMinWidth
		width = Min(width, 1000) /*= max width */
		if .mandatoryCols.Has?(field)
			width = Max(.mandatoryColMinWidth, width)
		return width
		}

	mandatoryColMinWidth: 50
	Header_TrackMinWidth(colIdx)
		{
		if .mandatoryCols.Has?(.colModel.Get(colIdx))
			return .mandatoryColMinWidth
		return 0
		}

	getter_mandatoryCols()
		{
		return .mandatoryCols = .Send('GetMandatoryCols') // once only
		}

	MoveColumnToFront(col)
		{
		if false is colIdx = .colModel.FindCol(col)
			return
		if colIdx is 0 // already at the front
			return
		swapIdx = 0
		if .checkBoxColumn isnt false and .colModel.FindCol(.checkBoxColumn) is 0
			swapIdx = 1
		.header.SwapItems(swapIdx, colIdx)

		.HeaderReorder(colIdx, swapIdx)
		}

	CONTEXTMENU(lParam)
		{
		if .colModel is false
			return 0

		SetFocus(.Hwnd)
		x = LOSWORD(lParam)
		y = HISWORD(lParam)
		ScreenToClient(.Hwnd, pt = Object(:x, :y))
		col = .colModel.GetColByX(pt.x)

		extraMenu = .buildSortMenu()

		.Send('VirtualListHeader_ContextMenu', col, x, y, :extraMenu)
		return 0
		}

	buildSortMenu()
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

	Destroy()
		{
		.header.Destroy()
		super.Destroy()
		}
	}
