// Copyright (C) 2013 Axon Development Corporation All rights reserved worldwide.
Test
	{
	Test_single_table()
		{
		table = .MakeTable('(field) key()', [field: "TEST"])
		x = Config(table)
		Assert(x.field is: 'TEST')

		// test invalidate
		Config.Invalidate(table)
		Assert(not Object?(Suneido.Config[table]))
		}

	Test_multi_query()
		{
		table = .MakeTable('(field) key(field)', [field: "ONE"], [field: "TWO"])

		Assert({ Config(table) } throws: "not unique")
		x = Config(table $ ' where field is "ONE"')
		Assert(x.field is: 'ONE')
		x = Config(table $ ' where field is "TWO"')
		Assert(x.field is: 'TWO')
		x = Config(table $ ' where field is "ONE"')
		Assert(x.field is: 'ONE')

		// test invalidate
		Config.Invalidate(table $ ' where field is "ONE"')
		Assert(not Object?(Suneido.Config[table $ ' where field is "ONE"']))
		}

	Test_OverrideRestore()
		{
		// test override with no members
		table = .MakeTable('(a, b, c) key()')
		Assert(Config(table) members: #("asof"))
		Config.Override(table, Object())
		Assert(Config(table) members: #("asof", "asof_override"))
		Assert(Config.ConfigCached?('config_original_' $ table))
		Config.Restore(table)
		Assert(Config(table) members: #("asof"))
		Assert(Config.ConfigCached?('config_original_' $ table) is: false)

		// test with members when table already exists in Suneido variable
		Config.Invalidate(table)
		QueryOutput(table, Record(a: 0, b: 0, c: 0))
		test_config = Config(table)
		Assert(test_config.a is: 0)
		Assert(test_config.b is: 0)
		Assert(test_config.c is: 0)
		Config.Override(table, Object(a: 1, b: 2, c: 3))
		test_config = Config(table)
		Assert(test_config.a is: 1)
		Assert(test_config.b is: 2)
		Assert(test_config.c is: 3)

		// test with members when table doesn't already exist in Suenido variable
		Config.Invalidate(table)
		Config.Override(table, Object(a: 1, b: 2, c: 3))
		test_config = Config(table)
		Assert(test_config.a is: 1)
		Assert(test_config.b is: 2)
		Assert(test_config.c is: 3)

		Config.Restore(table)
		test_config = Config(table)
		Assert(test_config.a is: 0)
		Assert(test_config.b is: 0)
		Assert(test_config.c is: 0)

		// test override with members twice
		test_config = Config(table)
		Config.Override(table, Object(a: 1, b: 2, c: 3))
		Config.Override(table, Object(a: 4, b: 5, d: 6))
		test_config = Config(table)
		Assert(test_config.a is: 4)
		Assert(test_config.b is: 5)
		Assert(test_config.c is: 3)
		Assert(test_config.d is: 6)

		Config.Restore(table)
		test_config = Config(table)
		Assert(test_config.a is: 0)
		Assert(test_config.b is: 0)
		Assert(test_config.c is: 0)
		Assert(test_config hasntMember: 'd')
		}

	Test_Cache()
		{
		table = .MakeTable('(a) key()', [a: 1])
		Assert(Config.ConfigCached?(table) is: false)
		test_config = Config(table)
		Assert(Config.ConfigCached?(table))

		// test that GetCachedConfig does not reset asof
		test_config2 = Config.GetCachedConfig(table)
		Assert(test_config.asof is: test_config2.asof)

		// test changing value in the cache
		Suneido.Config[table].a = 2
		// asof is not expired; should just use what is in Suneido
		Assert(Config(table).a is: 2)

		// set asof to be expired
		Suneido.Config[table].asof = Date().Plus(seconds: -11)
		Assert(Config.ConfigCached?(table))
		Assert(Config.GetCachedConfig(table).a is: 2)
		// asof is expired; should requery
		Assert(Config(table).a is: 1)
		}
	}