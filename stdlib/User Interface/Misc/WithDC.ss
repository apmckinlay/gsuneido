// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
function (hwnd, block)
	{
	dc = GetDC(hwnd)
	Finally({ block(dc) },
		{ ReleaseDC(hwnd, dc) })
	}