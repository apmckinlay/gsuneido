// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
function (@args)
	{
	if not Suneido.Member?('Query1Cache')
		Suneido.Query1Cache = LruCache(Query1, 200)
	return Suneido.Query1Cache.Get(@args)
	}
