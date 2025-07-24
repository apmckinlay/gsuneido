// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Addon_VirtualListViewBase
	{
	VirtualListThumb_Dragging(speed)
		{
		.Grid.VertScrolling(.Grid.GetRows(speed))
		}

	VirtualListThumb_EndDragging()
		{
		.stopScrolling()
		}

	On_VirtualListThumb_ArrowDown()
		{
		.Grid.VertScroll(1)
		}

	On_VirtualListThumb_ArrowDown_MouseHold()
		{
		.Grid.VertScrolling(1, notify?:)
		}

	On_VirtualListThumb_ArrowDown_MouseUp()
		{
		.stopScrolling()
		}

	On_VirtualListThumb_ArrowUp()
		{
		.Grid.VertScroll(-1)
		}

	On_VirtualListThumb_ArrowSelect()
		{
		if .SaveFirst() is false
			return

		if 0 isnt .Send('VirtualList_On_Select')
			return

		.Thumb.SetSelectPressed(pressed:)
		SelectControl(.Parent, .Model.ColModel.GetSelectMgr().Name(),
			okbutton:, defaultButton: "Select",
			noUserDefaultSelects?: not .Model.ColModel.UserDefaultSelectEnabled?())
		// Need to handle that this could be destroyed while the select control is open
		if .Model is false
			return
		.Thumb.SetSelectPressed(.Model.ColModel.HasSelectedVals?())
		}

	On_VirtualListThumb_ArrowUp_MouseHold()
		{
		.Grid.VertScrolling(-1, notify?:)
		}

	On_VirtualListThumb_ArrowUp_MouseUp()
		{
		.stopScrolling()
		}

	On_VirtualListThumb_ArrowHome()
		{
		.Grid.KEYDOWN(VK.HOME)
		}

	On_VirtualListThumb_ArrowEnd()
		{
		.Grid.KEYDOWN(VK.END)
		}

	VirtualListThumb_MouseHold(direction)
		{
		.Grid.VertScrolling(direction * (.Model.VisibleRows - 1), notify?:)
		}

	VirtualListThumb_MouseUp()
		{
		.stopScrolling()
		}

	VirtualListThumb_PageUp()
		{
		.Grid.VertScroll(-.Model.VisibleRows + 1)
		}

	VirtualListThumb_PageDown()
		{
		.Grid.VertScroll(.Model.VisibleRows - 1)
		}

	stopScrolling()
		{
		.Grid.StopVScrolling()
		.Thumb.SetThumbPosition(.Model.GetPosition())
		}

	VirtualListThumb_Expand(rowIndex, expand, keepPos? = false)
		{
		.Parent.ExpandedIndex = rowIndex
		.Grid.ToggleExpand(rowIndex, expand, :keepPos?)
		.Parent.ExpandedIndex = false
		}
	}