// Copyright (C) 2010 Axon Development Corporation All rights reserved worldwide.
Test
	{
	Test_main()
		{
		.test('/', [])
		.test('c:/', [])
		.test('c:/a/b/c', ['c:/a', 'c:/a/b'])
		.test('c:/a/b/', ['c:/a', 'c:/a/b'])

		.test('//server/a/b/', ['//server/a', '//server/a/b'])
		.test('//server/', [])
		}
	test(path, expected)
		{
		for ..2
			{
			cds = new CreateDirectories
				{ CreateDirectories_ensureDir(dir) { .Log.Add(dir) } }
			cds.Log = []
			cds.CallClass(path)
			Assert(cds.Log is: expected)
			path = Paths.ToWindows(path) // test both ways
			}
		}
	}