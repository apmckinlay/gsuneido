// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
function (@data)
	{
	n = data.Size(list:)
	for (i = 0; i < n; ++i)
		{
		x = data[i].Compile()
		if x isnt x
			return false
		for (j = i + 1; j < n; ++j)
			{
			y = data[j].Compile()
			if x >= y or y <= x
				return false
			}
		}
	return true
	}