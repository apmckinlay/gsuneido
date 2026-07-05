// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Memoize
	{
	Func(includeViews? = false)
		{
		tables = QueryList('tables', 'table')
		if includeViews?
			tables.Add(@QueryList('views', 'view_name'))
		return tables.Sort!()
		}
	}