// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
// INCOMPLETE!
Controller
	{
	Title: 'Library Diff'
	CallClass(_hwnd = 0)
		{
		ToolDialog(hwnd, [this], keep_size: false)
		}
	Controls:
		(Vert
			(Pair
				(Static Diff)
				(LibLocate, name: 'name1'))
			Skip
			(Pair
				(Static with)
				(LibLocate, name: 'name2'))
			Skip
			(HorzEqual
				Fill
				(Button Diff)
				Skip
				(Button Cancel)
				)
			xstretch: 0)
	On_Diff()
		{
		}
	LocateEscape()
		{
		.On_Cancel()
		}
	}