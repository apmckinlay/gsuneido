// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		files = [
			.MakeFile(''),
			.MakeFile('File 1 Content'),
			.MakeFile('File 2 Content'),
			.MakeFile('File 1 Content with Extra texts')
			]
		for f1 in files
			for f2 in files
				Assert(CompareFiles?(f1, f2, 3) is: f1 is f2)
		}
	}