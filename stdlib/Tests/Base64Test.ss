// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	data: (
		("", ""),
		(a, "YQ=="),
		(aa, "YWE="),
		(aaa, YWFh),
		(aaaaaa, YWFhYWFh),
		("hello world", "aGVsbG8gd29ybGQ="),
		("\x00z", "AHo="),
		("a\x00z", YQB6),
		("ab\x00z", "YWIAeg=="),
		("abc\x00z", "YWJjAHo="),
		("\xff\xff\xff", "////")
		)
	Test_main()
		{
		for x in .data
			{
			Assert(Base64.Encode(x[0]) is: x[1])
			Assert(Base64.Decode(x[1]) is: x[0])
			}
		Assert(Base64.Decode(#YQ) is: 'a')
		Assert(Base64.Decode(Base64.EncodeUrl('a')) is: 'a')
		Assert(Base64.Decode(Base64.EncodeUrl(#aa)) is: #aa)
		Assert(Base64.Decode(Base64.EncodeUrl(#aaa)) is: #aaa)
		Assert(Base64.Decode(Base64.EncodeUrl("\xff\xff\xff")) is: "\xff\xff\xff")
		Assert(Base64.EncodeUrl('a') is: #YQ)
		Assert(Base64.EncodeUrl(#aa) is: #YWE)
		Assert(Base64.EncodeUrl(#aaa) is: #YWFh)
		Assert(Base64.EncodeUrl("hello world") is: #aGVsbG8gd29ybGQ)
		Assert(Base64.EncodeUrl("\xff\xff\xff") is: "____")

		Assert(Base64.EncodeLines('x'.Repeat(100), linelen: 60)
			is: "eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4\r\n" $
				"eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4\r\n" $
				"eHh4eHh4eHh4eA==\r\n")

		Assert(Base64.EncodeLines('x'.Repeat(100))
			is: "eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4" $
				"eHh4\r\n" $
				"eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eA==\r\n")
		}
	}
