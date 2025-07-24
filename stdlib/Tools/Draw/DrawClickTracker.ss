// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
DrawTracker
	{
	New(.hwnd, .item, .canvas) {}
	MouseDown(x/*unused*/, y/*unused*/)
		{
		GetClipCursor(.prev_ClipCursor = Object())
		r = GetWindowRect(.hwnd)
		ClipCursor(r)
		}

	MouseUp(x, y)
		{
		if .prev_ClipCursor isnt false
			ClipCursor(.prev_ClipCursor)
		return .Point(x, y)
		}

	Point(x, y)
		{
		return (.item)(x, y, canvas: .canvas)
		}

	GetSuJSTracker()
		{
		return 'SuDrawClickTracker'
		}
	}
