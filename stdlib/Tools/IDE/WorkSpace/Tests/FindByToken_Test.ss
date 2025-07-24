// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ForEachMatch()
		{
		positions = function (code, find)
			{
			list = Object()
			FindByToken.ForEachMatch(code, find, {|f, t| list.Add([f, t]) })
			return list
			}

		list = positions("x.Sort!(); y.\nSort! () /* z.Sort!() */", #('.', 'Sort!', '('))
		Assert(list is: #((1, 8), (12, 21)))

		list = positions("a a a b", #(a, a, b))
		Assert(list is: #((2, 7)))

		list = positions("a b c", #(b))
		Assert(list is: #((2, 3)))
		}
	}
