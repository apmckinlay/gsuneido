// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(CommaList() is: "")
		Assert(CommaList("") is: "")
		Assert(CommaList("one") is: "one")
		Assert(CommaList("one", "two", "three") is: "one, two, three")
		Assert(CommaList("one", "", "three") is: "one, three")
		}
	}