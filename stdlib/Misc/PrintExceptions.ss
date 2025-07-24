// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
function (block)
	{
	try
		block()
	catch (e)
		{
		Print('ERROR:', e)
		Print(FormatCallStack(e.Callstack(), levels: 5, indent:))
		}
	}
