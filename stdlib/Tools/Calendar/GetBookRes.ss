// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
GetRes
	{
	CallClass(env)
		{
		args = env.queryvalues
		if not args.Member?(#name) or not String?(args.name)
			return .ReturnWithDelay()

		book = args.GetDefault('book', 'imagebook')
		if not GetContributions('ResBook').Has?(book)
			return .ReturnWithDelay()

		if env.GetDefault('If-Modified-Since', false) is #20000101.InternetFormat()
			return ['NotModified', #(), '']

		res = GetBookText(args.name, book)
		if res is false
			return .ReturnWithDelay()

		headers = .Headers(args.name)
		return ['OK', headers, res]
		}
	}