// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
ListBodyBaseComponent
	{
	contextOnly: false
	New(.parentEl, .header, noShading = false, .indicateHovered = false)
		{
		.scrollContainerEl = .Parent.El

		.tbody = CreateElement('tbody', parentEl,
			'su-listbody' $
			(noShading ? '' : ' su-shadinglist') $
			(.indicateHovered ? ' su-listbody-hover' : ''))

		.data = Object()

		.tbody.AddEventListener('mousedown', .mousedown)
		.tbody.AddEventListener('dblclick', .doubleClick)
		.tbody.AddEventListener('contextmenu', .contextMenu)

		.deleteMark = IconFontHelper.GetCode('delete')
		.selectColor = ToCssColor(0xd77600/*=COLOR.HIGHLIGHT*/)
		}

	Reset()
		{
		.tbody.innerText = ''
		.prepareFillRow()
		.rows = Object()
		}

	prepareFillRow()
		{
		.fillRow = CreateElement('tr', .tbody)
		.fillRow.SetStyle('height', '100%')
		.SetCellAttributes(.fillRow, y: -1, type: 'fill-row')

		td = CreateElement('td', .fillRow)
		.SetCellAttributes(td, x: 0, type: 'mark-cell')

		.header.ForEachHeadCol()
			{ |col|
			td = CreateElement('td', .fillRow)
			.SetCellAttributes(td, x: col, type: 'fill-cell')
			}
		if .header.StretchCol() isnt false
			return
		el = CreateElement('td', .fillRow)
		.SetCellAttributes(el, x: .header.GetColsNum(), type: 'fill-cell')
		}

	ClearData()
		{
		.data = Object()
		.Reset()
		}

	UpdateDataCell(row, col, newCell)
		{
		.updateDataCell(row, .header.ToWebColIndex(col), newCell)
		}

	updateDataCell(row, col, newCell)
		{
		el = .getCellEl(row, col).firstChild
		.SetCellValue(el, newCell)
		.data[row][.header.GetField(col)] = newCell
		}

	UpdateData(row, record)
		{
		.data[row] = record
		if .header.HasMarkCol?()
			.updateMarkCol(row, record)
		.header.ForEachHeadCol()
			{ |col, field|
			.updateDataCell(row, col, record[field])
			}
		}

	UpdateBatch(batch)
		{
		for row in batch.Members()
			{
			.UpdateData(row, batch[row][0])
			if batch[row][1]
				.SelectRow(row)
			else
				.DeSelectRow(row)
			}
		}

	updateMarkCol(row, record)
		{
		if not .header.HasMarkCol?()
			return

		el = .getCellEl(row, 0)
		if record.Member?(#listrow_deleted) and record.listrow_deleted.data is "true"
			el.innerText = .deleteMark.Chr()
		else
			el.innerText = ''
		}

	DeleteRows(rowsToDelete)
		{
		if rowsToDelete.Empty?()
			return

		rowsToDelete.Sort!()
		offset = 0
		i = 0
		row = rowsToDelete[i]

		while (row < .getNumRows())
			{
			if i < rowsToDelete.Size() and row is rowsToDelete[i]
				{
				offset++
				i++
				.rows[row].Remove()
				}
			else
				{
				.data[row - offset] = .data[row]
				.SetCellAttributes(.rows[row], y: row - offset)
				.rows[row - offset] = .rows[row]
				}
			row++
			}
		Assert(offset is: rowsToDelete.Size())
		.data = .data[..-offset]
		.rows = .rows[..-offset]
		}

	InsertData(row, record)
		{
		n = .getNumRows()
		for (i = n; i > row; i--)
			{
			.rows[i] = .rows[i - 1]
			.data[i] = .data[i - 1]
			.SetCellAttributes(.rows[i], y: i)
			}
		fragment = SuUI.GetCurrentDocument().CreateDocumentFragment()
		.addData(record, row, fragment)
		.tbody.InsertBefore(fragment, row >= n ? .fillRow : .rows[row + 1])
		}

	SwapRows(from, to)
		{
		.SetCellAttributes(.rows[from], y: to)
		.SetCellAttributes(.rows[to], y: from)
		if to - from is 1
			.tbody.InsertBefore(.rows[to], .rows[from])
		else
			{
			.tbody.InsertBefore(.rows[from], .rows[to])
			.tbody.InsertBefore(.rows[to],
				from + 1 >= .getNumRows() ? .fillRow : .rows[from + 1])
			}
		.data.Swap(from, to)
		.rows.Swap(from, to)
		}

	HeaderChanged(data)
		{
		.Reset()
		if data isnt false
			.data = data
		fragment = SuUI.GetCurrentDocument().CreateDocumentFragment()
		for rowIdx in .data.Members().Sort!()
			.addRow(.data[rowIdx], rowIdx, fragment)
		.tbody.InsertBefore(fragment, .fillRow)
		}

	AddBatch(batch)
		{
		fragment = SuUI.GetCurrentDocument().CreateDocumentFragment()
		origSize = .data.Size()
		for m in batch.Members().Sort!()
			{
			.addData(batch[m], origSize + m, fragment)
			}
		.tbody.InsertBefore(fragment, .fillRow)
		}

	addData(rec, at, container = false)
		{
		.data[at] = rec
		.addRow(rec, at, container)
		}

	addRow(rec, rowIdx, container)
		{
		if container is false
			container = .tbody

		row = CreateElement('tr', container)
		.SetStyles(Object('line-height': .GetRowHeight() $ 'px'), row)
		.SetCellAttributes(row, y: rowIdx, type: 'data-row')

		td = CreateElement('td', row)
		.SetCellAttributes(td, x: 0, type: 'mark-cell')

		.header.ForEachHeadCol()
			{ |col, field, width|
			cell = rec[field]
			td = CreateElement('td', row)
			.CreateCellElement(td, cell, col, 'su-listbody-cell')
			if width is 0
				td.SetStyle('display', 'none')
			}
		// empty col for filling remaining space
		el = CreateElement('td', row)
		// invisible character to take vertical space when the row is empty
		el.innerHTML = '&#8205;'
		.SetCellAttributes(el, x: .header.GetColsNum(), type: 'empty-cell')
		.rows[rowIdx] = row
		}

	getNumRows()
		{
		return .data.Size()
		}

	// ListEditWindowComponent
	GetCellEl(row, col)
		{
		return .getCellEl(row, .header.ToWebColIndex(col))
		}

	getCellEl(row, col)
		{
		return .rows[row].children.item(col)
		}

	mousedown?: false
	mouseEventId: 0
	mousedown(event)
		{
		if .contextOnly or event.button isnt 0 or .freeze is true
			return
		.handleMouseEvent(event, 'LBUTTONDOWN', {
			.mousedown? = true
			.StartMouseTracking(.mouseup, .mousemove)
			})
		}

	dragging: false
	focused: false
	origFocused: false
	mousemove(event)
		{
		if .dragging is false
			return

		newrow = .getRowFromY(event.clientY)
		if .focused isnt newrow
			{
			.SwapRows(.focused, newrow)
			.focused = newrow
			}
		}

	getRowFromY(y)
		{
		containerClientRect = SuRender.GetClientRect(.scrollContainerEl)
		offsetHeight = Max(
			Min(y - containerClientRect.top, .scrollContainerEl.clientHeight) -
				.header.GetOffsetHeight(),
			0)
		height = offsetHeight + .scrollContainerEl.scrollTop
		return Min((height / .GetRowHeight()).Floor(), .getNumRows() - 1)
		}

	List_AllowDragging(focused, mouseEventId)
		{
		if mouseEventId is .mouseEventId and .mousedown? is true
			{
			.dragging = true
			.origFocused = .focused = focused
			.tbody.classList.Add('su-list-dragging')
			}
		}

	freeze: false
	mouseup(event)
		{
		if .dragging is true
			{
			.tbody.classList.Remove('su-list-dragging')
			if .focused isnt .origFocused
				.Parent.Event('ListMoveRow', .origFocused, .focused)
			}
		.mousedown? = false
		.origFocused = .focused = .dragging = false
		.StopMouseTracking()
		.handleMouseEvent(event, 'LBUTTONUP')
		// Freeze to avoid sending the duplicate the mousedown and mouseup events
		// when double clicking. Somehow, dblclick doesn't fire when double clicking
		// the rect type cell without this.
		.freeze = true
		SuDelayed(100/*=cooldown*/, .releaseFreeze)
		}

	releaseFreeze()
		{
		.freeze = false
		}

	doubleClick(event)
		{
		if .contextOnly
			return

		.handleMouseEvent(event, 'LBUTTONDBLCLK', freeze?:)
		}

	handleMouseEvent(event, name, block = false, freeze? = false)
		{
		target = event.target
		// target can be a Text Node when this is triggered by dragstart
		if not .tbody.Contains(target) or target.GetDefault(#tagName, false) is false
			return
		if target.tagName is 'IMG'
			target = target.parentElement
		else if target.classList.Contains('su-listbody-cell-rect')
			target = target.parentElement

		switch (target.GetAttribute('data-type'))
			{
		case 'cell', 'empty-cell', 'mark-cell', 'fill-cell':
			if block isnt false
				block()

			row = .getRow(target)
			col = .header.ToControlColIndex(Number(target.GetAttribute('data-x')))
			.Parent.RunWhenNotFrozen()
				{
				if freeze?
					.Parent.EventWithFreeze(name, row, col,
						shift: event.shiftKey, control: event.ctrlKey,
						mouseEventId: ++.mouseEventId)
				else
					.Parent.Event(name, row, col,
						shift: event.shiftKey, control: event.ctrlKey,
						mouseEventId: ++.mouseEventId)
				}
		default:
			}
		}

	getRow(target)
		{
		if target.tagName is 'DIV'
			target = target.parentElement
		row = Number(target.parentElement.GetAttribute('data-y'))
		if row is -1
			row = .getNumRows()
		return row
		}

	SelectRow(row)
		{
		Assert(.data hasMember: row, msg: 'sujslib:ListBodyComponent.SelectRow')
		.rows[row].SetAttribute('data-selected', 'true')
		.SetStyles(Object(
			'background-color': .selectColor,
			'color': 'white'), .rows[row])
		}

	DeSelectRow(row)
		{
		Assert(.data hasMember: row, msg: 'sujslib:ListBodyComponent.DeSelectRow')
		.rows[row].SetAttribute('data-selected', 'false')
		.SetStyles(Object(
			'background-color': .getHighlight(row),
			'color': ''), .rows[row])
		}

	AddHighlightRow(row, color)
		{
		Assert(.data hasMember: row, msg: 'sujslib:ListBodyComponent.AddHighlightRow')
		.rows[row].SetAttribute('data-highlight', color = ToCssColor(color))
		if not .selected?(row)
			{
			.rows[row].SetStyle('background-color', color)
			}
		}

	RemoveHighlightRow(row)
		{
		Assert(.data hasMember: row, msg: 'sujslib:ListBodyComponent.RemoveHightlighRow')
		.rows[row].SetAttribute('data-highlight', '')
		if not .selected?(row)
			.rows[row].SetStyle('background-color', '')
		}

	getHighlight(row)
		{
		if not .rows[row].HasAttribute('data-highlight')
			return ''
		return .rows[row].GetAttribute('data-highlight')
		}

	selected?(row)
		{
		if not .rows[row].HasAttribute('data-selected')
			return false
		return .rows[row].GetAttribute('data-selected') is 'true'
		}

	AllowContextOnly(allow)
		{
		.contextOnly = allow
		}

	ScrollRowToView(row)
		{
		Assert(.rows hasMember: row, msg: 'sujslib:ListBodyComponent.ScrollRowToView')
		headerHeight = .header.GetOffsetHeight()
		rowHeight = .rows[row].offsetHeight
		rowOffsetTop = .rows[row].offsetTop
		scrollTop = .scrollContainerEl.scrollTop
		scrollHeight = .scrollContainerEl.clientHeight

		if scrollTop + headerHeight > rowOffsetTop
			.scrollContainerEl.scrollTop = rowOffsetTop - headerHeight
		else if scrollTop + scrollHeight < rowOffsetTop + rowHeight
			.scrollContainerEl.scrollTop = rowOffsetTop + rowHeight - scrollHeight
		}

	ScrollColToView(col)
		{
		colRect = .header.GetColRect(col)
		scrollLeft = .scrollContainerEl.scrollLeft
		scrollWidth = .scrollContainerEl.clientWidth
		if scrollLeft > colRect.left
			.scrollContainerEl.scrollLeft = colRect.left
		else if scrollLeft + scrollWidth < colRect.left + colRect.width
			.scrollContainerEl.scrollLeft = colRect.left + colRect.width - scrollWidth
		}

	GetScrollContainerEl()
		{
		return .scrollContainerEl
		}

	savedVScrollPos: false
	SaveVScrollPos()
		{
		.savedVScrollPos = .scrollContainerEl.scrollTop
		}
	RestoreVscrollPos()
		{
		if .savedVScrollPos is false
			return
		.scrollContainerEl.scrollTop = .savedVScrollPos
		}

	contextMenu(event)
		{
		target = event.target
		if target.tagName in ('IMG')
			target = target.parentElement
		else if target.tagName is 'DIV' and
			target.classList.Contains('su-listbody-cell-rect')
			target = target.parentElement

		switch (target.GetAttribute('data-type'))
			{
		case 'cell', 'empty-cell', 'mark-cell', 'fill-cell':
			row = .getRow(target)
			col = .header.ToControlColIndex(Number(target.GetAttribute('data-x')))
			.Parent.RunWhenNotFrozen()
				{
				.Parent.EventWithOverlay('CONTEXTMENU', event.clientX, event.clientY,
					:row, :col)
				}
		default:
			}
		event.PreventDefault()
		event.StopPropagation()
		}

	MeasureWidth(i)
		{
		width = 0
		for row in .rows
			{
			div = row.children.Item(i).firstChild
			m = SuRender().GetTextMetrics(div, div.innerText)
			width = Max(width, m.width)
			}
		return width
		}

	// Called by ListHeaderComponent
	HideCol(col)
		{
		for row in .rows
			row.children.item(col).SetStyle('display', 'none')
		.fillRow.children.item(col).SetStyle('display', 'none')
		}

	hoverRow: false
	HoverRow(row)
		{
		if .indicateHovered is false
			return
		if .RowHovered?(row) or not .rows.Member?(row)
			return

		.hoverRow = row
		.rows[row].classList.Add('su-listbody-hovered')
		}

	RowHovered?(row)
		{
		return row isnt false and .hoverRow is row
		}

	DeHoverRow()
		{
		if .hoverRow is false or not .rows.Member?(.hoverRow)
			{
			.hoverRow = false
			return
			}
		.rows[.hoverRow].classList.Remove('su-listbody-hovered')
		.hoverRow = false
		}

	ReorderList(newOrders)
		{
		for i in newOrders.Members()
			{
			if newOrders[i] < 0
				continue
			cur = i
			while newOrders[cur] >= 0
				{
				.SwapRows(i, newOrders[cur])
				temp = cur
				cur = newOrders[cur]
				newOrders[temp] = -1
				}
			}
		}
	}
