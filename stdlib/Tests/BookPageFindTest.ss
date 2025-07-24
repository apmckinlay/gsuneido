// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		table = .MakeTable('(name, path, text) key (name, path)')
		QueryOutput(table, Record(name: 'Configuration',
			path: '/' $ table $ '/Test'))
		page = BookPageFind(table, '/' $ table $ '/Test', 'Configuration')
		Assert(page isnt: false)
		Assert(page.name is: 'Configuration')
		Assert(page.path is: '/' $ table $ '/Test')

		page = BookPageFind(table, '/book/Test', 'Configuration')
		Assert(page isnt: false)
		Assert(page.name is: 'Configuration')
		Assert(page.path is: '/' $ table $ '/Test')

		QueryOutput(table, Record(name: 'Page',
			path: '/' $ table $ '/Test/Path'))
		page = BookPageFind(table, '/book/Test2/Path', 'Page')
		Assert(page isnt: false)
		Assert(page.name is: 'Page')
		Assert(page.path is: '/' $ table $ '/Test/Path')

		// test with Reference section in the path
		QueryOutput(table, Record(name: 'Equipment Statements',
			path: '/Trucking/Reference/Equipment and Driver Statements'))
		page = BookPageFind(table, '/Trucking/Equipment and Driver Statements',
			'Equipment Statements')
		Assert(page isnt: false)
		Assert(page.name is: 'Equipment Statements')
		Assert(page.path is: '/Trucking/Reference/Equipment and Driver Statements')
		}
	}