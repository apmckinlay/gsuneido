// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		fn = SvcDisabledTables.SvcDisabledTables_disabled

		table = .MakeTable('(one, two) key(one)')
		Assert(fn(table), msg: 'table')

		newlib = .MakeLibrary()
		Assert(fn(newlib), msg: 'newlib no svc columns')

		Database('ensure ' $ newlib $
			' (lib_before_hash, lib_before_text, lib_before_path)')
		Assert(fn(newlib) is: false, msg: 'newlib has svc columns')
		}
	}