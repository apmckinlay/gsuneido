// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
RackComposeBase
	{
	Call(env)
		{
		if -1 is result = .App(:env)
			return result
		ext = env.path.AfterLast('.')
		if false is type = MimeTypes.GetDefault(ext, false)
			return result
		if Object?(result) and result.Size() >= 2 and result[1].Member?(#Content_Type)
			return result
		code, headers, content = RackServer.Result(result)
		if not headers.Member?(#Content_Type)
			headers = headers.Copy().Add(type at: #Content_Type)
		return [code, headers, content]
		}
	}
