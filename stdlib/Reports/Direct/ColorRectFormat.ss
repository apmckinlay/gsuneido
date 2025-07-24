// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Format
	{
	New(.data = false, .width = .5, .height = .25)
		{
		}
	GetSize(data /*unused*/ = false)
		{
		return Object(w: .width, h: .height, d: 0)
		}
	Print(x, y, w, h, data = false)
		{
		if .data isnt false
			data = .data
		_report.AddRect(x, y, w, h, thick: 1, fillColor: data)
		}
	}