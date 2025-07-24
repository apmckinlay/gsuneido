// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
Memoize
	{
	Expires?: true
	Func(query)
		{
		return .getColumns(query)
		}

	getColumns(query)
		{
		query = QuerySuppress(QueryStripSort(query))
		try
			return WithQuery(query)
				{|q|
				q.Columns()
				}
		catch (unused, '*nonexistent table')
			return #()
		}
	}
