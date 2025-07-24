// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (page)
	{
	// If the server restarts while in edit mode, this variable will be cleared out.
	if not Suneido.Member?('WikiLock')
		Suneido.WikiLock = Object()
	Suneido.WikiLock.Delete(page)
	}