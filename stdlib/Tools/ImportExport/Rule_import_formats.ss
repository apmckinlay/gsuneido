// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	return Plugins().Contributions("ImportExport",
		'import_formats').Map!({|x| x.name })
	}