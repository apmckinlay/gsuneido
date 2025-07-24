// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(GetLastWriteTime('doest exist') is: false)

		startDate = Date().NoTime()
		testFile = .MakeFile('test')
		lastModTime = GetLastWriteTime(testFile)
		Assert(lastModTime greaterThanOrEqualTo: startDate)

		// because linux file last modified time does not have milliseconds
		Thread.Sleep(1000)

		AddFile(testFile, 'added some text')
		Assert(GetLastWriteTime(testFile) greaterThan: lastModTime)
		}
	}
