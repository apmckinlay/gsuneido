// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
WndProc
	{
	Name: 'VirtualListExpandBar'
	Ystretch: 1

	enableExpand: false
	model: false
	deleteImage: false
	delete: false
	editBtns: #()
	New(.preventCustomExpand?, .enableDeleteBar = false, .switchToForm = false)
		{
		.CreateWindow("SuBtnfaceArrowNoDblClks", windowName: "",
			style: WS.VISIBLE | WS.CLIPSIBLINGS)
		.SubClass()

		.SetFont(size: "", text: "W...")
		.Xmin = .enableDeleteBar ? .Ymin + 4 : 0 /*= padding */
		.deleteImage = ImageResource('delete.emf')
		.delete = CreateSolidBrush(CLR.RED)
		.editBtns = Object()
		}
	SetInfo(.model, .rowHeight, .headerYmin, .expandBtns)
		{
		.destroyEditBtns()
		.showExpandBar()
		.expandBtns.Reset()
		}

	showExpandBar()
		{
		if .Send('VirtualListGrid_Expand', []) isnt 0
			.enableExpand = true
		else if not .preventCustomExpand?
			{
			if 0 is expandInfo = .Send('Customizable_ExpandInfo')
				expandInfo = Object(availableFields: false, defaultLayout: '')
			if 0 is customKey = .Send('GetAccessCustomKey')
				customKey = ''
			table = .model.GetTableName()
			c = Customizable(table, defaultLayout: expandInfo.defaultLayout,
				user: Suneido.User, :customKey)
			.enableExpand = c.LayoutExists?(CustomizeExpandControl.LayoutName)
			}
		.Xmin = .enableExpand or .enableDeleteBar ? .Ymin + 4 : 0 /*= padding */
		.expandBtns.EnsureExpandButton(
			visible: .Xmin isnt 0 and .model.ExpandModel isnt false)
		}

	ShowExpand?()
		{
		return .enableExpand or .enableDeleteBar
		}

	MoveButtonTo(row_num, expanded, invalid, curLeft = 0)
		{
		if .hideBtn?(invalid, row_num)
			{
			.expandBtns.SetVisible(false)
			return false
			}
		else
			{
			.expandBtns.SetVisible(true)
			.expandBtns.MoveTo(
				row_num, .rowHeight, .headerYmin, minus: expanded, :curLeft)
			return true
			}
		}

	HideButton()
		{
		.expandBtns.SetVisible(false)
		}

	hideBtn?(invalid, row_num)
		{
		if .invalid(invalid, row_num)
			return true
		if .switchToForm is true
			return false
		if .Xmin is 0
			return true
		if .enableExpand is false and .enableDeleteBar is true
			return true
		return false
		}

	invalid(invalid, row_num)
		{
		if invalid or .model is false or row_num < 0
			return true
		rec = .model.GetRecord(row_num)
		return rec is false or rec.vl_deleted is true
		}

	ShowEditButtons()
		{
		.destroyEditBtns()
		if .model is false or not .model.EditModel.Editable?()
			return
		if false is .Controller.Send('VirtualList_ShowEditButton?')
			return
		for (i = 0; i < .model.VisibleRows; i++)
			{
			rec = .model.GetRecord(i)
			if rec isnt false and rec.vl_expanded_rows isnt ''
				if .needEditBtn?(rec)
					.showEditButton(rec, i)
			}
		}

	showEditButton(rec, i)
		{
		locked = .model.EditModel.RecordLocked?(rec)
		editBtn = .Construct('VirtualListEditButton', .Xmin)
		editBtn.MoveTo(i + 1, .rowHeight, .headerYmin)
		if locked
			editBtn.Pushed?(true)
		.editBtns.Add(editBtn)
		}

	needEditBtn?(rec)
		{
		// control data is only what is in the expand control (not the columns)
		if false isnt recordControl = .model.ExpandModel.GetRecordControl(rec)
			recordControl.GetControlData().Members().Each()
				{
				if not FieldProtected?(it, rec, .model.EditModel.ProtectField)
					return true
				}
		return false
		}

	RefreshEditState()
		{
		if .model is false or not .model.EditModel.Editable?()
			return
		for btn in .editBtns
			{
			pushed = .model.EditModel.RecordLocked?(btn.GetRecord())
			btn.Pushed?(pushed)
			}
		}

	RepaintRow(row_num)
		{
		if .model is false
			return
		y = row_num * .rowHeight + .headerYmin
		rect = Object(left: 0, top: y, right: .Ymin, bottom: y + .rowHeight)
		InvalidateRect(.Hwnd, rect, erase:)
		}

	ERASEBKGND()
		{
		return 1 // no erase, the window is completly redrawn by PAINT
		}

	PAINT()
		{
		hdc = BeginPaint(.Hwnd, ps = Object())
		if .model isnt false
			WithBkMode(hdc, TRANSPARENT, { .paint(hdc, ps.rcPaint) })
		EndPaint(.Hwnd, ps)
		return 0
		}

	paint(hdc, rc)
		{
		topRow = .getRowFromY(rc.top)
		top = topRow * .rowHeight + .headerYmin
		numRows = .model.VisibleRows
		.paintMarkCol(hdc, rc, top, topRow, numRows)
		}

	getRowFromY(y)
		{
		return Max(0, ((y - .headerYmin) / .rowHeight).Int())
		}

	GetRecordFromY(y)
		{
		return .model.GetRecord(.getRowFromY(y) - 1)
		}

	paintMarkCol(hdc, rc, top, topRow, numRows)
		{
		rcCell = Object(left: 0, top: top - .rowHeight,
			right: .Xmin, bottom: top)
		for (row = topRow; row < numRows and rcCell.top < rc.bottom; ++row)
			{
			rcCell.top += .rowHeight
			rcCell.bottom += .rowHeight
			FillRect(hdc, rcCell, GetSysColorBrush(COLOR.BTNFACE))
			rec = .model.GetRecord(row)
			if rec isnt false and rec.vl_deleted is true
				.drawDeleteMark(rcCell, hdc)
			}
		if rcCell.bottom < rc.bottom
			{
			rcCell.top = rcCell.bottom
			rcCell.bottom = rc.bottom
			FillRect(hdc, rcCell, GetSysColorBrush(COLOR.BTNFACE))
			}
		if rc.right < .Xmin
			return false
		rc.left = .Xmin
		ExcludeClipRect(hdc, 0, rc.top, .Xmin, rc.bottom)
		return true
		}

	drawDeleteMark(rcCell, hdc)
		{
		padding = ScaleWithDpiFactor(4)  /*= padding*/
		cellWidth = rcCell.right - rcCell.left
		cellHeight = rcCell.bottom - rcCell.top
		wh = Min(cellWidth, cellHeight) - padding * 2
		leftPad = (cellWidth - wh) / 2
		topPad = (cellHeight - wh) / 2
		.deleteImage.Draw(hdc, rcCell.left + leftPad, rcCell.top + topPad,
			wh, wh, .delete)
		}

	MOUSEMOVE(lParam)
		{
		.Send('VirtualListExpandBar_MouseMove', HISWORD(lParam))
		return 0
		}

	MOUSEWHEEL(wParam)
		{
		.Send('VirtualList_MouseWheel', wParam)
		return 0
		}

	destroyEditBtns()
		{
		for btn in .editBtns
			btn.Destroy()
		.editBtns = Object()
		}

	Destroy()
		{
		if .deleteImage isnt false
			.deleteImage.Close()
		DeleteObject(.delete)
		.destroyEditBtns()
		super.Destroy()
		}
	}