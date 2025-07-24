// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Format
	{
	Size: #(h: 240, w: 0, d: 0)
	New(ymin = false)
		{
		if ymin isnt false
			.Size = Object(h: ymin.InchesInTwips().Int(), w: 0, d: 0)
		}
	GetSize(data/*unused*/ = false)
		{
		return .Size
		}
	}