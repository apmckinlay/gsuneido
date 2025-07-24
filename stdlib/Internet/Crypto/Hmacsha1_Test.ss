// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(Hmacsha1("", "").ToHex()
			is: "fbdb1d1b18aa6c08324b7d64b71fb76370690e1d")
		Assert(Hmacsha1("The quick brown fox jumps over the lazy dog",
			"key").ToHex()
			is: "de7c9b85b8b78aa6bc8a7a36f70a90701c9db4d9")
		}
	}