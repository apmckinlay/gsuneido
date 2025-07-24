// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ValidFileName()
		{
		Assert(CheckDirectory.ValidFileName?(`test.txt`))
		Assert(CheckDirectory.ValidFileName?(`test`))
		Assert(CheckDirectory.ValidFileName?(`.txt`) is: false)
		Assert(CheckDirectory.ValidFileName?(`test | more`) is: false)
		Assert(CheckDirectory.ValidFileName?(`c:\work\test.txt`))
		Assert(CheckDirectory.ValidFileName?(`c:\work\test:again`) is: false)
		Assert(CheckDirectory.ValidFileName?(`c:\work\`) is: false)
		Assert(CheckDirectory.ValidFileName?(`c:\work\NUL`) is: false)
		Assert(CheckDirectory.ValidFileName?(`c:\work\NUL.txt`) is: false)
		Assert(CheckDirectory.ValidFileName?(`c:\work\NULlified.txt`))
		}

	Test_deleteFile?()
		{
		fn = CheckDirectory.CheckDirectory_deleteFile?
		fileDate = #20230814.150012345
		startTime = #20230814.160012345
		filePrefix = 'suneido_dir_test_'

		file = 'suneido_dir_test_20230814_151210000'
		Assert(fn(file, fileDate, startTime, filePrefix))
		Assert(fn(file, startTime, fileDate, filePrefix) is: false)

		file = 'suneido_dir_test_20230814_151210001'
		Assert(fn(file, fileDate, startTime, filePrefix))
		Assert(fn(file, startTime, fileDate, filePrefix) is: false)

		file = 'suneido_dir_test_20230814_1511'
		Assert(fn(file, fileDate, startTime, filePrefix) is: false)

		file = 'suneido_dir_test_20230814_1512345'
		Assert(fn(file, fileDate, startTime, filePrefix) is: false)

		file = 'suneido_dir_test_20230814_15121000012'
		Assert(fn(file, fileDate, startTime, filePrefix) is: false)

		file = 'suneido_dir_test_20230814_1512100001_'
		Assert(fn(file, fileDate, startTime, filePrefix) is: false)

		file = 'suneido_dir_test_20230814_151210000.txt'
		Assert(fn(file, fileDate, startTime, filePrefix) is: false)
		}

	Test_ReviewPaths()
		{
		mock = Mock(CheckDirectory)
		mock.CheckDirectory_table = table = .MakeTable('(cd_path, cd_TS) key (cd_path)')
		mock.When.FilePrefix().Return(prefix = 'suneido_dir_test_')
		mock.When.dirExists?([anyArgs:]).Return(true)
		mock.When.dir([anyArgs:]).Return([])
		mock.When.deleteFileApi([anyArgs:]).Return(false)
		mock.When.ReviewPaths().CallThrough()

		// Empty "checked_directories" table
		Assert(mock.ReviewPaths() is: 'SUCCEEDED')
		mock.Verify.Never().deleteOutstanding([anyArgs:])

		// Empty directory
		QueryOutput(table, [cd_path: `T:\FakeDir1`])
		Assert(QueryCount(table) is: 1)
		Assert(mock.ReviewPaths() is: 'SUCCEEDED')
		mock.Verify.deleteOutstanding(`T:\FakeDir1`, [isDate:], prefix)
		Assert(QueryCount(table) is: 0)

		// Directory has files, fails to delete all of them
		QueryOutput(table, [cd_path: `T:\FakeDir1`])
		mock.When.dir([anyArgs:]).Return([
			[name: 'suneido_dir_test_20230814_151210000', date: Date().Minus(minutes: 3)],
			[name: 'suneido_dir_test_20230815_151210000', date: Date().Minus(minutes: 3)]
			])
		Assert(QueryCount(table) is: 1)
		result = mock.ReviewPaths().Split('\r\n')
		Assert(result[0] is: 'WARNING')
		Assert(result[1] is: `- Path: T:\FakeDir1, Days Stuck: 0`)
		mock.Verify.Times(2).deleteOutstanding(`T:\FakeDir1`, [isDate:], prefix)
		Assert(QueryCount(table) is: 1)

		// Directory has files, fails to delete all of them, stuck several days
		mock.When.startTime().Return(Date().Plus(days: 3))
		mock.When.dir([anyArgs:]).Return([
			[name: 'suneido_dir_test_20230814_151210000', date: Date().Minus(minutes: 3)],
			[name: 'suneido_dir_test_20230815_151210000', date: Date().Minus(minutes: 3)]
			],
			[
			[name: 'suneido_dir_test_20230815_151210000', date: Date().Minus(minutes: 3)]
			])
		mock.When.deleteFileApi([anyArgs:]).Return(true, false)
		Assert(QueryCount(table) is: 1)
		result = mock.ReviewPaths().Split('\r\n')
		Assert(result[0] is: 'FAILED')
		Assert(result[1] is: `- Path: T:\FakeDir1, Days Stuck: 3`)
		mock.Verify.Times(3).deleteOutstanding(`T:\FakeDir1`, [isDate:], prefix)
		Assert(QueryCount(table) is: 1)

		// Directory has files, deletes all of them
		mock.When.startTime().Return(Date().Plus(days: 3))
		mock.When.dir([anyArgs:]).Return([
			[name: 'suneido_dir_test_20230814_151210000', date: Date().Minus(minutes: 3)],
			[name: 'suneido_dir_test_20230815_151210000', date: Date().Minus(minutes: 3)]
			],
			[])
		mock.When.deleteFileApi([anyArgs:]).Return(true)
		Assert(QueryCount(table) is: 1)
		Assert(mock.ReviewPaths() is: 'SUCCEEDED')
		mock.Verify.Times(4).deleteOutstanding(`T:\FakeDir1`, [isDate:], prefix)
		Assert(QueryCount(table) is: 0)

		// DirExists returns false (directory is potentially inaccessible or deleted).
		// path record is deleted (as we can no longer interact properly with the path)
		QueryOutput(table, [cd_path: `T:\FakeDir1`])
		mock.When.startTime().Return(Date())
		mock.When.dir([anyArgs:]).Return([
			[name: 'suneido_dir_test_20230814_151210000', date: Date().Minus(minutes: 3)],
			])
		mock.When.dirExists?([anyArgs:]).Return(false)
		Assert(QueryCount(table) is: 1)
		Assert(mock.ReviewPaths() is: 'SUCCEEDED')
		mock.Verify.Times(5).deleteOutstanding(`T:\FakeDir1`, [isDate:], prefix)
		Assert(QueryCount(table) is: 0)
		}
	}