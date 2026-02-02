// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_checkForDups()
		{
		checkForDups = Qc_DupCodeContinuousTests.Qc_DupCodeContinuousTests_checkForDups
		oldHash = #(00001: 'testlib:File1:8',
			00002: 'testlib:File2:9',
			00003: 'testlib:File3:4, testlib:File4:2',
			00004: 'testlib:File3:5, testlib:File5:5')
		newHash = #(00001: 'testlib:File1:8',
			00002: 'testlib:File2:9',
			00003: 'testlib:File3:4, testlib:File4:2',
			00004: 'testlib:File3:5, testlib:File5:5')
		newDups = Object()
		checkForDups(oldHash, newHash, newDups)
		Assert(newDups.Empty?())

		newHash = #(00001: 'testlib:File1:8',
			00002: 'testlib:File2:9, testlib:File6:14',
			00003: 'testlib:File3:4, testlib:File4:2, testlib:File6:20',
			00004: 'testlib:File3:5',
			00005: 'testlib:File6:14, testlib:File7:18')
		newDups = Object()
		checkForDups(oldHash, newHash, newDups)
		Assert(newDups isSize: 3)
		Assert(newDups has: 'testlib:File2:9, testlib:File6:14')
		Assert(newDups has: 'testlib:File3:4, testlib:File4:2, testlib:File6:20')
		Assert(newDups has: 'testlib:File6:14, testlib:File7:18')

		// check if just the line numbers changed
		oldHash = #(00001: 'testlib:File1:8, testlib:File2:16')
		newHash = #(00002: 'testlib:File1:9, testlib:File2:17')
		newDups = Object()
		checkForDups(oldHash, newHash, newDups)
		Assert(newDups.Empty?())

		// ensure we dont flag a removal as a new match
		oldHash = #(00001: 'testlib:File1:8, testlib:File2:16')
		newHash = #(00001: 'testlib:File1:9')
		newDups = Object()
		checkForDups(oldHash, newHash, newDups)
		Assert(newDups.Empty?())

		// ensure we dont flag a removal as a new match
		oldHash = #(00001: 'testlib:File1:8, testlib:File2:16')
		newHash = #(00001: 'testlib:File1:8, testlib:File2:16, testlib:File3:32')
		newDups = Object()
		checkForDups(oldHash, newHash, newDups)
		Assert(newDups isSize: 1)
		Assert(newDups has: 'testlib:File1:8, testlib:File2:16, testlib:File3:32')
		}
	}