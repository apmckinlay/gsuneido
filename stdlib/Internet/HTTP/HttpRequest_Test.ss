// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		r = HttpRequest(s = 'GET /Wiki?page HTTP/1.1')
		Assert(r.Request() is: s)
		Assert(r.Method() is: 'GET')
		Assert(r.Path() is: '/Wiki')
		Assert(r.Query() is: 'page')
		}
	Test_split()
		{
		r = HttpRequest.SplitRequestLine('GET /Wiki?page HTTP/1.1')
		Assert(r.method is: 'GET')
		Assert(r.path is: '/Wiki')
		Assert(r.query is: 'page')
		Assert(r.version is: 1.1)

		err = Catch()
			{
			r = HttpRequest.SplitRequestLine('GET /Wiki?page HTTP/1.1^M')
			}
		Assert(err is HttpRequest.BadRequest)
		}
	Test_default()
		{
		r = HttpRequest('')
		r.Headers.authorization = 'authorization'
		Assert(r.Authorization() is: 'authorization')
		r.Authorization('test')
		Assert(r.Authorization() is: 'test')
		}
	}
