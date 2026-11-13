// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// server side function that handles pack/unpack
	CallClass(func, env, auth? = false)
		{
		args = .getArgsFromEnv(env, func)
		if auth? is true and Auth.UseToken(args.GetDefault(#runOnHttpToken, "")) is false
			return ["401", #(), ""]

		result = func(@args)
		return Pack(result)
		}

	getArgsFromEnv(env, func)
		{
		if false is body = env.GetDefault(#body, { env.GetDefault(#entity_body, false) })
			return #()
		try
			result = Unpack(body)
		catch (err)
			{
			SuneidoLog("ERROR: (CAUGHT) Unable to Unpack body: " $ err,
				params: [env, func])
			return #()
			}
		return result
		}
	}
