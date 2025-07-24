// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
// ugly hack to get find and replace fields to line up
// Grid would be simpler but doesn't handle stretch
// see also: FindStaticControl and ReplaceBarHorzControl
StaticControl
	{
	New(@args)
		{
		super(@args)
		.Xmin = Max(.Xmin, .TextExtent('Find').x +
			.TextExtent('M').y + 10 + ScaleWithDpiFactor(10)) // 10 = skip
		}
	}