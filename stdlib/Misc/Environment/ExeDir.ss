// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	// using forward slashes ("ToUnix") to avoid usage of escape characters
	return Paths.ToUnix(Paths.ParentOf(ExePath()))
	}
