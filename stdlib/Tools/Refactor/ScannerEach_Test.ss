// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		src = 'function (arg) { return arg + 1 }'
		scan = Scanner(src)
		ScannerEach(src)
			{|prev2/*unused*/, prev/*unused*/, token, next/*unused*/|
			scan.Next2()
			Assert(token is: scan.Text())
			}
		Assert(scan.Next2() is: scan)
		}
	}