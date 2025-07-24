// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(model)
		{
		.model = model
		.lru = LruCache(.Getfn, cache_size: 50)
		}
	Getrecord(recnum)
		{
		return .lru.Get(recnum)
		}
	Getitem(recnum, field)
		{
		x = .lru.Get(recnum)
		if (x is false or x is Object())
			return ""
		return x[field]
		}
	Getfn(i)
		{
		return .model.Getrecord(i)
		}
	Default(@args)
		{
		return .model[args[0]](@+1args)
		}
	}
