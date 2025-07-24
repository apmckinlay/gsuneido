// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Format
	{
	New(format)
		{
		.format = _report.Construct(format)
		.Xmin = .format.Xmin
		.Ymin = .format.Ymin
		.Xstretch = .format.Xstretch
		.Ystretch = .format.Ystretch
		}
	GetSize(@args)
		{
		size = .format.GetSize(@args).Copy()
		size.d = size.h
		return size
		}
	Print(@args)
		{
		.format.Print(@args)
		}
	}
