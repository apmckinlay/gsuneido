// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		s = "Now is the time
			for all good men
			to come to the aid"
		file = .MakeFile(s)
		Assert(GetFile(file) is: s)
		Assert(GetFile(file, 10) is: "Now is the")
		Assert(GetFile(file, 1000) is: s)
		AddFile(file, "bye")
		Assert(GetFile(file) is: s $ "bye")
		}
	}
