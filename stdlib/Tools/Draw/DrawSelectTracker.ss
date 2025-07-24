// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
DrawRectTracker
	{
	New(hwnd, item/*unused*/, canvas)
		{
		super(hwnd, false, canvas)
		.hwnd = hwnd
		.canvas = canvas
		.rects = Object()
		}
	MouseDown(x, y)
		{
		.rects = Object()
		if false isnt i = .canvas.ItemAtPoint(x, y)
			{
			.item = .canvas.GetAllItems()[i]
			if not .canvas.GetSelected().Has?(.item)
				{
				.canvas.MaybeClearSelect()
				.canvas.Select(i)
				}
			else
				{
				if KeyPressed?(VK.CONTROL) or KeyPressed?(VK.SHIFT)
					.canvas.UnSelect(i)
				}

			for item in .canvas.GetSelected()
				{
				r = item.BoundingRect()
				rect = Object(left: r.x1, right: r.x2, top: r.y1, bottom: r.y2, :item)
				.rects.Add(rect)
				hdc = GetDC(.hwnd)
				DrawFocusRect(hdc, rect)
				}
			.x = .origx = x
			.y = .origy = y
			}
		else
			super.MouseDown(x, y)
		}
	dragging()
		{
		return .rects.NotEmpty?()
		}
	MouseMove(x, y)
		{
		if .dragging()
			{
			hdc = GetDC(.hwnd)
			dx = x - .x
			dy = y - .y

			move = DrawSelectTracker_CalcNextMove(dx, dy, .rects, .canvas)
			for rect in .rects
				{
				DrawFocusRect(hdc, rect) // erase previous rect
				rect.left += move.x
				rect.right += move.x
				rect.top += move.y
				rect.bottom += move.y
				DrawFocusRect(hdc, rect)
				}
			.x += move.x
			.y += move.y
			}
		else
			super.MouseMove(x, y)
		}


	MouseUp(x, y)
		{
		if .dragging()
			{
			hdc = GetDC(.hwnd)
			dx = .x - .origx
			dy = .y - .origy
			for rect in .rects
				{
				InvalidateRect(.canvas.Hwnd,
					.canvas.RectConversion(rect.item.BoundingRect()), true)
				DrawFocusRect(hdc, rect) // erase previous rect
				rect.item.Move(dx, dy)
				InvalidateRect(.canvas.Hwnd,
					.canvas.RectConversion(rect.item.BoundingRect()), true)
				}
			return false
			}
		else
			return super.MouseUp(x, y)
		}
	AddPoint(x, y)
		{
		selected = .canvas.GetSelected()
		for (item in selected)
			if (item.Contains(x, y))
				return
		.canvas.SelectPoint(x, y)
		for (item in selected)
			.canvas.Select(.canvas.GetAllItems().Find(item))
		}
	Point(x, y)
		{
		if KeyPressed?(VK.CONTROL) or KeyPressed?(VK.SHIFT)
			.AddPoint(x, y)
		else
			.canvas.SelectPoint(x, y)
		return false
		}
	Rect(x1, y1, x2, y2)
		{
		.canvas.SelectRect(Min(x1, x2), Min(y1, y2), Max(x1, x2), Max(y1, y2))
		return false
		}
	GetSuJSTracker()
		{
		return 'SuDrawSelectTracker'
		}
	}