// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		try Database("destroy dberrorstest")
		Database("create dberrorstest (a, b, c) key(a)")
		}
	admin:
		(
		("create dberrorstest (a, b, c) key(a)"
			"exist")
		("create dberrorstest2 (a b c)"
			"key required")

		("alter dberrorstest2 create (a b c)"
			"nonexistent table")
		("alter dberrorstest create (c)"
			"existing column")
		("alter dberrorstest rename y to z"
			"nonexistent column")
		("alter dberrorstest drop (nonex)"
			"nonexistent column")
		("alter tables rename table to tablenum"
			"system")

		("rename dberrorstest2 to abc"
			"nonexistent table")
		("rename tables to abc"
			"system")
		("rename dberrorstest to stdlib"
			"exist")

		("destroy tables"
			"can't")
		("destroy dberrorstest2"
			"nonexistent table")
		)
	Test_admin()
		{
		for (x in .admin)
			Assert({ Database(x[0]) } throws: x[1])
		}
	query:
		(
		("dberrorstest2",
			"nonexistent table")

		("dberrorstest project x,y,z"
			"nonexistent column")

		("dberrorstest rename w to x, y to z"
			"nonexistent column")
		)
	Test_query()
		{
		for (x in .query)
			Assert({ Cursor(x[0]) } throws: x[1])
		}
	Teardown()
		{
		Database("destroy dberrorstest")
		}
	}
