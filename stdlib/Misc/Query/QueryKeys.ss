// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
Memoize
	{
	Func(query)
		{
		return WithQuery(QuerySuppress(Query.StripSort(query)))
			{|q|
			q.Keys()
			}
		}
	}
