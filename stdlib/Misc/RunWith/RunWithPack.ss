// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// server side function that handles pack/unpack
	CallClass(func, env)
		{
		args = .getArgsFromEnv(env, func)
		result = func(@args)
		return Pack(result)
		}

	getArgsFromEnv(env, func)
		{
		if false is body = env.GetDefault(#body, { env.GetDefault(#entity_body, false) })
			return #()
		try
			{
			result = Unpack(body)
			}
		catch (err)
			{
			SuneidoLog('ERROR: (CAUGHT) Unable to Unpack body: ' $ err,
				params: Object(env, func))
			return #()
			}
		return result
		}
	}
