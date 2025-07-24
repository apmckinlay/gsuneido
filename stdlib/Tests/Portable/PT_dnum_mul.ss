// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
function (@args)
	{
	mul = function (x, y) { x * y }
	return Pt.Nums(args)
		{|x,y,z| Pt.Binary(mul, x, y, z) and Pt.Binary(mul, y, x, z) }
	}
