// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		fn = LibRecGetPath

		testLib = .MakeLibrary([name: 'TestFolder1', group: 0, num: 33301, parent: 0],
			[name: 'TestRecord0', group: -1, num: 33300, parent: 0],
			[name: 'TestRecord1', group: -1, num: 33302, parent: 33301],
			[name: 'TestRecord2', group: -1, num: 33303, parent: 33301],
			[name: 'TestFolder2', group: 33301, num: 33304, parent: 33301],
			[name: 'TestRecord3', group: -1, num: 33305, parent: 33304])

		Assert(fn([name: ''], testLib) is: testLib $ '/')

		Assert(fn([name: 'TestRecord0', group: -1, num: 33300, parent: 0], testLib) is:
			testLib $ '/TestRecord0')

		Assert(fn([name: 'TestRecord1', group: -1, num: 33302, parent: 33301], testLib)
			is: testLib $ '/TestFolder1/TestRecord1')
		Assert(fn([name: 'TestRecord2', group: -1, num: 33303, parent: 33301], testLib)
			is: testLib $ '/TestFolder1/TestRecord2')

		Assert(fn([name: 'TestRecord3', group: -1, num: 33305, parent: 33304], testLib)
			is: testLib $ '/TestFolder1/TestFolder2/TestRecord3')
		}
	}