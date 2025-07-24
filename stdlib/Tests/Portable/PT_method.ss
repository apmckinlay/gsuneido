// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(@args)
		{
		str? = args.Extract(#str?)
		ob = .toValue(str?, args, 0)
		method = args[1]
		expected = .toValue(str?, args, args.Size() - 1)
		argvals = Object()
		for (i = 2; i < args.Size() - 1; ++i)
			argvals.Add(.toValue(str?, args, i))
		result = ob[method](@argvals)
		ok = result is expected
		if not ok
			Print("\tgot: " + Display(result))
		return ok
		}
	toValue(str?, args, i)
		{
		return str?[i] ? args[i] : args[i].Compile()
		}
	}
