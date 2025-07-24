// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Addon_VirtualListViewBase
	{
	VirtualListHeader_HeaderTrack(col, width, movePix, resizeLeftSide?)
		{
		.Grid.ScrollClientRect(col, width, movePix, resizeLeftSide?)
		}

	VirtualListHeader_HeaderReorder()
		{
		.Grid.Repaint()
		}

	VirtualListHeader_ResetColumns()
		{
		.Grid.ScrollToLeft()
		.Grid.Repaint()
		}

	VirtualListHeader_MeasureWidth(field)
		{
		return VirtualListMeasureWidth(field, .Grid, .Model)
		}

	VirtualListHeader_ContextMenu(col, x, y, extraMenu = #())
		{
		contextMenu = .GetContextMenu()
		contextMenu.SetContext(rec: false, :col, columns: .Model.ColModel.GetColumns())
		contextMenu.Show(.Parent, x, y, :extraMenu)
		}
	}