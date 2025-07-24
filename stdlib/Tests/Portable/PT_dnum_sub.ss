// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
function (@args)
	{
	sub = function (x, y) { x - y }
	return Pt.Nums(args)
		{|x,y,z| Pt.Binary(sub, x, y, z) and Pt.Binary(sub, y, x, -z) }
	}
