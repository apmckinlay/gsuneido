// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.MakeLibraryRecord([name: "Rule_neg", text: "function () { return -.k }"])
		}
	test(query, expected)
		{
		// specify limit to avoid sorting in memory
		Assert(QueryAll(query, limit: 99) is: expected)
		}
	Test_single()
		{
		tbl = .MakeTable("(k, Neg) key(k)", [k: 1], [k: 2])
		.test(tbl, [[k: 1], [k: 2]])
		.test(tbl $ " sort k", [[k: 1], [k: 2]])
		.test(tbl $ " sort neg", [[k: 2], [k: 1]])
		}
	Test_multi()
		{
		tbl = .MakeTable("(k) key(k)", [k: 1], [k: 2])
		q = tbl $ " extend neg"
		.test(q, [[k: 1], [k: 2]])
		.test(q $ " sort k", [[k: 1], [k: 2]])
		.test(q $ " sort neg", [[k: 2], [k: 1]])

		tbl2 = .MakeTable("(k,x,Neg) key(k)", [k: 1, x: 11], [k: 2, x: 22])
		q = tbl $ " join by(k) " $ tbl2
		.test(q, [[k: 1, x: 11], [k: 2, x: 22]])
		.test(q $ " sort k", [[k: 1, x: 11], [k: 2, x: 22]])
		.test(q $ " sort neg", [[k: 2, x: 22], [k: 1, x: 11]])

		q = tbl $ " union " $ tbl2 $ '/* CHECKQUERY SUPPRESS: UNION NOT DISJOINT */'
		.test(q, [[k: 1], [k: 1, x: 11], [k: 2], [k: 2, x: 22]])
		.test(q $ " sort k, x", [[k: 1], [k: 1, x: 11], [k: 2], [k: 2, x: 22]])
		.test(q $ " sort neg, x", [[k: 2], [k: 2, x: 22], [k: 1], [k: 1, x: 11]])
		}
	Test_extend()
		{
		tbl = .MakeTable("(k) key(k)", [k: 1], [k: 2])
		q = tbl $ " extend neg = -k"
		.test(q, [[k: 1, neg: -1], [k: 2, neg: -2]])
		.test(q $ " sort k", [[k: 1, neg: -1], [k: 2, neg: -2]])
		.test(q $ " sort neg", [[k: 2, neg: -2], [k: 1, neg: -1]])
		}
	}