// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(Opt() is: "")
		Assert(Opt("") is: "")
		Assert(Opt("", "one") is: "")
		Assert(Opt("", "one", "two", "three") is: "")
		Assert(Opt("one", "two", "three", "") is: "")
		Assert(Opt("one", "", "two", "three") is: "")
		Assert(Opt("", "", "", "") is: "")
		Assert(Opt("abc") is: "abc")
		Assert(Opt("abc", "def") is: "abcdef")
		Assert(Opt("abc", "def", "ghi") is: "abcdefghi")
		}
	}