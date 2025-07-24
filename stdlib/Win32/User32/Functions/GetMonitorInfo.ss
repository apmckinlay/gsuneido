// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
function (hMonitor)
	{
	return GetMonitorInfoApi(hMonitor, mi = Object(cbSize: MONITORINFO.Size()))
		? mi
		: false
	}