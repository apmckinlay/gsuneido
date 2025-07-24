// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		try
			{
			.Outer()
			throw "shouldn't get here"
			}
		catch (e)
			{
			Assert(Display(e.Callstack()[0].fn), has: 'Inner')
			}
		}
	Outer()
		{
		try
			.Inner()
		catch (e)
			throw e // rethrow
		}
	Inner()
		{
		throw "boom"
		}
	}