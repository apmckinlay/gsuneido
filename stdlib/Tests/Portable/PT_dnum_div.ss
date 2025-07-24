// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
function (@args)
	{
	div = function (x, y) { x / y }
	return Pt.Nums(args)
		{|x,y,z| Pt.Binary(div, x, y, z) }
	}
