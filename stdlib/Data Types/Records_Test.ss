// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
// this is tests for methods defined in Records
// RecordTest is for built-in methods
Test
	{
	Test_Project()
		{
		Assert(Type([].Project(#(a))) is 'Record')
		Assert([a: 1, b: 2, c: 3].Project(#(a, c, d)) is: [a: 1, c: 3])
		Assert([a: 1, b: 2, c: 3].Project(#a, #c, #d) is: [a: 1, c: 3])
		}
	Test_queries()
		{
		.test_queries(Record())
		Transaction(read:)
			{|t|
			r = t.QueryFirst("stdlib sort num") // any record from the transaction
			.test_queries(r)
			}
		}
	test_queries(r)
		{
		args = #(stdlib, name: Records_Test)
		for m in #(Query1, QueryEmpty?)
			Assert(r[m](@args) is: Global(m)(@args))
		args = #("stdlib sort num", name: Records_Test)
		for m in #(QueryFirst, QueryLast)
			Assert(r[m](@args) is: Global(m)(@args))
		}
	}