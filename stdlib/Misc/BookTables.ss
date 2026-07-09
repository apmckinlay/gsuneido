// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// returns a list of tables that have the fields required of books
// SEE ALSO: LibraryTables
MemoizeSingle
	{
	Func()
		{
		query = 'columns
					summarize table, list column
					where HasBookColumns?(list_column)'
		return QueryList(query, 'table').SortWith!(#Lower)
		}
	}
