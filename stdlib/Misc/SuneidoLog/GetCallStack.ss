// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
function (skip = 0, limit = 999)
	{
	try
		throw ""
	catch (e)
		return e.Callstack()[skip + 1 :: limit]
	}