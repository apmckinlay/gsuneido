// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Memoize
	{
	Func(name)
		{
		for lib in LibraryTables()
			if false isnt Query1(lib, group: -1, :name)
				return lib
		return false
		}
	}