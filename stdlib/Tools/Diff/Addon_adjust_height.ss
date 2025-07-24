// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	Init()
		{
		.adjusted = false
		}
	Painted()
		{
		scrollBar? = HorzScrollBar?(.Hwnd)
		if .adjusted is true and not scrollBar?
			{
			.adjusted = false
			.Parent.Ymin -= GetSystemMetrics(SM.CXHSCROLL)
			.WindowRefresh()
			}
		else if .adjusted is false and scrollBar?
			{
			.adjusted = true
			.Parent.Ymin += GetSystemMetrics(SM.CXHSCROLL)
			.WindowRefresh()
			}
		}
	}