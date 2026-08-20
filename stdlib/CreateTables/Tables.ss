// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass()
		{
		tables = Object()
		.ForEachTable()
			{
			tables.Add(it)
			}
		return tables
		}

	ForEachTable(block)
		{
		for lib in .libraries()
			{
			tables = Object()
			for table in QueryList(lib $ ' where name > "Table_" and name < "Table_~"
				and group is -1', 'name')
				tables.AddUnique(LibraryTags.RemoveTagFromName(table))
			for table in tables
				{
				cl = Global(table)
				if cl.RenamedForCustomize? is false
					block(cl)
				}
			}
		}

	libraries()
		{
		return Libraries()
		}

	List()
		{
		tables = Object()
		.ForEachTable()
			{
			tables.Add(it.Table)
			}
		return tables
		}

	GetTable(table, member = false)
		{
		try
			{
			cl = Global('Table_' $ table)
			return member isnt false ? cl[member] : cl
			}
		catch (unused, "can't find|member not found")
			return false
		}

	RemoveColumns(table, columns)
		{
		if not TableExists?(table)
			return
		QueryColumns.ResetCache()
		removeCols = columns.Intersect(QueryColumns(table))
		if not removeCols.Empty?()
			Database('alter ' $ table $ ' drop (' $ removeCols.Join(',') $ ')')
		}
	}