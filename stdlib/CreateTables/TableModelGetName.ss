// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Memoize
	{
	Func(table)
		{
		fn = Tables.GetTable(table)
		if fn is false
			return false
		return fn().Name
		}
	}