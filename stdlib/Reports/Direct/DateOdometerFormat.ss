// Copyright (C) 2026 Axon Development Corporation. All rights reserved worldwide.
// for date or odometer
ShortDateFormat
	{
	Width() // ensure width is wide enough
		{
		minWidth = 11
		return Max(minWidth, super.Width())
		}
	}
