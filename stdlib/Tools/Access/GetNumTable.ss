// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
// finds the first table that has num/name/abbrev with key on num
Memoize
	{
	Func(base_name)
		{
		cols = #(_num, _name, _abbrev).Map({ base_name $ it })
		x = QueryFirst("columns
			where column in " $ Display(cols)[1 ..] $ "
			summarize table, count
			where count = 3 sort table")
		return x is false or
			QueryEmpty?("indexes", table: x.table, key: true, columns: cols[0])
			? false : x.table
		}
	}
