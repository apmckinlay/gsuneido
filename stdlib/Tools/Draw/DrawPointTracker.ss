// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
DrawTracker
	{
	New(hwnd, item)
		{
		.item = item
		}
	MouseUp(x, y)
		{
		return (.item)(x, y)
		}
	}