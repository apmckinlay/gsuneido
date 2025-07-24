// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
TextFormat
	{
	New()
		{
		super(width: 3)
		}
	Print(x, y, w, h, data/*unused*/ = "")
		{
		super.Print(x, y, w, h, _report.GetPage())
		}
	}
