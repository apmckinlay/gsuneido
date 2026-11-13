// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// returns a list of tables that have the fields required of libraries
// SEE ALSO: BookTables
MemoizeSingle
	{
	Func()
		{
		query = "tables where HasLibraryColumns?(TableColumns(table))"
		return QueryList(query, #table).SortWith!(#Lower)
		}
	}
