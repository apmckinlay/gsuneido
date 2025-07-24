// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
RackComposeBase
	{
	Call(env)
		{
		if -1 is result = .App(:env)
			return result
		result = RackServer.ResultOb(result)
		if result[1].Member?(#Content_Type)
			return result
		ext = env.path.AfterLast('.')
		if false isnt type = MimeTypes.GetDefault(ext, false)
			result = [result[0],
				result[1].Copy().Add(type at: #Content_Type), result[2]]
		return result
		}
	}
