// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		super.Setup()
		tableName = .MakeTable("(col1, col2, col3, col4, col5, col6) key (col1)")
		.MakeLibraryRecord([name: "Table_" $ tableName,
			text: `TableModel
				{
				Table: '` $ tableName $ `'
				Columns: '(col1, col2, col3)'
				UserInterfaceExclude() { return #(col2) }
				ReporterExcludes() { return #(col1, col2) }
				}`])
		contribName = .TestLibName() $ "_Table_" $ tableName
		.MakeLibraryRecord([name: contribName,
			text: `#(
				Columns: '(col4, col5)'
				UserInterfaceExclude: #(
					"col6",
					function () { return #(col4) }),
				ReporterExcludes: #(
					"col6",
					function () { return #(col5) }
				))`])
		.model = Global("Table_" $ tableName)()
		}

	Test_GetDuplicateFieldInfo()
		{
		noDups = .TempName()
		.MakeLibraryRecord([name: noDups, text: `TableModel { }`])
		info = Global(noDups)().GetDuplicateFieldInfo()
		Assert(info is: false, msg: #nodups)

		dups = .TempName()
		.MakeLibraryRecord([name: dups,
			text: `TableModel {
				Table: testing
				DuplicateFieldInfo: function ()
					{
					return Object(
						fields: #('test', 'test2'),
						configTable: 'configuration',
						configField: 'checkdupfields')
					}}`])
		info = Global(dups)().GetDuplicateFieldInfo()
		sbe = Object(fields: #(test, test2),
			configTable: #configuration, configField: #checkdupfields)
		Assert(info equalsSet: sbe, msg: #dups)

		parent = .TempName()
		.MakeLibraryRecord(
			[name: "Table_" $ parent, text: `TableModel { ` $ `Table: ` $ parent $ `}`])
		contribDups = .TestLibName() $ "_Table_" $ parent
		.MakeLibraryRecord([name: contribDups,
			text: `#(
			DuplicateFieldInfo: #(function ()
				{
				return Object(
					fields: #('test', 'test2'),
					configTable: 'configuration',
					configField: 'checkdupfields')
				}))`])
		info = Global("Table_" $ parent)().GetDuplicateFieldInfo()
		Assert(info equalsSet: sbe, msg: "contrib dups")

		parent = .MakeTable("(test, test2, test3, test4, test5, custom_0, custom_1)
			key (test)")

		sbe = Object(fields: [#test, #test2, #custom_1, #test4],
				excludeFields: [#test3, #test5, #custom_0],
				configTable: #configuration, configField: #checkdupfields).
			Copy().
			Sort!()
		.MakeLibraryRecord(
			[name: name = "Table_" $ parent,
				text: `TableModel
			{
			Table: '` $ parent $ `'
			DuplicateFieldInfo: function ()
				{
				return Object(
					fields: #('test', 'test2'),
					excludeFields: #('test3', 'test5', 'custom_0')
					configTable: 'configuration',
					configField: 'checkdupfields')
				}
			}`], table: #stdlib)
		.AddTeardown({ QueryDo("delete stdlib where name is " $ Display(name)) })
		extra = .TestLibName() $ "_Table_" $ parent
		.MakeLibraryRecord([name: extra,
			text: `#(
			DuplicateFieldInfo: (function ()
				{
				return Object(
					fields: #('test3', 'test4'))
				}))`])
		info = Global("Table_" $ parent)().GetDuplicateFieldInfo()
		Assert(info equalsSet: sbe, msg: "contrib dups")
		}

	Test_GetIndexes()
		{
		tableName = .TempName()
		configTable = .MakeTable("(checkdupfields) key ()")
		QueryOutput(configTable, [checkdupfields: #(test1, test2, test3)])
		.MakeLibraryRecord(
			[name: name = "Table_" $ tableName,
				text: `TableModel
			{
			Table: '` $ name $ `'
			Columns: #(test1, test2, test3)
			Keys: #(test1)
			UniqueIndexes: #(test2)
			DuplicateFieldInfo: function ()
				{
				return Object(
					fields: #('test', 'test2', 'test3'),
					configTable: '` $ configTable $ `',
					configField: 'checkdupfields')
				}
			}`], table: #stdlib)
		.AddTeardown({ QueryDo("delete stdlib where name is " $ Display(name)) })
		Assert(Global(name)().GetIndexes() is: #(test3))
		}

	Test_UserInterfaceColumns()
		{
		userInterfaceCols = .model.UserInterfaceColumns()
		Assert(userInterfaceCols isnt: #())
		Assert(userInterfaceCols isSize: 3)
		Assert(userInterfaceCols equalsSet: #(col1, col3, col5))
		}

	Test_ReporterExcludeColumns()
		{
		reporterExcludes = .model.ReporterExcludeColumns()
		Assert(reporterExcludes isnt: #())
		Assert(reporterExcludes isSize: 4)
		Assert(reporterExcludes equalsSet: #(col1, col2, col5, col6))
		}

	Test_UserInterfaceExcludeColumns()
		{
		userInterfaceExcludes = .model.UserInterfaceExcludeColumns()
		Assert(userInterfaceExcludes isnt: #())
		Assert(userInterfaceExcludes isSize: 3)
		Assert(userInterfaceExcludes equalsSet: #(col2, col4, col6))
		}
	}
