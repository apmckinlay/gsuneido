// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(Md5("").ToHex() is: "d41d8cd98f00b204e9800998ecf8427e")
		Assert(Md5("\x00\xff").ToHex() is: "d07d34efac6328007ad67c7e0a985e00")

		cksum = "fc5e038d38a57032085441e7fe7010b0"
		Assert(Md5("helloworld").ToHex() is: cksum)
		Assert(Md5("hello", "world").ToHex() is: cksum)
		Assert(Md5().Update("hello").Update("world").Value().ToHex() is: cksum)
		}
	}