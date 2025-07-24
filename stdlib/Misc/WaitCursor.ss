// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// NOTE: don't use this unless you are going to "block"
function (block)
	{
	oldcursor = SetCursor(LoadCursor(NULL, IDC.WAIT))
	try
		block()
	catch (x)
		{
		SetCursor(oldcursor)
		throw x
		}
	SetCursor(oldcursor)
	}