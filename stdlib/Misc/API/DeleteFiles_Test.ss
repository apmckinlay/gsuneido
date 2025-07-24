// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
// TAGS: client
Test
	{
	Test_main()
		{
		dir = .MakeDir()
		subDir = .MakeDir(dir)
		.PutFile(dir $ '/testfile.txt', '')
		.AddTeardown({ DeleteFile(dir $ '/testfile.txt') })
		DeleteFiles(dir $ '/*.*')
		Assert(not FileExists?(dir $ '/testfile.txt'))
		Assert(DirExists?(subDir))
		}
	}
