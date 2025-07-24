// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
// TAGS: !client
Test
	{
	delete_tables()
		{
		QueryDo('delete ' $ .table1)
		QueryDo('delete ' $ .table2)
		QueryDo('delete ' $ .table3)
		QueryDo('delete ' $ .table4)
		}
	Test_main()
		{
		if Sys.Client?()
			return
		.table1 = .MakeTable("(a,b,c) key(a)")
		.table2 = .MakeTable("(a,b,c) key(a)")
		.table3 = .MakeTable("(a,b,c) key(a)")
		.table4 = .MakeTable("(a,b,c) key(a)")

		.MakeLibraryRecord(
			Object(
				name: 'Trigger_' $ .table1,
				text: "function (t, oldrec, newrec)
					{
					if oldrec is false
						QueryOutput(" $ Display(.table3) $ ", Record(a:1, b:2, c:3))
					}")
			Object(
				name: 'Trigger_' $ .table2,
				text: "function (t, oldrec, newrec)
					{
					if oldrec is false
						QueryOutput(" $ Display(.table4) $ ", Record(a:1, b:2, c:3))
					}"))

		// make sure triggers are working
		QueryOutput(.table1, Record(a: 'a', b: 'b', c: 'c'))
		Assert(Query1(.table3) isnt: false)
		QueryOutput(.table2, Record(a: 'a', b: 'b', c: 'c'))
		Assert(Query1(.table4) isnt: false)
		.delete_tables()

		// test with empty object, trigger should still work
		DoWithoutTriggers(#())
			{
			QueryOutput(.table1, Record(a: 'a', b: 'b', c: 'c'))
			}
		Assert(Query1(.table3) isnt: false)
		.delete_tables()

		// test with 1 table
		DoWithoutTriggers(Object(.table1))
			{
			QueryOutput(.table1, Record(a: 'a', b: 'b', c: 'c'))
			}
		Assert(Query1(.table3) is: false)
		.delete_tables()

		// test with 2 tables
		DoWithoutTriggers(Object(.table1, .table2))
			{
			QueryOutput(.table1, Record(a: 'a', b: 'b', c: 'c'))
			QueryOutput(.table2, Record(a: 'a', b: 'b', c: 'c'))
			}
		Assert(Query1(.table3) is: false)
		Assert(Query1(.table4) is: false)
		.delete_tables()
		}
	}