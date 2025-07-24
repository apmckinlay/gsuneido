// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(Sha1("").ToHex() is: "da39a3ee5e6b4b0d3255bfef95601890afd80709")
		Assert(Sha1("\x00\xff").ToHex() is: "aa3e5dcdd77b153f2e59bd0d8794fde33cb4e486")

		cksum = "6adfb183a4a2c94a2f92dab5ade762a47889a5a1"
		Assert(Sha1("helloworld").ToHex() is: cksum)
		Assert(Sha1("hello", "world").ToHex() is: cksum)
		Assert(Sha1().Update("hello").Update("world").Value().ToHex() is: cksum)
		}
	}