// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
GetRes
	{
	CallClass(env)
		{
		args = env.queryvalues
		if not args.Member?(#name) or not String?(args.name)
			return .ReturnWithDelay()

		name = args.name
		if not GetContributions('ResCode').Has?(name)
			return .ReturnWithDelay()

		lib = args.GetDefault('lib', 'stdlib')
		if not Libraries().Has?(lib)
			return .ReturnWithDelay()

		if false is rec = Query1Cached(lib, :name)
			return .ReturnWithDelay()

		headers = .Headers(name)
		return ['OK', headers, rec.text]
		}
	}