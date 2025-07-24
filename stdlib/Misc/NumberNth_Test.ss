// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(NumberNth(0) is: '0th')
		Assert(NumberNth(1) is: '1st')
		Assert(NumberNth(2) is: '2nd')
		Assert(NumberNth(3) is: '3rd')
		Assert(NumberNth(4) is: '4th')
		Assert(NumberNth(9) is: '9th')
		Assert(NumberNth(11) is: '11th')
		Assert(NumberNth(12) is: '12th')
		Assert(NumberNth(13) is: '13th')
		Assert(NumberNth(14) is: '14th')
		Assert(NumberNth(21) is: '21st')
		Assert(NumberNth(22) is: '22nd')
		Assert(NumberNth(23) is: '23rd')
		Assert(NumberNth(24) is: '24th')
		}
	}