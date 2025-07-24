// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		table = .MakeTable("(a,b,c,d) key(a,b)",
			rec = [a: 1, b: 2, c: 3, d: 4])
		Assert(QueryOutputIfNew(table, rec) is: false)
		rec = [a: 1, b: 22, c: 3, d: 4]
		Assert(QueryOutputIfNew(table, rec) is: true)
		}
	Test_empty_key()
		{
		table = .MakeTable("(a,b,c) key()")
		Assert(QueryOutputIfNew(table, []) is: true)
		Assert(QueryOutputIfNew(table, []) is: false)
		}
	}