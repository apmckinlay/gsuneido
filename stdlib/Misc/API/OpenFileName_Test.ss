// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_multipleFiles()
		{
		multi = OpenFileName.OpenFileName_multipleFiles
		s = `C:\Users` $ '\x00a.jpg\x00b.jpg\x00\x00'
		Assert(multi(s) is: #(`C:\Users/a.jpg`, `C:\Users/b.jpg`))

		s = `C:\Users` $ '\x00a.jpg\x00\x00\x00'
		Assert(multi(s) is: #(`C:\Users/a.jpg`))
		}
	}