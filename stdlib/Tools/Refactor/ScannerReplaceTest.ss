// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(ScannerReplace("", "", "") is: "")
		Assert(ScannerReplace("hello", "hell", "heck") is: "hello")
		Assert(ScannerReplace("hello", "hello", "howdy") is: "howdy")
		Assert(ScannerReplace("hello world", "hello", "howdy") is: "howdy world")
		Assert(ScannerReplace("hello world", "world", "there") is: "hello there")
		Assert(ScannerReplace("my hello world", "hello", "howdy") is: "my howdy world")
		Assert(ScannerReplace("my 'hello' world", "hello", "howdy")
			is: "my 'hello' world")
		Assert(ScannerReplace("my /*hello*/ world", "hello", "howdy")
			is: "my /*hello*/ world")
		}
	Test_multi()
		{
		Assert(
			ScannerReplace("[ 'string{}[]' { a: 1} ]",
				['[', ']', '{', '}'], ['#(', ')', '#(', ')'])
			is: "#( 'string{}[]' #( a: 1) )")
		Assert({ ScannerReplace('', [1], [1, 2]) }
			throws: "Assert FAILED:")
		}
	}