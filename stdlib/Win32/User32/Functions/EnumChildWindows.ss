// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
function (hwnd, block)
	{
	EnumChildWindowsApi(hwnd, b = {|childHwnd, unused| block(childHwnd) }, lParam: 0)
	ClearCallback(b)
	}
