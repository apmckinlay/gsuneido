// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Format
	{
	radiusHeightRatio: 0.4
	New(.data = false, .radius = 0.25)
		{
		}
	GetSize()
		{
		}
	Print(x, y, w/*unused*/, h, data = false)
		{
		if .data isnt false
			data = .data
		_report.AddCircle(x + h/2, y + h/2, .radiusHeightRatio * h, 1, fillColor: data)
		}
	}