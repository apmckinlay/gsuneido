// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
MatcherWas
	{
	Match(value, args)
		{
		m = args.Members()[0]
		v = args[m]
		return Object?(value) and value.Member?(m) and value[m] is v
		}
	Expected(args)
		{
		return "the value to be an object containing " $ Display(args)
		}
	}