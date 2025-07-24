// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Format
	{
	Size: #(h: 0, w: 240, d: 0)
	New(xmin = false)
		{
		if (xmin isnt false)
			.Size = Object(h: 0, w: xmin.InchesInTwips(), d: 0)
		}
	GetSize(data/*unused*/ = false)
		{
		return .Size
		}
	}
