// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
DrawTracker
	{
	New(.hwnd, .item) { }
	x1: false
	MouseDown(x, y)
		{
		SetCapture(.hwnd)
		.x0 = x
		.y0 = y
		}
	MouseMove(x, y)
		{
		hdc = GetDC(.hwnd)
		if .x1 isnt false
			.line(hdc) // erase previous line
		.x1 = x
		.y1 = y
		.line(hdc)
		ReleaseDC(.hwnd, hdc)
		}
	MouseUp(x/*unused*/, y/*unused*/)
		{
		ReleaseCapture()
		item = false
		if .x1 isnt false
			{
			hdc = GetDC(.hwnd)
			.line(hdc) // erase previous line
			item = .Line(.x0, .y0, .x1, .y1)
			ReleaseDC(.hwnd, hdc)
			.x1 = false
			}
		return item
		}
	line(hdc)
		{
		oldrop = SetROP2(hdc, R2.NOTXORPEN)
		MoveTo(hdc, .x0, .y0)
		LineTo(hdc, .x1, .y1)
		SetROP2(hdc, oldrop)
		}

	Line(x0, y0, y1, y2)
		{
		return (.item)(x0, y0, y1, y2)
		}

	GetSuJSTracker()
		{
		return 'SuDrawLineTracker'
		}
	}