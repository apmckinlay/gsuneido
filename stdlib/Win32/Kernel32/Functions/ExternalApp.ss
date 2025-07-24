// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Memoize
	{
	Func(app)
		{
		return .findApp(app, FindProgram)
		}

	findApp(app, func)
		{
		for contrib in GetContributions('ExternalApp').Reverse!()
			if contrib.app is app
				{
				result = func(contrib.path)
				if result isnt false and result isnt 0
					return result
				}
		return false
		}
	}