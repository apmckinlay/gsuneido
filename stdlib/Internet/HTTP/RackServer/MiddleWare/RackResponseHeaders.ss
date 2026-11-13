// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
RackComposeBase
	{
	Call(env)
		{
		if -1 is result = .App(:env)
			return result
		code, headers, content = RackServer.Result(result)
		newHeaders = Object()
		for field in headers.Members()
			HttpResponse.ResponseHeaderField(field, headers[field], newHeaders)
		.setDefaultResponseHeaders(newHeaders)
		return [code, newHeaders, content]
		}

	setDefaultResponseHeaders(headers)
		{
		date = .date()
		.setIfNotExist(headers, #Expires, date)
		.setIfNotExist(headers, "Last-Modified", date)
		.setIfNotExist(headers, "Cache-Control", "no-cache")
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
