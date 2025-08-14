// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
WndProc
	{
	Name: 		"List"
	Xmin: 		100
	Ymin: 		50
	Xstretch:	1
	Ystretch:	1

	New(.columns = false, .style = 0, .exStyle = 0, .stretch = false)
		{
		.CreateWindow(WC_LISTVIEW, "", .style | WS.VISIBLE, WS_EX.CLIENTEDGE)
		.SetExtendedStyle(.exStyle)
		.SubClass()
		.SetFont()

		// make sure it's on top
		SetWindowPos(.Hwnd, HWND.TOP, .x, .y, .w, .h, 0)

		.cols = Object()
		.ncols = .nrows = 0
		.EnsureVisible(0)
		}
	DeleteAll()
		{
		.nrows = 0
		SendMessage(.Hwnd, LVM.DELETEALLITEMS, 0, 0)
		}
	SetStyle(.style)
		{
		SetWindowLong(.Hwnd, GWL.STYLE, WS.CHILD | WS.VISIBLE | .style)
		}
	SetExtendedStyle(.exStyle)
		{
		SendMessage(.Hwnd, LVM.SETEXTENDEDLISTVIEWSTYLE, 0, exStyle | LVS_EX.DOUBLEBUFFER)
		}
	AddItem(label, image = 0, lParam = 0)
		{
		item = Object(
			mask: LVIF.TEXT | LVIF.IMAGE | LVIF.PARAM,
			pszText: label,
			iImage: image,
			iItem: .nrows,
			iSubItem: 0,
			:lParam)
		SendMessageListItem(.Hwnd, LVM.INSERTITEM, 0, item)
		return .nrows++
		}
	CheckAll(checked = true)
		{
		for (i = 0; i < .nrows; ++i)
			.SetCheckState(i, checked)
		}
	stateImageIndexBits: 12
	SetCheckState(i, state)
		{
		if (state is '')
			state = LVCS.DISABLE
		else
			state =  (state is true) ? LVCS.CHECKED : LVCS.UNCHECKED
		item = Object(
			mask: LVIF.STATE,
			state: state << .stateImageIndexBits,
			stateMask: LVIS.STATEIMAGEMASK)
		SendMessageListItem(.Hwnd, LVM.SETITEMSTATE, i, item)
		}
	GetCheckState(i)
		{
		state = SendMessage(.Hwnd, LVM.GETITEMSTATE, i, LVIS.STATEIMAGEMASK)
		state = state >> .stateImageIndexBits
		return LVCS.CHECKED is state
		}
	GetColWidth(col)
		{
		return SendMessage(.Hwnd, LVM.GETCOLUMNWIDTH, col, 0)
		}
	SetMaxWidth(column)
		{
		if String?(column)
			column = .columnNameToIndex(column)
		maxWidth = 0
		for rec in .ToObject()
			maxWidth = Max(maxWidth, .TextExtent(rec[.cols[column]]).x)
		SendMessage(.Hwnd, LVM.SETCOLUMNWIDTH, column,
			maxWidth + .headerBorder + 5 /*= inner content padding*/)
		}
	SetColWidth(col, width)
		{
		if width is false
			width = .getWidth(.cols[col], .cols[col])
		SendMessage(.Hwnd, LVM.SETCOLUMNWIDTH, col, width)
		}
	SetColumnWidth(column, width)
		{
		if String?(column)
			column = .columnNameToIndex(column)
		if (column is false)
			return
		SendMessage(.Hwnd, LVM.SETCOLUMNWIDTH, column, width)
		}
	columnNameToIndex(column)
		{
		column_index = false
		for (i = 0; i < .ncols; ++i)
			{
			if (.cols[i] is column)
				{
				column_index = i
				break
				}
			}
		return column_index
		}
	Addrow(ob, lParam = 0)
		{
		// Note: All the columns must exist, or the object must have a default.
		item_index = .nrows
		item = Object(
			mask: LVIF.TEXT | LVIF.PARAM,
			pszText: "label",
			iItem: item_index,
			iSubItem: 0,
			:lParam)
		SendMessageListItem(.Hwnd, LVM.INSERTITEM, 0, item)
		for (i = 0; i < .ncols; ++i)
			{
			item = Object(
				mask: LVIF.TEXT,
				pszText: .display(ob[.cols[i]]),
				iItem: item_index,
				iSubItem: i)
			SendMessageListItem(.Hwnd, LVM.SETITEM, 0, item)
			}
		.SelectItem(.nrows)
		.nrows++
		return item_index
		}
	AddColumn(name)
		{
		header = PromptOrHeading(name)
		header = header is "" ? name : header
		width = .getWidth(name, header)
		col = Object(
			mask: LVCF.FMT | LVCF.WIDTH | LVCF.TEXT | LVCF.SUBITEM,
			fmt: LVCFMT.LEFT,
			cx: width,
			pszText: header,
			iSubItem: .ncols)
		SendMessageListColumn(.Hwnd, LVM.INSERTCOLUMN, .ncols, col)
		.cols[.ncols] = name
		++.ncols
		}
	headerBorder: 20 // determined by experimentation
	maxFieldWidth: 300
	getWidth(field, text)
		{
		heading_width = text is "" ? 0 : .TextExtent(text).x
		hdc = GetDC(.Window.Hwnd)
		format_width = FieldFormatWidth(field, .GetAveCharWidth(), hdc)
		ReleaseDC(.Window.Hwnd, hdc)

		return Min(Max(heading_width, format_width) + .headerBorder, .maxFieldWidth)
		}
	GetAveCharWidth()
		{
		hdc = GetDC(.Window.Hwnd)
		GetTextMetrics(hdc, tm = Object())
		ReleaseDC(.Window.Hwnd, hdc)
		return tm.AveCharWidth
		}
	GetRecord(item)
		{
		item_object = Object()
		for (i = 0; i < .ncols; ++i)
			item_object[.cols[i]] = SendMessageListItemOut(.Hwnd, item, i)
		return item_object
		}
	SelectItem(i)
		{
		item = Object()
		item.stateMask = LVIS.SELECTED | LVIS.FOCUSED
		item.state = LVIS.SELECTED | LVIS.FOCUSED
		SendMessageListItem(.Hwnd, LVM.SETITEMSTATE, i, item)
		}
	UnSelectItem(i)
		{
		item = Object()
		item.stateMask = LVIS.SELECTED
		item.state = 0
		SendMessageListItem(.Hwnd, LVM.SETITEMSTATE, i, item)
		}
	UnSelect()
		{
		.UnSelectItem(.GetSelected())
		}
	EnsureVisible(i)
		{
		SendMessage(.Hwnd, LVM.ENSUREVISIBLE, i, false)
		}
	GetSelected()
		{
		result = SendMessage(.Hwnd, LVM.GETNEXTITEM, -1, LVNI.SELECTED)
		return (result > -1) ? result : false
		}
	display(value, limit = 500)
		{
		value = not String?(value) ? Display(value) : value
		ellipsis = ' ...'
		if value.Size() > limit
			value = value[.. limit - ellipsis.Size()] $ ellipsis
		return value
		}
	ToObject()
		{
		listob = Object()
		for (i = 0; i < .nrows; ++i)
			listob[i] = .GetRecord(i)
		return listob
		}
	x: 0, y: 0, w: 0, h: 0
	Resize(.x, .y, .w, .h)
		{
		super.Resize(x, y, w, h)
		.adjust()
		}
	adjust()
		{
		if .stretch is false
			return
		else if .stretch is true
			col = .ncols - 1
		else if String?(.stretch)
			col = .columns.Find(.stretch)
		else
			col = .stretch
		GetClientRect(.Hwnd, rect = Object())
		w = rect.right - rect.left
		for (i = 0; i < .ncols; ++i)
			if i isnt col
				w -= .GetColWidth(i)
		minColWidth = 30
		w = Max(minColWidth, w)

		// need the "if" to avoid infinite recursion
		if .GetColWidth(col) isnt w
			.SetColWidth(col, w)
		}
	}
