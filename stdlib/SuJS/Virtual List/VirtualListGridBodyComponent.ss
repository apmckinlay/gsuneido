// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
ListBodyBaseComponent
	{
	contextOnly: false
	New(.parentEl, .header)
		{
		.scrollContainerEl = .Parent.El
		LoadCssStyles('vlist-body-control.css', VirtualListBodyStyles)
		.tbody = CreateElement('tbody', parentEl, 'su-vlistbody su-vshadinglist')

		.data = Object()

		.tbody.AddEventListener('mousedown', .mousedown)
		.tbody.AddEventListener('mousemove', .mousemove)
		.tbody.AddEventListener('dblclick', .doubleClick)
		.tbody.AddEventListener('contextmenu', .contextMenu)

		.observer = SuUI.MakeWebObject('IntersectionObserver',
			.observerFn,
			Object(root: .scrollContainerEl/*, rootMargin: '100px 0px 100px 0px'*/))
		.observerStatus = Object()

		.deleteMark = IconFontHelper.GetCode('delete')

		.expandedCtrls = Object()
		.recycledExpands = Object()
		}

	observerFn(entries, observer/*unused*/)
		{
		for entry in entries
			{
			row = Number(entry.target.GetAttribute('data-y'))
			if not .observerStatus.Member?(row) // has been removed
				continue

			load? = .observerStatus[row] isnt true and entry.isIntersecting is true
			.observerStatus[row] = entry.isIntersecting
			if not load?
				continue

			.removeObserver(row)
			.Parent.Event('VirtualListGridComponent_Load', row)
			Print(:row)
			}
		}

	addObserver(rowIdx)
		{
		Assert(.rows hasMember: rowIdx,
			msg: 'sujslib:VirtualListGridBodyComponent.addObserver')

		rowRect = .rows[rowIdx].GetBoundingClientRect()
		viewRect = .scrollContainerEl.GetBoundingClientRect()
		if rowRect.height is 0 or // invisible
			rowRect.top > viewRect.bottom or
			rowRect.bottom < viewRect.top
			{
			.observer.Observe(.rows[rowIdx])
			.observerStatus[rowIdx] = ''
			}
		else
			{
			// Wait for the finish of all the actions
			// in case the rowIdx is outdated due to insert or delete
			.addDelayedLoad(.rows[rowIdx], .Parent)
			}
		}

	delayedLoads: false
	addDelayedLoad(row, parent)
		{
		if .delayedLoads is false
			.delayedLoads = Object()
		delay = SuDelayed(0)
			{
			parent.Event('VirtualListGridComponent_Load',
				Number(row.GetAttribute('data-y')))
			}
		.delayedLoads.Add(delay)
		}

	cancelDelayedLoads()
		{
		if .delayedLoads isnt false
			{
			.delayedLoads.Each(#Kill)
			.delayedLoads = false
			}
		}

	removeObserver(rowIdx)
		{
		Assert(.rows hasMember: rowIdx,
			msg: 'sujslib:VirtualListGridBodyComponent.removeObserver')
		.observer.Unobserve(.rows[rowIdx])
		.observerStatus.Delete(rowIdx)
		}

	isObserving(rowIdx)
		{
		return .rows.Member?(rowIdx) and .observerStatus.Member?(rowIdx)
		}

	Reset()
		{
		.tbody.innerText = ''
		.prepareFillRow()

		.rows = Object()
		.top = .bottom = 0
		}

	FillRowY: 999999
	prepareFillRow()
		{
		.fillRow = CreateElement('tr', .tbody)
		.fillRow.SetStyle('height', '100%')
		.SetCellAttributes(.fillRow, y: .FillRowY, type: 'fill-row')

		td = CreateElement('td', .fillRow)
		td.SetStyle('width', '100%')
		.SetCellAttributes(td, x: 0, type: 'mark-cell')

		.header.ForEachHeadCol()
			{ |col|
			td = CreateElement('td', .fillRow)
			.SetCellAttributes(td, x: col, type: 'fill-cell')
			}
		el = CreateElement('td', .fillRow)
		.SetCellAttributes(el, x: .header.GetColsNum(), type: 'fill-cell')
		}

	ClearData()
		{
		.data = Object()
		.cancelDelayedLoads()
		.Reset()
		}

	updateDataCell(row, col, newCell)
		{
		el = .getCellEl(row, col).firstChild
		.SetCellValue(el, newCell)
		.data[row][.header.GetField(col)] = newCell
		}

	UpdateData(row, record, keepPos? = false)
		{
		.data[row] = record
		if false is .rows.Member?(row)
			return
		if .header.HasMarkCol?()
			.updateMarkCol(row, record)

		.setHighlight(.rows[row], record.vl_brush)

		.header.ForEachHeadCol()
			{ |col, field|
			.updateDataCell(row, col, record[field])
			}
		if keepPos?
			.Parent.UpdateScroll()
		}

	updateMarkCol(row, record)
		{
		if not .header.HasMarkCol?()
			return

		el = .getCellEl(row, 0)
		.setMarkColValue(row, el, record)
		}

	setHighlight(row, color)
		{
		if color isnt false
			{
			row.classList.Add('su-vlistbody-row-highlighted')
			row.style.setProperty('--su-vlistbody-row-color',
				ToCssColor(color))
			}
		else
			row.classList.Remove('su-vlistbody-row-highlighted')
		}

	setMarkColValue(row, td, record)
		{
		if record.vl_deleted is true
			td.textContent = .deleteMark.Chr()
		else
			{
			td.textContent = ''
			if .expandAttachedRow is row
				AttachElement(.expandDiv, td, false)
			}
		}

	InsertData(row, record, shiftTop?)
		{
		addObserver? = shiftTop? is true
			? .insertRowShiftTop(row, record)
			: .insertRowShiftBottom(row, record)

		if addObserver? is true
			.addObserver(row)
		}

	insertRowShiftTop(row, record)
		{
		addObserver? = .clearObserver(row, -1)
		for (i = .top; i <= row and i < .bottom; i++)
			{
			.rows[i - 1] = .rows[i]
			.data[i - 1] = .data[i]
			.SetCellAttributes(.rows[i - 1], y: i - 1)
			.updateExpandRow(i, i - 1)
			}
		fragment = SuUI.GetCurrentDocument().CreateDocumentFragment()
		.addData(record, row, fragment)
		.tbody.InsertBefore(fragment, row is .bottom - 1 ? .fillRow : .rows[row + 1])

		if .observerStatus.Member?(.top)
			{
			.observerStatus[.top - 1] = .observerStatus[.top]
			.observerStatus.Delete(.top)
			}
		.top--
		return addObserver?
		}

	insertRowShiftBottom(row, record)
		{
		addObserver? = .clearObserver(row)
		for (i = .bottom - 1; i >= row; i--)
			{
			.rows[i + 1] = .rows[i]
			.data[i + 1] = .data[i]
			.SetCellAttributes(.rows[i + 1], y: i + 1)
			.updateExpandRow(i, i + 1)
			}
		fragment = SuUI.GetCurrentDocument().CreateDocumentFragment()
		.addData(record, row, fragment)
		.tbody.InsertBefore(fragment, row is .bottom ? .fillRow : .rows[row + 1])

		if .observerStatus.Member?(.bottom - 1)
			{
			.observerStatus[.bottom] = .observerStatus[.bottom - 1]
			.observerStatus.Delete(.bottom - 1)
			}
		.bottom++
		return addObserver?
		}

	clearObserver(row, adj = 0)
		{
		addObserver? = false
		if row is .bottom + adj and true is addObserver? = .isObserving(.bottom - 1)
			.removeObserver(.bottom - 1)
		else if row is .top + adj and true is addObserver? = .isObserving(.top)
			.removeObserver(.top)
		return addObserver?
		}

	DeleteRecord(rowNum, shiftTop?)
		{
		Assert(.rows hasMember: rowNum,
			msg: 'sujslib:VirtualListGridBodyComponent.DeleteRecord')
		newObserve = shiftTop? is true
			? .deleteRowShiftTop(rowNum)
			: .deleteRowShiftBottom(rowNum)
		if .rows.Size() > 0
			newObserve.Each({ .addObserver(it is #top ? .top : .bottom - 1) })
		}

	deleteRowShiftTop(rowNum)
		{
		newObserve = Object()
		if .isObserving(.top)
			{
			.removeObserver(.top)
			newObserve.Add(#top)
			}
		if rowNum is .bottom - 1 and .isObserving(.bottom - 1)
			{
			.removeObserver(.bottom - 1)
			newObserve.Add(#bottom)
			}

		.rows[rowNum].Remove()
		for (i = rowNum; i > .top; i--)
			{
			.rows[i] = .rows[i - 1]
			.data[i] = .data[i - 1]
			.SetCellAttributes(.rows[i], y: i)
			.updateExpandRow(i - 1, i)
			}
		.rows.Erase(.top)
		.data.Erase(.top)
		.top++
		return newObserve
		}

	deleteRowShiftBottom(rowNum)
		{
		newObserve = Object()
		if .isObserving(.bottom - 1)
			{
			.removeObserver(.bottom - 1)
			newObserve.Add(#bottom)
			}
		if rowNum is .top and .isObserving(.top)
			{
			.removeObserver(.top)
			newObserve.Add(#top)
			}

		.rows[rowNum].Remove()
		for (i = rowNum; i < .bottom - 1; i++)
			{
			.rows[i] = .rows[i + 1]
			.data[i] = .data[i + 1]
			.SetCellAttributes(.rows[i], y: i)
			.updateExpandRow(i + 1, i)
			}
		.rows.Erase(.bottom - 1)
		.data.Erase(.bottom - 1)
		.bottom--
		return newObserve
		}

	HeaderChanged()
		{
		// ClearData will be called
		}

	top: 0
	bottom: 0 // exclusive
	lastLoadOnTop?: false
	AddBatch(batch, newTop, topEnded?, newBottom, bottomEnded?, loadOnTop? = false)
		{
		.updateTopBottom(newTop, newBottom)
		.lastLoadOnTop? = loadOnTop?
		pos = loadOnTop?
			? .scrollContainerEl.scrollHeight - .scrollContainerEl.scrollTop
			: false
		fragment = SuUI.GetCurrentDocument().CreateDocumentFragment()
		for m in batch.Members().Sort!()
			{
			.addData(batch[m], m, fragment)
			}
		if loadOnTop?
			AttachElement(fragment, .tbody, 0)
		else
			.tbody.InsertBefore(fragment, .fillRow)

		.handlerObserver(newTop, topEnded?)
		.handlerObserver(newBottom, bottomEnded?)

		if pos isnt false
			.scrollContainerEl.scrollTop = .scrollContainerEl.scrollHeight - pos
		}

	atBottom?: false
	BeforeResize()
		{
		.atBottom? = (.lastLoadOnTop? or .scrollContainerEl.scrollTop isnt 0) and
			.isInViewWithinParent(.fillRow, .scrollContainerEl)
		}
	AfterResize()
		{
		if .atBottom?
			{
			.atBottom? = false
			el = .scrollContainerEl
			SuDelayed(0, { el.scrollTop = el.scrollHeight })
			}
		}

	isInViewWithinParent(child, parent)
		{
		childRect = child.GetBoundingClientRect();
		parentRect = parent.GetBoundingClientRect();

		return childRect.bottom > parentRect.top and childRect.top < parentRect.bottom
		}

	handlerObserver(pos, ended?)
		{
		if pos isnt false and ended? is false
			.addObserver(pos)
		}

	updateTopBottom(newTop, newBottom)
		{
		if newTop isnt false
			{
			if .isObserving(.top)
				.removeObserver(.top)
			.top = newTop
			}
		if newBottom isnt false
			{
			if .isObserving(.bottom - 1)
				.removeObserver(.bottom - 1)
			.bottom = newBottom + 1
			}
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
		.SetCellAttributes(row, y: rowIdx, type: 'data-row')
		.SetStyles(Object('line-height': .GetRowHeight() $ 'px'), row)

		td = CreateElement('td', row)
		.SetCellAttributes(td, x: 0, type: 'mark-cell')
		.setMarkColValue(row, td, rec)

		.setHighlight(row, rec.vl_brush)
		.header.ForEachHeadCol()
			{ |col, field, width|
			cell = rec[field]
			td = CreateElement('td', row)
			el = .CreateCellElement(td, cell, col, 'su-vlistbody-cell')
			if width is 0
				td.SetStyle('display', 'none')
			}
		// empty col for filling remaining space
		el = CreateElement('td', row)
		// invisible character to take vertical space when the row is empty
		el.innerHTML = '&#8205;'
		.SetCellAttributes(el, x: .header.GetColsNum(), type: 'empty-cell')
		row.AddEventListener('mouseenter', .mouseEnterRow)
		row.AddEventListener('mouseleave', .mouseLeaveRow)
		.rows[rowIdx] = row
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
	mousedown? : false
	mouseEventId: 0
	mousedown(event)
		{
		if .contextOnly or event.button isnt 0 or .freeze is true
			return
		.handleMouseEvent(event, 'LBUTTONDOWN')
			{
			.mousedown? = true
			.StartMouseTracking(.mouseup, .mousemove)
			}
		}

	freeze: false
	mouseup(event)
		{
		if .dragging is true
			{
			.tbody.classList.Remove('su-vlist-dragging')
			if .focused isnt .origFocused
				.Parent.Event('MoveRow', .origFocused, .focused)
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
		if .Destroyed?() or .expandedCtrls.Any?({ it.ctrl.El.Contains(event.target) })
			return
		event.StopPropagation()
		target = event.target
		// target can be a Text Node when this is triggered by dragstart
		if  not .tbody.Contains(target) or target.GetDefault(#tagName, false) is false
			return
		if target.tagName is 'IMG'
			target = target.parentElement
		else if target.classList.Contains('su-listbody-cell-rect')
			target = target.parentElement

		switch (target.GetAttribute('data-type'))
			{
		case 'cell', 'empty-cell', 'fill-cell':
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
		if target.tagName is 'TD'
			target = target.parentElement
		return Number(target.GetAttribute('data-y'))
		}

	SelectRow(row)
		{
		Assert(.data hasMember: row,
			msg: 'sujslib:VirtualListGridBodyComponent.SelectRow')
		.rows[row].classList.Add('su-vlistbody-row-selected')
		}

	DeSelectRow(row)
		{
		if .rows.Member?(row)
			.rows[row].classList.Remove('su-vlistbody-row-selected')
		}

	AllowContextOnly(allow)
		{
		.contextOnly = allow
		}

	ScrollRowToView(row)
		{
		if false is rowEl = row is .FillRowY ? .fillRow : .rows.GetDefault(row, false)
			return
		headerHeight = .header.GetOffsetHeight()
		rowHeight = rowEl.offsetHeight
		rowOffsetTop =rowEl.offsetTop
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

	Release(releaseTo, releaseTop?)
		{
		if releaseTop? is true
			{
			.removeObserver(.top)
			for (i = .top; i <= releaseTo; i++)
				{
				.rows[i].Remove()
				.rows.Erase(i)
				.data.Erase(i)
				}
			.top = releaseTo + 1
			.addObserver(.top)
			}
		else
			{
			.removeObserver(.bottom - 1)
			for (i = .bottom - 1; i >= releaseTo; i--)
				{
				.rows[i].Remove()
				.rows.Erase(i)
				.data.Erase(i)
				}
			.bottom = releaseTo
			.addObserver(.bottom - 1)
			}
		}

	contextMenu(event)
		{
		target = event.target
		switch (target.GetAttribute('data-type'))
			{
		case 'cell', 'empty-cell', 'mark-cell', 'fill-cell':
			row = .getRow(target)
			col =  .header.ToControlColIndex(Number(target.GetAttribute('data-x')))
			.Parent.RunWhenNotFrozen()
				{
				.Parent.EventWithOverlay('CONTEXTMENU', event.clientX, event.clientY,
					:row, :col)
				}
		default:
			}
		event.StopPropagation()
		event.PreventDefault()
		}

	VirtualListExpand_ContructAt(rowIdx, ctrl, rows = false)
		{
		expandRow = CreateElement('tr')
		.SetCellAttributes(expandRow, y: rowIdx, type: 'expanded-row')

		markCell = CreateElement('td', expandRow)
		.SetCellAttributes(markCell, x: 0, type: 'expand-mark-cell')
		.createEditButton(markCell)

		td = CreateElement('td', expandRow)
		td.SetAttribute('colspan', '9999')
		insertBefore = not .rows.Member?(rowIdx + 1)
			? .fillRow
			: .rows[rowIdx + 1]
		.tbody.InsertBefore(expandRow, insertBefore)

		.rows[rowIdx].SetAttribute('data-expanded', '')

		c = .construct(ctrl, td)
		.SetStyles(Object('overflow': 'auto', 'user-select': 'text'), td)
		.SetStyles(Object(
			'position': 'absolute',
			'top': '0px',
			'left': '0px'), c.El)
		DoStartup(c)

		.expandedCtrls[rowIdx] = [ctrl: c, :expandRow, recordRow: .rows[rowIdx], :rows]
		.WindowRefresh()
		}

	createEditButton(markCell)
		{
		if .showEditButton? is false
			return
		edit = CreateElement('div', markCell, className: 'su-vlist-edit-button')
		edit.innerText = IconFontHelper.GetCode('edit').Chr()
		edit.AddEventListener('click', .editClicked)
		}

	editClicked(event)
		{
		row = .getRow(event.target)
		.Parent.Event('ExpandButton_EditClicked', row)
		event.StopPropagation()
		}

	construct(ctrl, td)
		{
		if Number?(ctrl)
			{
			c = .recycledExpands.FindOne({ it.UniqueId is ctrl })
			td.AppendChild(c.El)
			return c
			}

		.TargetEl = td
		c = .Construct(ctrl)
		.Delete(#TargetEl)
		if c.Xstretch > 0
			c.SetStyles(#(width: '100%'))
		return c
		}

	VirtualListExpand_Recycle(uniqueId)
		{
		ctrl = .VirtualListExpand_Destroy(uniqueId)
		.recycledExpands.Add(ctrl)
		}

	VirtualListExpand_Destroy(uniqueId)
		{
		if false is i = .expandedCtrls.FindIf({ it.ctrl.UniqueId is uniqueId })
			{
			throw "cannot find the expanded control to recycle: " $ Display(uniqueId)
			}
		ctrl = .expandedCtrls[i].ctrl
		.expandedCtrls[i].ctrl.El.Remove()
		.expandedCtrls[i].expandRow.Remove()
		.expandedCtrls[i].recordRow.RemoveAttribute('data-expanded')
		.expandedCtrls.Erase(i)
		return ctrl
		}

	updateExpandRow(oldIdx, newIdx)
		{
		if not .expandedCtrls.Member?(oldIdx)
			return

		.expandedCtrls[newIdx] = .expandedCtrls[oldIdx]
		.expandedCtrls.Erase(oldIdx)
		.SetCellAttributes(.expandedCtrls[newIdx].expandRow, y: newIdx)
		}

	VirtualListExpandEditPushed(updates)
		{
		for row in updates.Members()
			if updates[row]
				.expandedCtrls[row].expandRow.SetAttribute('data-editing', '')
			else
				.expandedCtrls[row].expandRow.RemoveAttribute('data-editing')
		}

	showEditButton?: true
	SetShowEditButton?(.showEditButton?) {}

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
		}

	Recalc()
		{
		for ctrl in .expandedCtrls
			{
			if ctrl.rows isnt false
				ctrl.ctrl.BottomUp(#OverrideHeight, (.GetRowHeight() * ctrl.rows ))
			ctrl.ctrl.BottomUp(#Recalc)
			ctrl.ctrl.El.parentElement.SetStyle('height',
				(ctrl.rows isnt false
					? .GetRowHeight() * ctrl.rows
					: ctrl.ctrl.Ymin) $ 'px')
			}
		}

	expandDiv: false
	expandAttachedRow: false
	expandButtons: false
	showSwitch: false
	SetExpandButtons(buttons)
		{
		.expandButtons = Object()
		.expandDiv = CreateElement('div', className: 'su-vlist-expand-buttons')
		.showSwitch = true
		for button in buttons
			{
			el = CreateElement('span', .expandDiv, className: 'su-vlist-expand-button')
			el.textContent = Number?(button[1]) ? button[1].Chr() : button[1][1].Chr()
			el.SetAttribute('data-cmd', button[0])
			el.SetAttribute('translate', 'no')
			el.AddEventListener('click', .expandButtonClicked)
			.expandButtons.Add(Object(:el, :button))
			}
		}

	expandButtonClicked(event)
		{
		if .expandAttachedRow is false
			return
		cmd = event.target.GetAttribute('data-cmd')
		.Parent.EventWithOverlay(cmd, .expandAttachedRow)
		}

	mouseEnterRow(event)
		{
		if .expandDiv is false
			return
		target= event.target
		row = .getRow(target)
		if .data[row].vl_deleted is true
			return
		for buttonOb in .expandButtons
			if not Number?(buttonOb.button[1])
				buttonOb.el.textContent = target.HasAttribute(buttonOb.button[1][0])
					? buttonOb.button[1][1].Chr()
					: buttonOb.button[1][2].Chr()
		AttachElement(.expandDiv, target.firstChild, false)
		.expandAttachedRow = row
		}

	mouseLeaveRow(event/*unused*/)
		{
		if .expandDiv is false
			return
//		.expandDiv.Remove()
//		.expandAttachedRow = false
		}

	dragging: false
	focused: false
	DistanceToShowSwitchBtn: 100
	mousemove(event)
		{
		if .expandDiv is false
			return
		rect = SuRender.GetClientRect(.scrollContainerEl)
		cusorX = event.clientX - rect.left
		show? = cusorX < .DistanceToShowSwitchBtn
		if show? isnt .showSwitch
			{
			for item in .expandButtons
				{
				alwaysDisplay? = item.button.GetDefault(#alwaysDisplay?, false)
				item.el.SetStyle('display', alwaysDisplay? or show? ? '' : 'none')
				}
			.showSwitch = show?
			}
		.moveRowAround(event)
		}

	moveRowAround(event)
		{
		if .dragging is false
			return

		focusedPos = .getRowPos(.focused)
		y = event.clientY
		if y < focusedPos.top // going up
			{
			if false isnt previousPos = .getRowPos(.focused - 1)
				{
				if previousPos.top <= y and y <= (previousPos.top + .GetRowHeight())
					{
					// .swapRows doesn't handle moving the expand ctrl
					// the from and to args are assigned delicately to avoid the need of moving expand ctrls
					.swapRows(.focused - 1, .focused)
					.updateExpandRow(.focused - 1, .focused)
					.focused = .focused - 1
					}
				}
			}
		else if y > focusedPos.bottom // going down
			{
			if false isnt nextPos = .getRowPos(.focused + 1)
				{
				if nextPos.bottom >= y and y >= (nextPos.bottom - .GetRowHeight())
					{
					// .swapRows doesn't handle moving the expand ctrl
					// the from and to args are assigned delicately to avoid the need of moving expand ctrls
					.swapRows(.focused + 1, .focused)
					.updateExpandRow(.focused + 1, .focused)
					.focused = .focused + 1
					}
				}
			}
		}

	getRowPos(index)
		{
		if not .rows.Member?(index)
			return false

		row = .rows[index]
		rect = SuRender.GetClientRect(row)
		top = rect.top
		bottom = rect.bottom
		if .expandedCtrls.Member?(index)
			{
			expandRect = SuRender.GetClientRect(.expandedCtrls[index].expandRow)
			bottom = expandRect.bottom
			}
		return Object(:top, :bottom)
		}

	getNumRows()
		{
		return .data.Size()
		}

	// this method doesn't handle expand controls yet
	swapRows(from, to)
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

	VirtualList_AllowDragging(focused, mouseEventId)
		{
		if mouseEventId is .mouseEventId and .mousedown? is true
			{
			.dragging = true
			.origFocused = .focused = focused
			.tbody.classList.Add('su-vlist-dragging')
			}
		}

	Destroy()
		{
		.cancelDelayedLoads()
		.observer.Disconnect()
		super.Destroy()
		}
	}
