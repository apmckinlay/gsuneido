// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Memoize
	{
	Func(book, option)
		{
		QueryFirst(BookOptionQuery(book, option) $ ' sort path')
		}
	}