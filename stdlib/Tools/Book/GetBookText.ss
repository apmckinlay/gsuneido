// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
Memoize
	{
	Func(name, book = 'imagebook')
		{
		if Number?(name)
			{
			x = Query1(book, num: name)
			SuneidoLog("ERROR: GetBookText with number")
			}
		else
			{
			path = '/res'
			if name.Has?('/')
				{
				path $= '/' $ name.BeforeLast('/')
				name = name.AfterLast('/')
				}
			x = Query1(book, :path, :name)
			}
		return x is false ? false : x.text
		}
	}