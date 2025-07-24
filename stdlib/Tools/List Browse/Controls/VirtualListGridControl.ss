// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
// TODO: move focusedRow into selection
WndProc
	{
	Ystretch:		1
	Xstretch:		1
	Name: 			'VirtualListGrid'
	model: 			false
	colModel: 		false
	formatting:		false
	readOnly:		false

	New()
		{
		.CreateWindow("SuWhiteArrow", "",
			WS.TABSTOP | WS.HSCROLL | WS.VISIBLE | WS.CLIPSIBLINGS, WS_EX.CONTROLPARENT)
		.SubClass()

		.SetFont(size: "", text: "W...")
		.rowHeight = .Ymin += 4 /*= row margin*/
		.focusedRow = false

		.painter = VirtualListGridPaint(this, .rowHeight)
		.tooltip = VirtualListGridTooltip(this)
		}

	SetModel(model)
		{
		.model = model
		.colModel = model.ColModel
		model.AutoSave? = .Controller.Send('VirtualList_AutoSave?') isnt false
		.clearSelects()
		.selection = model.InitSelection()
		.painter.SetModel(.model)
		}

	Getter_RowHeight()
		{
		return .rowHeight
		}

	GetChildren()
		{
		if .model is false or .model.ExpandModel is false
			return #()
		return .model.ExpandModel.GetControls()
		}

	Resize(x, y, w, h)
		{
		super.Resize(x, y, w, h)
		.recalcGrid()
		if .model.ExpandModel isnt false
			.updateExpand()
		}

	recalcGrid()
		{
		VirtualListGridScroll.UpdateHorzScrollPos(this, .colModel)
		// in case horz scroll bar changed the rows
		if .model isnt false
			.model.UpdateVisibleRows(.GetRows(.GetClientRect().GetHeight()))
		}

	updateExpand()
		{
		viewChanged = false
		.model.ExpandModel.UpdateExpand(.Hwnd, .rowHeight)
			{ |newRows, oldRows, y|
			viewChanged = true
			rowIndex = .GetRows(y) - 1
			.model.SetRecordCollapsed(rowIndex, keepPosition?:)
			.model.SetRecordExpanded(rowIndex, newRows)
			.scroll(0, (newRows - oldRows) * .rowHeight, rowIndex + oldRows + 1)
			// won't handle if non-focused row changed at same time
			if .focusedRow isnt false and .model.Offset < 0
				.focusedRow -= newRows - oldRows
			.selection.DecreaseShiftStart(newRows - oldRows)
			if newRows < oldRows
				for (i = 0; i < oldRows - newRows; ++i)
					.repaintRow(rowIndex + .model.Offset + oldRows + i)
			}
		if viewChanged
			.Send('VirtualListGrid_ViewChanged')
		}

	scroll(dx, dy, fromRowIndex = false)
		{
		VirtualListGridScroll(this, .model, .rowHeight, dx, dy, fromRowIndex)
		}

	// called by VirtualListViewControl
	ScrollToLeft()
		{
		.HSCROLL(SB.LEFT)
		}

	ScrollColToView(col)
		{
		colRect = .getColRect(col)
		offset = Max(0, colRect.GetX() + colRect.GetWidth() - .GetClientRect().GetWidth())
		offset = Min(offset, colRect.GetX())
		if (offset isnt 0)
			.HSCROLL(MAKELONG(SB.THUMBTRACK, (offset + .colModel.Offset).Int()))
		}

	// called by VirtualListViewControl
	ScrollClientRect(col, width, movePix, resizeLeftSide? = false)
		{
		oldRect = .getColRect(col)

		.colModel.SetColWidth(col, width)
		newRect = .getColRect(col)

		if (movePix < 0)
			oldRect = newRect

		left = Max(0, oldRect.GetX() + oldRect.GetWidth())
		if resizeLeftSide? // when the view port is at right end
			VirtualListGridScroll(this, .model, .rowHeight, movePix, 0, fromRight: left)
		else
			{
			VirtualListGridScroll(this, .model, .rowHeight, movePix, 0, fromLeft: left)
			newRect.Set(y: 0)
			InvalidateRect(.Hwnd, newRect.ToWindowsRect(), false)
			}

		.recalcGrid()

		.Send('VirtualListGrid_ViewChanged')
		}

	getColRect(col)
		{
		for (x = 0, c = col; c > 0;)
			x += .colModel.GetColWidth(--c)
		return Rect(x - .colModel.Offset, 0, .colModel.GetColWidth(col),
			.GetClientRect().GetHeight())
		}

	VertScroll(dy)
		{
		.tooltip.Activate(false)
		.vertScroll(dy, notify?:)
		}

	vertScroll(dy, notify? = false) // dy > 0 : down; dy < 0 : down
		{
		if dy is 0
			return

		dy = .model.UpdateOffset(dy, .saveAndCollapseRelease)
		if dy is 0
			return

		diff = -dy * .rowHeight
		clientRect = .GetClientRect()
		clientRect.Set(height: .rowHeight * .model.VisibleRows)
		.scroll(0, diff)

		.Send('VirtualListGrid_ViewChanged', ignoreThumb: notify? is false)
		}

	saveAndCollapseRelease(rec, row_num)
		{
		.SetFocus()
		freshRec = false
		if .model.EditModel.RecordLocked?(rec) or .model.EditModel.RecordChanged?(rec)
			{
			if false is freshRec = .Send('VirtualListGrid_SaveRecord', rec)
				{
				.SelectRecord(rec)
				.AlertInfo('Save',
					'The information on the current record is invalid. ' $
					'Please correct it first.')
				return false
				}
			}

		targetRec = freshRec is false ? rec : freshRec
		if rec.vl_expanded_rows isnt ''
			{
			.setRecordCollapsed(row_num)
			.model.ExpandModel.Collapse(targetRec)
			}
		if .model.Selection.HasSelectedRow?(targetRec)
			{
			.model.Selection.ClearSelect(false)
			.model.Selection.ClearShiftStart()
			}
		return true
		}

	vsSpeed: 		0
	notify?:		false
	scrollingTimer: false
	VertScrolling(dy, notify? = false)
		{
		.vsSpeed = dy
		.notify? = notify?
		if .vsSpeed is 0
			{
			.StopVScrolling()
			return
			}

		if .scrollingTimer is false and .vsSpeed isnt 0
			.scrollingTimer = .Defer(.scrolling, uniqueID: 'vl_scrolling')
		}

	scrolling()
		{
		if .vsSpeed isnt 0
			{
			.vertScroll(.vsSpeed, .notify?)
			.scrollingTimer = .Delay(100, /*= delay to keep scrolling */
				.scrolling, uniqueID: 'vl_scrolling')
			}
		}

	StopVScrolling()
		{
		.vsSpeed = 0
		.notify? = false
		if .scrollingTimer isnt false
			.scrollingTimer.Kill()
		.scrollingTimer = false
		}

	HSCROLL(wParam)
		{
		if .model is false
			return 0
		return VirtualListGridScroll.HSCROLL(this, .model, .colModel, wParam, .rowHeight)
		}

	InsertRow(record = false, pos = 'current', force = false) // pos: current, start, end
		{
		if false is .okayToInsert?(force)
			return false
		if pos is 'current' and .focusedRow is false
			pos = 'end'
		rowIndex = .insertRowIndex(pos)
		if not .validRowIndex?(rowIndex)
			return false
		if false is newRec = .model.InsertNewRecord(record, rowIndex, :force)
			return false

		.Send('VirtualListGrid_NewRowAdded', newRec)
		if not force
			{
			.selectRowAfterInsert(pos, rowIndex)
			.EditField(.GetSelectedRecord(), .model.ColModel.Get(0))
			}
		else
			.Repaint()
		return true
		}

	okayToInsert?(force)
		{
		if not .editable?() and not force
			return false
		if false is .Controller.Send('VirtualList_AllowNewRecord')
			return false
		if not .Controller.SaveOutstandingChanges()
			return false
		return true
		}

	insertRowIndex(pos)
		{
		if pos is 'end'
			return false
		if pos is 'start'
			return 0
		return Number?(pos)	? pos : .focusedRow - .model.Offset // current
		}

	validRowIndex?(rowIndex)
		{
		if Number?(rowIndex) and rowIndex < 0 and .model.GetStartLast() is false
			{
			SuneidoLog('ERROR: (CAUGHT) row_num should not be negative ' $
				'when inserting row', calls:, params: [:rowIndex, focus: .focusedRow,
					offset: .model.Offset, size: .model.GetLoadedData().Size()])
			return false
			}
		return true
		}

	selectRowAfterInsert(pos, rowIndex)
		{
		if pos isnt 'end'
			{
			.repaintSpecifiedVisibleRows(.clearSelected?)
			.clearSelects()
			.scroll(0, 1 * .rowHeight, rowIndex)
			.selectRow(rowIndex + .model.Offset)
			}
		else
			.selectRow(.model.GetLastVisibleRowIndex())
		}

	clearSelected?(rec)
		{
		if not .selection.HasSelectedRow?(rec)
			return false
		.selection.ClearSelect(rec)
		return true
		}

	DeleteRow(rec, rowNum)
		{
		if rec isnt false
			.selection.ClearSelect(rec)
		.scroll(0, -1 * .rowHeight, rowNum)
		.model.DeleteRecord(rec)
		.focusedRow = .selection.AdjustFocusedRow(.focusedRow, rowNum)
		.Send('VirtualListGrid_RowDeleted')
		}

	SelectFocusedRow()
		{
		.focusedRow = .model.ValidateRow(.focusedRow, returnBoundary?:)
		.selectRow(.focusedRow, moveDown:)
		}

	toggleExpandWithHotkeys(ctrl, expand)
		{
		if .model isnt false and
			ctrl and .model.ExpandModel isnt false and .focusedRow isnt false
			.ToggleExpand(.focusedRow - .model.Offset, :expand)
		}

	ToggleExpand(rowIndex, expand, keepPos? = false)
		{
		if false is rec = .model.GetRecord(rowIndex)
			return
		if rec.vl_expand? is true
			return
		ctrl = false
		if expand
			{
			if rec.vl_expanded_rows isnt ''
				return
			if 0 is layoutOb = .Send('VirtualListGrid_Expand', rec)
				layoutOb = .model.ExpandModel.CustomizableLayout()
			_expandRec = rec
			.model.ExpandModel.ConstructAt(layoutOb, rowIndex, this, .model, .rowHeight)
			readOnly? = .Controller.Send('VirtualList_ReadOnly?')
			if not Boolean?(readOnly?)
				readOnly? = not .model.EditModel.RecordLocked?(rec)
			.model.ExpandModel.Expand(rec, layoutOb, .model, :readOnly?)
			ctrl = layoutOb.ctrl
			.setRecordExpanded(rowIndex, layoutOb.rows, keepPos?)
			.ScrollToLeft()
			}
		else
			{
			if rec.vl_expanded_rows is ''
				return
			.setRecordCollapsed(rowIndex)
			.model.ExpandModel.Collapse(rec)
			.CommitRecord(rec)
			}
		.Send('VirtualListGrid_AfterExpand', :rec, :ctrl, :expand)
		}

	CommitRecord(rec, highlighErr? = false)
		{
		if not .model.EditModel.RecordChanged?(rec)
			{
			.model.UnlockRecord(rec)
			return true
			}
		if false isnt .Send('VirtualListGrid_SaveRecord', rec, :highlighErr?)
			return true
		return false
		}

	setRecordExpanded(rowIndex, rows, keepPos?)
		{
		.model.SetRecordExpanded(rowIndex, rows)
		.moveExpandControls(rowIndex, rows, expand?:)
		absRowIndex = .updateForExpand(rowIndex, rows)
		rows = Min(.model.VisibleRows-1, rows)
		if not keepPos?
			.scrollRowToView(absRowIndex + rows)
		}

	setRecordCollapsed(rowIndex)
		{
		rec = .model.GetRecord(rowIndex)
		rows = rec.vl_expanded_rows // get current expanded rows first
		.model.SetRecordCollapsed(rowIndex)
		.moveExpandControls(rowIndex, rows)
		.updateForExpand(.updateCollapsedRowIndex(rowIndex, rows), -rows)
		}

	moveExpandControls(rowIndex, rows, expand? = false)
		{
		if expand? is true
			.scroll(0, rows * .rowHeight, rowIndex + 1)
		else
			{
			topIndex = Min(rowIndex, 0)
			bottomIndex = Max(rowIndex + rows, 0)
			// move up all expand ctrls that are underneath the collpased row
			// if the collpased row has parts under the view port
			if bottomIndex isnt 0
				.scroll(0, -Min(bottomIndex, rows) * .rowHeight, bottomIndex + 1)
			// move down all expand ctrls that are above the collpased row
			// if the collpased row has parts above the view port
			// and make sure the first visible record stays
			if topIndex isnt 0
				VirtualListGridScroll.MoveExpandControls(this, .model, .rowHeight,
					topIndex - 1, Min(-topIndex, rows))
			}
		}

	updateCollapsedRowIndex(rowIndex, rows)
		{
		if rowIndex >= 0
			return rowIndex
		return Min(rowIndex + rows, 0)
		}

	updateForExpand(rowIndex, rows)
		{
		absRowIndex = rowIndex + .model.Offset
		if .model.Offset < 0
			{
			if .focusedRow isnt false and .focusedRow <= absRowIndex + rows
				.focusedRow -= rows
			}
		else
			{
			if .focusedRow isnt false and .focusedRow > absRowIndex
				.focusedRow += rows
			}
		.model.LogInvalidFocus(.focusedRow)
		.selection.UpdateShiftStart(absRowIndex, rows)
		.repaintExpandedRows(absRowIndex, rows.Abs())
		return absRowIndex
		}

	repaintExpandedRows(rowIndex, rows)
		{
		for(i=1; i<=rows; i++)
			.repaintRow(rowIndex + i)
		if .focusedRow isnt false
			.repaintRow(.focusedRow)
		}

	GETDLGCODE(@unused)
		{
		return DLGC.WANTCHARS | DLGC.WANTARROWS | DLGC.WANTALLKEYS
		}

	//TODO: handle mouse "One screen at a time" option
	MOUSEWHEEL(wParam)
		{
		if .model is false
			return 0
		scroll = GetWheelScrollInfo(wParam)
		prevOffset = .model.Offset
		// limit the reading lines to be within [50, -50] to avoid slow query
		lines = Max(Min(-scroll.lines, 50), -50/*=scroll limit*/)
		.VertScroll(lines)
		if prevOffset is .model.Offset
			.Send('VirtualListGrid_MouseWheel', wParam)
		return 0
		}

	MOUSEMOVE(lParam)
		{
		if .model is false
			return 0

		TrackMouseEvent(Object(cbSize: TRACKMOUSEEVENT.Size(),
			dwFlags: TME.LEAVE hwndTrack: .Hwnd))

		x = LOSWORD(lParam)
		y = HISWORD(lParam)
		row_num = .GetRows(y)
		col = .model.ColModel.GetColByX(x)
		.tooltip.UpdateToolTip(row_num, col, .model)
		.ShowExpandButton(row_num)

		if not .dragging
			return 0
		data = .model.GetLoadedData()
		if false is newRow = .getNewRowToSwap(row_num, data)
			return 0
		// ensure Swap() and Send() is 1 row a time
		for (inc = newRow > .focusedRow ? 1 : -1;
			.focusedRow isnt newRow; .focusedRow += inc)
			{
			movedRec = .model.GetRecord(.focusedRow + inc - .model.Offset)
			if .model.ExpandModel isnt false and movedRec.vl_expanded_rows isnt ''
				{
				expand = .model.ExpandModel.GetExpandedControl(movedRec)
				expand.ctrl.Resize(
					-.model.ColModel.Offset + .rowHeight,
					(.focusedRow + 1 - .model.Offset) * .rowHeight,
					20000, /*= max width */
					expand.rows * .rowHeight)
				}
			data.Swap(.focusedRow, .focusedRow + inc)
			.repaintRow(.focusedRow)
			.Controller.Send("VirtualList_Move")
			}
		.model.LogInvalidFocus(.focusedRow)
		.repaintRow(.focusedRow)
		return 0
		}

	getNewRowToSwap(row_num, data)
		{
		newRow = Max(0, row_num + .model.Offset)
		newRow = Min(newRow, .model.Offset + .model.VisibleRows - 1)
		newRow = Min(newRow, data.Size() - 1)
		newRow = Max(newRow, .model.Offset)
		newRowRec = .model.GetRecord(newRow - .model.Offset)
		SetCursor(LoadCursor(ResourceModule(), IDC.DRAG1))
		if newRow < .focusedRow and newRowRec.vl_expand? is true
			return false
		if newRow > .focusedRow and .mouseOnExpand?(newRowRec)
			return false
		if .focusedRow is false or .focusedRow is newRow or newRow < 0
			return false
		return newRow
		}

	mouseOnExpand?(newRowRec)
		{
		return newRowRec.vl_expanded_rows isnt '' or
			newRowRec.vl_expand? is true and
				newRowRec.vl_expand_index isnt newRowRec.vl_rows - 1
		}

	ShowExpandButton(row_num = false)
		{
		GetCursorPos(pt = Object())
		if row_num is false
			row_num = .GetRows(pt.y)
		ScreenToClient(.Hwnd, pt)
		curLeft = pt.x / .GetClientRect().GetWidth()
		rec = .model.GetRecord(row_num)
		invalid = rec is false or rec.vl_expand? is true
		expanded = rec isnt false and .model.ExpandModel isnt false and
			.model.ExpandModel.GetExpandedControl(rec) isnt false
		.Send('VirtualListGrid_MouseMove', row_num, expanded, invalid, :curLeft)
		}

	TTN_SHOW(lParam)
		{
		.tooltip.TTN_SHOW(lParam, .getCellRect)
		}

	getCellRect(row_num, col)
		{
		rect = .getRowRect(row_num)
		rect.Set(x: .model.ColModel.GetColumnOffset(col),
			width: .model.ColModel.GetColumnWidth(col))
		return rect.ToWindowsRect()
		}

	MOUSELEAVE()
		{
		.tooltip.ResetLast()
		.Send('VirtualListGrid_MouseLeave')
		return 0
		}

	KEYDOWN(wParam)
		{
		if .model is false or .dragging
			return 0
		ctrl = KeyPressed?(VK.CONTROL)
		shift = KeyPressed?(VK.SHIFT)
		if .keydown_fns.Member?(wParam)
			(.keydown_fns[wParam])(:ctrl, :shift)
		return 0
		}

	getter_keydown_fns()
		{
		fns = Object()
		fns[VK.LEFT] = 		{|ctrl|		.HSCROLL(ctrl ? SB.LEFT : SB.LINELEFT) 			}
		fns[VK.RIGHT] = 	{|ctrl|		.HSCROLL(ctrl ? SB.RIGHT : SB.LINERIGHT)		}
		fns[VK.UP] = 		{|shift|	.selectRow(.focusedRow - 1, :shift)				}
		fns[VK.DOWN] = 		{|shift|	.selectRow(.focusedRow + 1, :shift, moveDown:)	}
		fns[VK.PRIOR] = 	{|shift| 	.selection.PageKey(
											.focusedRow, shift, .selectRow, up?:)		}
		fns[VK.NEXT] = 		{|shift|	.selection.PageKey(
											.focusedRow, shift, .selectRow)				}
		fns[VK.HOME] = 		{			.setStartLast(false)							}
		fns[VK.END] = 		{			.setStartLast(true)								}
		fns[VK.RETURN] = 	{			.Send('VirtualListGrid_Return')					}
		fns[VK.ESCAPE] = 	{			.Send('VirtualListGrid_Escape')					}
		fns[VK.TAB] = 		{			.tabThrough()									}
		fns[VK.F5] = 		{			.Repaint()										}
		fns[VK.SPACE] = 	{			.Send('VirtualListGrid_Space')					}
		fns[VK.ADD]	=		{|ctrl| 	.toggleExpandWithHotkeys(ctrl, expand:) 		}
		fns[VK.SUBTRACT] = 	{|ctrl|		.toggleExpandWithHotkeys(ctrl, expand: false)	}
		fns[VK.OEM_PLUS] = 	{|ctrl|		.toggleExpandWithHotkeys(ctrl, expand:)			}
		fns[VK.OEM_MINUS] = {|ctrl|		.toggleExpandWithHotkeys(ctrl, expand: false)	}
		fns[VK.DELETE]	= 	{			.deleteFromKeyboard()							}
		fns[VK.INSERT]	= 	{			.insertFromKeyboard()							}
		fns[VK.F8] =
			{
			if Suneido.User is 'default'
				Inspect(this)
			}
		return .keydown_fns = fns // once only
		}

	setStartLast(startLast)
		{
		if false is .Send('VirtualListGrid_SetStartLast')
			return
		moveDown = startLast
		.clearSelects()
		if .model.AllRead?
			{
			if .model.GetStartLast()
				.selectRow(startLast ? -1 : -.model.GetLoadedData().Size(), :moveDown)
			else
				.selectRow(startLast ? .model.GetLoadedData().Size() - 1 : 0, :moveDown)
			return
			}

		if .model.SetStartLast(startLast)
			{
			.Repaint()
			.Send('VirtualListGrid_ViewChanged')
			}

		.selectRow(startLast ? -1 : 0, :moveDown)
		}

	deleteFromKeyboard()
		{
		if 0 isnt .Controller.Send('On_Context_DeleteUndelete', rec: .GetSelectedRecord())
			return

		.Send('Keyboard_Delete', rec: .GetSelectedRecord())
		}

	insertFromKeyboard()
		{
		if .focusedRow isnt false
			.scrollRowToView(.focusedRow)
		.Send('VirtualListGrid_Insert')
		}

	tabThrough()
		{
		if .model.ExpandModel is false or .focusedRow is false
			return .Send('VirtualListGrid_Tab')

		rowIndex = .focusedRow
		while (false isnt rec = .model.GetRecord(rowIndex - .model.Offset))
			{
			if rec.vl_expanded_rows isnt ''
				{
				.FocusFirst(.model.ExpandModel.GetExpandedControl(rec).ctrl.Hwnd)
				return 0
				}
			rowIndex++
			}
		hwnd = GetNextDlgTabItem(.Hwnd, NULL, false)
		SetFocus(hwnd)
		return 0
		}

	SETFOCUS(wParam)
		{
		.RepaintSelectedRows()
		if (.edit isnt false)
			{
			editorHwnd = .edit.GetEditorHwnd()
			if not .edit.ClosingListEdit() and
				editorHwnd isnt 0 and not .edit.ChildOf?(wParam)
				{
				SetActiveWindow(editorHwnd) // put the editor on top
				return 0
				}
			.edit.Return()						// or end it
			}
		return 0
		}
	KILLFOCUS(wParam)
		{
		.RepaintSelectedRows()
		.endDrag()
		if .saveOnLeaving?(wParam)
			.rowChanged(false) // save
		return 'callsuper'
		}
	// NOTE: we don't auto save when focus leaves
	// 		 from expand controls to anywhere outside virtual list
	// 		 so it behaves same as AccessControl with separate edit button
	saveOnLeaving?(curFocus)
		{
		if .edit isnt false and .edit.Editing is true
			return false
		// killfocus is still triggered when the window is re-activated
		// even though the focus is set to itself, like finishing cell editing
		if curFocus is .Hwnd
			return false
		if .model isnt false and .model.ExpandModel isnt false and
			false isnt .model.ExpandModel.GetCurrentFocusedRecord(curFocus)
			return false

		return not .popUpFromChildWindow?(curFocus)
		}

	popUpFromChildWindow?(curFocus)
		{
		p = curFocus
		parentWindow = 0
		while 0 isnt p = GetParent(p)
			{
			parentWindow = p
			}
		return parentWindow isnt .Window.Hwnd and
			GetWindow(parentWindow, GW.OWNER) is .Window.Hwnd
		}
	PAINT()
		{
		return .painter.PAINT()
		}

	HighlightValues(member, values, color)
		{
		.painter.HighlightValues(member, values, color)
		.Repaint()
		}

	HighlightRecords(recs, color)
		{
		.painter.HighlightRecords(recs, color)
		.Repaint()
		}

	ClearHighlightRecord(rec)
		{
		.painter.ClearHighlightRecord(rec)
		}

	ClearHighlight()
		{
		.painter.ClearHighlight()
		.Repaint()
		}

	Repaint(keepPos? /*unused*/ = false)
		{
		.recalcGrid()
		rc = .GetClientRect().ToWindowsRect()
		InvalidateRect(.Hwnd, rc, false)
		}

	clearSelects()
		{
		if .focusedRow is false
			return
		oldRow = .focusedRow
		.focusedRow = false
		.repaintRow(oldRow)
		}

	GetSelectedRecord()
		{
		if .focusedRow isnt false
			return .model.GetRecord(.focusedRow - .model.Offset)
		recs = .GetSelectedRecords()
		if not recs.Empty?()
			return recs[0]
		return false
		}

	GetSelectedRecords()
		{
		if .selection is false
			return #()
		return .selection.GetSelectedRecords()
		}

	SelectTopRecords(field, values)
		{
		row = 0
		while(false isnt data = .model.GetRecord(row))
			{
			if not values.Has?(data[field])
				break
			.selection.SelectRows(ctrl: true, shift: false, :row)
			++row
			}
		}

	selection: false
	ClearSelect(rec = false)
		{
		if .selection is false
			return
		needToRepaint = rec isnt false or .focusedRow isnt false or .selection.NotEmpty?()
		.focusedRow = false
		.selection.ClearSelect(rec)
		if needToRepaint
			.Repaint()
		}

	SetFocusedRow(.focusedRow)
		{
		}

	SelectRow(row)
		{
		.selectRow(row)
		}

	selectRow(row, ctrl = false, shift = false, moveDown = false)
		{
		if .selection.NotEmpty?()
			.RepaintSelectedRows()

		.scrollRowToView(row)

		if false is .model.ValidateRow(row)
			return false

		rec = .model.GetRecord(row - .model.Offset)
		if rec isnt false and rec.vl_expand? is true
			return .selectRow(row + (moveDown ? 1 : -1), :ctrl, :shift, :moveDown)

		if .model.LogInvalidFocus(row)
			return false

		try
			.selection.SelectRows(ctrl, shift, row)
		catch (err, 'Cannot select more than')
			.AlertWarn(.Title, err)
		.rowChanged(row)

		.focusedRow = row
		.RepaintSelectedRows()
		.Send('VirtualListGrid_ItemSelected', rec)
		return true
		}
	rowChanged(row)
		{
		if .focusedRow is false or .focusedRow is row or not .model.AutoSave?
			return
		if false is prevRec = .model.GetRecord(.focusedRow - .model.Offset)
			return
		.CommitRecord(prevRec)
		}

	SelectNextRow(rowChange)
		{
		return .selectRow(.focusedRow + rowChange)
		}

	SelectRecord(rec)
		{
		.selectRow(.model.GetRecordRowNum(rec) + .model.Offset)
		}

	RepaintSelectedRows()
		{
		if .selection is false
			return
		.repaintSpecifiedVisibleRows(.selection.HasSelectedRow?)
		}

	RepaintRecord(rec)
		{
		.repaintSpecifiedVisibleRows({ it is rec })
		}

	repaintSpecifiedVisibleRows(block)
		{
		for(i = 0; i < .model.VisibleRows; ++i)
			if block(.model.GetRecord(i))
				.repaintRow(i + .model.Offset)
		}

	scrollRowToView(row)
		{
		offset = Max(Min(0, row - .model.Offset),
			row - .model.Offset - .model.VisibleRows + 1)
		.VertScroll(offset)
		}

	repaintRow(row_num)
		{
		if .model is false or row_num is false
			return

		row_num = row_num - .model.Offset
		clientRect = .GetClientRect()
		clientRect.Set(y: 0, height: clientRect.GetHeight())
		if (clientRect.Overlaps?(rect = .getRowRect(row_num)))
			InvalidateRect(.Hwnd, rect.ToWindowsRect(), false)
		.Send('VirtualListGrid_RepaintingRow', row_num)
		}

	getRowRect(row_num)
		{
		return .GetClientRect().Set(y: row_num * .rowHeight, height: .rowHeight)
		}

	edit: false
	LBUTTONDBLCLK(lParam)
		{
		if .model is false
			return 0

		x = LOSWORD(lParam)
		y = HISWORD(lParam)
		row_num = .selectRowByClick(y)
		if false is col = .model.ColModel.GetColByX(x)
			return 0
		if false is rec = .model.GetRecord(row_num)
			.InsertRow(pos: 'end')
		else
			.EditField(rec, col)
		rec = .GetSelectedRecord()
		.Send('VirtualListGrid_DoubleClick', rec, col)
		return 0
		}

	editable?()
		{
		return .readOnly is false and .model.EditModel.Editable?()
		}

	EditField(rec, col)
		{
		if not .editable?() or rec.GetDefault('vl_deleted', false) is true
			return
		SetFocus(.Hwnd)

		.edit = VirtualListEdit(this, .model)
		colIndex = .model.ColModel.FindCol(col)
		if rec isnt false and not .readOnly
			{
			.ScrollColToView(colIndex)
			screenCellRect = .GetCurrentRowCellRect(col)
			if false isnt rec = .edit.EditCell(rec, col, screenCellRect)
				.RepaintRecord(rec)
			}
		.Send('VirtualListGrid_Edit', :rec, :col)
		}

	GetCurrentRowCellRect(col)
		{
		row_num = .focusedRow - .model.Offset
		cellRect = .getCellRect(row_num, col)
		width = cellRect.right - cellRect.left
		height = cellRect.bottom - cellRect.top

		ClientToScreen(.Hwnd, p = [x: cellRect.left, y: cellRect.top])
		cellRect.top = p.y
		cellRect.left = p.x
		screenCellRect = Object(
			left: p.x, top: p.y, right: p.x + width, bottom: p.y + height)
		return screenCellRect
		}

	dragging: false
	LBUTTONDOWN(lParam)
		{
		if .model is false
			return 0

		SetFocus(.Hwnd)
		ctrl = KeyPressed?(VK.CONTROL)
		shift = KeyPressed?(VK.SHIFT)
		rowIndex = .GetRows(HISWORD(lParam))
		selected = .selectRow(rowIndex + .model.Offset, ctrl, shift)
		if selected and .draggable?(shift, ctrl)
			{
			.dragging = true
			SetCapture(.Hwnd)
			}
		return 0
		}

	draggable?(shift, ctrl)
		{
		return not .readOnly and not shift and not ctrl and
			not .currentRowExpanded?() and
			.Controller.Send("VirtualList_AllowMove", rec: .GetSelectedRecord()) is true
		}

	currentRowExpanded?()
		{
		return .model.ExpandModel isnt false and
			.GetSelectedRecord().vl_expanded_rows isnt ''
		}

	LBUTTONUP(lParam)
		{
		.endDrag()
		rowIndex = .GetRows(HISWORD(lParam))
		col = .model.ColModel.GetColByX(LOSWORD(lParam))
		if rowIndex + .model.Offset is .focusedRow
			.Send('VirtualListGrid_LeftClick', .model.GetRecord(rowIndex), col)
		return 0
		}

	endDrag()
		{
		if .dragging
			{
			ReleaseCapture()
			.dragging = false
			}
		}

	GetRows(height)
		{
		return (height / .rowHeight).Floor()
		}

	CONTEXTMENU(lParam)
		{
		if .model is false or .dragging
			return 0

		SetFocus(.Hwnd)
		x = LOSWORD(lParam)
		y = HISWORD(lParam)
		if x is -1 and y is -1 // from keyboard
			{
			col = .model.ColModel.GetColByX(0)
			ClientToScreen(.Hwnd, pt = Object())
			.Send('VirtualListHeader_ContextMenu', col, pt.x, pt.y)
			return 0
			}
		ScreenToClient(.Hwnd, pt = Object(:x, :y))
		row_num = .GetRows(pt.y)
		rec = .model.GetRecord(row_num)
		if rec isnt false and rec.vl_expand? is true
			return 0
		if not .selection.HasSelectedRow?(rec)
			row_num = .selectRowByClick(pt.y)
		col = .model.ColModel.GetColByX(pt.x)
		.Send('VirtualListGrid_ContextMenu', rec, col, x, y, :row_num)
		return 0
		}

	selectRowByClick(y)
		{
		row_num = .GetRows(y)
		.selectRow(row_num + .model.Offset)
		return row_num
		}

	SetReadOnly(readOnly)
		{
		.readOnly = readOnly
		.painter.SetBackground(.readOnly ? COLOR.BTNFACE : COLOR.WINDOW)
		if .model isnt false and .model.ExpandModel isnt false
			.model.ExpandModel.SetReadOnly(readOnly)
		.Repaint()
		}

	GetReadOnly()
		{
		return .readOnly
		}

	AlwaysHighlightSelected?()
		{
		return .Controller.Send("VirtualList_AlwaysHighlightSelected?") is true
		}

	FinishEdit()
		{
		if .edit isnt false
			.edit.Return()
		.endDrag()
		}

	Destroy()
		{
		.FinishEdit()
		.painter.Destroy()
		.tooltip.Destroy()
		super.Destroy()
		}
	}
