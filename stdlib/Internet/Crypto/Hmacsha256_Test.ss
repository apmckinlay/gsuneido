// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(Hmacsha256("", "").ToHex()
			is: "b613679a0814d9ec772f95d778c35fc5ff1697c493715653c6c712144292c5ad")
		Assert(Hmacsha256("The quick brown fox jumps over the lazy dog",
			"key").ToHex()
			is: "f7bc83f430538424b13298e6aa6fb143ef4d59a14946175997479dbc2d1a3cd8")
		Assert(Hmacsha256("Another test for hmac256.\r\nWith a second line",
			"different_key").ToHex()
			is: "12ce0627c8a3f0bf07b762c85f14545a9f5fa6bec85e83e843c67f006066d614")

		// key is bigger than block size
		Assert(Hmacsha256("This is a test",
			"abc123abc123abc123abc123abc123abc123abc123abc123abc123abc123abc123").ToHex()
			is: "216ac10a80834eefeee448d1ed27a2ff571a7cd13b62b1f41764e5f807d9cd37")
		}
	}