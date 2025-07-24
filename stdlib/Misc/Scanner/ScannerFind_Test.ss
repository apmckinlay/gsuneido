// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(ScannerFind("", "") is: 0)
		Assert(ScannerFind("x", "x") is: 0)
		Assert(ScannerFind("x", "y") is: 1)
		Assert(ScannerFind("'x' x", "x") is: 4)
		Assert(ScannerFind("sort_field", "sort") is: 10)
		Assert(ScannerFind("sort,field", "sort") is: 0)
		}
	}