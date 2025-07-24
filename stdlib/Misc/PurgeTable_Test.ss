// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		table = .MakeTable("(num, date) key (num)")
		// the default is 3 months
		// will create records before and after that point
		// Create one record a month
		.outputRecords(table, Date().Plus(minutes: -5), 10)
		Assert(QueryCount(table) is: 10)

		PurgeTable(table, 'date', Date().Plus(months: -12)) // shouldn't delete anything
		Assert(QueryCount(table) is: 10)

		PurgeTable(table, 'date', monthsToKeep: 6) // Delete older than 6 months
		Assert(QueryCount(table) is: 6)

		PurgeTable(table, 'date') // Delete older than 3 months (default)
		Assert(QueryCount(table) is: 3)

		PurgeTable(table, 'date', Date()) // delete everything
		Assert(QueryCount(table) is: 0)
		}

	outputRecords(table, date, count)
		{
		for i in .. count
			QueryOutput(table, Object(num: i, date: date.Plus(months: -i)))
		}
	}