// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(QuerySelectColumns('tables project table') is: #(table))

		Assert(QuerySelectColumns('columns project column
			/* CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE */
			sort column') is: #(column))

		view = .MakeView('(tables union (views rename view_name to table)
			project table)')
		Assert(QuerySelectColumns(view $
			'/* CHECKQUERY SUPPRESS: UNION NOT DISJOINT */') is: #(table))

		Assert(QuerySelectColumns('(columns project column)
			/* CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE */')
			is: #(column))
		}

	Test_customPermissions()
		{
		.SpyOn(Customizable.GetNonPermissableFields).Return(#(table))
		Assert(QuerySelectColumns('tables project table') is: #())
		}

	Test_deletedCustomField()
		{
		table = .MakeTable('(keyColumn, column1, column2, column3, custom_555599)
			key (keyColumn)')

		columns = QuerySelectColumns(table)
		Assert(columns isSize: 5)
		Assert(columns has: #keyColumn)
		Assert(columns has: #column1)
		Assert(columns has: #column2)
		Assert(columns has: #column3)
		Assert(columns has: #custom_555599)

		.SpyOn(Datadict).Return(#(Internal:))

		columns = QuerySelectColumns(table)
		Assert(columns isSize: 4)
		Assert(columns has: #keyColumn)
		Assert(columns has: #column1)
		Assert(columns has: #column2)
		Assert(columns has: #column3)
		Assert(columns hasnt: #custom_555599)
		}
	}