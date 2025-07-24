// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		table = .MakeTable('(a,b,c,d) key(a,b) key(c)')
		x = [a: 1, b: 2, c: 3, d: 4]
		QueryEnsure(table, x) // should output
		Assert(Query1(table) is: x)
		QueryEnsure(table, x) // same so do nothing
		Assert(Query1(table) is: x)
		x.d = 5
		QueryEnsure(table, x) // should delete/output
		Assert(Query1(table) is: x)
		}
	}