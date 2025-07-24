// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
function (@data)
	{
	n = data.Size(list:)
	for (i = 0; i < n; ++i)
		{
		x = Pt.Num(data[i])
		if x isnt x
			return false
		for (j = i + 1; j < n; ++j)
			{
			y = Pt.Num(data[j])
			if x >= y or y <= x
				return false
			}
		}
	return true
	}
