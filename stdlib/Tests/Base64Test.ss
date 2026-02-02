// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	data: (
		("", "")
		("a", "YQ==")
		("aa", "YWE=")
		("aaa", "YWFh")
		("aaaaaa", "YWFhYWFh")
		("hello world", "aGVsbG8gd29ybGQ=")
		("\x00z", "AHo=")
		("a\x00z", "YQB6")
		("ab\x00z", "YWIAeg==")
		("abc\x00z", "YWJjAHo=")
		("\xff\xff\xff", "////")
		)
	Test_main()
		{
		for x in .data
			{
			Assert(Base64.Encode(x[0]) is: x[1])
			Assert(Base64.Decode(x[1]) is: x[0])
			}
		Assert(Base64.Decode("YQ") is: "a")

		Assert(Base64.EncodeLines('x'.Repeat(100), linelen: 70)
			is: "eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eH
h4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eA==
")

Assert(Base64.EncodeLines('x'.Repeat(100))
			is: "eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eH"$
			"h4\r\neHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eA==\r\n")
		}
	}