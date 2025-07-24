// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Set Order'
	CallClass(hwnd, oldorder)
		{
		OkCancel(Object(this, oldorder), .Title, hwnd)
		}
	New(oldorder)
		{ .Vert.Number.Set(.oldorder = oldorder) }

	Controls: (Vert
		(Static "Enter Item's Order")
		(Skip 3)
		(Number mask: '###.###'))

	OK()
		{
		return .Vert.Number.Get()
		}

	Cancel()
		{
		return .oldorder
		}

	}