// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
_Addon_VirtualListView_Header
	{
	VirtualListHeader_DisableSort?()
		{
		.Send('VirtualList_DisableSort?')
		}

	VirtualListHeader_HeaderReorder()
		{
		.GetViewControls().header.RefreshSort(.Model.GetPrimarySort())
		super.VirtualListHeader_HeaderReorder()
		}

	VirtuallistHeader_HeaderDividerDoubleClick(col, width)
		{
		.GetViewControls().header.HeaderDividerDoubleClick(col, width)
		}

	VirtualListHeader_HeaderResize(col, width)
		{
		.GetViewControls().header.HeaderResize(col, width)
		}

	VirtualListHeader_ResetColumns()
		{
		.Repaint()
		}

	VirtualListHeader_ContextMenu(col, x, y)
		{
		extraMenu = .GetViewControls().header.BuildSortMenu()
		super.VirtualListHeader_ContextMenu(col, x, y, :extraMenu)
		}
	}