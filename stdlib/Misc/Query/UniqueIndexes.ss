// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
Memoize
	{
	Func(table)
		{
		return QueryList("indexes where table is " $ Display(table) $ " and key is 'u'",
			'columns')
		}
	}