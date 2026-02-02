// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Addon_VirtualListViewBase
	{
	VirtualListGrid_HorzScroll(movePix)
		{
		.Header.ScrollClientRect(movePix)
		}

	VirtualListGrid_ViewChanged(ignoreThumb = false)
		{
		if not ignoreThumb
			.Thumb.SetThumbPosition(.Model.GetPosition())
		.RepaintExpandBar()
		}

	VirtualListGrid_DoubleClick(rec, col)
		{
		result = .Send('VirtualList_DoubleClick', :rec, :col)

		if result is 0 and .Model.EditModel.ProtectField is false and rec isnt false and
			col isnt false and rec[col] isnt '' and
			false isnt ctrlClass = GetControlClass.FromField(col)
			ctrlClass.ZoomReadonly(FormatValue(rec[col], col))
		}

	VirtualListGrid_LeftClick(rec, col)
		{
		checkbox = .GetCheckBoxField()
		if col is checkbox and false is .Send('VirtualList_AllowCheckRecord?', rec)
			return
		.Model.CheckRecord(rec, col)
		.Send('VirtualList_LeftClick', rec, :col)
		.UpdateTotalSelected(recalc:)
		if not .Destroyed?()
			.Grid.RepaintSelectedRows()
		}

	VirtualListGrid_Edit(rec, col)
		{
		.Send('VirtualList_Edit', rec, col)
		}

	VirtualListGrid_ItemSelected(rec)
		{
		.Send('VirtualList_ItemSelected', rec)
		.Grid.RepaintSelectedRows()
		if .Model.ExpandModel isnt false
			.Model.ExpandModel.ClearAllSelections()
		.RefreshValid(rec)
		.GetViewControls().expandBar.RefreshEditState()
		}

	VirtualListGrid_AfterExpand(rec, ctrl, expand)
		{
		if .Model.ExpandModel is false
			return

		c = VirtualListViewExtra
		if c.Member?('VirtualListGrid_AfterExpand')
			return c.VirtualListGrid_AfterExpand(rec, ctrl, :expand, view: .Parent)

		.repaintAndResize()
		if expand is true
			{
			rowNum = .Model.ExpandModel.FindIfOverLimit(rec, .Model.GetRecordRowNum)
			if rowNum isnt false
				{
				.Grid.ToggleExpand(rowNum, expand: false)
				.repaintAndResize()
				}
			}

		.Send('VirtualList_AfterExpand', :rec, :ctrl)
		}

	repaintAndResize()
		{
		.RepaintExpandBar()
		.Scroll.ResizeWindow()
		}

	VirtualListGrid_MouseLeave()
		{
		if .GetViewControls().expandBtns isnt false and
			not .GetViewControls().expandBtns.MouseOver?()
			.ExpandBar.HideButton()
		}

	VirtualListGrid_Return()
		{
		.Send('VirtualList_Return')
		}

	VirtualListGrid_Tab()
		{
		.Send('VirtualList_Tab')
		}

	VirtualListGrid_Escape()
		{
		.Send('VirtualList_Escape')
		}

	VirtualListGrid_Space()
		{
		checkBoxColumn = .GetCheckBoxField()
		if checkBoxColumn isnt false and
			false is .Send('VirtualList_AllowCheckRecord?', .GetSelectedRecord())
			return
		for rec in .Grid.GetSelectedRecords()
			.Model.CheckRecord(rec, forceCheck:)
		.Send('VirtualList_Space')
		.UpdateTotalSelected(recalc:)
		.Grid.RepaintSelectedRows()
		}

	VirtualListGrid_ContextMenu(rec, col, x, y, row_num)
		{
		.GetContextMenu().ShowMenu(.Parent, rec, col, row_num, [:x, :y])
		}

	VirtualListGrid_SetStartLast()
		{
		return .IsLinked?() is true ? true : .SaveFirst()
		}

	VirtualListGrid_Expand(rec)
		{
		return .Send('VirtualList_Expand', rec)
		}

	VirtualListGrid_MouseMove(row_num, expanded, invalid, curLeft = 0)
		{
		.ExpandBar.MoveButtonTo(row_num, expanded, invalid, :curLeft)
		}

	VirtualListGrid_MouseWheel(wParam)
		{
		.Send('VirtualList_MouseWheel', wParam)
		}

	VirtualListGrid_NewRowAdded(newRec)
		{
		.Send('VirtualList_NewRowAdded', newRec)
		.ExpandBar.Repaint()
		.Scroll.ResizeWindow()
		}

	VirtualListGrid_RowDeleted()
		{
		.Scroll.ResizeWindow()
		}

	VirtualListGrid_Insert()
		{
		if false is .Send('VirtualList_AllowInsert')
			return

		.Grid.InsertRow(pos: 'current')
		}

	VirtualListGrid_RepaintingRow(row_num)
		{
		.ExpandBar.RepaintRow(row_num)
		}
	}
