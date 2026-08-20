// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ForEachTable()
		{
		if not LibraryTags.GetTagsInUse().Has?('__trial')
			return

		.MakeLibraryRecord(
			// table 1 has no __trial override
			[name: 'Table_TablesTest1', text: 'class {
				RenamedForCustomize?: false
				Name: "Table_TablesTest1"}', group: -1],

			// table 2 has __trial override
			[name: 'Table_TablesTest2', text: 'class {
				RenamedForCustomize?: false
				Name: "Table_TablesTest2"}', group: -1],
			[name: 'Table_TablesTest2__trial', text: 'class {
				RenamedForCustomize?: false
				Name: "Table_TablesTest2__trial"}', group: -1],

			// table 3 has only __trial override
			[name: 'Table_TablesTest3__trial', text: 'class {
				RenamedForCustomize?: false
				Name: "Table_TablesTest3__trial"}', group: -1])

		cl = Tables
			{
			Tables_libraries()
				{
				return Object(Test.TestLibName())
				}
			}
		tables = Object()
		cl.ForEachTable({ tables.Add(it.Name) })
		Assert(tables equalsSet: #(Table_TablesTest1,
			Table_TablesTest2__trial, Table_TablesTest3__trial))
		}
	}