// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
WndProc
	{
	Name:			"List"
	Ymin:			0
	Xmin:			0
	Ystretch:		1
	Xstretch:		1
	lcf:			#(SELECTED: 0x0100, HIGHLIGHTED:0x0200, FOCUSED: 0x0400,
						HOVERED: 0x0800)
	fullEditMode:	false
	editor:			false
	readOnly:		false
	allowTab:		true
	contextOnly:	false
	header:			false
	columns:		false
	widths:			false
	data:			false
	multiSelect:	false
	focused:		false
	shiftAnchor:	false
	sortCol:		false
	forceDelete:	false
	deleteImage:	false
	alwaysHighlightSelected: false
	horzMargin:		5
	rowOffset:		0
	horzOffset:		0
	markColWidth:	0

	New(columns = false, data = false, .defWidth = 100, .noShading = false,
		noDragDrop = false, highlightColor = 0x00FF9957, noHeaderButtons = false,
		.headerSelectPrompt = false, booleansAsBox = false, fontSize = "",
		.resetColumns = false, .customizeColumns = false,
		.alwaysHighlightSelected = false, .indicateHovered = false,
		.columnsSaveName = false, .checkBoxColumn = false, .sortSaveHandler = false,
		.trackValid = false)
		{
		Assert((Integer?(defWidth) and (defWidth > -1)) or (defWidth is false))
		.origColumns = columns isnt false ? columns : Object()
		.CreateWindow("SuWhiteArrow", "", .WindowStyle(), WS_EX.CLIENTEDGE)
		.SubClass()
		.SetFont(size: fontSize, text: "W...")
		.colOverlap = .Xmin + .horzMargin
		.rowHeight = .Ymin += 4 /*= row margin*/
		.header = .Construct(Object('Header'
			style: .hdrStyle(noDragDrop, noHeaderButtons), :headerSelectPrompt))
		.Ymin = 3 * .rowHeight + .header.Ymin		// min 3 rows + headerheight
		.formatting = new ListFormatting(fontSize isnt "", booleansAsBox)
		.SetColumns(.origColumns)
		.Set(data isnt false ? data : Object())
		.brushes = Object()
		.highlight_colors = Object()
		.AddHighlight(false, highlightColor)	// set default at index 0
		.brushes.background = GetSysColorBrush(COLOR.WINDOW)
		.brushes.shaded = CreateSolidBrush(CLR.azure)
		.brushes.focused = GetSysColorBrush(COLOR.HIGHLIGHT)
		.brushes.invalid = CreateSolidBrush(CLR.LIGHTRED)
		.brushes.delete = CreateSolidBrush(CLR.RED)
		.deleteImage = ImageResource('delete.emf')
		.SetVisible(true)
		if ReadOnlyAccess(this)
			.SetReadOnly(true)
		.lastEdit = Date().Minus(seconds: 5)
		if .columnsSaveName isnt false
			UserColumns.Load(.GetColumns(), .columnsSaveName, this, deletecol: false,
				initialized?: true, load_visible?:)
		if .trackValid
			.Window.AddValidationItem(this)
		}
	Startup()
		{
		.header.Startup()
		}
	WindowStyle()
		{
		return WS.TABSTOP | WS.HSCROLL | WS.VSCROLL
		}
	hdrStyle(noDragDrop, noHeaderButtons)
		{
		return (noDragDrop ? 0 : HDS.DRAGDROP) | (noHeaderButtons ? 0 : HDS.BUTTONS)
		}
	Resize(x, y, w, h)		// override to size and move header, update scrolling
		{
		super.Resize(x, y, w, h)
		.rowOffset = Max(0, Min(.rowOffset, .GetNumRows() - .GetNumVisibleRows() + 2))
		.horzOffset = Max(0,
			Min(.horzOffset, .GetTotalColWidths() - .GetClientRect().GetWidth()))
		.UpdateScrollbars()
		.synchHeader()
		}
	Repaint()				// override to avoid repainting header
		{
		rc = .GetClientRect().ToWindowsRect()
		rc.top += .header.Ymin
		InvalidateRect(.Hwnd, rc, false)
		}
	Clear()
		{
		.Set(Object())
		}
	Get()
		{
		return .data
		}
	Set(data)
		{
		if .data is data
			return
		.data = Object()
		.rowOffset = 0
		.focused = false
		.AddRows(data)
		.select(0, true, false, false)
		}
	AddRows(data)
		{
		Assert(Object?(data) and not data.HasNamed?(), "data must be a list")
		for (dataRow in data)
			{
			.checkRow(dataRow)
			.data.Add(dataRow.Set_default(""))
			}
		.sortCol = false
		.UpdateScrollbars()
		.setSortIndicator(false)
		.Repaint()
		}
	FindRowIdx(field, value)
		{
		return .data.FindIf({ it[field] is value })
		}
	GetRow(row)
		{
		return .data[row]
		}
	SetRow(row, value)
		{
		Assert(.data.Member?(row))
		.checkRow(value)
		.data[row] = value.Set_default("")
		.sortCol = false
		.RepaintRow(row)
		}
	AddRow(value)
		{
		.InsertRow(.GetNumRows(), value)
		}
	InsertRow(row, value)
		{
		Assert(.data.Member?(row) or row is .GetNumRows())
		.checkRow(value)
		.data.Add(value.Set_default(""), at: row)
		.sortCol = false
		.UpdateScrollbars()
		if row is .GetNumRows() - 1
			.RepaintRow(row)
		else if row < .rowOffset + .GetNumVisibleRows()
			.Repaint()
		}
	CheckAndInsertRow(row, newRecord, useDefaultsIfEmpty? = false)
		{
		if false is .Send("List_WantNewRow", prevRow: .focused,
			record: newRecord, :useDefaultsIfEmpty?)
			return false;			// not allowed by parent
		.InsertRow(row, newRecord)
		.Send("List_NewRowAdded", row, record: newRecord)
		return .data[row]
		}
	GetNumRows()
		{
		return .data.Size()
		}
	getRowRect(row)
		{
		return .GetClientRect().Set(
			y: .header.Ymin + (row - .rowOffset) * .rowHeight,
			height: .rowHeight)
		}
	GetDataRowFromY(y)
		{
		return .getRowFromY(y) + .rowOffset
		}
	getRowFromY(y)
		{
		// May return a value >= .GetNumRows
		return Max(0, ((y - .header.Ymin) / .rowHeight).Int())
		}

	getNumVisibleRowsWithFraction()
		{
		return ((.GetClientRect().GetHeight() - .header.Ymin) / .rowHeight)
		}

	GetNumVisibleRows()
		{
		// returns number of rows which are at least partially
		// visible in the current client area
		return .getNumVisibleRowsWithFraction().Ceiling()
		}
	GetNumFullyVisibleRows()
		{
		return .getNumVisibleRowsWithFraction().Floor()
		}
	DeleteSelection()
		{
		// selected rows are marked for deletion / deleted
		// returns actual number of rows deleted
		return .DeleteRows(@ .getSelection())
		}
	DeleteAll() // WARNING: bad on large lists
		{
		return .DeleteRows(@ .data.Members())
		}
	SetForceDelete()
		{
		// next DeleteRows() DELETES rows
		// wo confimation or notification of controller !
		.forceDelete = true
		}
	DeleteRows(@args)
		{
		// args is a list of row indices to delete
		// rows are marked for deletion / deleted
		// returns actual number of rows deleted
		rowsToDelete = .getDeletedRow(args)
		if not rowsToDelete.Empty?()
			{
			dataRowsDeleted = Object()
			for (row in rowsToDelete)
				{
				dataRowsDeleted.Add(.data[row])
				.data.Delete(row)
				if row < .rowOffset
					.rowOffset--
				if .focused isnt false and row < .focused
					.focused--
				}
			.checkAndSetFocusedRow()
			.rowOffset = Max(0, Min(.rowOffset, .GetNumRows() - .GetNumVisibleRows() + 2))
			.UpdateScrollbars()
			.Repaint()
			// notify controller of deletions
			if .forceDelete is false
				.Send("List_Deletions", deletions: dataRowsDeleted.Reverse!())
			}
		.forceDelete = false
		return rowsToDelete.Size()
		}
	getDeletedRow(args)
		{
		rowsToDelete = Object()					// build valid list of rows to delete
		for (row in args.Sort!().Reverse!())			// traverse in decending order
			if row >= 0 and row < .GetNumRows()	// filter valid rownumbers
				if not rowsToDelete.Has?(row)		// ignore duplicates
					if true is .forceDelete or
						false isnt .Send("List_DeleteRecord", .data[row])
						rowsToDelete.Add(row)
		return rowsToDelete
		}
	checkAndSetFocusedRow()
		{
		if .focused isnt false
			if .GetNumRows() is 0
				.focused = false
			else if .focused >= .GetNumRows()
				.focused = .GetNumRows() - 1
		}
	ScrollRowToView(row)
		{
		Assert(.data.Member?(row))
		offset = Max(
			Min(0, row - .rowOffset),
			row - .rowOffset - .GetNumVisibleRows() + 2)
		if offset isnt 0
			.scroll(0, offset)
		else
			.Update()
		}
	ScrollToBottom()
		{
		rows = .GetNumRows()
		if rows > 1
			.ScrollRowToView(rows - 1)
		}
	DoWithCurrentVScrollPos(block)
		{
		savedVScrollPos = .rowOffset
		block()
		.scroll(0, savedVScrollPos - .rowOffset)
		}
	RepaintRow(row)
		{
		Assert(.data.Member?(row))
		clientRect = .GetClientRect()
		if 0 > height = clientRect.GetHeight() - .header.Ymin
			return
		clientRect.Set(y: .header.Ymin, :height)
		if clientRect.Overlaps?(rect = .getRowRect(row))
			InvalidateRect(.Hwnd, rect.ToWindowsRect(), false)
		}
	GetRowHeight()
		{
		return .rowHeight / GetDpiFactor()
		}
	SetRowHeight(newHeight)
		{
		.rowHeight = newHeight
		.Repaint()
		}
	SetColWidth(col, width, fromHeader = false)
		{
		Assert(.columns.Member?(col))
		if width is false
			width = .header.GetDefaultColumnWidth(.columns[col])
		if ((col is 0 and .markColWidth isnt 0) or 0 is movePix = width - .widths[col])
			return
		oldRect = .getColRect(col)
		.widths[col] = width
		newRect = .getColRect(col)
		if movePix < 0
			oldRect = newRect
		updateRect = .GetClientRect()
		updateRect.Set(
				x:		Max(.markColWidth, oldRect.GetX() + oldRect.GetWidth()),
				y:		.header.Ymin
				width:	updateRect.GetWidth() - oldRect.GetX() - oldRect.GetWidth()
				)
		.setNewRect(newRect, fromHeader, col, width, oldRect)
		.synchHeader()
		.Update()		// force update of clientarea
		.drawFocus()	// hide focus rectangle
		ScrollWindowEx(.Hwnd, movePix, 0,
			updateRect.ToWindowsRect(), updateRect.ToWindowsRect(),
			0, 0, SW.INVALIDATE)
		.drawFocus()	// show focus rectangle
		InvalidateRect(.Hwnd, newRect.ToWindowsRect(), false)
		.scroll(0, 0)
		.Update()		// force update of clientarea
		}
	setNewRect(newRect, fromHeader, col, width, oldRect)
		{
		newRect.Set(y: .header.Ymin)
		if not fromHeader
			.header.SetItemWidth(.markColWidth is 0 ? col : col - 1, width)
		else if .formatting.AssumeLeftJust?(.columns[col])
			{			// shortcut if only right portion has to be redrawn
			w = Min(oldRect.GetWidth(), .colOverlap)
			newRect.Set(x: oldRect.GetX() + oldRect.GetWidth() - w, width: w)
			}
		}
	GetColWidth(col)
		{
		return .widths[col]
		}
	SetCol(col /*unused*/, text /*unused*/)
		{
		throw "not implemented"
		}
	GetCol(col)
		{
		// returns the name of column at index col
		return .columns.GetDefault(col, false)
		}
	SetColumns(columns, reset = false)
		{
		Assert(Object?(columns) and not columns.HasNamed?())
		if .columns is columns and not reset
			return
		.columns = columns.Copy()
		.widths = Object()
		.markColWidth = 0
		.horzOffset = 0
		.header.Clear()
		.formatting.SetFormats(.columns) // create formats for displaying columns

		colno = 0
		for (col in .columns.Members())
			if col is 0 and columns[0] is "listrow_deleted"
				{
				.markColWidth = GetSystemMetrics(SM.CXHSCROLL)
				.widths.Add(.markColWidth)
				}
			else
				.addHeader(columns[col], colno++)
		if reset
			.header_changed? = true
		.updateForColumnsChange()
		}
	AppendColumns(newColumns)
		{
		.formatting.SetFormats(.columns.Copy().Append(newColumns))
		offset = .columns.Has?("listrow_deleted") ? -1 : 0
		for column in newColumns
			{
			.addHeader(column, .columns.Size() + offset)
			.columns.Add(column)
			}
		.updateForColumnsChange()
		}
	addHeader(column, colno)
		{
		.header.AddItem(column, .defWidth)
		.header.SetItemFormat(colno, .formatting.GetHeaderAlign(column))
		.widths.Add(.header.GetItemWidth(colno))
		}
	updateForColumnsChange()
		{
		.UpdateScrollbars()
		.Repaint()
		.synchHeader()
		}
	GetColumns()
		{
		return .columns.Copy()
		}
	GetVisibleColumns()
		{
		cols = Object()
		for (i = 0; i < .columns.Size(); i++)
			if .GetColWidth(i) > 0
				cols.Add(.columns[i])
		return cols
		}
	GetNumCols()
		{
		return .columns.Size()
		}
	getColRect(col)
		{
		for (x = 0, c = col; c > 0; )
			x += .widths[--c]
		return Rect(x - .horzOffset, .header.Ymin,
			.widths[col], .GetClientRect().GetHeight() - .header.Ymin)
		}
	GetTotalColWidths()
		{
		return .widths.Sum()
		}
	AllowTab(allow)
		{
		Assert(Boolean?(allow))
		.allowTab = allow
		}
	AllowContextOnly(allow)
		{
		Assert(Boolean?(allow))
		.contextOnly = allow
		}
	ScrollColToView(col)
		{
		Assert(.columns.Member?(col))
		colRect = .getColRect(col)
		offset = Max(0, colRect.GetX() + colRect.GetWidth() - .GetClientRect().GetWidth())
		offset = Min(offset, colRect.GetX() - .markColWidth)
		if offset isnt 0
			.scroll(offset, 0)
		}
	getNextCell(col, row, amt)
		{
		if row is 0 and col is 0 and amt < 0 //no cells left to get
			return Object(col: 0, row: 0)

		numCols = .GetNumCols()
		numRows = .GetNumRows()
		forever
			{
			col += amt
			row = Max(0, row + (col / numCols).Floor())
			col %= numCols
			if col < 0
				col += numCols
			if row >= numRows		// will give non-existent row at end for edit
				return Object(col: 0, row: numRows)

			if .skip?(row, col)
				continue

			return Object(:col, :row)
			}
		}
	skip?(row, col)
		{
		// skip deleted rows
		if .data[row].GetDefault("listrow_deleted", false) is true
			return true
		if .hiddenInvalidCol?(col, row)
			return true
		if true is .Send('List_Tabover?', .GetCol(col))
			return true
		return false
		}
	hiddenInvalidCol?(col, row)
		{
		minWidth = 5
		return (.GetColWidth(col) < minWidth and
			not (.data[row].Member?("list_invalid_cells") and
				.data[row].list_invalid_cells.Member?(.columns[col])) and
			.widths[1..].Max() >= minWidth)
		}
	SetBrush(brush, color)
		{
		Assert(String?(brush) and Number?(color))
		DeleteObject(.brushes[brush])
		.brushes[brush] = CreateSolidBrush(color)
		}
	SetReadOnly(readOnly, grayOut = true)
		{
		Assert(Boolean?(readOnly) and Boolean?(grayOut))
		if readOnly
			.FinishEdit()
		.readOnly = readOnly
		DeleteObject(.brushes.background)
		.brushes.background = GetSysColorBrush(
			.readOnly and grayOut ? COLOR.BTNFACE : COLOR.WINDOW)
		.Repaint()
		}
	GetReadOnly()
		{
		return .readOnly
		}
	SetMultiSelect(multi)
		{
		Assert(Boolean?(multi))
		if .multiSelect isnt multi and false is (.multiSelect = multi)
			.select(0, false, false, true)
		}
	GetMultiSelect()
		{
		return .multiSelect
		}
	SetFullEditMode(mode)
		{
		// false restricts continuous editing (clicking) to the same row
		Assert(Boolean?(mode))
		.fullEditMode = mode
		}

	GetCheckBoxField()
		{
		return .checkBoxColumn
		}

	GetField(field)
		{
		if .getSelection().Size() isnt 1
			throw "List.GetField requires a single row selection"
		return .GetCurrentRecord()[field]
		}

	SetField(field, value, idx = false)
		{
		if idx is false
			{
			selected = .getSelection()
			if selected.Size() isnt 1
				throw "List.SetField requires a single row selection"
			idx = selected[0]
			}
		record = .GetRow(idx)
		record[field] = value
		}

	GetCurrentRecord()
		{
		selected = .GetSelection()
		return selected.Size() is 1 ? .GetRow(selected[0]) : false
		}

	GetSelection()
		{
		return .getFocusSelection()
		}

	getFocusSelection()
		{
		selected = .getSelection()
		if selected.Size() is 0 and .focused isnt false
			{
			.SetSelection(.focused)
			selected = Object(.focused)
			}
		return selected
		}

	getSelection()
		{
		selection = Object()
		for (row in .data.Members())
			if .RowSelected?(row)
				selection.Add(row)
		return selection
		}
	SetSelection(row)
		{
		Assert(.data.Member?(row))
		.selectFocus(row)
		}
	setRowFlags(rowOb, flags)
		{
		if Record?(rowOb)
			rowOb.PreSet("listrow_flags", flags)
		else
			rowOb.listrow_flags = flags
		}
	SelectRow(row)
		{
		if not .data.Member?(row) or .RowSelected?(row)
			return

		.setRowFlags(.data[row], .data[row].listrow_flags | .lcf.SELECTED)
		.RepaintRow(row)
		}
	DeSelectRow(row)
		{
		if not .data.Member?(row) or not .RowSelected?(row)
			return

		.FinishEdit()
		.setRowFlags(.data[row], .data[row].listrow_flags ^ .lcf.SELECTED)
		.RepaintRow(row)
		}
	HoverRow(row)
		{
		if .indicateHovered is false
			return
		if .RowHovered?(row) or not .data.Member?(row)
			return

		.hoverRow = row
		.setRowFlags(.data[row], .data[row].listrow_flags | .lcf.HOVERED)
		.RepaintRow(row)
		}
	DeHoverRow()
		{
		if .hoverRow is false or not .data.Member?(.hoverRow)
			return
		.setRowFlags(.data[.hoverRow], .data[.hoverRow].listrow_flags ^ .lcf.HOVERED)
		.RepaintRow(.hoverRow)
		.hoverRow = false
		}
	ClearSelectFocus()
		{
		.ClearSelect()
		.focused = false
		}
	ClearSelect()
		{
		.Send("List_BeforeClearSelect")
		for row in .data.Members().Filter(.RowSelected?).Copy()
			.DeSelectRow(row)
		.Send("List_AfterClearSelect")
		}
	AddInvalidCell(col, row)
		{
		rec = .data[row]
		Assert(.columns.Member?(col))
		field = .columns[col]
		if not rec.Member?("list_invalid_cells")
			rec.list_invalid_cells = Object()
		rec.list_invalid_cells[field] = rec[field]
		.RepaintRow(row)
		}
	HasInvalidCell?(record, member)
		{
		return record.Member?("list_invalid_cells") and
			record.list_invalid_cells.Member?(member)
		}
	RowHasInvalidCell?(record)
		{
		return record.Member?("list_invalid_cells") and
			not record.list_invalid_cells.Empty?()
		}
	InvalidCellValue(record, member)
		{
		return record.list_invalid_cells[member]
		}
	RemoveInvalidCell(col, row)
		{
		rec = .data[row]
		Assert(.columns.Member?(col))
		if not rec.Member?("list_invalid_cells")
			return
		rec.list_invalid_cells.Delete(.columns[col])
		.RepaintRow(row)
		}
	AddHighlight(row, color = false, rowData = false)
		{
		// color is added to the colors/brushes if not already present
		// the order in which colors are used is preserved for sorting/grouping
		// the default highlightcolor is at index 0 !
		// returns the index number of the highlightcolor/brush
		if color is false
			cidx = 0
		else if false is cidx = .highlight_colors.Find(color)
			{
			cidx = .highlight_colors.Size()
			.highlight_colors.Add(color)
			.brushes.Add(CreateSolidBrush(color))
			}
		// if row is false, just add color to force order for sorting/grouping
		if row isnt false
			{
			.addHeightLightFlag(.data[row], cidx)
			.RepaintRow(row)
			}
		// rowData is used to set a highlight prior to row construction
		else if rowData isnt false
			.addHeightLightFlag(rowData, cidx)
		return cidx
		}
	addHeightLightFlag(rec, cidx)
		{
		.setRowFlags(rec,
			LOWORD(rec.listrow_flags) | .lcf.HIGHLIGHTED + (cidx << 16)) /*= add first hex
				as color index*/
		}
	GetHighlighted()
		{
		highlight = Object()
		for (row in .data.Members())
			if .RowHighlighted?(row)
				highlight.Add(row)
		return highlight
		}
	HighlightValues(member, values, color = false, sortHighlight = false, group = false)
		{	// allways add new colors, so these are eventually set for sorting/grouping
		cidx = color is false ? 0 : .AddHighlight(false, color)
		for (rec in .data)
			if values.Has?(rec[member])
				.addHeightLightFlag(rec, cidx)
		if sortHighlight
			.SortHighlight(group)
		.Repaint()
		}
	ClearHighlight(row = false)
		{
		if row isnt false
			{
			Assert(.data.Member?(row))
			if .RowHighlighted?(row)
				{
				.setRowFlags(.data[row],
					.data[row].listrow_flags & (0xffff - .lcf.HIGHLIGHTED))
				.RepaintRow(row)
				}
			}
		else
			for (row in .data.Members())
				.ClearHighlight(row)
		}
	RowHighlighted?(row)
		{
		return (.data[row].listrow_flags & .lcf.HIGHLIGHTED) isnt 0
		}
	RowSelected?(row)
		{
		return (.data[row].listrow_flags & .lcf.SELECTED) isnt 0
		}
	RowHovered?(row)
		{
		return row isnt false and .hoverRow is row
		}
	drawFocus()
		{
		// shows or hides the focus rectangle so that horizontal
		// window scrolling doesn't create garbage
		if .focused is false or
			.focused < .rowOffset or
			.focused >= .rowOffset + .GetNumVisibleRows()
			return
		.drawLineFocus(hdc = GetDC(.Hwnd), .getRowRect(.focused).ToWindowsRect())
		ReleaseDC(.Hwnd, hdc)
		}
	selectFocus(row, ctrl = false, shift = false, select = true)
		{
		row = Min(row, .GetNumRows() - 1)
		oldFocus = .focused
		.focused = row
		if oldFocus is false
			.select(0, ctrl, shift, select)
		else
			{
			if oldFocus isnt .focused
				.RepaintRow(oldFocus)
			.select(shift ? .focused - oldFocus : 0, ctrl, shift, select)
			}
		}
	select(amt, ctrl = false, shift = false, select = false)
		{
		if not .selectable?()
			return
		if not shift
			.shiftAnchor = false
		newFocus = Max(0, Min(.focused + amt, .GetNumRows() - 1))
		oldSelection = .getSelection()
		newSelection = oldSelection.Copy()
		if ctrl
			newSelection = .ctrlSelect(select, newFocus, newSelection)
		else if shift
			newSelection = .shiftSelect(newFocus)
		else if select
			newSelection = Object(newFocus)

		.RepaintRow(newFocus)
		.RepaintRow(.focused)
		.focused = newFocus
		.handleSelectionChanged(newSelection, oldSelection)
		.ScrollRowToView(.focused)		// ensure focus is in view
		}
	selectable?()
		{
		if .GetNumRows() < 1
			{
			.focused = false
			.Send("List_Selection", selection: false)
			return false
			}
		if .focused is false
			.focused = 0
		if .GetNumVisibleRows() < 1
			return false
		return true
		}
	handleSelectionChanged(newSelection, oldSelection)
		{
		if newSelection.Sort!() isnt oldSelection	// if selection changed
			{
			for (row in oldSelection.Difference(newSelection))
				.DeSelectRow(row)
			for (row in newSelection.Difference(oldSelection))
				.SelectRow(row)
			.Send("List_Selection",
				selection: newSelection.Empty?() ? false : newSelection)
			}
		}
	ctrlSelect(select, newFocus, newSelection)
		{
		if select
			{
			if .RowSelected?(newFocus)
				newSelection.Remove(newFocus)
			else
				newSelection = .multiSelect
					? newSelection.Add(newFocus) : Object(newFocus)
			}
		return newSelection
		}
	shiftSelect(newFocus)
		{
		if .shiftAnchor is false
			.shiftAnchor = .focused
		inc = .shiftAnchor < newFocus ? 1: -1
		newSelection = Object(newFocus)
		for (row = .shiftAnchor; row isnt newFocus; row += inc)
			newSelection.Add(row)
		return newSelection
		}
	repaintSelection(selection = false)
		{
		for (row in (selection is false) ? .getSelection() : selection)
			.RepaintRow(row)
		}
	checkRow(dataRow)
		{
		Assert(Object?(dataRow), "row must be an object")
		Assert(not dataRow.Readonly?(), "row cannot be read-only object")
		}
	customFields: false
	SetCustomFields(customFields)
		{
		.customFields = customFields
		}
	ContextNew()
		{
		// simulate insert key
		.KEYDOWN(VK.INSERT, 0)
		}
	edit(col, row)
		// attempts to edit a cell
		// if the row specified by row does not exist, a new row is added
		// returns true if
		// 		editing begun successfully OR readOnly OR
		// 		no cols OR callback disallowed row adding
		// false otherwise
		{
		if row >= .GetNumRows()	// add a row?
			{
			row = .GetNumRows()
			newRow = Record()
			if false is .Send("List_WantNewRow", prevRow: row - 1, record: newRow)
				return true
			.AddRow(newRow)
			.Send("List_NewRowAdded", row, record: newRow)
			}
		if false is .Send("List_AllowCellEdit", col, row) or
			false is control = .Send("List_WantEditField",
				:col, :row, data: .data[row][.columns[col]])
			return false						// no editing allowed by parent

		// The following line MUST ALWAYS be done (ie. when inserting rows)
		.selectFocus(row)
		.ScrollColToView(col)

		cellRect = .getColRect(col)
		rowRect = .getRowRect(row)
		clientToScreenResult = ClientToScreen(.Hwnd, pt = Object(x: 0, y: 0))
		cellRect.Set(
			x:		cellRect.GetX() + pt.x + .horzMargin,
			y:		rowRect.GetY() + pt.y,
			width:	Min(.Window.GetClientRect().GetWidth() - .markColWidth,
				Max(cellRect.GetWidth(), 2 * .horzMargin + 13 /*= minimum width*/)) -
				2 * .horzMargin,
			height:	rowRect.GetHeight())
		if .getSelection().Size() isnt 1 // should be assert but don't want to annoy users
			SuneidoLog('ERROR: ListControl edit selection size isnt 1. It is ' $
				Display(.getSelection().Size()) $
				', ClientToScreenResult: ' $ Display(clientToScreenResult), calls:)
		custom = .customFields isnt false
			? .customFields.GetDefault(.GetCol(col), false)
			: false
		.editor = new ListEditWindow(control,
			.Send('List_EditFieldReadonly', col, row) is true,
			col, row, this, cellRect.ToWindowsRect(), :custom,
			customFields: .customFields)
		return true;							// field editing begun!
		}
	ContextEdit(pt)
		{
		ScreenToClient(.Hwnd, pt)
		row = Min(.GetDataRowFromY(pt.y), .GetNumRows())
		col = Min(.GetColFromX(pt.x), .GetNumCols() - 1)
		if .data[row].listrow_deleted isnt true
			.Edit(col, row, 1)
		}
	Edit(col, row, amt, canMoveRows = false)
		{
		// attempts to edit cell at (col, row)
		// if this is not possible, attempts to edit next cell amt cells away
		// if canMoveRows is false,
		// it will not attempt to edit a cell in another row
		while .edit(col, row) is false
			{
			nextCell = .getNextCell(col, row, amt)	// try next cell
			if not .nextCellAvailable?(nextCell, col, row, amt, canMoveRows)
				return false
			col = nextCell.col
			row = nextCell.row
			}
		return true
		}
	nextCellAvailable?(nextCell, col, row, amt, canMoveRows)
		{
		if nextCell.col is col and nextCell.row is row
			return false	// no available next cell
		if nextCell.row isnt row and not canMoveRows
			return false	// can't move (or create) rows
		// only check next cell available when tabbing to next row
		if nextCell.row isnt row and amt is 1 and
			true is .Send('List_NextCellNotAvailable?', row)
			return false
		return true
		}
	FinishEdit()
		{
		if .editor isnt false
			.editor.Return()

		// ensure any drag process being done by the user is ended. This helps prevent
		// situations like AccessControl leaving edit mode while the user is still
		// dragging a record, which can lead to Access not in edit mode when saving
		.EndDrag()
		}
	ListEditWindow_Commit(col, row, dir, data, valid?, unvalidated_val= '')
		{
		.editor = false
		if not .commitable?()
			return

		if unvalidated_val isnt '' and unvalidated_val is
			.GetInvalidFieldData(.data[row], .columns[col])
			return
		_committing = .columns[col]
		.commit(col, row, data, valid?, unvalidated_val)

		if .allowTab and dir isnt 0
			{
			nextCell = .getNextCell(col, row, dir)
			.Edit(nextCell.col, nextCell.row, dir, canMoveRows:)
			}
		}

	commitable?()
		{
		if .data is false or .data.Empty?()
			return false
		if .Destroyed?() or .readOnly is true
			{
			SuneidoLog("ERRATIC: ListEditCommit on destroyed or readonly list",
				params: Record(destroyed: .Destroyed?(), readonly: .readOnly))
			return false
			}
		return true
		}

	GetInvalidFieldData(rec, field)
		{
		if rec.Member?('List_InvalidData') and Object?(rec.List_InvalidData)
			return rec.List_InvalidData.GetDefault(field, '')
		return ''
		}

	commit(col, row, data, valid?, unvalidated_val)
		{
		// from invalid to valid
		corrected? = valid? is true and
			'' isnt .GetInvalidFieldData(.data[row], .columns[col])
		if .SetInvalidFieldData(.data[row], .columns[col], unvalidated_val)
			.Send('List_InvalidDataChanged', .data[row])
		valueChanged? = data isnt .data[row][.columns[col]]
		.updateCellValid(col, row, valid?, unvalidated_val)
		if valueChanged? and false isnt .Send("List_CellEdit", :col, :row, :data, :valid?)
			{
			.data[row][.columns[col]] = data
			.RepaintRow(row)
			.Send("List_CellValueChanged", :col, :row, :data)
			.Send("List_AfterEdit", :col, :row, :data, :valid?)
			}
		.Send('List_AfterEditWindowCommit', :col, :row, :data, :valid?,
			valueChanged?: valueChanged? or unvalidated_val isnt '' or corrected?)
		}

	updateCellValid(col, row, valid?, unvalidated_val)
		{
		if not .trackValid
			return

		if not valid? and unvalidated_val isnt ''
			.AddInvalidCell(col, row)
		else
			.RemoveInvalidCell(col, row)
		}

	SetInvalidFieldData(rec, field, val)
		{
		if not rec.Member?('List_InvalidData')
			rec.List_InvalidData = Object().Set_default('')
		if rec.List_InvalidData[field] isnt val
			{
			rec.List_InvalidData[field] = val
			return true
			}
		return false
		}

	// interface (header)
	header_changed?: false
	HeaderChanged?()
		{
		return .header_changed?
		}
	SetHeaderChanged(status)
		{
		.header_changed? = status
		}
	Header_AllowTrack(col)
		{
		return .Send("List_AllowHeaderResize", .markColWidth is 0 ? col : col + 1)
		}
	HeaderTrack(col, width)
		{
		// sent by header control
		Assert(.columns.Member?(col))
		.SetColWidth(.markColWidth is 0 ? col : col + 1, width, true)
		.header_changed? = true
		}

	GetColNum(col)
		{
		return .markColWidth is 0 ? col : col + 1
		}

	HeaderClick(col, button)
		{
		// sent by header control
		Assert(.columns.Member?(col))
		if button isnt 0 or .GetNumRows() < 1
			return
		scol = .markColWidth isnt 0 ? col + 2 : col + 1
		scol = .sortCol is false or .sortCol.Abs() isnt scol ? scol : -.sortCol
		.SortListData(scol)
		.Send("List_AfterSort")
		}
	SortListData(col = false)
		{
		if col isnt false
			.sortCol = col
		if .sortCol is false
			return
		compareFunc = .Send("List_GetCompareFunc", col: .columns[.sortCol.Abs() - 1])
		if compareFunc is false
			return // parent doesn't allow sorting on this column
		if compareFunc is 0			// use default comparison if no callback
			compareFunc = .defCompareFunc	// or parent has no dedicated function
		if .focused isnt false		// mark focused row
			.setRowFlags(.data[.focused], .data[.focused].listrow_flags | .lcf.FOCUSED)
		cmpfn = .sortCol < 0 ? {|x, y| compareFunc(y, x) } : compareFunc
		.data.Sort!(cmpfn)
		.restoreFocused()
		.setSortIndicator(.sortCol)
		.Repaint()
		}

	showSortIndicator: true
	setSortIndicator(sortCol)
		{
		if not .showSortIndicator
			return

		sortDown? = sortCol < 0
		if sortCol isnt false
			sortCol = sortCol.Abs() - (.markColWidth isnt 0 ? 2 : 1)
		colOffset = .columns.Has?("listrow_deleted") ? 1 : 0
		for (idx = 0; idx < .header.GetItemCount(); idx++)
			{
			fmt = .formatting.GetHeaderAlign(.columns[idx + colOffset])
			if idx is sortCol
				fmt |= sortDown? ? HDF.SORTDOWN : HDF.SORTUP
			.header.SetItemFormat(idx, fmt)
			}
		}

	signedSort(savedSort)
		{
		if savedSort.Has?(',')
			savedSort = savedSort.BeforeFirst(',')
		sortDown? = false
		if savedSort.Has?('reverse')
			{
			savedSort = savedSort.AfterFirst(' ')
			sortDown? = true
			}
		if false is colNum = .columns.Find(savedSort)
			return false
		sortCol = (colNum + 1) * (sortDown? ? -1 : 1)
		.SetSortCol(sortCol)
		}

	SavedSortIndicator(savedSort)
		{
		if savedSort is false
			return
		.signedSort(savedSort)
		}

	GetSortCol()
		{
		return .sortCol
		}
	GetSort(nonMarkExtraCol? = false)
		{
		if not Number?(sortCol = .GetSortCol())
			return ''

		col = .markColWidth isnt 0 or nonMarkExtraCol? ? sortCol.Abs() - 1 : sortCol
		if false is fieldname = .GetCol(col)
			return ''

		return sortCol.Sign() is -1 ? 'reverse ' $ fieldname : fieldname
		}
	SetSortCol(col)
		{
		.sortCol = col
		.SortListData()
		}
	defCompareFunc(x, y)
		{
		return .formatting.CompareRows(.columns[.sortCol.Abs() - 1], x , y)
		}
	Header_AllowDrag(col)
		{
		return .Send("List_AllowHeaderReorder", .markColWidth is 0 ? col : col + 1)
		}
	HeaderReorder(col, newIdx)
		{
		// sent by header control
		Assert(.columns.Member?(col))
		Assert(.columns.Member?(newIdx))
		if col is newIdx
			return
		if .markColWidth isnt 0
			{
			col++
			newIdx++
			}
		oldColumns = .columns.Copy()
		.columns.Delete(col).Add(oldColumns[col], at: newIdx)
		org = .widths[col]
		.widths.Delete(col).Add(org, at: newIdx)
		if .sortCol isnt false		// adjust index sortcolumn
			{
			col = .columns.Find(oldColumns[.sortCol.Abs() - 1]) + 1
			.sortCol = .sortCol > 0 ? col : -col
			}
		.header_changed? = true
		.Repaint()
		}
	synchHeader()
		{
		width = Max(.GetClientRect().GetWidth() - .markColWidth, .header.Xmin)
		offset = .markColWidth - .horzOffset
		.header.Resize(offset, 0, width, .header.Ymin)
		.header.Update()
		if offset > 0
			{
			rect = Object(left:0 , top:0, right: offset, bottom: .header.Ymin)
			FillRect(hdc = GetDC(.Hwnd), rect, GetSysColorBrush(COLOR.BTNFACE))
			ReleaseDC(.Hwnd, hdc)
			}
		}
	restoreFocused()
		{
		// try to restore screenposition as much as possible
		if .focused is false
			return
		for (row in .data.Members()) // eventually changed after sorting
			if 0 isnt (.data[row].listrow_flags & .lcf.FOCUSED)
				{
				.setRowFlags(.data[row], .data[row].listrow_flags ^ .lcf.FOCUSED)
				if .focused >= .rowOffset and
					.focused <= .rowOffset + .GetNumVisibleRows() + 2
					.rowOffset = Max(0,	Min(row - .focused + .rowOffset,
						.GetNumRows() - .GetNumVisibleRows() + 2))
				.focused = row
				.Send("List_SelectedRowPositionChanged", selection: Object(.focused))
				.UpdateScrollbars()
				return
				}
		}
	SortHighlight(group = false)
		{
		if .focused isnt false
			.setRowFlags(.data[.focused], .data[.focused].listrow_flags | .lcf.FOCUSED)
		if group is false	// just sort all highlighted (on top)
			.data.Sort!({|x, y|
				(x.listrow_flags & .lcf.HIGHLIGHTED) >
				(y.listrow_flags & .lcf.HIGHLIGHTED) })
		else					// group colors
			.data.Sort!({|x, y|
				((x.listrow_flags & .lcf.HIGHLIGHTED) is 0 ?
					65536 : HIWORD(x.listrow_flags)) < /*= 0x10000 no hightlights */
				((y.listrow_flags & .lcf.HIGHLIGHTED) is 0 ?
					65536 : HIWORD(y.listrow_flags)) }) /*= 0x10000 no hightlights */
		.sortCol = false
		.restoreFocused()
		.Repaint()
		}
	// interface (scrolling):
	UpdateScrollbars()
		{
		if .data is false
			return
		sif = Object(
			cbSize:	SCROLLINFO.Size(),
			fMask:	SIF.RANGE | SIF.POS | SIF.PAGE,
			nPage:	.GetClientRect().GetWidth(),
			nMin:	0,
			nMax:	.GetTotalColWidths(),
			nPos:	.horzOffset
			)
		.UpdateHorzScroll(sif)
		sif.nPage = .GetNumVisibleRows() - 1
		sif.nMax = .GetNumRows()
		sif.nPos = .rowOffset
		SetScrollInfo(.Hwnd, SB.VERT, sif, true)
		}
	UpdateHorzScroll(sif)
		{
		SetScrollInfo(.Hwnd, SB.HORZ, sif, true)
		}
	scroll(dx, dy)
		{
		// dy
		// plus 2 (1 to handle going from 0 based to 1 based, and 1 for blank row at end)
		newOffset = Max(0, Min(.rowOffset + dy, .GetNumRows() - .GetNumVisibleRows() + 2))
		dy = .rowOffset - newOffset
		.rowOffset = newOffset
		// dx
		newOffset = Max(0,
			Min(.horzOffset + dx, .GetTotalColWidths() - .GetClientRect().GetWidth()))
		dx = .horzOffset - newOffset
		.horzOffset = newOffset
		.UpdateScrollbars()
		// shortcut if Repaint() probably faster (smoother PGUP/PGDN)
		if dy.Abs() > .GetNumVisibleRows() / 2
			.Repaint()
		else if dx isnt 0 or dy isnt 0
			{
			dy *= .rowHeight
			rcScroll = .GetClientRect()
			if dx isnt 0
				{
				if .markColWidth - .horzOffset > 0
					.synchHeader()
				else
					{
					rcScroll.Set(y: 0, height: .header.Ymin)
					ScrollWindowEx(.Hwnd, dx, dy, rcScroll.ToWindowsRect(),
						rcScroll.ToWindowsRect(), NULL, NULL,
							SW.INVALIDATE | SW.SCROLLCHILDREN)
					rcScroll = .GetClientRect()
					}
				rcScroll.Set(x: rcScroll.GetX() + .markColWidth)
				rcScroll.Set(width: rcScroll.GetWidth() - .markColWidth)
				.drawFocus()	// hide focus rect
				}
			rcScroll.Set(y: .header.Ymin)
			ScrollWindowEx(.Hwnd, dx, dy, rcScroll.ToWindowsRect(),
				rcScroll.ToWindowsRect(), NULL, NULL, SW.INVALIDATE)
			if dx isnt 0
				.drawFocus()	// show focus rect again
			}
		}

	// windows messages

	ERASEBKGND()
		{
		return 1 // no erase, the window is completly redrawn by PAINT
		}

	hideContent: false
	HideContent(hide)
		{
		.hideContent = hide
		}

	PAINT(lParam /*unused*/)
		{
		hdc = BeginPaint(.Hwnd, ps = Object())

		// need the BeginPaint and EndPaint even if the content is to be hidden
		if not .hideContent
			.paint(hdc, ps, ps.rcPaint)

		EndPaint(.Hwnd, ps)
		return 0
		}

	paint(hdc, ps, rc)
		{
		topRow = .getRowFromY(ps.rcPaint.top)
		top = topRow * .rowHeight + .header.Ymin
		topRow += .rowOffset
		numRows = .GetNumRows()

		if false is .paintMarkCol(hdc, rc, top, topRow, numRows)
			return

		numCols = .GetNumCols()
		for (col = 0, left = -.horzOffset;
			col < numCols and left + .widths[col] <= rc.left; )
			left += .widths[col++]
		WithHdcSettings(hdc, [.GetFont(), SetBkMode: TRANSPARENT])
			{
			.formatting.SetDC(hdc)
			top = .paintVisibleRows(hdc, rc, numRows, numCols,
				Object(:topRow, :top, :left, :col))
			if top < rc.bottom
				{
				rcSel = Object(left: rc.left, :top, right: rc.right, bottom: rc.bottom)
				FillRect(hdc, rcSel, .brushes.background)
				}
			}
		}

	paintMarkCol(hdc, rc, top, topRow, numRows)
		{
		if rc.left >= .markColWidth
			return true

		rcCell = Object(left: 0, top: top - .rowHeight,
			right: .markColWidth, bottom: top)
		for (row = topRow; row < numRows and rcCell.top < rc.bottom; ++row)
			{
			rcCell.top += .rowHeight
			rcCell.bottom += .rowHeight
			FillRect(hdc, rcCell, GetSysColorBrush(COLOR.BTNFACE))
			if row is .focused
				.drawLineFocus(hdc, Object(top: rcCell.top, bottom: rcCell.bottom))
			if .data[row].listrow_deleted is true
				.drawDeleteMark(rcCell, hdc)
			}
		if rcCell.bottom < rc.bottom
			{
			rcCell.top = rcCell.bottom
			rcCell.bottom = rc.bottom
			FillRect(hdc, rcCell, GetSysColorBrush(COLOR.BTNFACE))
			}
		if rc.right < .markColWidth
			return false
		rc.left = .markColWidth
		ExcludeClipRect(hdc, 0, rc.top, .markColWidth, rc.bottom)
		return true
		}
	drawDeleteMark(rcCell, hdc)
		{
		padding = ScaleWithDpiFactor(4) /*= padding*/
		cellWidth = rcCell.right - rcCell.left
		cellHeight = rcCell.bottom - rcCell.top
		wh = Min(cellWidth, cellHeight) - padding * 2
		leftPad = (cellWidth - wh) / 2
		topPad = (cellHeight - wh) / 2
		.deleteImage.Draw(hdc, rcCell.left + leftPad, rcCell.top + topPad,
			wh, wh, .brushes.delete)
		}
	paintVisibleRows(hdc, rc, numRows, numCols, startFrom)
		{
		topRow = startFrom.topRow
		top = startFrom.top
		left = startFrom.left
		col = startFrom.col
		// for each row needing painting
		for (row = topRow; row < numRows and top < rc.bottom; top += .rowHeight, ++row)
			{
			rec = .data[row]
			.prepareRow(rc, top, row, rec, hdc)

			// for each column needing painting...
			for (c = col, x = left + .horzMargin; c < numCols and x < rc.right;
				x += .widths[c], ++c)
				if .widths[c] > 2 * .horzMargin
					.formatting.PaintCell(.columns[c], x, top + 2,
						.widths[c] - 2 * .horzMargin, .rowHeight - 2, rec)
			}
		return top
		}

	prepareRow(rc, top, row, rec, hdc)
		{
		rcSel = Object(left: rc.left, :top, right: rc.right, bottom: top + .rowHeight)
		.drawBackground(row, rec, hdc, rcSel)

		.drawFocusedRow(row, hdc, rcSel)

		.paintInvalidCells(rec, rcSel, hdc)

		SetTextColor(hdc, .getTextColor(row))
		}
	drawBackground(row, rec, hdc, rcSel)
		{
		if .RowSelected?(row) and (.alwaysHighlightSelected or .HasFocus?())
			brush = .brushes.focused
		else if .RowHighlighted?(row)
			brush = .brushes[HIWORD(rec.listrow_flags)]
		else
			brush = .noShading or row % 2 is 1 ? .brushes.background : .brushes.shaded
		.formatting.SetBackgroundBrush(brush)
		FillRect(hdc, rcSel, brush)
		}
	paintInvalidCells(rec, rcSel, hdc)
		{
		// paint invalid cells for row
		if rec.Member?("list_invalid_cells")
			for (column in rec.list_invalid_cells.Members())
				{
				if false is cc = .columns.Find(column)
					continue
				rect = .getColRect(cc)
				rcSel.left = rect.GetX()
				rcSel.right = rect.GetX() + rect.GetWidth()
				FillRect(hdc, rcSel, .brushes.invalid)
				}
		}
	drawFocusedRow(row, hdc, rcSel)
		{
		if row is .focused
			{
			SetTextColor(hdc, GetSysColor(COLOR.WINDOWTEXT))
			.drawLineFocus(hdc, rcSel)
			}
		}
	drawLineFocus(hdc, rcSel)
		{
		rcSel.left = -2  // only draw top and bottom dotted lines for focus
		rcSel.right = 20000
		// DrawFocusRect() is influenced by text color
		DrawFocusRect(hdc, rcSel)
		}
	getTextColor(row)
		{
		if .RowHovered?(row)
			return CLR.silver

		highlight = .editor is false and
			.RowSelected?(row) and
			(.HasFocus?() or .alwaysHighlightSelected)
		return GetSysColor(highlight ? COLOR.HIGHLIGHTTEXT : COLOR.WINDOWTEXT)
		}
	SETFOCUS(wParam)
		{
		if .editor isnt false and not .editor.Destroyed?()
			{
			if not .editor.ClosingListEdit and not .editor.ChildOf?(wParam)
				{
				SetActiveWindow(.editor.Hwnd)		// put the editor on top
				return 0
				}
			.editor.Return()						// or end it
			.lastEdit = Date()
			if GetActiveWindow() isnt .Window.Hwnd
				.Defer({ SetActiveWindow(.Window.Hwnd) })
			}
		.repaintSelection()
		return 0
		}
	KILLFOCUS()
		{
		.repaintSelection()
		return 'callsuper'
		}
	GETDLGCODE(wParam, lParam)
		{
		keys = .Send('List_GetDlgCode', :wParam, :lParam)
		return keys is 0 ? DLGC.WANTCHARS | DLGC.WANTARROWS : keys
		}
	SetListFocus()
		{
		.FinishEdit()		// tell editor window to close
		SetFocus(.Hwnd)
		}
	move(direction)
		{
		sels = .getSelection()
		if .readOnly or sels.Size() isnt 1 or true isnt .Send("List_AllowMove")
			return
		fromRow = sels[0]
		toRow = fromRow + direction
		if toRow < 0 or toRow >= .GetNumRows()
			return
		.data.Swap(fromRow, toRow)
		.focused = toRow
		.RepaintRow(fromRow)
		.RepaintRow(toRow)
		.ScrollRowToView(toRow)
		.sortCol = false
		.Send("List_Move", fromRow, toRow)
		}
	dragging: false
	RBUTTONDOWN()
		{
		SetFocus(.Hwnd)
		if not .contextOnly
			.Send("List_RightClick")
		return 0
		}
	LBUTTONDOWN(lParam)
		{
		if .contextOnly
			return 0
		row = .GetDataRowFromY(HISWORD(lParam))
		col = .GetColFromX(LOSWORD(lParam))
		shift = KeyPressed?(VK.SHIFT)
		control = KeyPressed?(VK.CONTROL)
		if .editing?(shift, control, row, col)
			return 0
		SetFocus(.Hwnd)

		// when another ctrl loses focus, it could trigger list destroy (dynamic layouts)
		if .Destroyed?() is true
			return 0
		if 0 isnt .Send("List_SingleClick", .data.Member?(row) ? row : false, col)
			return 0

		return .updateSelection(row, shift, control)
		}

	editing?(shift, control, row, col)
		{
		return not shift and not control and .allowTab and
			row < .GetNumRows() and col < .GetNumCols() and
			.editingRecent(row, col)
		}

	editingRecent(row, col)
		{
		return (.fullEditMode or row is .focused) and
			Date().MinusSeconds(.lastEdit) < .12 and /*= determines if from recent edit */
			.edit(col, row)
		}

	updateSelection(row, shift, control)
		{
		if row >= .GetNumRows()
			{
			if not shift and not control
				.ClearSelect()
			return 0
			}
		if not shift or .multiSelect is false
			.selectFocus(row, control, false, true)
		else
			.select(row - .focused, false, true, true)
		if .dragging?(shift, control)
			{
			.dragging = true
			SetCapture(.Hwnd)
			}
		return 0
		}

	dragging?(shift, control)
		{
		return not .readOnly and not shift and
			not control and true is .Send("List_AllowMove")
		}

	MOUSEMOVE(lParam)
		{
		if not .dragging and not .indicateHovered
			return 0

		TrackMouseEvent(Object(cbSize: TRACKMOUSEEVENT.Size(),
			dwFlags: TME.LEAVE, hwndTrack: .Hwnd))

		newRow = Min(.getRowFromY(HISWORD(lParam)), .GetNumVisibleRows() - 1)
		newRow = Min(newRow + .rowOffset, .GetNumRows() - 1)

		.DeHoverRow()
		.HoverRow(newRow)
		if not .dragging
			return 0

		SetCursor(LoadCursor(ResourceModule(), IDC.DRAG1))
		if .focused isnt newRow
			{						// ensure Swap() and Send() is 1 row a time
			for (inc = newRow > .focused ? 1 : -1; .focused isnt newRow; .focused += inc)
				{
				.RepaintRow(.focused)
				.data.Swap(.focused, .focused + inc)
				.Send("List_Move", .focused, .focused + inc)
				}
			.RepaintRow(.focused)
			.sortCol = false
			}
		return 0
		}

	hoverRow: false
	MOUSELEAVE()
		{
		.DeHoverRow()
		return 0
		}

	EndDrag()
		{
		if .dragging
			{
			ReleaseCapture()
			.dragging = false
			}
		}
	LBUTTONUP(lParam)
		{
		if .contextOnly
			return 0

		.EndDrag()

		// get row/col from coordinates in lParam
		row = .GetDataRowFromY(HISWORD(lParam))
		col = .GetColFromX(LOSWORD(lParam))
		.Send("List_LButtonUp", .data.Member?(row) ? row : false, col)
		return 'callsuper'
		}
	LBUTTONDBLCLK(lParam)
		{
		if .contextOnly
			return 0
		// get row/col from lParam
		row = Min(.GetDataRowFromY(HISWORD(lParam)), .GetNumRows())
		col = Min(.GetColFromX(LOSWORD(lParam)), .GetNumCols() - 1)
		datarow = .data.Member?(row) ? row : false
		// update selection and focused row
		if datarow is false
			.ClearSelect()
		else
			.selectFocus(row)
		dbl_click_result = .Send("List_DoubleClick", datarow, col)
		// added zoom here to make sure the logic still follows the same sequence
		if .Send("List_AllowZoom", datarow, col) is true
			.zoomOnField(datarow, col, .GetCol(col))
		.editFromDblClick(dbl_click_result, datarow, row, col)
		return 0
		}
	zoomOnField(row, col, field)
		{
		if row isnt false and col isnt false and
			String?(value = .GetField(field)) and value isnt "" and
			false isnt ctrl = GetControlClass.FromField(field)
			ctrl.ZoomReadonly(value)
		}
	editFromDblClick(dbl_click_result, datarow, row, col)
		{
		if .edtiable?(dbl_click_result, datarow, row)
			{
			col = (col is 0 and .markColWidth isnt 0) ? 1 : col
			if not .Edit(col, row, 1, canMoveRows: false) and col > 0 // begin edit
				.Edit(col - 1, row, -1, canMoveRows: false)
			}
		}
	edtiable?(dbl_click_result, datarow, row)
		{
		return (dbl_click_result is 0 and not .readOnly) and
			(datarow is false or .data[row].listrow_deleted isnt true)
		}
	GetColFromX(x)
		{
		numCols = .GetNumCols()
		for (col = 0, left = -.horzOffset; col < numCols; col++)
			if x < left += .widths[col]
				break
		return col
		}
	CONTEXTMENU(lParam)
		{
		.SetListFocus()

		// when another ctrl loses focus, it could trigger list destroy (dynamic layouts)
		if .Destroyed?() is true
			return 0
		x = LOSWORD(lParam)
		y = HISWORD(lParam)
		ScreenToClient(.Hwnd, pt = Object(:x, :y))
		if pt.y < .header.Ymin and not .contextOnly
			{
			result = .Send("List_HeaderContextMenu", x, y)
			if result is 0
				.buildHeaderContextMenu(x, y)
			}
		else
			.contextMenuFromRow(pt, x, y)
		return 0
		}

	buildHeaderContextMenu(x, y)
		{
		contextItems = Object()
		if .resetColumns
			contextItems.Add('Reset Columns')
		if .customizeColumns
			contextItems.Add('Customize Columns...')
		if .sortSaveHandler isnt false
			contextItems.Append(.buildSortMenu())
		if contextItems.Empty?() is false
			ContextMenu(contextItems).ShowCallCascade(this, x, y)

		}

	buildSortMenu()
		{
		sortMenu = Object('Reset Sort to System Default', '')
		if .sortCol is false
			sortMenu.Add(Object(name: 'Set as Default Sort for Current User',
				state: MFS.DISABLED))
		else
			{
			field = .columns[.sortCol.Abs() - 1]
			header = .header.GetHeaderText(field)
			sortMenu.Add(Object(
				name: 'Set (' $ header $ ') as Default Sort for Current User',
				cmd: 'Set as Default Sort'))
			}
		return sortMenu
		}

	contextMenuFromRow(pt, x, y)
		{
		row = .GetDataRowFromY(pt.y)
		sel = .getSelection()
		if not .contextOnly or (sel isnt #() and sel[0] is row)
			{
			if row >= .GetNumRows()
				.ClearSelectFocus()
			else if not .RowSelected?(row)
				.SetSelection(row)	// new selection if not selected
			.Send("List_ContextMenu", x, y)
			}
		}

	On_Context_Reset_Columns()
		{
		if not .resetColumns
			return
		reverse = .sortCol < 0
		fieldName = Number?(.sortCol)
			? .columns[.sortCol.Abs() - 1]
			: false

		if .customizeColumns is true
			UserColumns.Reset(
				this, .columnsSaveName, .origColumns, deletecol: false, load_visible?:)
		else
			.SetColumns(columns: .origColumns, reset:)
		if false is colNum = .columns.Find(fieldName)
			return
		.sortCol = (colNum + 1) * (reverse ? -1 : 1)
		.setSortIndicator(.sortCol)
		}

	On_Context_Customize_Columns()
		{
		if 0 is mandatory = .Send('List_MandatoryColumns')
			mandatory = #()
		CustomizeColumnsDialog(.Hwnd, this, .origColumns, .columnsSaveName, mandatory,
			headerSelectPrompt: .headerSelectPrompt)
		}

	On_Context_Reset_Sort_to_System_Default()
		{
		if .sortSaveHandler is false or false is .sortCol
			return

		defaultSort = (.sortSaveHandler)(reset:)
		.signedSort(defaultSort)
		}
	On_Context_Set_as_Default_Sort()
		{
		if .sortSaveHandler is false or false is .sortCol
			return

		fieldName = .columns[.sortCol.Abs() - 1]
		save = ((.sortCol < 0) ? 'reverse ' : '') $ fieldName
		(.sortSaveHandler)(:save)
		}
	HSCROLL(wParam)
		{
		return .scroll_list(wParam)
		}
	VSCROLL(wParam)
		{
		return .scroll_list(wParam, vert:)
		}
	scroll_list(wParam, vert = false)
		{
		SetFocus(.Hwnd)
		fns = vert ? .vscrollfns : .hscrollfns
		(fns[LOWORD(wParam)])(:wParam)
		return 0
		}
	getter_hscrollfns()
		{
		ob = Object().Set_default(function () { })
		ob[SB.LEFT]				= { 		.scroll(-.horzOffset, 0)					}
		ob[SB.RIGHT]			= { 		.scroll(.GetTotalColWidths(), 0)			}
		ob[SB.LINELEFT]			= { 		.scroll(-10, 0) /*= move distance */		}
		ob[SB.LINERIGHT]		= { 		.scroll(10, 0) /*= move distance */			}
		ob[SB.PAGELEFT]			= { 		.scroll(-.GetClientRect().GetWidth(), 0)	}
		ob[SB.PAGERIGHT]		= { 		.scroll(.GetClientRect().GetWidth(), 0)		}
		ob[SB.THUMBTRACK]		= {|wParam| .scroll(HIWORD(wParam) - .horzOffset, 0)	}
		return .hscrollfns = ob // once only
		}
	getter_vscrollfns()
		{
		ob = Object().Set_default(function () { })
		ob[SB.LEFT]				= { 		.scroll(0, -.rowOffset) 					}
		ob[SB.RIGHT]			= { 		.scroll(0, .bottom())						}
		ob[SB.LINELEFT]			= { 		.scroll(0, -1) 								}
		ob[SB.LINERIGHT]		= { 		.scroll(0, 1) 								}
		ob[SB.PAGELEFT]			= { 		.scroll(0, -.GetNumVisibleRows() + 1) 		}
		ob[SB.PAGERIGHT]		= { 		.scroll(0, .GetNumVisibleRows() - 1) 		}
		ob[SB.THUMBTRACK]		= {|wParam| .scroll(0, HIWORD(wParam) - .rowOffset) 	}
		return .vscrollfns = ob // once only
		}
	bottom()
		{
		// plus 2 (1 to handle going from 0 based to 1 based, and 1 for blank row at end)
		return .GetNumRows() - .GetNumVisibleRows() - .rowOffset + 2
		}
	MOUSEWHEEL(wParam)
		{
		scroll = GetWheelScrollInfo(wParam)
		lines = scroll.lines
		if scroll.page? or lines.Abs() > .GetNumVisibleRows() - 1
			lines = scroll.down? ? -.GetNumVisibleRows() + 1 : .GetNumVisibleRows() - 1
		.scroll(0, -lines)
		return 0
		}
	KEYDOWN(wParam, lParam)
		{
		if .contextOnly
			return 0
		ctrl = KeyPressed?(VK.CONTROL)
		shift = KeyPressed?(VK.SHIFT)
		.keydown(wParam, shift, ctrl)
		if .focused isnt false	// if actual rows, .focused should be set !
			{
			// two methods because switch on key is too big for one method...
			.keydown_focused_scroll_keys(wParam, shift, ctrl)
			.keydown_focused_selection_keys(wParam, shift, ctrl)
			}
		.Send('List_KeyDown', :wParam, :lParam)
		return 0
		}
	keydown(wParam, shift, ctrl) // independent of actual rows/focused
		{
		(.keydown_fns[wParam])(:shift, :ctrl)
		}
	getter_keydown_fns()
		{
		ob = Object().Set_default(function () { })
		ob[VK.LEFT] = 	{|ctrl| .HSCROLL(ctrl ? SB.LEFT : SB.LINELEFT) }
		ob[VK.RIGHT] =  {|ctrl|	.HSCROLL(ctrl ? SB.RIGHT : SB.LINERIGHT) }
		ob[VK.F2] = // edit focused row
			{
			if .readOnly isnt true and .focused isnt false
				.Edit(0, .focused, 1, canMoveRows: false)
			}
		ob[VK.INSERT] = .insertNewRow
		ob[VK.F5] =	// refresh display
			{
			.UpdateScrollbars()
			.Repaint()
			}
		ob[VK.F8] =
			{
			if Suneido.User is 'default'
				Inspect(this)
			}
		ob[VK.SPACE] = .listToggle
		return .keydown_fns = ob // once only
		}
	insertNewRow(shift)
		{
		if .readOnly isnt true
			{
			if .getSelection() is #() or .focused is false
				{
				insertAt = .GetNumRows()
				prevRow = insertAt - 1
				}
			else
				{
				insertAt = shift ? .focused + 1 : .focused
				prevRow = .focused
				}
			newRow = Record()
			if false isnt .Send("List_WantNewRow", :prevRow, record: newRow)
				{ // parent didn't disallow
				.InsertRow(insertAt, newRow)
				.Send("List_NewRowAdded", insertAt, record: newRow)
				.Edit(0, insertAt, 1, canMoveRows: false)
				}
			}
		}
	listToggle(shift = false, ctrl = false)
		{
		if .checkBoxColumn is false or
			shift isnt false or ctrl isnt false
			return

		sel = .getSelection()
		if sel.Empty?()
			return

		for row in sel
			{
			data = .GetRow(row)
			if false isnt .Send('List_AllowToggle', :data, :row)
				{
				data[.checkBoxColumn] = data[.checkBoxColumn] isnt true
				.Send('List_AfterToggle', :data, :row)
				.RepaintRow(row)
				}
			}
		}

	keydown_focused_scroll_keys(wParam, shift, ctrl) // based on focused row
		{
		(.keydown_focused_fns[wParam])(:shift, :ctrl)
		}
	getter_keydown_focused_fns()
		{
		ob = Object().Set_default(function () { })
		ob[VK.PRIOR] = .vk_prior
		ob[VK.NEXT] = .vk_next
		ob[VK.HOME] = .vk_home
		ob[VK.END] = .vk_end
		return .keydown_focused_fns = ob // once only
		}
	vk_prior(shift, ctrl)
		{
		if .focused is .rowOffset and .rowOffset > 0
			.VSCROLL(SB.PAGEUP)
		if not shift or .multiSelect is false
			.selectFocus(.rowOffset, ctrl, false, not ctrl)
		else
			.select(.rowOffset - .focused, false, true, true)
		}
	vk_next(shift, ctrl)
		{
		if .focused is .rowOffset + .GetNumVisibleRows() - 2
			.VSCROLL(SB.PAGEDOWN)
		if not shift or .multiSelect is false
			.selectFocus(.rowOffset + .GetNumVisibleRows() - 2, ctrl, false, not ctrl)
		else
			{
			row = (.rowOffset + .GetNumVisibleRows() - 2) - .focused
			.select(row, false, true, true)
			}
		}
	vk_home(shift, ctrl)
		{
		.VSCROLL(SB.TOP)
		if not shift or .multiSelect is false
			.selectFocus(0, ctrl, false, not ctrl)
		else
			.select(-.focused, false, true, true)
		}
	vk_end(shift, ctrl)
		{
		.VSCROLL(SB.BOTTOM)
		if not shift or .multiSelect is false
			.selectFocus(.GetNumRows() - 1, ctrl, false, not ctrl)
		else
			.select(.GetNumRows() - .focused - 1, false, true, true)
		}
	keydown_focused_selection_keys(wParam, shift, ctrl) // based on focused row
		{
		(.keydown_selection_fns[wParam])(:ctrl, :shift)
		}
	getter_keydown_selection_fns()
		{
		ob = Object().Set_default(function () { })
		ob[VK.UP] = {|ctrl, shift| .vk_up_down(ctrl, shift, dir: -1) }
		ob[VK.DOWN] = {|ctrl, shift| .vk_up_down(ctrl, shift, dir: 1) }
		ob[VK.SPACE] =
			{ |ctrl, shift|
			if ctrl
				.select(0, ctrl, shift, true)
			}
		ob[VK.DELETE] =
			{
			if .readOnly isnt true and false isnt .Send("List_DeleteKeyDown")
				.DeleteSelection()
			}
		.keydown_selection_fns = ob // once only
		}
	vk_up_down(ctrl, shift, dir)
		{
		if ctrl
			.move(dir)
		else
			.select(dir, ctrl, shift and .multiSelect isnt false, not ctrl)
		}

	QueueDeleteAttachmentFile(newFile, oldFile, name, action)
		{
		return .Send('QueueDeleteAttachmentFile', newFile, oldFile, name, action)
		}

	trackValid: false
	Valid?()
		{
		if .trackValid
			{
			// loop through the lines and check if any invalid data
			for rec in .data
				if .RowHasInvalidCell?(rec)
					return false
			}
		return true
		}

	ConfirmDestroy()
		{
		return not .trackValid ? true : .Valid?()
		}

	Destroy()
		{
		.FinishEdit()
		if .deleteImage isnt false
			.deleteImage.Close()
		if .columnsSaveName isnt false
			UserColumns.Save(.columnsSaveName, this, .origColumns)
		for (brush in .brushes)
			DeleteObject(brush)
		.header.Destroy()
		.formatting.Destroy()
		if .trackValid
			.Window.RemoveValidationItem(this)
		super.Destroy()
		}
	}
