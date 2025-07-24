// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
function (@args)
	{
	add = function (x, y) { x + y }
	return Pt.Nums(args)
		{|x,y,z| Pt.Binary(add, x, y, z) and Pt.Binary(add, y, x, z) }
	}
