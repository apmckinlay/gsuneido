// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
// TAGS: !client
// Can't run on a client because clients get batches of timestamps
// and may not be in sequence with the server.
Test
	{
	Test_main()
		{
		.table = .MakeTable('(a, b_TS, c) key(a)')
		.test1(0)
			{ QueryOutput(.table, x = #(a: 0)) }
		Assert(x is: #(a: 0)) // readonly so _TS not added

		.test1(1)
			{ QueryOutput(.table, x = Record(a: 1)) }
		.test1(1)
			{ QueryDo('update ' $ .table $ ' set c = 2') }
		.test1(2)
			{ QueryDo('insert { a: 2 } into ' $ .table) }
		.test1(2)
			{
			Transaction(update:)
				{ |t|
				x = t.Query1(.table $ ' where a = ' $ 2)
				x.c = 'c'
				x.Update()
				}
			}
		}
	test1(a, block)
		{
		before = Timestamp()
		block()
		x = Query1(.table $ ' where a = ' $ a)
		Assert(x.b_TS isDate:)
		Assert(x.b_TS greaterThan: before)
		Assert(x.b_TS lessThan: Timestamp())
		}
	}