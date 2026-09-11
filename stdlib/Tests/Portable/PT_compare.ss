// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
function(@data)
	{
	n = data.Size(list:)
	for (i = 0; i < n; ++i)
		{
		x = Suneido.Compile(data[i])
		if x isnt x
			return false
		for (j = i + 1; j < n; ++j)
			{
			y = Suneido.Compile(data[j])
			if x >= y or y <= x
				return false
			}
		}
	return true
	}
