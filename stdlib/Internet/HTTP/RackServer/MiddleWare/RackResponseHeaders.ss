// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
RackComposeBase
	{
	Call(env)
		{
		if -1 is result = .App(:env)
			return result
		result = RackServer.ResultOb(result)
		headers = Object()
		for field in result[1].Members()
			HttpResponse.ResponseHeaderField(field, result[1][field], headers)
		.setDefaultResponseHeaders(headers)
		return [result[0], headers, result[2]]
		}

	setDefaultResponseHeaders(headers)
		{
		date = .date()
		.setIfNotExist(headers, 'Expires', date)
		.setIfNotExist(headers, 'Last-Modified', date)
		.setIfNotExist(headers, 'Cache-Control', 'no-cache')
		}

	// for testing
	date()
		{
		return Date()
		}

	setIfNotExist(headers, field, value)
		{
		if not headers.Member?(field)
			HttpResponse.ResponseHeaderField(field, value, headers)
		}
	}
