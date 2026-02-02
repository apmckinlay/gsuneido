// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		tbl = .MakeTable("(a,b,c,d,e,f) key(a) index(d) key(b,c) index(e,f)")
		test = {|@args| {|expected|
			Assert(Query.Strategy1(@args), like: expected.Replace("tbl", tbl)) } }

		// Query1/Empty?
		test(tbl)("no select: tbl^(a)")
		test(tbl, a: 1)("key: tbl^(a)")
		test(tbl, b: 1, c: 2, d: 3)("key: tbl^(b,c)") // key takes priority
		test(tbl, d: 1)("just index: tbl^(d)")
		test(tbl, d: 1, e: 2)("multiple indexes: tbl (d) (e,f)")

		// QueryFirst/Last
		test(tbl $ " sort e")(
			"{0} tbl^(e,f)
			[nrecs~ 0 cost~ 0]")
		test(tbl $ " where a = 1")(
			"{0.000x 0 0+500} tbl^(a)
			{0 0+500} where*1 a is 1
			[nrecs~ 0 cost~ 500]")

		tbl = .MakeTable("(a,b,c) key()")
		test(tbl)("no select: tbl^()")
		test(tbl, a: 1)("key: tbl^()")
		}
	}