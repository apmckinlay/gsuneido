// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Format
	{
	New(.height = 1440, .thick = 10, .before = 120, .after = 120)
		{
		}
	GetSize(data /*unused*/ = false)
		{
		return Object(h: .height, w: .before + .thick + .after, d: 0)
		}
	Print(x, y, w /*unused*/, h, data /*unused*/ = false)
		{
		x += .before
		_report.AddLine(x, y, x, y + h, .thick)
		}
	}