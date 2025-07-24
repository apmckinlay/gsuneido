// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		table = .MakeTable('(a,b) key(a)')
		for i in ..20
			QueryOutput(table, [a: i, b: i + 100])

		fn = QueryListLimited

		Assert(fn(table, 'a', 20).Sort!() is: Seq(20))
		Assert(fn(table $ ' where a > 100 sort a', 'a', 20) is: #())
		Assert(fn(table, 'a', 19) is: false)
		Assert(fn(table, 'a', 19, returnDetailIfOverLimit?:)
			is: Object(count: 20, list: Seq(19)))
		Assert(fn(table, 'a', 1, returnDetailIfOverLimit?:)
			is: Object(count: 20, list: #(0)))
		}
	}
