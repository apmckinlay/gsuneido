// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		s = "for (i = 0; i < 10; ++i)"
		Assert(ScannerFindIf(s, { it is "i" }) is: 5)
		Assert(ScannerFindIf(s, { it is 123 }) is: s.Size())
		}
	}