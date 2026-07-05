// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(env)
		{
		if .invalidToken?(env)
			return Object('Unauthorized', [], 'Your session is invalid or expired')

		env.path = env.path.AfterFirst('/rally')
		router = RackRouter(GetContributions('RackRoutes').
			Filter({ it.GetDefault(#rally, false) }))
		return router(env)
		}

	invalidToken?(env)
		{
		return JsSessionToken.Validate(env) is false
		}
	}