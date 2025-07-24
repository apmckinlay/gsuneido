// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_duplicate_key()
		{
		table = .MakeTable('(a,b) key(a) key(b)')
		Assert(.nrows(table) is: 0)
		QueryOutput(table, #(a: 1, b: 2))
		Assert(.nrows(table) is: 1)
		Assert({ QueryOutput(table, #(a: 1, b: 3)) } throws: "duplicate key")
		Assert(.nrows(table) is: 1)
		Assert({ QueryOutput(table, #(a: 3, b: 2)) } throws: "duplicate key")
		Assert(.nrows(table) is: 1)
		}
	Test_trigger_exception()
		{
		if Sys.Client?()
			return
		table = .MakeTable('(a) key(a)')
		Assert(.nrows(table) is: 0, msg: "table should be empty initially")
		.MakeLibraryRecord(Record(name: 'Trigger_' $ table,
			text: 'function (t, oldrec, newrec) { x }'))
		Transaction(update:)
			{ |t|
			try t.QueryOutput(table, Record(a: 'a'))
			Assert(QueryEmpty?(table), msg: 'inside transaction')
			t.Rollback()
			}
		Assert(QueryEmpty?(table), msg: 'outside transaction')
		Assert(.nrows(table) is: 0,
			msg: "table should still be empty after aborting output tran")

		// output should succeed despite trigger exception
		// (if exception is caught and transaction completed)
		Transaction(update:)
			{ |t|
			try t.QueryOutput(table, Record(a: 'a'))
			}
		Assert(not QueryEmpty?(table), msg: 'record missing')
		Assert(.nrows(table) is: 1)
		}
	nrows(table)
		{
		return QueryCount(table)
		}
	}