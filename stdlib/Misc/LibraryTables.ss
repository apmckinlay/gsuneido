// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// returns a list of tables that have the fields required of libraries
// SEE ALSO: BookTables
MemoizeSingle
	{
	Func()
		{
		query = 'columns
					summarize table, list column
					where HasLibraryColumns?(list_column)'
		return QueryList(query, 'table').SortWith!(#Lower)
		}
	}