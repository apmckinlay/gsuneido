// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		s = "hello world".Repeat(100)
		c = Zlib.Compress(s)
		Assert(c.Size() lessThan: s.Size())
		s2 = Zlib.Uncompress(c)
		Assert(s2 is: s)
		}
	}