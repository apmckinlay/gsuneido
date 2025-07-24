// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		check = function (s)
			{
			for k in #('k', 'key', 'abracadabra')
				{
				t = StringXor(s, k)
				Assert(t.Size() is: s.Size())
				Assert(StringXor(t, k) is: s)
				}
			}
		check("")
		check("x")
		check("\x00\xff")
		check("hello world")
		check("now is the time for all good men to come to the aid of their party")
		}
	}