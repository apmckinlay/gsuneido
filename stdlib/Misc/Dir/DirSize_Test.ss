// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.rootDir = .MakeDir()
		.PutFile(.rootDir $ "/f1", "This is the first file, in dirsize_test")

		.depth2 = .MakeDir(.rootDir)
		.PutFile(.depth2 $ "/f2",
			"This is the second file that must be recursivly found in Depth2")
		.PutFile(.depth2 $ "/f3",
			"This is the third file that is being added to f2 in Depth2")
		}
	Test_main()
		{
		f1Size = FileSize(.rootDir $ "/f1")
		f2Size = FileSize(.depth2  $ "/f2")
		f3Size = FileSize(.depth2  $ "/f3")
		folder2Size = f2Size + f3Size

		test =
			{
			Assert(DirSize(.rootDir) is: f1Size + folder2Size,
				msg: 'f1Size + folder2Size')
			Assert(DirSize(.depth2) is: folder2Size, msg: 'depth2 is folder2Size')
			}
		try
			test()
		catch
			test() // retry since it fails intermittently

		//Remove dirsize_test's files so now dirsize_test size should equal Depth2's size
		DeleteFile(.rootDir $ "/f1")
		Assert(DirSize(.rootDir) is: folder2Size, msg: 'rootDir is folder2Size')
		}
	}