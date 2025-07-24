// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
// checks that all the parents exist
function (book)
	{
	check = function (book, path)
		{
		result = true
		orig_path = path
		while path isnt ''
			{
			name = path.AfterLast('/')
			path = path.BeforeLast('/')
			if QueryEmpty?(book, :path, :name)
				{
				Print("ERROR: bad path: " $ orig_path)
				result = false
				break
				}
			}
		result
		}
	cache = LruCache(check)
	QueryApply(book $ ' sort path')
		{ |x|
		cache.Get(book, x.path)
		}
	}