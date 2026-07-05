// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_checkArgs()
		{
		m = AccessMarkTab.AccessMarkTab_checkArgs

		// invalid args
		Assert({ m(#()) } throws: 'missing argument')
		Assert({ m(#(access: 'access')) } throws: 'missing argument')

		// valid args
		m(#(access: 'access', plugin: 'plugin'))
		}

	Test_tabHasData?()
		{
		m = AccessMarkTab.AccessMarkTab_tabHasData?
		plugin = Object(fields: #(test_a,test_b,test_c))
		access = class { GetData() { return Record() } }
		Assert(m(access, plugin) is: false)
		access = class { GetData() { return Record(test_d: 'd') } }
		Assert(m(access, plugin) is: false)
		access = class { GetData() { return Record(test_b: false) } }
		Assert(m(access, plugin) is: false)
		access = class { GetData() { return Record(test_c: true) } }
		Assert(m(access, plugin))
		access = class { GetData() { return Record(test_b: 'abc') } }
		Assert(m(access, plugin))
		access = class { GetData()
			{ return Record(test_a: 123, test_b: 'abc', test_c: true) } }
		Assert(m(access, plugin))
		}
	}