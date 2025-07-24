// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
function (@fields)
	{
	return {|x, y|
		c = 0
		for f in fields
			if 0 isnt c = Cmp(x[f], y[f])
				break
		c < 0
		}
	}