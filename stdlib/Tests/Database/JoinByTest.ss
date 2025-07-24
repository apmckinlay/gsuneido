// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		t1 = .MakeTable('(a, b, c) key(a)')
		t2 = .MakeTable('(b, c, d) key(b, c)')
		Assert(QueryStrategy(t1 $ ' join by (b,c)' $ t2) has: "by(b,c)")
		QueryStrategy(t1 $ ' join by (b,c)' $ t2)
		Assert({ QueryStrategy(t1 $ ' join by (a,b,c)' $ t2) }
			throws: 'join: by does not match')
		Assert({ QueryStrategy(t1 $ ' join by (b)' $ t2) }
			throws: 'join: by does not match')
		QueryStrategy(t1 $ ' join by (c,b)' $ t2)
		Assert({ QueryStrategy(t1 $ ' join by ()' $ t2) }
			throws: 'syntax error')
		}
	}
