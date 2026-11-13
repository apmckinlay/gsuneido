// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
function(table)
	{
	// remove ServerEval < BuiltDate 20261030
	s = ServerEval("Database.Schema", table) // faster than querying columns
	columns = s.AfterFirst(' ').BeforeFirst(' ')
	columns = columns[1..-1].Split(',').Remove('-')
	return columns
	}
