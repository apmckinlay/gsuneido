// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	book = Suneido.GetDefault(#CurrentBook, fallback = AccessPermissions.GetDefaultBook())
	return book isnt ''
		? book
		: fallback
	}
