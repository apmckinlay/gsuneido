// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_findClosestNonBlankChar()
		{
		fn = Addon_stepping_debugger.Addon_stepping_debugger_findClosestNonBlankChar
		s = "   \r\nfunction ()\r\n\t{\r\n\t1 + 1     \r\n\t}"
		Assert(fn(s, 0) is: false)
		Assert(fn(s, 5) is: 5)

		Assert(fn(s, 36) is: 36)
		Assert(fn(s, 37) is: false)

		Assert(fn(s, 16) is: 15)
		Assert(fn(s, 22) is: 23)
		Assert(fn(s, 31) is: 27)
		Assert(fn(s, 24) is: 25)
		}
	}
