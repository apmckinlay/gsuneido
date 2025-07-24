// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(Unprivatize("Name.Name_priv") is: "Name.priv")
		Assert(Unprivatize("Name.priv") is: "Name.priv")
		Assert(Unprivatize("Name.Public") is: "Name.Public")
		Assert(Unprivatize("123.456_789") is: "123.456_789")
		Assert(Unprivatize("") is: "")
		}
	}