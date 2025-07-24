// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
function (field, grid, model)
	{
	WithDC(grid.Hwnd)
		{|dc|
		DoWithHdcObjects(dc, [grid.GetFont()])
			{
			model.ColModel.SetDC(dc)
			maxSize = false
			for (i = 0; i < model.VisibleRows; i++)
				{
				if false is rec = model.GetRecord(i)
					break
				s = model.ColModel.MeasureWidth(field, rec)
				if s.w is 0
					continue
				maxSize = Max(s.w + 2 * VirtualListGridPaint.HorzMargin, maxSize)
				}
			}
		}
	return maxSize
	}
