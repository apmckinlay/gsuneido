// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Name: 			'VirtualListGrid'
	ComponentName:	'VirtualListGrid'
	ComponentArgs: 	#()
	rowHeight: 		0 // Fake
	model: 			false
	focusedRow: 	false
	readOnly:		false
	New()
		{
		.Controller = .Controller.Controller // skip VirtualListScrollControl
		.focusedRow = false
		}

	virtualVisibleRows: 10
	SetModel(model)
		{
		.model = model
		.colModel = model.ColModel
		model.AutoSave? = .Controller.Send('VirtualList_AutoSave?') isnt false
		.clearSelects()
		.selection = model.InitSelection()
		.selection.ClearSelect()

		.Act(#ClearData)
		.loadHelper = VirtualListLoadAdapter(this, .model, .paintRow,
			.saveAndCollapseRelease)
		.model.UpdateVisibleRows(.virtualVisibleRows)
		if not .model.EditModel.Editable?() or
			false is .Controller.Send('VirtualList_ShowEditButton?')
			.Act(#SetShowEditButton?, false)

		.init()
		}

	init()
		{
		.Act(#ClearData)
		.loadHelper.InitLoad()
		}

	Repaint(keepPos? = false)
		{
		if .model is false or .Destroyed?()
			return
		if .loadHelper.ModelIsReset?()
			.init()
		else
			.loadHelper.ForLoadedData()
				{ |rowNum, rec|
				.CancelAct(#UpdateData, { it[0] is rowNum })
				.Act(#UpdateData, rowNum, .paintRow(rec), :keepPos?)
				}
		}

	paintRow(rec)
		{
		info = .model.EditModel.GetInvalidInfo(rec)
		rec.PreSet(#list_invalid_row, info.validRule isnt '')
		rec.PreSet(#list_invalid_cells, Object())
		for field in info.invalidCols
			rec.list_invalid_cells[field] = rec[field]

		if '' isnt .model.EditModel.GetWarningMsg(rec)
			.HighlightRecords([rec], CLR.WarnColor, true)

		row = .colModel.GetFormatting().PaintRow(rec, .colModel)
		row.vl_deleted = rec.GetDefault(#vl_deleted, false)
		row.vl_brush = .brushMgr is false ? false : .brushMgr.GetBrush(rec)
		return row
		}

	clearSelects()
		{
		if .focusedRow is false
			return

		.Act(#DeSelectRow, .focusedRow)
		.focusedRow = false
		}

	saveAndCollapseRelease(rec, row_num)
		{
		.SetFocus()
		freshRec = false
		if .model.EditModel.RecordLocked?(rec) or .model.EditModel.RecordChanged?(rec)
			{
			if false is freshRec = .Send('VirtualListGrid_SaveRecord', rec)
				{
				unsavedRec = rec
				if OkCancel(
					'The information on the editting record is invalid.\r\n' $
					'Choose OK button to go back and fix it.\r\n' $
					'Choose Cancel button to stay on the current page.',
					'Save')
					.SelectRecord(unsavedRec)
				return false
				}
			}

		targetRec = freshRec is false ? rec : freshRec
		if rec.vl_expanded_rows isnt ''
			{
			.model.SetRecordCollapsed(row_num)
			.model.ExpandModel.Collapse(targetRec, this)
			}
		if .model.Selection.HasSelectedRow?(targetRec)
			{
			.model.Selection.ClearSelect(false)
			.model.Selection.ClearShiftStart()
			}
		return true
		}

	GetChildren()
		{
		if .model is false or .model.ExpandModel is false
			return #()
		return .model.ExpandModel.GetControls()
		}

	VirtualListGridComponent_Load(row)
		{
		.loadHelper.LoadFrom(row)
		}

	LBUTTONDOWN(row, col/*unused*/, shift, control, mouseEventId = false)
		{
		if .model is false
			return 0

		SetFocus(.Hwnd)
		selected = .selectRow(row, control, shift)
		if selected and .draggable?(shift, control)
			{
			.dragging = true
			.Act('VirtualList_AllowDragging', .focusedRow, mouseEventId)
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

	MoveRow(focused, newRow)
		{
		Assert(focused is: .focusedRow)
		for (inc = newRow > .focusedRow ? 1 : -1; .focusedRow isnt newRow;
			.focusedRow += inc)
			{
			rec = .model.GetLoadedData()
			rec.Swap(.focusedRow, .focusedRow + inc)
			.Controller.Send("VirtualList_Move")
			.repaintRow(.focusedRow, rec[.focusedRow])
			.repaintRow(.focusedRow + inc, rec[.focusedRow + inc])
			}
		.sortCol = false
		}

	selection: false
	selectRow(row, ctrl = false, shift = false, moveDown/*unused*/ = false)
		{
		oldSelections = .getSelections()

		.scrollRowToView(row)

		if false is .model.ValidateRow(row)
			return false

		rec = .model.GetRecord(row - .model.Offset)
//		if rec isnt false and rec.vl_expand? is true
//			return .selectRow(row + (moveDown ? 1 : -1), :ctrl, :shift, :moveDown)

		if .model.LogInvalidFocus(row)
			return false

		try
			.selection.SelectRows(ctrl, shift, row)
		catch (err, 'Cannot select more than')
			.AlertWarn(.Title, err)
		.rowChanged(row)

		.focusedRow = row
		newSelections = .getSelections()
		.handleSelectionChanged(newSelections, oldSelections)
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

	getSelections()
		{
		if .selection is false
			return Object()

		selections = Object()
		data = .model.GetLoadedData()
		for rec in .selection.GetSelectedRecords()
			{
			// m can be false if the record has been deleted
			if false isnt m = data.FindIf({ Same?(rec, it) })
				selections.Add(m)
			}
		return selections
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

	handleSelectionChanged(newSelection, oldSelection)
		{
		if newSelection.Sort!() isnt oldSelection	// if selection changed
			{
			for (row in oldSelection.Difference(newSelection))
				.Act(#DeSelectRow, row)
			for (row in newSelection.Difference(oldSelection))
				.Act(#SelectRow, row)
			}
		}

	SelectTopRecords(field, values)
		{
		oldSelections = .getSelections()
		row = 0
		while(false isnt data = .model.GetRecord(row))
			{
			if not values.Has?(data[field])
				break
			.selection.SelectRows(ctrl: true, shift: false, :row)
			.scrollRowToView(row)
			++row
			}
		newSelections = .getSelections()
		.handleSelectionChanged(newSelections, oldSelections)
		}

	ClearSelect(rec = false)
		{
		if .selection is false
			return
		oldSelections = .getSelections()
		.focusedRow = false
		.selection.ClearSelect(rec)
		newSelections = .getSelections()
		.handleSelectionChanged(newSelections, oldSelections)
		}

	SetFocusedRow(.focusedRow)
		{
		}

	committing: false
	CommitRecord(rec, highlighErr? = false)
		{
		if .committing is true
			return false

		Finally(
			{
			.committing = true
			return .commitRecord(rec, :highlighErr?)
			},
			{
			.committing = false
			})
		}

	commitRecord(rec, highlighErr? = false)
		{
		if not .model.EditModel.RecordChanged?(rec)
			{
			.model.UnlockRecord(rec)
			.RepaintRecord(rec)
			return true
			}
		if false isnt .Send('VirtualListGrid_SaveRecord', rec, :highlighErr?)
			return true

		.RepaintRecord(rec)
		return false
		}

	RepaintSelectedRows()
		{
		if .selection is false
			return

		selections = .getSelections()
		for sel in selections
			.repaintRow(sel, .model.GetRecord(sel - .model.Offset))
		}

	RepaintRecord(rec)
		{
		// i can be false if the rec is a new record and hasn't been inserted
		if false isnt i = .model.GetLoadedData().FindIf({ Same?(rec, it) })
			.repaintRow(i, rec)
		}

	repaintRow(rowNum, rowRec)
		{
		if .model is false or rowNum is false
			return

		.Act(#UpdateData, rowNum, .paintRow(rowRec))
		.Act(.selection.HasSelectedRow?(rowRec) ? #SelectRow : #DeSelectRow, rowNum)
		.Send('VirtualListGrid_RepaintingRow', rowNum)
		}

	LBUTTONUP(row, col)
		{
		if row is .focusedRow
			.Send('VirtualListGrid_LeftClick', .model.GetRecord(row - .model.Offset),
				.model.ColModel.Get(col))
		}

	edit: false
	LBUTTONDBLCLK(row, col)
		{
		if .model is false
			return

		if not .HasFocus?()
			return

		.selectRow(row)
		if false is col = .model.ColModel.Get(col)
			return
		if false is rec = .model.GetRecord(row - .model.Offset)
			.InsertRow(pos: 'end')
		else
			.EditField(rec, col)
		rec = .GetSelectedRecord()
		.Send('VirtualListGrid_DoubleClick', rec, col)
		}

	InsertRow(record = false, pos = 'current', force = false) // pos: current, start, end
		{
		if false is .okayToInsert?(force)
			return false

		if pos is 'current' and .focusedRow is false
			pos = 'end'
		rowIndex = pos is 'end'
			? false
			: pos is 'start'
				? 0
				: Number?(pos)
					? pos
					: .focusedRow - .model.Offset // current

		if false is newRecOb = .loadHelper.InsertNewRecord(record, rowIndex, :force)
			return false

		.Send('VirtualListGrid_NewRowAdded', newRecOb.newRec)
		if not force
			{
			.selectRowAfterInsert(newRecOb.newRowNum)
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

	selectRowAfterInsert(rowNum)
		{
		.clearSelects()
		.selectRow(rowNum)
		}

	DeleteRow(rec, rowNum)
		{
		if rec isnt false
			.selection.ClearSelect(rec)
		.loadHelper.DeleteRecord(rowNum + .model.Offset)
		.model.DeleteRecord(rec)
		.focusedRow = .selection.AdjustFocusedRow(.focusedRow, rowNum)
		.Send('VirtualListGrid_RowDeleted')
		}

	SelectFocusedRow()
		{
		.focusedRow = .model.ValidateRow(.focusedRow, returnBoundary?:)
		.selectRow(.focusedRow, moveDown:)
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
			if false isnt rec = .edit.EditCell(rec, col)
				.RepaintRecord(rec)
			}
		.Send('VirtualListGrid_Edit', :rec, :col)
		}

	editable?()
		{
		return .readOnly is false and .model.EditModel.Editable?()
		}

	SetReadOnly(readOnly)
		{
		.readOnly = readOnly
		if .model isnt false and .model.ExpandModel isnt false
			.model.ExpandModel.SetReadOnly(readOnly)
		.Act("SetReadOnly", .readOnly)
		}

	GetReadOnly()
		{
		return .readOnly
		}

	ScrollColToView(col)
		{
		.Act('ScrollColToView', col)
		}

	scrollRowToView(row)
		{
		.loadHelper.EnsureRow(row)
		.Act('ScrollRowToView', row)
		}

	SelectRow(row)
		{
		.selectRow(row)
		}

	focused: false
	SETFOCUS()
		{
		.focused = true
		if .edit isnt false
			.edit.Return()						// or end it
		.RepaintSelectedRows()
		return 0
		}

	KILLFOCUS(wParam)
		{
		if .focused is false
			return 'callsuper'

		.focused = false
		.RepaintSelectedRows()
		if .saveOnLeaving?(wParam)
			.rowChanged(false) // save
		return 'callsuper'
		}

	// NOTE: we don't auto save when focus leaves
	// 		 from expand controls to anywhere outside virtual list
	// 		 so it behaves same as AccessControl with separate edit button
	saveOnLeaving?(curFocus)
		{
		if .committing is true
			return false
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
		if false is control = SuRenderBackend().GetRegisteredControl(curFocus)
			return false

		return not Same?(control.Window, .Window) and
			control.Window.ComponentName in (#ModalWindow, #Dialog)
		}

	CONTEXTMENU(x, y, row, col)
		{
		if .model is false// or .dragging
			return 0

		SetFocus(.Hwnd)
		row_num = row - .model.Offset
		rec = .model.GetRecord(row_num)
		if rec isnt false and rec.vl_expand? is true
			return 0
		if not .selection.HasSelectedRow?(rec)
			.selectRow(row)
		column = .model.ColModel.Get(col)
		.Send('VirtualListGrid_ContextMenu', rec, column, x, y, :row_num)
		return 0
		}

	CONTEXTMENU_HEADER(x, y, col)
		{
		if .colModel is false
			return 0

		SetFocus(.Hwnd)
		.Send('VirtualListHeader_ContextMenu', .model.ColModel.Get(col), x, y)
		return 0
		}

	KEYDOWN(wParam, ctrl = false, shift = false)
		{
		if .model is false// or .dragging
			return 0

		if .keydown_fns.Member?(wParam)
			(.keydown_fns[wParam])(:ctrl, :shift)
		return 0
		}

	getter_keydown_fns()
		{
		fns = Object()
//		fns[VK.LEFT] = 		{|ctrl|		.HSCROLL(ctrl ? SB.LEFT : SB.LINELEFT) 			}
//		fns[VK.RIGHT] = 	{|ctrl|		.HSCROLL(ctrl ? SB.RIGHT : SB.LINERIGHT)		}
		fns[VK.UP] = 		{|shift|	.selectRow(.focusedRow - 1, :shift)				}
		fns[VK.DOWN] = 		{|shift|	.selectRow(.focusedRow + 1, :shift, moveDown:)	}
//		fns[VK.PRIOR] = 	{|shift| 	.selection.PageKey(
//											.focusedRow, shift, .selectRow, up?:)		}
//		fns[VK.NEXT] = 		{|shift|	.selection.PageKey(
//											.focusedRow, shift, .selectRow)				}
		fns[VK.HOME] = 		{			.setStartLast(false)							}
		fns[VK.END] = 		{			.setStartLast(true)								}
		fns[VK.RETURN] = 	{			.Send('VirtualListGrid_Return')					}
		fns[VK.ESCAPE] = 	{			.Send('VirtualListGrid_Escape')					}
//		fns[VK.TAB] = 		{			.tabThrough()									}
		fns[VK.F5] = 		{			.Repaint()										}
		fns[VK.SPACE] = 	{			.Send('VirtualListGrid_Space')					}
//		fns[VK.ADD]	=		{|ctrl| 	.toggleExpandWithHotkeys(ctrl, expand:) 		}
//		fns[VK.SUBTRACT] = 	{|ctrl|		.toggleExpandWithHotkeys(ctrl, expand: false)	}
//		fns[VK.OEM_PLUS] = 	{|ctrl|		.toggleExpandWithHotkeys(ctrl, expand:)			}
//		fns[VK.OEM_MINUS] = {|ctrl|		.toggleExpandWithHotkeys(ctrl, expand: false)	}
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

		.clearSelects()
		if .model.AllRead?
			{
			if .model.GetStartLast()
				.selectRow(startLast ? -1 : -.model.GetLoadedData().Size())
			else
				.selectRow(startLast ? .model.GetLoadedData().Size() - 1 : 0)
			return
			}

		if .model.SetStartLast(startLast) is true
			.init()

		.selectRow(startLast ? -1 : 0)
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

	Getter_RowHeight()
		{
		return .rowHeight
		}

	ExpandButton_LBUTTONDOWN(row)
		{
		rowIndex = row - .model.Offset
		if false isnt rec = .model.GetRecord(rowIndex)
			.Send('VirtualListThumb_Expand', rowIndex, expand: rec.vl_expanded_rows is '')
		}

	VirtualListExpand_SwitchToForm(row)
		{
		.Send('VirtualListExpand_SwitchToForm', row)
		}

	ToggleExpand(rowIndex, expand, keepPos? /*unused*/ = false)
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
				layoutOb = Object(ctrl: Object('Record',
					Object('Customizable', tabName: CustomizeExpandControl.LayoutName)))
			_expandRec = rec
			.model.ExpandModel.ConstructAt(layoutOb, rowIndex, this, .model, .rowHeight)
			readOnly? = .Controller.Send('VirtualList_ReadOnly?')
			if not Boolean?(readOnly?)
				readOnly? = not .model.EditModel.RecordLocked?(rec)
			.model.ExpandModel.Expand(rec, layoutOb, .model, :readOnly?)
			ctrl = layoutOb.ctrl
			.model.SetRecordExpanded(rowIndex, layoutOb.rows)
			.ScrollToLeft()
			}
		else
			{
			if rec.vl_expanded_rows is ''
				return
			.model.SetRecordCollapsed(rowIndex)
			.model.ExpandModel.Collapse(rec, this)
			.CommitRecord(rec)
			}
		.Send('VirtualListGrid_AfterExpand', :rec, :ctrl, :expand)
		}

	ExpandButton_EditClicked(row)
		{
		.Send('On_Edit',
			source: VirtualListDummyEditButton(.model.GetRecord(row - .model.Offset)))
		}

	HeaderResize(col, width)
		{
		.Send('VirtualListHeader_HeaderResize', col, width)
		}

	HeaderClick(col)
		{
		.Send('VirtualListHeader_HeaderClick', col: .colModel.Get(col))
		}

	HeaderReorder(oldIdx, newIdx)
		{
		if oldIdx isnt newIdx
			.colModel.ReorderColumn(oldIdx, newIdx)
		.Send('VirtualListHeader_HeaderReorder')
		.colModel.SetHeaderChanged()
		}

	HeaderDividerDoubleClick(col, width)
		{
		.Send('VirtuallistHeader_HeaderDividerDoubleClick', col, width)
		}

	brushMgr: false
	HighlightValues(member, values, color)
		{
		if .brushMgr is false
			.brushMgr = VirtualListBrushes()
		.brushMgr.HighlightValues(member, values, color)
		.Defer(.Repaint, uniqueID: 'Repaint')
		}

	HighlightRecords(recs, color, skipRepaint? = false)
		{
		if .brushMgr is false
			.brushMgr = VirtualListBrushes()
		.brushMgr.HighlightRecords(recs, color)
		if skipRepaint?
			return
		.Defer(.Repaint, uniqueID: 'Repaint')
		}

	ClearHighlightRecord(rec)
		{
		if .brushMgr isnt false
			{
			.brushMgr.ClearHighlightRecord(rec)
			.Defer(.Repaint, uniqueID: 'Repaint')
			}
		}

	ClearHighlight()
		{
		if .brushMgr isnt false
			{
			.brushMgr.Destroy()
			.Defer(.Repaint, uniqueID: 'Repaint')
			}
		.brushMgr = false
		}

	Default(@unused) {	}
	}
