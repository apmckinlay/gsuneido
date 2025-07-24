// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
function (hdc, obs, block)
	{
	olds = Object()
	for ob in obs
		olds.Add(SelectObject(hdc, ob))
	Finally(block,
		{
		for ob in olds
			SelectObject(hdc, ob)
		})
	}