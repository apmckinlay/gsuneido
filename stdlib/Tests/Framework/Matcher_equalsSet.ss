// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
MatcherWas
	{
	Match(value, args)
		{
		return value.EqualSet?(args)
		}
	Actual(value, args)
		{
		that = args.Copy()
		missing = Object()
		for x in value
			if false isnt i = that.Find(x)
				that.Delete(i)
			else
				missing.Add(x)
		maxchars = 1000
		msg = Object()
		if not that.Empty?()
			msg.Add("missing: " $ that.Join(', ').Ellipsis(maxchars))
		if not missing.Empty?()
			msg.Add("extra: " $ missing.Join(', ').Ellipsis(maxchars))
		return "had " $ msg.Join('; ')
		}
	Expected(args)
		{
		return "the value to be same set as " $ Display(args)
		}
	}
