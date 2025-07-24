// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
function (rect = false)
	{
	if false is rect
		rect = Object()
	info = GetMonitorInfo(MonitorFromRect(rect, MONITOR.DEFAULTTOPRIMARY))
	return info.rcWork
	}