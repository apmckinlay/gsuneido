// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		fn = CurlOrInternalError?
		Assert(fn('curl: (79) Could not open remote file for reading: Operation failed'),
			msg: 'could not open remote file')
		Assert(fn('Http.POST failed: HTTP/1.1 500 Internal Server Error'),
			msg: 'internal server error')
		Assert(fn('Invalid HTTP response code in: (empty header)'), msg: 'empty header')
		Assert(fn('Invalid HTTP response code in: ...') is: false, msg: '...')
		Assert(fn('Program Error') is: false, msg: 'program error')
		Assert(fn('curl: Http failed to get header file content'),
			msg: 'header file content')
		Assert(fn('Http.POST failed: HTTP/1.1 504 GATEWAY_TIMEOUT'), msg: '504')
		Assert(fn('Http.POST failed: HTTP/1.133 5046 NOT VALID ERROR') is: false,
			msg: 'not valid error')
		}
	}