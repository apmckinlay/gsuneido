// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ResponseCode()
		{
		s = ''
		Assert({ Http.ResponseCode(s) }, throws: "(empty header)")

		s = '\r\n'
		Assert({ Http.ResponseCode(s) }, throws: "(empty header)")

		s = '\r\n\r\n'
		Assert({ Http.ResponseCode(s) }, throws: "(empty header)")

		s = '\r\nContent-Length: 50'
		Assert({ Http.ResponseCode(s) }, throws: "\r\nContent-Length: 50")

		s = 'Content-Length: 50\r\n'
		Assert({ Http.ResponseCode(s) }, throws: "Content-Length: 50")

		s = 'HTTP/1.1 200 OK\r\n'
		Assert(Http.ResponseCode(s) is: '200')

		s = 'HTTP/1.0 505 Internal Server Error\r\n'
		Assert(Http.ResponseCode(s) is: '505')

		s = 'HTTP/1.1 505 Internal Server Error\r\n'
		Assert(Http.ResponseCode(s) is: '505')
		}
	}
