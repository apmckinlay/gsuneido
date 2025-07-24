// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Generator
	{
	New(lib1, lib2)
		{
		.list = LibCompare(lib1, lib2)
		.i = -1
		}
	Header()
		{
		return #(Vert (Text 'Differences') Vskip)
		}
	Next()
		{
		++.i
		if .i >= .list.Size()
			return false
		return _report.Construct(
			Object('Code',
				.list[.i][0] $ .list[.i][2],
				w: _report.GetWidth()))
		}
	}
