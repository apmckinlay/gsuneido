// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		test = { |n|
			Assert(StrToInt(IntToStr(n)) is: n)
			}
		test(0)
		test(1)
		test(300)
		test(999932785)
		}
	}