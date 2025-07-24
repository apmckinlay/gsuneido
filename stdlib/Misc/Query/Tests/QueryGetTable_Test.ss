// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		spy = .SpyOn(ProgrammerError).Return('')
		Assert(QueryGetTable("", nothrow:) is: "")
		Assert(QueryGetTable("stdlib") is: "stdlib")
		Assert(QueryGetTable("stdlib sort test_field") is: "stdlib")
		Assert(QueryGetTable("stdlib\nsort test_field") is: "stdlib")
		Assert(QueryGetTable("stdlib\r\n\t\tsort test_field") is: "stdlib")
		Assert(QueryGetTable(" \t
			stdlib\r\n\t\tsort test_field") is: "stdlib")

		Assert(QueryGetTable("(stdlib where num is 1)") is: "stdlib")
		Assert(QueryGetTable("(stdlib) where num is 1)") is: "stdlib")
		Assert(QueryGetTable("(stdlib
			where num is 1)") is: "stdlib")
		Assert(spy.CallLogs() isSize: 0)
		QueryGetTable(" (
			(stdlib
			where num is 1) join by(x) foo) where y is 1")
		Assert(spy.CallLogs() isSize: 1)

		Assert(QueryGetTable(" (
			( /* tableHint: */ stdlib
			where num is 1) join by(x) foo) where y is 1")
			is: "stdlib")

		Assert(QueryGetTable(" (
			( /* tableHint: stdlib */ bar
			where num is 1) join by(x) foo) where y is 1")
			is: "stdlib")

		QueryGetTable(" (
			( /* this is a comment */ stdlib
			where num is 1) join by(x) foo) where y is 1")
		Assert(spy.CallLogs() isSize: 2)

		Assert({ QueryGetTable(" (
			(stdlib
			where num is 1) join by(x) foo) where y is 1
			/* tableHint: */") }
			throws: 'QueryGetTable failed')

		Assert({ QueryGetTable(" (
			(bar
			where num is 1) join by(x) foo /* tableHint: */)
			where column is 'num'") }
			throws: 'QueryGetTable failed')

		Assert(QueryGetTable(" (
			(bar
			where num is 1) join by(x) foo /* tableHint: stdlib */)
			where column is 'num'")
			is: 'stdlib')

		Assert({ QueryGetTable("") } throws: "QueryGetTable failed")
		Assert(spy.CallLogs() isSize: 2)
		}
	Test_orView()
		{
		name = .MakeView("tables")
		Assert(QueryGetTable(name, nothrow:) is: "")
		Assert(QueryGetTable(name $ " where foo", orview:) is: name)
		Assert(QueryGetTable('/* tableHint: */ ' $ name $ " where foo", nothrow:) is: "")
		}
	}