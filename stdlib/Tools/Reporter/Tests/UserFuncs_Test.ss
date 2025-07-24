// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		lib = .MakeLibraryRecord()
		uf = UserFuncs('UserFuncs_Test', 2, lib)

		name = uf.NeedFunc('123' $ '/**/'.Repeat(20))
		Assert(name is: "UserFuncs_Test000001")
		Assert(Global(name) is: 123)

		name = uf.NeedFunc('456')
		Assert(name is: "UserFuncs_Test000002")
		Assert(Global(name) is: 456)

		name = uf.NeedFunc('123' $ '/**/'.Repeat(20))
		Assert(name is: "UserFuncs_Test000001")
		Assert(Global(name) is: 123)

		name = uf.NeedFunc('789')
		Assert(name is: "UserFuncs_Test000003")
		Assert(Global(name) is: 789)

		query = lib $ ' where name > "UserFuncs_Test" and name < "UserFuncs_Test999999"'
		Assert(QueryCount(query) is: 2)
		}

	// This sets up the scenario where a record that has been read already (name1/name2)
	// gets deleted by the next lookup (name3/name4), and is then unavailable when
	// it gets used -
	Test_deletion()
		{
		lib = .MakeLibraryRecord()
		uf = UserFuncs('UserFuncs_DeleteTest', 35, lib)
		uf2 = UserFuncs('UserFuncs_BDeleteTest', n: 3, :lib)

		func1 = 'function () { return 1 }'
		func2 = 'function () { return 2 }'
		func3 = 'function () { return 3 }'
		name1 = uf.NeedFunc(func1)
		uf2.NeedFunc(func1) // insterting other records that do not match prefix
		// DON'T Global 'name1' here, that will put it in memory and the code below will
		// not fail corretly
		for i in ..34
			uf.NeedFunc('/*this is filler ' $ Display(i) $ '*/')

		uf2.NeedFunc(func2) // insterting other records that do not match prefix

		name2 = uf.NeedFunc(func1)
		// record should have been re-output as it was not at the top of the stack
		Assert(name1 isnt: name2)
		name3 = uf.NeedFunc(func2)
		name4 = uf.NeedFunc(func3)
		Assert(Global(name2)() is: 1)
		Assert(Global(name3)() is: 2)
		Assert(Global(name4)() is: 3)

		Assert(QueryCount(lib $ ' where name =~ "^UserFuncs_DeleteTest"') is: 35)
		}

	Test_max()
		{
		lib = .MakeLibraryRecord()
		uf = UserFuncs('UserFuncs_MaxTest', n: 3, :lib)
		uf2 = UserFuncs('UserFuncs_BMaxTest', n: 3, :lib)

		func1 = 'function() { return "one" }'
		name = uf.NeedFunc(func1)
		QueryDo('update ' $ lib $ ' where name is ' $ Display(name) $
			' set name = "UserFuncs_MaxTest999998"')

		name2 = uf2.NeedFunc(func1)
		Assert(name2 is: 'UserFuncs_BMaxTest000001')
		QueryDo('update ' $ lib $ ' where name is ' $ Display(name2) $
			' set name = "UserFuncs_BMaxTest999998"')

		func2 = 'function() { return "two" }'
		name = uf.NeedFunc(func2)
		Assert(name is: 'UserFuncs_MaxTest000001')

		name2 = uf2.NeedFunc(func2)
		Assert(name2 is: 'UserFuncs_BMaxTest000001')

		func3 = 'function() { return "three" }'
		name = uf.NeedFunc(func3)
		Assert(name is: 'UserFuncs_MaxTest000002')

		func4 = 'function() { return "four" }'
		name = uf.NeedFunc(func4)
		Assert(name is: 'UserFuncs_MaxTest000003')

		// orig should be removed
		Assert(QueryEmpty?(lib, name: "UserFuncs_MaxTest999998"),
			msg: '999998 not recycled')

		Assert(not QueryEmpty?(lib, name: "UserFuncs_BMaxTest999998"),
			msg: '999998 B exists')
		Assert(not QueryEmpty?(lib, name: name2), msg: 'uf2 rec should not delete')
		}
	}