// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	cl: Reporter_table
		{
		fakeTables: #(
			fullValid: #(BookLocation: '/Book/Screen', Permission: 'fullAuth'),
			valuesInSource: #(BookLocation: '/Book/Table', Permission: 'fromTable'),
			valuesNotInTable: #()
			noTable: false
			noBookLocation: #()
			)
		Reporter_table_getTable(source)
			{
			return .fakeTables[source.table]
			}
		}
	fullValid: #(table: 'fullValid')
	valuesInSource: #(table: 'valuesInSource', bookLocation: '/book/source',
		auth: 'fromSource')
	valuesNotInTable: #(table: 'valuesNotInTable', bookLocation: '/book/notFromTable',
		auth: 'notTheTable')
	noTable: #(table: 'noTable', name: 'No Table Found')
	tableValuesNoBookLocation: #(table: 'noBookLocation', auth: 'fromSource')
	Test_BookLocation()
		{
		Assert(.cl.BookLocation(.fullValid) is: '/Book/Screen')
		Assert(.cl.BookLocation(.valuesInSource) is: '/book/source')
		Assert(.cl.BookLocation(.valuesNotInTable) is: '/book/notFromTable')
		Assert(.cl.BookLocation(.noTable) is: false)
		Assert(.cl.BookLocation(.tableValuesNoBookLocation) is: false)
		}

	Test_Authorization()
		{
		Assert(.cl.Authorization(.fullValid) is: 'fullAuth')
		Assert(.cl.Authorization(.valuesInSource) is: 'fromSource')
		Assert(.cl.Authorization(.valuesNotInTable) is: 'notTheTable')
		Assert(.cl.Authorization(.noTable) is: 'No Table Found')
		}
	}