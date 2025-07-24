// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
class
	{
	DefaultSize: 9
	PtSize(lfsize, logPixel)
		{
		if lfsize is false
			lfsize = Suneido.logfont.lfHeight
		return -(lfsize * PointsPerInch / logPixel).Round(0)
		}
	LfSize(ptsize, logPixel)
		{
		if ptsize is false
			ptsize = .DefaultSize
		return -(ptsize * logPixel / PointsPerInch).Round(0)
		}
	}