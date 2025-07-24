// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
function (hwnd)
	{
	if GetWindowRectApi(hwnd, r = [])
		return r
	SuneidoLog('ERROR: GetWindowRect failed!', calls:,
		params: Record(:hwnd, :r, isWindow: IsWindow(hwnd)))
	return []
	}