// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		cl = RackResponseHeaders
			{
			RackResponseHeaders_date()
				{
				return #20170727.010203
				}
			}
		app = function(env)
			{
			if not env.Member?(#testFields)
				return 'No testFields'
			return ['200 OK', env.testFields, ""]
			}
		mw = new cl(app)

		// only default headers
		date = #20170727.010203
		Assert(mw(Object()) is: ['200 OK', Object(
			'Expires': date.InternetFormat(),
			'Last-Modified': date.InternetFormat(),
			'Cache-Control': 'no-cache'), 'No testFields'])

		// default should not override the specified
		expireDate = #20170727.050505
		Assert(mw(Object(testFields: Object(
			'Cache_Control': 'public',
			'Expires': expireDate,
			'Connection': 'keep-alive'))),
			is: ['200 OK', Object(
				'Expires': expireDate.InternetFormat(),
				'Last-Modified': date.InternetFormat(),
				'Cache-Control': 'public',
				'Connection': 'keep-alive'), ''])

		// invalid header
		Assert({ mw(Object(testFields: Object(invalidField: true))) }
			throws: "HttpResponse: method not found: invalidField")
		}
	}
