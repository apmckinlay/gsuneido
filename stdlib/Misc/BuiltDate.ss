// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	if not Suneido.Member?(#BuiltDate)
		Suneido.BuiltDate = Date(Built().BeforeFirst(' ('))
	return Suneido.BuiltDate
	}