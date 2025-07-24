// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
// ugly hack to get find and replace fields to line up
// Grid would be simpler but doesn't handle stretch
// see also: ReplaceStaticControl and ReplaceBarHorzControl
StaticControl
	{
	New(@args)
		{
		super(@args)
		.Xmin = Max(.Xmin, .TextExtent('Replace').x -
			.TextExtent('M').y - 10 - ScaleWithDpiFactor(10 /*= skip */))
		}
	}