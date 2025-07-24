// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Control
	{
	New(min = 0, fill = 1)
		{
		if (not .Parent.Member?("Dir"))
			.Xstretch = .Ystretch = fill
		else if (.Parent.Dir is "vert")
			.Ystretch = fill
		else if (.Parent.Dir is "horz")
			.Xstretch = fill
		if (min > 0)
			{
			min = ScaleWithDpiFactor(min)
			if (not .Parent.Member?("Dir"))
				.Xmin = .Ymin = min
			else if (.Parent.Dir is "vert")
				.Ymin = min
			else if (.Parent.Dir is "horz")
				.Xmin = min
			}
		}
	GetReadOnly() // read-only not applicable to fill
		{
		return true
		}
	}
