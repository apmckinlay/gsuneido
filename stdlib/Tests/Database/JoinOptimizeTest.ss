// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.big = .MakeTable('(a, b, c) key(a) index(b)')
		for i in .. 100
			QueryOutput(.big,
				Record(a: i, b: i % 7, c: "helloworld".Repeat(10)))

		.small_one = .MakeTable('(a, d) key(a)')
		n = 0
		for i in .. 10
			QueryOutput(.small_one,
				Record(a: n += 13, d: "helloworld".Repeat(10)))

		.small_n = .MakeTable('(a, e, f) index(a) key(e)')
		n = 0
		for i in .. 10
			QueryOutput(.small_n,
				Record(a: n += 9, e: n, f: "helloworld".Repeat(10)))
		}
	Test_one_one()
		{
		query = .big $ ' join by(a) ' $ .small_one
		Assert(QueryStrategy(query).Tr('()').Replace("bya", "on a")
			is: .small_one $ '^a join 1:1 on a ' $ .big $ '^a')
		}
	Test_one_n()
		{
		query = .big $ ' join by(a) ' $ .small_n
		Assert(QueryStrategy(query).Tr('()').Replace("bya", "on a")
			is: .small_n $ '^a join n:1 on a ' $ .big $ '^a')
		}
	}
