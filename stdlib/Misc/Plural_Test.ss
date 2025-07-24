// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(Plural(0, "dog") is: "0 dogs")
		Assert(Plural(1, "dog") is: "1 dog")
		Assert(Plural(3, "dog") is: "3 dogs")
		}
	}