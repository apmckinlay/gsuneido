// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
function (src, type, displayed)
	{
	Assert(src isString:)
	try
		x = src.Compile()
	catch (e)
		return type is 'Exception' and e.Has?(displayed)
	return Type(x) is type and Display(x) is displayed
	}