// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_WORD()
		{
		for lo in #(0, 1, 5, 123, 0xffff)
			for hi in #(0, 1, 5, 123, 0xffff)
				{
				x = MAKELONG(lo, hi)
				Assert(LOWORD(x) is: lo)
				Assert(HIWORD(x) is: hi)
				}
		}
	Test_SWORD()
		{
		for lo in #(-5, -1, 0, 1, 5, 0x7fff)
			for hi in #(-5, -1, 0, 1, 5, 0x7fff)
				{
				x = MAKELONG(lo, hi)
				Assert(LOSWORD(x) is: lo)
				Assert(HISWORD(x) is: hi)
				}
		x = 0xffff_ffde_1234
		Assert(LOSWORD(x) is: 0x1234)
		Assert(HISWORD(x) is: -34)
		}
	}