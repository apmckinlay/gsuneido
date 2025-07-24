// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
function (block)
	{
	try
		block()
	catch (e)
		return throw e
	return throw "expected exception not thrown"
	}