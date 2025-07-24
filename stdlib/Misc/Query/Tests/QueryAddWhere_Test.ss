// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(QueryAddWhere("test", "where x > 5")
			sameText: "test where x > 5 ")
		Assert(QueryAddWhere("test sort name", "where x > 5")
			sameText: BuiltDate() < #20250422
				? "test where x > 5 sort name"
				: "test sort name where x > 5")
		Assert(QueryAddWhere("test where s = 'sort' sort test1", "where x = 123")
			sameText: BuiltDate() < #20250422
				? "test where s = 'sort' where x = 123 sort test1"
				: "test where s = 'sort' sort test1 where x = 123")
		}
	}