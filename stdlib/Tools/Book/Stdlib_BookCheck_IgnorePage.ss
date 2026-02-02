// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
#(function (book, x)
	{
	return x.text is "" or BookContent.Match(book, x.text)
	}
)