// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// TODO: reset tooltips if user scrolls
// TODO: reset tooltips if user resizes column
WndProc
	{
	Name: 		"List"
	Xmin: 		100
	Ymin: 		50
	Xstretch:	1
	Ystretch:	1
	use_menu: false
	hdr_hwnd: 0
	columns: false

	New(.query = false, columns = false, display = false,
		.reverse = false, .no_prompts = false,
		.font = "", .size = "", .realtime_scrolling = true, .handle_leftclick = false,
		.style = 0, .exStyle = 0, .doubleclick_edit = false, .column_widths = false,
		.allowDrag = false, .stretch_to_fill = false, .usercolumns = false,
		.tips = false, menu = ('Reset Columns'), headers = false)
		{
		.FieldMap = Object()
		.display? = display is true
		if menu isnt false
			.SetMenu(menu)
		// set up hidden columns object
		.hidecolumns = Object()
		if (query isnt false) // Virtual listview
			{
			.style = .style | LVS.OWNERDATA | LVS.REPORT

			if (query > "")
				{
				.create()
				.setup_query(columns, headers)
				}
			else
				{
				.model = false
				.create()
				.SetListSize(.nrows = 0)
				}
			}
		else
			{
			.model = false
			.columns = columns
			.create()
			}
		.olditem = .oldindex = .newitem = false
		.EnsureVisible(0)
		.setup_header()

		if .usercolumns isnt false and .columns isnt false
			UserColumns.Load(.columns, .usercolumns, this)

		// tooltip offsets by trial & error - may not always be right
		.offset = [x: 6, y: 2]
		}
	setup_header()
		{
		if 0 isnt hwnd = .GetHeader()
			{
			.AddHwnd(.hdr_hwnd = hwnd)
			.hdr = ListViewHeader(this) // to get header tooltips
			}
		}
	create()
		{
		.CreateWindow(WC_LISTVIEW, "", .style | WS.VISIBLE, WS_EX.CLIENTEDGE)
		.SetExtendedStyle(.exStyle)
		.SubClass()
		.SetFont(.font, .size)

		// make sure it's on top
		SetWindowPos(.Hwnd, HWND.TOP, .x, .y, .w, .h, 0)

		.initializeNotificationMap()

		.cols = Object()
		.ncols = .nrows = 0

		if .tips is true or .tips is 'auto'
			{
			maxTipWidth = 400
			.Map[TTN.SHOW] = 'TTN_SHOW'
			.tip = .Construct(ToolTipControl)
			.tip.SendMessage(TTM.SETMAXTIPWIDTH, 0, maxTipWidth)
			.tip.Activate(false)
			.tip.AddTool(.Hwnd, '???')
			.SetRelay(.tip.RelayEvent)
			.Defer(.set_tip_font)
			}
		}
	initializeNotificationMap()
		{
		.Map = Object()
		.Map[NM.CLICK] = 'NM_CLICK'
		.Map[NM.RCLICK] = 'NM_RCLICK'
		.Map[NM.RETURN] = 'NM_RETURN'
		.Map[NM.DBLCLK] = 'NM_DBLCLK'
		.Map[LVN.ITEMCHANGED] = 'ITEMCHANGED'
		.Map[LVN.BEGINDRAG] = 'BEGINDRAG'
		.Map[LVN.GETDISPINFO] = 'GETDISPINFO'
		.Map[LVN.ODCACHEHINT] = 'LVN_ODCACHEHINT'
		.Map[LVN.KEYDOWN] = 'LVN_KEYDOWN'
		.Map[LVN.COLUMNCLICK] = 'LVN_COLUMNCLICK'
		.Map[HDN.BEGINTRACK] = 'HDN_BEGINTRACK'
		.Map[HDN.BEGINTRACKW] = 'HDN_BEGINTRACK'
		.Map[HDN.ENDTRACK] = 'HDN_ENDTRACK'
		.Map[HDN.ENDTRACKW] = 'HDN_ENDTRACK'
		.Map[HDN.ITEMCHANGED] = 'HDN_ITEMCHANGED'
		.Map[HDN.ITEMCHANGEDW] = 'HDN_ITEMCHANGED'
		.Map[HDN.ITEMCLICK] = 'HDN_ITEMCLICK'
		.Map[HDN.ITEMCLICKW] = 'HDN_ITEMCLICK'
		}
	set_tip_font()
		{
		if .tip isnt false
			.tip.SendMessage(WM.SETFONT, .SendMessage(WM.GETFONT), false)
		}
	tip: false
	Reset()
		{
		.close_model()
		if .tip isnt false
			{
			.tip.DestroyWindow()
			.tip = false
			}
		.DestroyWindow()
		.remove_header()

		.create()
		.previous_columns = false
		.setup_header()
		}
	hdr: false
	remove_header()
		{
		if .hdr isnt false
			{
			.hdr.Destroy()
			.hdr = false
			}
		if .hdr_hwnd isnt 0
			.DelHwnd(.hdr_hwnd)
		.hdr_hwnd = 0
		}
	setup_query(columns = false, headers = false, ignoreColumns? = false)
		{
		// model provides data for listview
		.close_model()
		.model = ListViewModelCached(
			ListViewModel(.query, reverse: .reverse))
		if not ignoreColumns?
			{
			if (columns is false)
				.columns = .model.Getcolumns()
			else
				.columns = columns
			.SetColumns(.columns, headers)
			}
		.nrows = .model.Getnumrows()
		.SetListSize(.nrows)
		}
	GetQuery()
		{
		return .query
		}
	ResetQuery(query, columns = false)
		{
		if .usercolumns isnt false
			UserColumns.Save(.usercolumns, this)
		.Reset()
		.model = .olditem = .newitem = false
		.query = query
		.setup_query(columns)
		if .usercolumns isnt false
			UserColumns.Load(.columns, .usercolumns, this)
		.adjust()
		}
	ResetQueryWithoutDestroy(query)
		{
		.close_model()
		.model = .olditem = .newitem = false
		.query = query
		.setup_query(.GetColumns(), ignoreColumns?:)
		.adjust()
		}
	GetModel()
		{
		return .model
		}
	SetListSize(num_items)
		{
		SendMessage(.Hwnd, LVM.SETITEMCOUNT, num_items, 0)
		}
	DeleteAll()
		{
		.nrows = 0
		SendMessage(.Hwnd, LVM.DELETEALLITEMS, 0, 0)
		}
	SetStyle(style)
		{
		.remove_header()
		.style = style
		if (.query isnt false)
			.style = .style | LVS.OWNERDATA | LVS.REPORT
		SetWindowLong(.Hwnd, GWL.STYLE, WS.CHILD | WS.VISIBLE | .style)
		.setup_header()
		}
	// for setting extended styles
	exStyle: 0
	SetExtendedStyle(ex_style)
		{
		.exStyle = ex_style
		SendMessage(.Hwnd, LVM.SETEXTENDEDLISTVIEWSTYLE,
			0, ex_style | LVS_EX.DOUBLEBUFFER)
		}
	SetView(view)
		{
		style = GetWindowLong(.Hwnd, GWL.STYLE)
		if ((style & LVS.TYPEMASK) isnt view)
			SetWindowLong(.Hwnd, GWL.STYLE, (style & ~LVS.TYPEMASK) | view)
		}
	SetImageList(images, which)
		{
		SendMessage(.Hwnd, LVM.SETIMAGELIST, which, images)
		}
	DeleteItem(item = false)
		{
		if (item is false)
			item = .GetSelected()
		SendMessageListItem(.Hwnd, LVM.DELETEITEM, item, 0)
		.nrows--
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
	GetItemParam(index)
		// pre:	index is an integer indentifying a valid list view item
		// post:	returns lParam associated with list view item
		{
		item = Object(mask: LVIF.PARAM, iItem: index, iSubItem: 0)
		SendMessageListItem(
			.Hwnd,
			LVM.GETITEM,
			0,
			item
			)
		return item.lParam
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
	GetCheckedItems()
		{
		checked = Object()
		for (i = 0; i < .nrows; ++i)
			if .GetCheckState(i)
				checked.Add(i)
		return checked
		}
	GetCheckState(i)
		{
		state = SendMessage(.Hwnd, LVM.GETITEMSTATE, i, LVIS.STATEIMAGEMASK)
		state = state >> .stateImageIndexBits
		return LVCS.CHECKED is state
		}
	SetColumnValue(column, text, index = false)
		{
		if (index is false)
			index = .GetSelected()
		column_index = .columnNameToIndex(column)
		if (column_index is false)
			return
		item = Object(
			mask: LVIF.TEXT,
			pszText: .display(text),
			iItem: index,
			iSubItem: column_index)
		SendMessageListItem(.Hwnd, LVM.SETITEM, 0, item)
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
	GetColumnWidth(column)
		{
		if String?(column)
			column = .columnNameToIndex(column)
		if (column is false)
			return 0
		return SendMessage(.Hwnd, LVM.GETCOLUMNWIDTH, column, 0)
		}
	GetColumnWidths()
		{
		col_widths = Object()
		for (column in .cols)
			col_widths[column] = .GetColumnWidth(column)
		return col_widths
		}
	HideColumn(column)
		{
		.SetColumnWidth(column, 0)
		}
	SetColumnName(column, text)
		{
		column_index = .columnNameToIndex(column)
		if (column_index is false)
			return
		col = Object()
		col.mask = LVCF.TEXT
		col.pszText = text
		SendMessageListColumn(.Hwnd, LVM.SETCOLUMN, column_index, col)
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
	DeleteAllColumns()
		{
		for (i = .ncols - 1; i >= 0; --i)
			SendMessageListColumn(.Hwnd, LVM.DELETECOLUMN, i, 0)
		.hidecolumns = Object()
		.cols = Object()
		.ncols = 0
		}
	// SetColumns and Addrow are for report view
	previous_columns: false
	SetColumns(cols, headers = false)
		{
		if (cols is .previous_columns)
			return
		.previous_columns = cols.Copy()
		.DeleteAllColumns()
		cols = cols.Copy()
		for (n = cols.Size(), i = 0; i < n; ++i)
			{
			if (cols[i].Prefix?("-"))
				.hidecolumns.Add(cols[i] = cols[i][1 ..])
			.AddColumn(cols[i],
				headers is false ? false : headers[i])
			}
		// hide columns specified as hidden
		for (c in .hidecolumns)
			.HideColumn(c)
		}
	OrderColumns(column_list)
		{
		// the column list must contain less than 256 columns
		// (which is all SETCOLUMNORDERARRAY can handle)
		Assert(column_list.Size() < 256) /*= max columns */
		for (outer_column in column_list.Members())
			{
			// each column in the column list must be a member of .cols
			Assert( .cols.Member?( column_list[outer_column] ) )
			// each column in the column list must be distinct
			for (inner_column in column_list.Members())
				{
				if ( outer_column is inner_column )
					continue
				Assert( column_list[outer_column] isnt column_list[inner_column] )
				}
			}
		// send message
		SendMessageListColumnOrder( .Hwnd, LVM.SETCOLUMNORDERARRAY, .cols.Size(),
			Object( order: column_list.Copy() ) )
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
	AddColumn(name, header = false)
		{
		if header is false
			{
			header = PromptOrHeading(name)
			header = (header is "" or .no_prompts) ? name : header
			}
		if (Object?(.column_widths) and .column_widths.Member?(name))
			width = .column_widths[name]
		else
			width = .getWidth(name, header)
		col = Object(
			mask: LVCF.FMT | LVCF.WIDTH | LVCF.TEXT | LVCF.SUBITEM,
			fmt: LVCFMT.LEFT,
			cx: width,
			pszText: header,
			iSubItem: .ncols)
		SendMessageListColumn(.Hwnd, LVM.INSERTCOLUMN, .ncols, col)
		.FieldMap[header] = name
		.cols[.ncols] = name
		++.ncols
		}
	ResetColumnWidths()
		{
		.column_widths = false
		for (column in .cols)
			{
			if (.hidecolumns.Has?(column))
				.HideColumn(column)
			else
				.SetColumnWidth(column, .getWidth(column, column))
			}
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
	Getrow(item)
		{
		item_object = Object()
		for (i = 0; i < .ncols; ++i)
			item_object[.cols[i]] = SendMessageListItemOut(.Hwnd, item, i)
		return item_object
		}
	GetRecord(item)
		{
		return .model is false ? .Getrow(item) : .model.Getrecord(item)
		}
	GetColumns()
		{
		return .cols
		}
	FindItem(field, value)
		{
		value = .display(value)
		subitem = .cols.Find(field)
		n = .SendMessage(LVM.GETITEMCOUNT)
		for (i = 0; i < n; ++i)
			if SendMessageListItemOut(.Hwnd, i, subitem) is value
				return i
		return false
		}
	// drag support
	dragging: false
	BEGINDRAG(lParam)
		{
		if .allowDrag is false
			return 0
		lv = NMLISTVIEW(lParam)
		item = lv.iItem
		p = lv.ptAction
		image = SendMessagePoint(.Hwnd, LVM.CREATEDRAGIMAGE, item, pt = Object())
		ImageList_BeginDrag(image, 0, p.x - pt.x, p.y - pt.y)
		.lx = p.x + GetSystemMetrics(SM.CXEDGE)
		.ly = p.y + GetSystemMetrics(SM.CYEDGE)
		ImageList_DragEnter(.Hwnd, .lx, .ly)
		SetCapture(.Hwnd)
		.dragging = item
		.offset = Object(x: p.x - pt.x, y: p.y - pt.y)
		return 0
		}
	curItem: false
	curSubItem: false
	MOUSEMOVE(lParam)
		{
		x = LOSWORD(lParam)
		y = HISWORD(lParam)
		if .tips is true or .tips is 'auto'
			.update_tip(x, y)
		if (.dragging isnt false)
			ImageList_DragMove(x - .x, y - .y)
		return 'callsuper'
		}
	update_tip(x, y)
		{
		ht = .SubItemHitTest(x, y)
		if .curItem isnt ht.iItem or .curSubItem isnt ht.iSubItem
			{
			.tip.Activate(false)
			if ht.iItem isnt -1 and ht.iSubItem isnt -1 and
				false isnt text = .need_tip(ht.iItem, ht.iSubItem)
				{
				.tip.UpdateTipText(.Hwnd, text)
				.tip.Activate(true)
				}
			.curItem = ht.iItem
			.curSubItem = ht.iSubItem
			}
		}
	need_tip(row, col)
		{
		if Record() is rec = .GetRecord(row)
			return false

		text = .display(rec[.cols[col]])
		strwid = .GetStringWidth(text)
		colwid = .GetColWidth(col)
		return strwid + .padding > colwid ? text : false
		}
	padding: 12 // from trial & error
	LVN_BEGINSCROLL()
		{
		.tip.Activate(false)
		return 0
		}
	TTN_SHOW(lParam)
		{
		r = .GetSubItemRect(.curItem, .curSubItem)
		.tip.AdjustRect(false, r)
		p = [x: r.left - .offset.x, y: r.top - .offset.y]
		ClientToScreen(.Hwnd, p)
		nmhdr = NMHDR(lParam)
		SetWindowPos(nmhdr.hwndFrom, 0,
			p.x, p.y, 0, 0, // rect
			SWP.NOACTIVATE | SWP.NOSIZE | SWP.NOZORDER)
		return true
		}
	HitTest(x, y)
		{
		SendMessageLVHITTESTINFO(.Hwnd, LVM.HITTEST, 0,
			ht = Object(pt: Object(:x, :y), flags: LVHT.ONITEM))
		return ht
		}
	SubItemHitTest(x, y)
		{
		SendMessageLVHITTESTINFO(.Hwnd, LVM.SUBITEMHITTEST, 0,
			ht = Object(pt: Object(:x, :y), flags: LVHT.ONITEM))
		return ht
		}

	LBUTTONDOWN(lParam)
		{
		if (.handle_leftclick is false)
			return 'callsuper'
		if (false is .Send('ListView_MouseClick', lParam))
			return 0
		return 'callsuper'
		}
	LBUTTONUP(lParam)
		{
		if (.dragging is false)
			return 0
		ImageList_DragLeave(.Hwnd)
		ImageList_EndDrag()
		ReleaseCapture()

		dx = LOWORD(lParam) - .lx
		dy = HIWORD(lParam) - .ly
		p = .GetItemPosition(.dragging)
		x = p.x + dx - .x
		y = p.y + dy - .y
		.SetItemPosition(.dragging, x, y)
		.Send('ItemMoved', .dragging, x, y)

		.dragging = false
		return 0
		}
	SetItemPosition(i, x, y)
		{
		SendMessage(.Hwnd, LVM.SETITEMPOSITION, i, MAKELONG(x, y))
		}
	GetItemPosition(i)
		{
		SendMessagePoint(.Hwnd, LVM.GETITEMPOSITION, i, p = Object())
		return p
		}
	GetFirstVisibleItem()
		{
		return .SendMessage(LVM.GETTOPINDEX)
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
	Seek(field, prefix)
		{
		if false is i = .model.Seek(field, prefix)
			return false
		.EnsureVisible(0)
		.EnsureVisible(Min(.nrows, i + .GetCountPerPage()) - 1)
		return i
		}
	GetCountPerPage()
		{
		return .SendMessage(LVM.GETCOUNTPERPAGE)
		}

	SetMenu(menu)
		{
		.menu = menu
		.use_menu = true
		}
	GetMenu()
		{
		return .use_menu ? .menu : false
		}
	rx: false
	RBUTTONDOWN(lParam)
		{
		.rx = LOWORD(lParam)
		.ry = HIWORD(lParam)
		return "callsuper"
		}
	NM_RCLICK()
		{
		if .use_menu and .rx isnt false
			{
			ClientToScreen(.Hwnd, p = Object(x: .rx, y: .ry))
			ContextMenu(.menu).ShowCall(this, p.x, p.y)
			}
		.rx = false
		return 0
		}
	Default(@args)
		{
		event = args[0]
		if event.Prefix?('On_Context_')
			{
			item = event.AfterFirst('On_Context_')
			event = 'On_' $ item
			if .Method?(event)
				this[event]()
			else
				.Send('On_' $ item) // specific
			.Send('On_Context', args.item) // generic
			}
		else
			super.Default(@args)
		}
	On_Large_Icons()
		{
		.SetView(LVS.ICON)
		}
	On_Small_Icons()
		{
		.SetView(LVS.SMALLICON)
		}
	On_List()
		{
		.SetView(LVS.LIST)
		}
	On_Details()
		{
		.SetView(LVS.REPORT)
		}
	On_Arrange()
		{
		SendMessage(.Hwnd, LVM.ARRANGE, LVA.DEFAULT, 0)
		}
	On_Line_Up()
		{
		SendMessage(.Hwnd, LVM.ARRANGE, LVA.SNAPTOGRID, 0)
		}
	On_Copy()
		{
		.with_value_under_mouse()
			{|record/*unused*/, value|
			if not String?(value)
				value = Display(value)
			ClipboardWriteString(value)
			}
		}
	On_Zoom()
		{
		.with_value_under_mouse()
			{|record/*unused*/, value|
			if not String?(value)
				value = Display(value)
			ZoomControl(.Window.Hwnd, value, true)
			}
		}
	On_Inspect()
		{
		.On_Inspect_Value()
		}
	On_Inspect_Value()
		{
		.with_value_under_mouse()
			{|record/*unused*/, value|
			Inspect.Window(value, 'Inspect Value')
			}
		}
	On_Inspect_Record()
		{
		.with_value_under_mouse()
			{|record, value/*unused*/|
			Inspect.Window(record, 'Inspect Record')
			}
		}
	On_Reset_Columns()
		{
		.ResetColumnWidths()
		}
	with_value_under_mouse(block)
		{
		if .rx is false // e.g. On_Copy from CTRL+C
			return
		ht = .SubItemHitTest(.rx, .ry)
		if ht.iItem is -1 or ht.iSubItem is -1
			{
			Beep()
			return
			}
		if Record() is data = .GetRecord(ht.iItem)
			return
		value = data[.cols[ht.iSubItem]]
		try
			block(data, value)
		catch
			Beep()
		}

	// used for saving and loading columns
	HeaderChanged?()
		{
		return true
		}
	GetStringWidth(s)
		{
		return SendMessageTextIn(.Hwnd, LVM.GETSTRINGWIDTH, 0, s)
		}
	GetItemRect(item)
		{
		SendMessageRect(.Hwnd, LVM.GETITEMRECT, item, r = Object())
		return r
		}
	GetSubItemRect(item, subitem)
		{
		return 0 is SendMessageRect(.Hwnd, LVM.GETSUBITEMRECT, item,
			r = Object(top: subitem, left: LVIR.BOUNDS))
			? false : r
		}

	GetNext(i, flags = false)
		{
		if (flags is false)
			flags = LVNI.ALL
		return SendMessage(.Hwnd, LVM.GETNEXTITEM, i, flags)
		}
	GetSelected()
		{
		result = SendMessage(.Hwnd, LVM.GETNEXTITEM, -1, LVNI.SELECTED)
		return (result > -1) ? result : false
		}
	GetSelectedItems()
		{
		selected = Object()
		start = -1
		do
			{
			result = SendMessage(.Hwnd, LVM.GETNEXTITEM, start, LVNI.SELECTED)
			if result > start
				selected.Add(result)
			start = result
			} while (result > -1)
		return selected
		}
	GetPreviousItem()
		{
		return .oldindex
		}
	ITEMCHANGED(lParam)
		{
		lv = NMLISTVIEW(lParam)
		item = lv.iItem

		// don't do regular processing if item checkbox is being disabled
		if .checkBoxDisabled(lv)
			return 0

		// check if item checkbox state has changed
		if .checkBoxChecked(lv)
			.Send("ItemChecked", item, true)
		if .checkBoxUnchecked(lv)
			.Send("ItemChecked", item, false)

		// get the item and subitems if item is losing focus
		if ((LVIS.FOCUSED & lv.uOldState) isnt 0)
			{
			.oldindex = item
			.olditem = .Getrow(item)
			}
		// get the item and subitems if focused
		if ((LVIS.FOCUSED & lv.uNewState) isnt 0)
			{
			.newitem = .Getrow(item)
			.Send("SelectChanged", .olditem, .newitem)
			}
		return 0
		}
	itemStateBitShift: 12
	checkBoxDisabled(lv)
		{
		return (lv.uOldState >> .itemStateBitShift isnt LVCS.DISABLE and
			lv.uNewState >> .itemStateBitShift is LVCS.DISABLE)
		}
	checkBoxChecked(lv)
		{
		return (lv.uOldState >> .itemStateBitShift is LVCS.UNCHECKED and
			lv.uNewState >> .itemStateBitShift is LVCS.CHECKED)
		}
	checkBoxUnchecked(lv)
		{
		return (lv.uOldState >> .itemStateBitShift is LVCS.CHECKED and
			lv.uNewState >> .itemStateBitShift is LVCS.UNCHECKED)
		}

	NM_CLICK()
		{
		.Send('ItemClicked', .GetSelected())
		return 0
		}

	GETDISPINFO(lParam)
		{
		if (.model is false)
			return 0
		dispinfo = NMLVDISPINFO(lParam)
		if ((LVIF.TEXT & dispinfo.item.mask) is LVIF.TEXT)
			{
			if not .cols.Member?(dispinfo.item.iSubItem)
				return 0
			try
				s = .model.Getitem(dispinfo.item.iItem, .cols[dispinfo.item.iSubItem])
			catch (err /*unused*/,
				'cannot use a completed Transaction|cannot use ended transaction')
				{
				.Send("On_Cancel")
				return 0
				}
			s = .display(s, dispinfo.item.cchTextMax - 1) $ '\x00'
			CopyMemory(dispinfo.item.pszText, s, s.Size())

			rows = .model.Getnumrows()
			if (rows isnt .nrows)
				.SetListSize(.nrows = rows)
			}
		return 0
		}
	display(value, limit = 500)
		{
		value = .display? or not String?(value)
			? Display(value) : value
		ellipsis = ' ...'
		if value.Size() > limit
			value = value[.. limit - ellipsis.Size()] $ ellipsis
		return value
		}
	UpdateCachedRow(i, rec)
		{
		r =  .model.Getrecord(i)
		for (member in rec.Members())
			r[member] = rec[member]
		InvalidateRect(.Hwnd, 0, true)
		}

	EnableMenu(enabled = true)
		{
		.use_menu = enabled
		}
	GetLastRow()
		{
		return (.nrows - 1)
		}
	GetRowCount()
		{
		return .nrows
		}
	ToObject()
		{
		listob = Object()
		for (i = 0; i < .nrows; ++i)
			listob[i] = .Getrow(i)
		return listob
		}
	NM_DBLCLK()
		{
		if (false is .Send('ListDoubleClick', .GetSelected()))
			return 0
		if (.doubleclick_edit is true and ((sel = .GetSelected()) isnt false))
			SendMessage(.Hwnd, LVM.EDITLABEL, sel, 0)
		.Send('NM_DBLCLK')
		return 0
		}
	VSCROLL(wParam, lParam)
		{
		if (.realtime_scrolling)
			return 'callsuper'
		switch (LOWORD(wParam))
			{
		case SB.THUMBTRACK :
			return 0
		case SB.THUMBPOSITION :
			return .Callsuper(.Hwnd, WM.VSCROLL,
				MAKELONG(SB.THUMBTRACK, HIWORD(wParam)), lParam)
		default :
			return 'callsuper'
			}
		}
	HSCROLL(wParam, lParam)
		{
		if (.realtime_scrolling)
			return 'callsuper'
		switch (LOWORD(wParam))
			{
		case SB.THUMBTRACK :
			return 0
		case SB.THUMBPOSITION :
			return .Callsuper(.Hwnd, WM.HSCROLL,
				MAKELONG(SB.THUMBTRACK, HIWORD(wParam)), lParam)
		default :
			return 'callsuper'
			}
		}
	LVN_COLUMNCLICK(lParam)
		{
		nm = NMLISTVIEW(lParam)
		.Send('ListView_ColumnClick', nm.iSubItem)
		return 0
		}

	GetHeader()
		{
		return .SendMessage(LVM.GETHEADER)
		}

	HDN_ITEMCLICK(lParam)
		{
		nmh = NMHEADER(lParam)
		.Send("ListView_HeaderItemClick", nmh.iItem)
		return 0
		}
	HDN_BEGINTRACK()
		{
		return .Send('ListView_HeaderBeginTrack')
		}
	HDN_ENDTRACK()
		{
		.Send('ListView_HeaderEndTrack')
		return 0
		}
	HDN_ITEMCHANGED()
		{
		.adjust()
		.Send('ListView_HeaderItemChanged')
		return 'callsuper'
		}

	Adjust()
		{
		.adjust()
		}

	adjust()
		{
		if .stretch_to_fill is false
			return
		else if .stretch_to_fill is true
			col = .ncols - 1
		else if String?(.stretch_to_fill)
			col = .columns.Find(.stretch_to_fill)
		else
			col = .stretch_to_fill
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

	Scroll(dx, dy)
		{
		.SendMessage(LVM.SCROLL, dx, dy)
		}

	x: 0, y: 0, w: 0, h: 0
	Resize(.x, .y, .w, .h)
		{
		super.Resize(x, y, w, h)
		.adjust()
		}

	model: false
	close_model()
		{
		if .model is false
			return
		.model.Close()
		.model = false
		}
	Destroy()
		{
		if .usercolumns isnt false
			UserColumns.Save(.usercolumns, this)
		.close_model()
		if .tip isnt false
			{
			.tip.Destroy()
			.tip = false
			}
		.remove_header()
		super.Destroy()
		}
	}
