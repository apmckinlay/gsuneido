// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_WriteHeader()
		{
		.test({|unused| }, 'HTTP/1.0 200 OK\r\n')
		.test({|r| r.ResponseCode(404) }, 'HTTP/1.0 404 Not Found\r\n')
		.test({|r| r.Date(#20130531) }, 'HTTP/1.0 200 OK\r\n' $
			'Date: Fri, 31 May 2013 00:00:00 GMT\r\n')
		}
	c: class
		{
		New() { .Log = "" }
		Write(s) { .Log $= s }
		Writeline(s) { .Log $= s $ '\r\n' }
		}
	test(block, result)
		{
		response = HttpResponse()
		block(response)
		x = new .c
		response.WriteHeader(x)
		Assert(x.Log is: result)
		}
	}