// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_select()
		{
		tbl2 = .MakeTable("(b,s) key(b)",
			[b: 2, s: "two"], [b: 4, s: "four"])
		tbl = .MakeTable('(a, b, tbl2, CursorTranTest) key(a)',
			[a: 1, b: 2, tbl2: tbl2],  r2 = [a: 3, b: 4, tbl2: tbl2])

		.MakeLibraryRecord(
			[name: 'Rule_cursorTranTest',
			text: 'function () { .Transaction().Query1(.tbl2, b: .b).s }'])
		query = tbl $ ' where cursorTranTest = "four"'
		r = QueryFirst(query $ ' sort a')
		Assert(r is: r2)
		Cursor(query)
			{|c|
			Transaction(read:)
				{|t|
				r = c.Next(t)
				Assert(r is: r2)
				}
			}
		}
	}
