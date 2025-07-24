// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
// part of AutoChooseList
ListBoxControl
	{
	LBUTTONDOWN(lParam /*unused*/)
		{
		return 0
		}
	LBUTTONUP(lParam)
		{
		result = SendMessage(.Hwnd, LB.ITEMFROMPOINT, 0, lParam)
		item = LOWORD(result)
		if (HIWORD(result) is 0 and item >= 0)
			.Send('AutoListBox_Click', item)
		return 0
		}
	MOUSEMOVE(lParam)
		{
		result = SendMessage(.Hwnd, LB.ITEMFROMPOINT, 0, lParam)
		item = LOWORD(result)
		if (HIWORD(result) is 0 and item >= 0)
			SendMessage(.Hwnd, LB.SETCURSEL, item, 0)
		return 0
		}
	}
