// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	c: class
		{
		New() { .Log = "" }
		Write(s) { .Log $= s }
		}
	Test_main()
		{
		.test("", #(),
			"Date: <date>\r\nContent-Length: 0\r\n\r\n")
		.test("", #('Content-Type': 'text/plain')
			"Date: <date>\r\nContent-Length: 0\r\nContent-Type: text/plain\r\n\r\n")
		.test("hello world", #(),
			"Date: <date>\r\nContent-Length: 11\r\n\r\nhello world")
		.test("hello", #('Content-Type': 'text/plain'),
			"Date: <date>\r\n" $
			"Content-Length: 5\r\n" $
			"Content-Type: text/plain\r\n\r\nhello")
		.test("", #('Content-Length': 6),
			"Date: <date>\r\n" $
			"Content-Length: 6\r\n" $
			"\r\n")
		.test("", #('Content_Length': 7),
			"Date: <date>\r\n" $
			"Content-Length: 7\r\n" $
			"\r\n")
		.test("", #('Set_Cookie': ['a=1', 'b=2']),
			"Date: <date>\r\n" $
			"Content-Length: 0\r\n" $
			"Set-Cookie: a=1\r\n" $
			"Set-Cookie: b=2\r\n" $
			"\r\n")
		}
	test(content, header, result)
		{
		x = new .c
		HttpSend(x, "HTTP/1.0", header, content)
		x.Log = x.Log.Replace("^Date: .*", 'Date: <date>')
		Assert(x.Log is: "HTTP/1.0\r\n" $ result)
		}
	}