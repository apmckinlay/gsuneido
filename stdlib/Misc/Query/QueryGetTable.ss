// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
// tableHint usage: place before the base table of the query
// 		/* tableHint: */ tables join by(table) columns)
// when the base of the query is a view, need to specify the base table of the view
// 		/* tableHint: tables */ view_table_indexes
class
	{
	CallClass(query, nothrow = false, orview = false)
		{
		table = .getTableName(query)
		if table isnt false and (TableExists?(table) or (orview and ViewExists?(table)))
			return table
		return .handleThrow("QueryGetTable failed on: " $ query, nothrow)
		}

	getTableName(query)
		{
		if false isnt table = .findTableHint?(query)
			return table
		// add other complexity checks
		if .hasMultipleTables?(query)
			ProgrammerError("Complex Query requires table hint " $
				query.Ellipsis(100 /*=maxsize*/))
		return .firstTable(query)
		}

	findTableHint?(q)
		{
		table = false
		regex = `/\*\s+tableHint:\s+`
		scan = ScannerWithContext(q, wantComments:)
		while scan isnt (tok = scan.Next()) and table is false
			{
			if scan.Type() is #COMMENT and tok =~ regex
				if false is table = tok.Extract(regex $ '(\w+)')
					table = scan.Ahead()
			}
		return table
		}

	hasMultipleTables?(query)
		{
		for tok in (scan = QueryScanner(query))
			if scan.Keyword?() and
				tok in ('join','union','times','minus','leftjoin','intersect')
				return true
		return false
		}

	firstTable(query)
		{
		for tok in (scan = QueryScanner(query))
			if scan.Keyword?() is false and scan.Type() is 'IDENTIFIER' and
				(TableExists?(scan.Value()) or ViewExists?(tok))
				return tok
		return false
		}

	handleThrow(msg, nothrow)
		{
		if nothrow
			return ""
		throw msg
		}
	}