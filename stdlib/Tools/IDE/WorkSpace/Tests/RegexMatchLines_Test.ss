// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		s = "now is the time
			for all good men
			to come to the aid
			of their party"
		Assert(RegexMatchLines(s, 'x') is: #())
		Assert(RegexMatchLines(s, 'n') is: #((0), (1)))
		Assert(RegexMatchLines(s, 'the') is: #((0), (2), (3)))
		}
	}