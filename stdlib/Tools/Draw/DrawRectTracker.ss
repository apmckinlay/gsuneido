// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
DrawTracker
	{
	x1: false
	y1: false
	New(hwnd, item, .canvas)
		{
		.hwnd = hwnd
		.item = item
		}
	rect: ()
	rectRequiredSize: 4
	prev_ClipCursor: false
	MouseDown(x, y)
		{
		GetClipCursor(.prev_ClipCursor = Object())
		r = GetWindowRect(.hwnd)
		ClipCursor(r)
		.rect = Object()
		.x1 = x
		.y1 = y
		}
	MouseMove(x, y)
		{
		if .x1 is false or .y1 is false
			return
		.x2 = x
		.y2 = y
		hdc = GetDC(.hwnd)
		if .rect.Size() is .rectRequiredSize
			DrawFocusRect(hdc, .rect) // erase previous rect
		.rect.left = Min(.x1, x)
		.rect.right = Max(.x1, x)
		.rect.top = Min(.y1, y)
		.rect.bottom = Max(.y1, y)
		DrawFocusRect(hdc, .rect)
		ReleaseDC(.hwnd, hdc)
		}
	MouseUp(x, y)
		{
		if .prev_ClipCursor is false
			return false
		ClipCursor(.prev_ClipCursor)
		item = false
		if .rect.Size() is .rectRequiredSize
			{
			hdc = GetDC(.hwnd)
			DrawFocusRect(hdc, .rect) // erase previous rect
			ReleaseDC(.hwnd, hdc)
			if ((.x1 - .x2).Abs() > 2 or (.y1 - .y2).Abs() > 2)
				item = .Rect(.x1, .y1, .x2, .y2)
			else
				item = .Point(.x1, .y1)
			.rect = Object()
			}
		else
			item = .Point(x, y)
		return item
		}
	Point(x /*unused*/, y /*unused*/)
		{ return false } // to be overridden by derived classes
	Rect(x1, y1, x2, y2)
		{
		return (.item)(x1, y1, x2, y2, canvas: .canvas)
		}
	resize_rect: ()
	varyx: 'none'
	varyy: 'none'
	ResizeDown(item, x, y)
		{
		r = item.BoundingRect()
		xmid = (r.x1 + r.x2) / 2
		ymid = (r.y1 + r.y2) / 2
		if item.InHandleArea(r.x1, r.y1, x, y)
			{ .varyx = 'left'; .varyy = 'top' }
		else if item.InHandleArea(r.x1, r.y2, x, y)
			{ .varyx = 'left'; .varyy = 'bottom' }
		else if item.InHandleArea(r.x2, r.y1, x, y)
			{ .varyx = 'right'; .varyy = 'top' }
		else if item.InHandleArea(r.x2, r.y2, x, y)
			{ .varyx = 'right'; .varyy = 'bottom' }
		else if item.InHandleArea(r.x1, ymid, x, y)
			{ .varyx = 'left'; .varyy = 'none' }
		else if item.InHandleArea(r.x2, ymid, x, y)
			{ .varyx = 'right'; .varyy = 'none' }
		else if item.InHandleArea(xmid, r.y1, x, y)
			{ .varyx = 'none'; .varyy = 'top' }
		else if item.InHandleArea(xmid, r.y2, x, y)
			{ .varyx = 'none'; .varyy = 'bottom' }
		.resize_rect = Object(left: r.x1, right: r.x2, top: r.y1, bottom: r.y2)
		hdc = GetDC(.hwnd)
		DrawFocusRect(hdc, .resize_rect)
		}
	ResizeMove(item, x, y)
		{
		hdc = GetDC(.hwnd)
		DrawFocusRect(hdc, .resize_rect) // erase previous rect
		if not item.Method?(#DoResizeMove)
			{
			if .varyx isnt 'none'
				.resize_rect[.varyx] = x
			if .varyy isnt 'none'
				.resize_rect[.varyy] = y
			}
		else
			item.DoResizeMove(x, y, .varyx, .varyy, .resize_rect)
		.uninvert(.resize_rect)
		DrawFocusRect(hdc, .resize_rect)
		ReleaseDC(.hwnd, hdc)
		ClipCursor(NULL)
		}
	uninvert(r)
		{
		if r.Empty?()
			return
		if r.left > r.right
			{
			r.Swap('left', 'right')
			.varyx = #(left: right, right: left, none: 'none')[.varyx]
			}
		if r.top > r.bottom
			{
			r.Swap('top', 'bottom')
			.varyy = #(top: bottom, bottom: top, none: 'none')[.varyy]
			}
		}
	ResizeUp(item /*unused*/, x /*unused*/, y /*unused*/)
		{
		hdc = GetDC(.hwnd)
		DrawFocusRect(hdc, .resize_rect) // erase previous rect
		}
	GetSuJSTracker()
		{
		return 'SuDrawRectTracker'
		}
	}
