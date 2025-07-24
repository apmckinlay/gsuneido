// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		gnt = GetNumTable.Func // bypass caching
		uniquePrefix = .TempTableName()
		Assert(gnt(uniquePrefix) is: false) // no table

		table = .MakeTable(
			"(k, " $ uniquePrefix $ "_num, " $ uniquePrefix $ '_name) key(k)')
		Assert(gnt(uniquePrefix) is: false) // no _abbrev

		Database("alter " $ table $ " create (" $ uniquePrefix $ "_abbrev)")
		Assert(gnt(uniquePrefix) is: false) // no key on num

		Database("alter " $ table $ " create key(" $ uniquePrefix $ "_num)")
		Assert(gnt(uniquePrefix) is: table)
		}
	}