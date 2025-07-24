// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
MatcherWas
	{
	Match(value, args)
		{
		return Type(value) is args
		}
	Expected(args)
		{
		return 'a type of ' $ args
		}
	}