// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.book = .MakeBook()
		QueryOutput(.book, [name: 'Menu1', order: 1, num: QueryMax(.book, 'num') + 1])
		QueryOutput(.book, [name: 'Menu2', order: 2, num: QueryMax(.book, 'num') + 1])
		.standalone = .MakeBookRecord(.book, 'Standalone', '/Menu1')
		.dup1 = .MakeBookRecord(.book, 'Dupe', '/Menu1')
		.dup2 = .MakeBookRecord(.book, 'Dupe', '/Menu2')
		}

	Test_main()
		{
		origUser = Suneido.User
		Suneido.User = 'default'
		// 0 matches
		Assert(QueryMatchBookOption.Func(.book, 'notExist', 'test') is: false)
		// 1 match
		rec = QueryMatchBookOption.Func(.book, 'Standalone', 'Menu1')
		Assert(rec.name is: .standalone.name)
		Assert(rec.path is: .standalone.path)

		// exact match
		rec2 = QueryMatchBookOption.Func(.book, 'Dupe', 'Menu1')
		Assert(rec2.name is: .dup1.name)
		Assert(rec2.path is: .dup1.path)

		// exact match
		rec3 = QueryMatchBookOption.Func(.book, 'Dupe', 'Menu2')
		Assert(rec3.name is: .dup2.name)
		Assert(rec3.path is: .dup2.path)

		// find best match
		rec4 = QueryMatchBookOption.Func(.book, 'Dupe', 'noMenu')
		Assert(rec4.name is: .dup1.name)
		Assert(rec4.path is: .dup1.path)

		.SpyOn(AccessPermissions).Return(false, true, false, false,
			'readOnly', true, 'readOnly', false, 'readOnly', 'readOnly')
		// find best match with permissions
		rec5 = QueryMatchBookOption.Func(.book, 'Dupe', 'noMenu')
		Assert(rec5.name is: .dup2.name)
		Assert(rec5.path is: .dup2.path)

		// no permissions
		rec6 = QueryMatchBookOption.Func(.book, 'Dupe', 'noMenu')
		Assert(rec6 is: false)

		// full permission gets priority over readOnly permission
		rec7 = QueryMatchBookOption.Func(.book, 'Dupe', 'noMenu')
		Assert(rec7.name is: .dup2.name)
		Assert(rec7.path is: .dup2.path)

		// only 1 readOnly permission
		rec8 = QueryMatchBookOption.Func(.book, 'Dupe', 'noMenu')
		Assert(rec8.name is: .dup1.name)
		Assert(rec8.path is: .dup1.path)

		// readOnly permission for both, first match gets priority
		rec9 = QueryMatchBookOption.Func(.book, 'Dupe', 'noMenu')
		Assert(rec9.name is: .dup1.name)
		Assert(rec9.path is: .dup1.path)

		Suneido.User = origUser
		}
	}