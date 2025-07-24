// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
#(
function (book)
	{
	book = book.RemoveSuffix('Help')
	return ApplicationDir() $ '/' $ book $ '.ico'
	}
)