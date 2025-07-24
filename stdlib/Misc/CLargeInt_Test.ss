// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(CLargeInt(0, 0) is: 0)
		Assert(CLargeInt(1, 1) is: 4294967297)
		Assert(CLargeInt(1, 100) is: 429496729601)
		Assert(CLargeInt(-1, -1) is: 1.844674407370955e19)
		Assert(CLargeInt(4.Gb(), 4.Gb()) is: 1.844674407800452e19)
		}
	}